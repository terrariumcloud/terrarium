package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/terrariumcloud/terrarium/internal/restapi"

	"github.com/terrariumcloud/terrarium/internal/module/services"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	defaultEndpoint = "0.0.0.0:3001"
)

var endpoint = defaultEndpoint
var awsAccessKey string
var awsSecretKey string
var awsRegion string

var rootCmd = &cobra.Command{
	Use:   "terrarium",
	Short: "Terrarium Services",
	Long:  "Runs backend that exposes Terrarium Services",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", defaultEndpoint, "Endpoint")
	rootCmd.PersistentFlags().StringVarP(&awsAccessKey, "aws-access-key-id", "k", "", "AWS Access Key (required)")
	rootCmd.MarkPersistentFlagRequired("aws-access-key-id")
	rootCmd.PersistentFlags().StringVarP(&awsSecretKey, "aws-secret-access-key", "s", "", "AWS Secret Key (required)")
	rootCmd.MarkPersistentFlagRequired("aws-secret-access-key")
	rootCmd.PersistentFlags().StringVarP(&awsRegion, "aws-region", "r", "", "AWS Region (required)")
	rootCmd.MarkPersistentFlagRequired("aws-region")
}

var (
	lsEndpoint    = "ingest.lightstep.com:443"
	lsToken       = ""
	lsEnvironment = "dev"
)

func newExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	var headers = map[string]string{
		"lightstep-access-token": lsToken,
	}
	fmt.Printf("Created exporter")
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithHeaders(headers),
		otlptracegrpc.WithEndpoint(lsEndpoint),
	)

	exp, err := otlptrace.New(ctx, client)
	if err != nil {
		log.Fatal(err)
	}

	return exp, nil
}

func newResource(name string) *resource.Resource {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name),
			semconv.ServiceVersionKey.String("1.0.2"),
			attribute.String("environment", lsEnvironment),
		),
	)
	if err != nil {
		log.Fatalf("The SchemaURL of the resources is not merged.")
	}
	log.Printf("Returning newResource")
	return r
}

func startGRPCService(name string, service services.Service) {

	log.Printf("Starting %s", name)

	ctx := context.Background()

	traceExporter, err := newExporter(ctx)
	if err != nil {
		log.Fatal(err)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(traceExporter),
		trace.WithResource(newResource(name)),
	)
	defer func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	otel.SetTracerProvider(tracerProvider)

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

	ctx := context.Background()

	traceExporter, err := newExporter(ctx)
	if err != nil {
		log.Fatal(err)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(traceExporter),
		trace.WithResource(newResource(name)),
	)
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	otel.SetTracerProvider(tracerProvider)

	log.Println(fmt.Printf("Listening on %s", endpoint))
	log.Fatal(http.ListenAndServe(endpoint, rootHandler.GetHttpHandler(mountPath)))
}

// Execute root command
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
