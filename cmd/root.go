package cmd

import (
	"context"
	"fmt"
	"github.com/terrariumcloud/terrarium/internal/storage"
	"log"
	"net"
	"net/http"
	"os"

	"go.opentelemetry.io/otel/propagation"

	"github.com/terrariumcloud/terrarium/internal/common/grpc_service"
	"github.com/terrariumcloud/terrarium/internal/restapi"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	noop "go.opentelemetry.io/otel/trace"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	buildVersion string
)

const (
	defaultEndpoint = "0.0.0.0:3001"
)

var (
	endpoint            = defaultEndpoint
	awsSessionConfig    = storage.AWSSessionConfig{}
	opentelemetryInited = false
	rootCmd             = &cobra.Command{
		Use:   "terrarium",
		Short: "Terrarium Services",
		Long:  "Runs backend that exposes Terrarium Services",
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", defaultEndpoint, "Endpoint")
	rootCmd.PersistentFlags().StringVarP(&awsSessionConfig.Key, "aws-access-key-id", "k", "", "AWS Access Key")
	rootCmd.PersistentFlags().StringVarP(&awsSessionConfig.Secret, "aws-secret-access-key", "s", "", "AWS Secret Key")
	rootCmd.PersistentFlags().StringVarP(&awsSessionConfig.Region, "aws-region", "r", "", "AWS Region")
	rootCmd.PersistentFlags().BoolVar(&awsSessionConfig.UseLocalStack, "use-localstack", false, "Connect to a localstack instance rather than AWS.")
}

func newTraceExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	if otelEndpoint, found := os.LookupEnv("OTEL_EXPORTER_OTLP_ENDPOINT"); found {
		fmt.Printf("Created exporter")
		client := otlptracegrpc.NewClient(
			otlptracegrpc.WithEndpoint(otelEndpoint),
			otlptracegrpc.WithInsecure(),
		)

		exp, err := otlptrace.New(ctx, client)
		if err != nil {
			return nil, err
		}
		return exp, nil
	} else {
		return nil, nil
	}
}

func newServiceResource(name string) *resource.Resource {
	versionInfo := buildVersion
	if serviceVersion, found := os.LookupEnv("OTEL_SERVICE_VERSION"); found {
		log.Println("Warning: build time version overriden by environment variable")
		versionInfo = serviceVersion
	}

	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(name),
		semconv.ServiceVersionKey.String(versionInfo),
	)
	return resources
}

func initOpenTelemetry(name string) func() {
	opentelemetryInited = true

	ctx := context.Background()

	traceExporter, err := newTraceExporter(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if traceExporter != nil {
		tracerProvider := trace.NewTracerProvider(
			trace.WithSampler(trace.AlwaysSample()),
			trace.WithBatcher(traceExporter),
			trace.WithResource(newServiceResource(name)),
		)
		otel.SetTracerProvider(tracerProvider)
		return func() {
			if err := tracerProvider.Shutdown(context.Background()); err != nil {
				log.Fatal(err)
			}
		}
	} else {
		otel.SetTracerProvider(noop.NewNoopTracerProvider())
	}
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return func() {
		// No-op
	}
}

func startGRPCService(name string, service grpc_service.Service) {
	log.Printf("Starting %s", name)
	if !opentelemetryInited {
		otelShutdown := initOpenTelemetry(name)
		defer otelShutdown()
	}

	listener, err := net.Listen("tcp4", endpoint)
	if err != nil {
		log.Fatalf("Failed to start: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	if err := service.RegisterWithServer(grpcServer); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}

	log.Printf("Listening at %s", endpoint)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed: %v", err)
	}
}

func startRESTAPIService(name, mountPath string, rootHandler restapi.RESTAPIHandler) {
	log.Printf("Starting %s", name)
	if !opentelemetryInited {
		otelShutdown := initOpenTelemetry(name)
		defer otelShutdown()
	}

	log.Printf("Listening on %s", endpoint)
	if err := http.ListenAndServe(endpoint, rootHandler.GetHttpHandler(mountPath)); err != nil {
		log.Fatalf("Failed: %v", err)
	}

}

// Execute root command
func Execute(serviceVersion string) {
	buildVersion = serviceVersion
	cobra.CheckErr(rootCmd.Execute())
}
