package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/terrariumcloud/terrarium/internal/restapi"

	"github.com/terrariumcloud/terrarium/internal/module/services"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
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
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(name)),
	)
	if err != nil {
		log.Fatalf("The SchemaURL of the resources is not merged.")
	}
	versionInfo := buildVersion
	if serviceVersion, found := os.LookupEnv("OTEL_SERVICE_VERSION"); found {
		log.Println("Warning: build time version overriden by environment variable")
		versionInfo = serviceVersion

	}
	res, err = resource.Merge(
		res,
		resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceVersionKey.String(versionInfo)),
	)
	if err != nil {
		log.Fatalf("The SchemaURL of the resources is not merged.")
	}

	return res
}

func startGRPCService(name string, service services.Service) {

	log.Printf("Starting %s", name)

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
		defer func() {
			if err := tracerProvider.Shutdown(context.Background()); err != nil {
				log.Fatal(err)
			}
		}()
		otel.SetTracerProvider(tracerProvider)
	} else {
		otel.SetTracerProvider(noop.NewNoopTracerProvider())
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
		defer func() {
			if err := tracerProvider.Shutdown(context.Background()); err != nil {
				log.Fatal(err)
			}
		}()
		otel.SetTracerProvider(tracerProvider)
	} else {
		otel.SetTracerProvider(noop.NewNoopTracerProvider())
	}

	log.Println(fmt.Printf("Listening on %s", endpoint))
	log.Fatal(http.ListenAndServe(endpoint, rootHandler.GetHttpHandler(mountPath)))
}

// Execute root command
func Execute(serviceVersion string) {
	buildVersion = serviceVersion
	cobra.CheckErr(rootCmd.Execute())
}
