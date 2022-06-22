package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	pb "terrarium-grpc-gateway/pkg/terrarium"

	"google.golang.org/grpc"
)

type TerrariumGrpcGateway struct {
	pb.UnimplementedPublisherServer
	pb.UnimplementedConsumerServer
}

func (s *TerrariumGrpcGateway) Configure(ctx context.Context, request *pb.ModuleConfigurationRequest) (*pb.TransactionStatusResponse, error) {
	return &pb.TransactionStatusResponse{
		Status: pb.TransactionStatusResponse_OK,
		StatusMessage: "All is good",
	}, nil
}

func (s *TerrariumGrpcGateway) BeginVersion(ctx context.Context, request *pb.BeginVersionRequest) (*pb.BeginVersionResponse, error) {
	return &pb.BeginVersionResponse{
		SessionKey: "1234",
	}, nil
}

func (s *TerrariumGrpcGateway) RegisterModuleDependencies(ctx context.Context, request *pb.RegisterModuleDependenciesRequest) (*pb.TransactionStatusResponse, error) {
	return &pb.TransactionStatusResponse{
		Status: pb.TransactionStatusResponse_OK,
		StatusMessage: "All is good",
	}, nil
}

func (s *TerrariumGrpcGateway) RegisterContainerDependencies(ctx context.Context, request *pb.RegisterContainerDependenciesRequest) (*pb.TransactionStatusResponse, error) {
	return &pb.TransactionStatusResponse{
		Status: pb.TransactionStatusResponse_OK,
		StatusMessage: "All is good",
	}, nil
}

func (s *TerrariumGrpcGateway) UploadSourceZip(server pb.Publisher_UploadSourceZipServer) error {
	return nil
}

func (s *TerrariumGrpcGateway) EndVersion(ctx context.Context, request *pb.EndVersionRequest) (*pb.TransactionStatusResponse, error) {
	return &pb.TransactionStatusResponse{
		Status: pb.TransactionStatusResponse_OK,
		StatusMessage: "All is good",
	}, nil
}

func (s *TerrariumGrpcGateway) DownloadSourceZip(request *pb.DownloadSourceZipRequest, server pb.Consumer_DownloadSourceZipServer) error {
	return nil
}

func (s *TerrariumGrpcGateway) RetrieveContainerDependencies(request *pb.RetrieveContainerDependenciesRequest, server pb.Consumer_RetrieveContainerDependenciesServer) error {
	return nil
}

func (s *TerrariumGrpcGateway) RetrieveModuleDependencies(request *pb.RetrieveModuleDependenciesRequest, server pb.Consumer_RetrieveModuleDependenciesServer) error {
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

	gatewayServer := &TerrariumGrpcGateway{}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterPublisherServer(grpcServer, gatewayServer)
	pb.RegisterConsumerServer(grpcServer, gatewayServer)
	grpcServer.Serve(lis)
}
