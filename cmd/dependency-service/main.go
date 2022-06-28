package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"terrarium-grpc-gateway/internal/services"
	"github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium"

	"google.golang.org/grpc"
)

type DependencyService struct {
	services.UnimplementedDependencyResolverServer
}

func (s *DependencyService) RegisterModuleDependencies(ctx context.Context, request *terrarium.RegisterModuleDependenciesRequest) (*terrarium.TransactionStatusResponse, error) {
	return nil, nil
}

func (s *DependencyService) RegisterContainerDependencies(ctx context.Context, request *terrarium.RegisterContainerDependenciesRequest) (*terrarium.TransactionStatusResponse, error) {
	return nil, nil
}

func (s *DependencyService) RetrieveContainerDependencies(request *terrarium.RetrieveContainerDependenciesRequest, server services.DependencyResolver_RetrieveContainerDependenciesServer) error {
	return nil
}

func (s *DependencyService) RetrieveModuleDependencies(request *terrarium.RetrieveModuleDependenciesRequest, server services.DependencyResolver_RetrieveModuleDependenciesServer) error {
	return nil
}

func main() {
	fmt.Println("Welcome to Terrarium GRPC API Gateway")
	flag.Parse()
	lis, err := net.Listen("tcp", "0.0.0.0:9443")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	dependencyServiceServer := &DependencyService{}
	// Need TLS
	grpcServer := grpc.NewServer(opts...)
	services.RegisterDependencyResolverServer(grpcServer, dependencyServiceServer)
	grpcServer.Serve(lis)
}