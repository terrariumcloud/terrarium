package cmd

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/terrariumcloud/terrarium/internal/storage"

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

func getEndPointOptions() []otlptracegrpc.Option {
	for _, envVar := range []string{"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "OTEL_EXPORTER_OTLP_ENDPOINT"} {
		if _, found := os.LookupEnv(envVar); found {
			// No need for options let the SDK pull out the values from the Environment variables
			return []otlptracegrpc.Option{}
		}
	}
	return []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint("jaeger:4317"),
		otlptracegrpc.WithInsecure(),
	}
}

func newTraceExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	options := getEndPointOptions()
	client := otlptracegrpc.NewClient(options...)
	return otlptrace.New(ctx, client)
}

func newServiceResource(name string) *resource.Resource {
	versionInfo := buildVersion
	if serviceVersion, found := os.LookupEnv("OTEL_SERVICE_VERSION"); found {
		log.Println("Warning: build time version overriden by environment variable")
		versionInfo = serviceVersion
	}

	if serviceName, found := os.LookupEnv("OTEL_SERVICE_NAME"); found {
		log.Println("Warning: service name overriden by environment variable")
		name = serviceName
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
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
		return func() {
			if err := tracerProvider.Shutdown(context.Background()); err != nil {
				log.Fatal(err)
			}
		}
	}

	otel.SetTracerProvider(noop.NewNoopTracerProvider())
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
