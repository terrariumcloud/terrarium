package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/terrariumcloud/terrarium/internal/restapi"

	services "github.com/terrariumcloud/terrarium/internal/module/services"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk_tracy "go.opentelemetry.io/otel/sdk/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

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
// 	exporter, err := stdoutTrace.New(stdoutTrace.WithPrettyPrint())
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

func newExporter(w io.Writer) (sdk_tracy.SpanExporter, error) {
	log.Printf("Returning newExporter")
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		// Do not print timestamps for the demo.
	)
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

	f, err := os.Create("traces.txt")
	if err != nil {
		log.Fatal(err)
		log.Printf("Failed 117")
	}
	defer f.Close()

	exp, err := newExporter(f)
	if err != nil {
		log.Printf("Failed 122")
		log.Fatal(err)
	}

	tp := sdk_tracy.NewTracerProvider(
		sdk_tracy.WithSyncer(exp),
		sdk_tracy.WithResource(newResource(name)),
	)
	defer func() {
		log.Printf("Failed 132")
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	log.Printf("Will execute TP")
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
	log.Println(fmt.Sprintf("Listening on %s", endpoint))
	log.Fatal(http.ListenAndServe(endpoint, rootHandler.GetHttpHandler(mountPath)))

}

// Execute root command
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
