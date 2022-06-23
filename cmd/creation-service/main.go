package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"terrarium-grpc-gateway/internal/services"

	"google.golang.org/grpc"
)

type CreationService struct {
	services.UnimplementedCreatorServer
}

func (s *CreationService) SetupModule(ctx context.Context, request *services.SetupModuleRequest) (*services.SetupModuleResponse, error) {
	return nil, nil
}

func main() {
	fmt.Println("Welcome to Terrarium GRPC API Gateway")
	flag.Parse()
	lis, err := net.Listen("tcp", "0.0.0.0:9443")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	creationServiceServer := &CreationService{}
	// Need TLS
	grpcServer := grpc.NewServer(opts...)
	services.RegisterCreatorServer(grpcServer, creationServiceServer)
	grpcServer.Serve(lis)
}