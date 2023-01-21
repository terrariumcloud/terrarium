package services

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Service interface {
	RegisterWithServer(grpcServer grpc.ServiceRegistrar) error
}

// createGRPCConnection takes an endpoint and returns a grpc connection
func CreateGRPCConnection(target string) (*grpc.ClientConn, error) {
	return grpc.Dial(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
}
