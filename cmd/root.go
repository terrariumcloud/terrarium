package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/terrariumcloud/terrarium/internal/restapi"

	services "github.com/terrariumcloud/terrarium/internal/module/services"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk_tracy "go.opentelemetry.io/otel/sdk/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const OtelTracerName = "terrarium-tracer"

const (
	defaultEndpoint = "0.0.0.0:3001"
)

var endpoint string = defaultEndpoint
var awsAccessKey string
var awsSecretKey string
var awsRegion string

var rootCmd = &cobra.Command{
	Use:   "terrarium",
	Short: "Terrarium Services",
	Long:  "Runs backend that exposes Terrarium Services",
}

// func InitOTELTracer() (*sdktrace.TracerProvider, error) {
// 	exporter, err := otlptrace.New(otlptrace.WithPrettyPrint())
// 	if err != nil {
// 		return nil, err
// 	}
// 	tp := sdktrace.NewTracerProvider(
// 		sdktrace.WithSampler(sdktrace.AlwaysSample()),
// 		sdktrace.WithBatcher(exporter),
// 	)
// 	otel.SetTracerProvider(tp)
// 	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

// 	return tp, nil
// }

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
	lsEndpoint    = ""
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
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "dev"),
		),
	)
	log.Printf("Returning newResource")
	return r
}

func startGRPCService(name string, service services.Service) {

	log.Printf("Starting %s", name)

	ctx := context.Background()

	exp, err := newExporter(ctx)
	if err != nil {
		log.Fatal(err)
	}

	tp := sdk_tracy.NewTracerProvider(
		sdk_tracy.WithSyncer(exp),
		// sdk_tracy.WithBatcher(exp),
		sdk_tracy.WithResource(newResource(name)),
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	otel.SetTracerProvider(tp)

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

	exp, err := newExporter(ctx)
	if err != nil {
		log.Fatal(err)
	}

	tp := sdk_tracy.NewTracerProvider(
		sdk_tracy.WithSyncer(exp),
		// sdk_tracy.WithBatcher(exp),
		sdk_tracy.WithResource(newResource(name)),
	)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	otel.SetTracerProvider(tp)

	log.Println(fmt.Printf("Listening on %s", endpoint))

	log.Fatal(http.ListenAndServe(endpoint, rootHandler.GetHttpHandler(mountPath)))

}

// Execute root command
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
