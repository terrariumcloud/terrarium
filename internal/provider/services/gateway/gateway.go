package gateway

import (
	"context"
	"log"

	"github.com/terrariumcloud/terrarium/internal/provider/services"
	terrarium "github.com/terrariumcloud/terrarium/pkg/terrarium/provider"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ConnectToVersionManagerError     = status.Error(codes.Unavailable, "Failed to connect to Version manager service.")
	UnknownVersionManagerActionError = status.Error(codes.InvalidArgument, "Unknown Version manager action requested.")
)

type TerrariumGrpcGateway struct {
	terrarium.UnimplementedPublisherServer
	versionManagerClient services.VersionManagerClient
}

func New(versionManagerClient services.VersionManagerClient) *TerrariumGrpcGateway {
	return &TerrariumGrpcGateway{
		versionManagerClient: versionManagerClient,
	}
}

// RegisterWithServer registers TerrariumGrpcGateway with grpc server
func (gw *TerrariumGrpcGateway) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	terrarium.RegisterPublisherServer(grpcServer, gw)
	return nil
}

// Register new provider with Registrar service
func (gw *TerrariumGrpcGateway) Register(ctx context.Context, request *terrarium.RegisterProviderRequest) (*terrarium.Response, error) {
	return gw.RegisterWithClient(ctx, request, gw.versionManagerClient)
}

// RegisterWithClient calls Register on version manager client
func (gw *TerrariumGrpcGateway) RegisterWithClient(ctx context.Context, request *terrarium.RegisterProviderRequest, client services.VersionManagerClient) (*terrarium.Response, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("provider.name", request.GetName()),
		attribute.String("provider.version", request.GetVersion()),
	)

	if res, delegateError := client.Register(ctx, request); delegateError != nil {
		log.Printf("Failed: %v", delegateError)
		span.RecordError(delegateError)
		return nil, delegateError
	} else {
		log.Println("Done <= Registered")
		span.AddEvent("Successful call to Register on Version Manager client.")
		return res, nil
	}
}

func (gw *TerrariumGrpcGateway) EndProvider(ctx context.Context, request *terrarium.EndProviderRequest) (*terrarium.Response, error) {
	return gw.EndProviderWithClient(ctx, request, gw.versionManagerClient)
}

// EndVersionWithClient calls AbortProvider/AbortProviderVersion/PublishVersion on Version Manager client
func (gw *TerrariumGrpcGateway) EndProviderWithClient(ctx context.Context, request *terrarium.EndProviderRequest, client services.VersionManagerClient) (*terrarium.Response, error) {

	terminateProviderVersion := services.TerminateVersionRequest{
		Provider: request.GetProvider(),
	}

	publishRequest := services.TerminateVersionRequest{
		Provider: request.GetProvider(),
	}

	span := trace.SpanFromContext(ctx)

	if request.GetAction() == terrarium.EndProviderRequest_DISCARD_VERSION {
		log.Println("Abort Provider Version => Version Manager")
		if res, delegateError := client.AbortProviderVersion(ctx, &terminateProviderVersion); delegateError != nil {
			log.Printf("Failed: %v", delegateError)
			span.RecordError(delegateError)
			return nil, delegateError
		} else {
			log.Println("Done <= Version Manager")
			span.AddEvent("Successfully aborted provider version", trace.WithAttributes(attribute.String("Provider Name", request.Provider.GetName()), attribute.String("Provider Version", request.Provider.GetVersion())))
			return res, nil
		}
	} else if request.GetAction() == terrarium.EndProviderRequest_PUBLISH {
		log.Println("Pubish Provider => Version Manager")
		if res, delegateError := client.PublishVersion(ctx, &publishRequest); delegateError != nil {
			log.Printf("Failed: %v", delegateError)
			span.RecordError(delegateError)
			return nil, delegateError
		} else {
			log.Println("Done <= Version Manager")
			span.AddEvent("Successfully published version", trace.WithAttributes(attribute.String("Provider Name", request.Provider.GetName()), attribute.String("Provider Version", request.Provider.GetVersion())))
			return res, nil
		}
	} else {
		log.Printf("Unknown Version manager action requested: %v", request.GetAction())
		span.AddEvent("Unknown action.")
		return nil, UnknownVersionManagerActionError
	}
}
