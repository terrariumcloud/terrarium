package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "terrarium-grpc-gateway/pkg/terrarium"
)

type TerrariumGrpcGateway struct {
	pb.UnimplementedPublisherServer
	pb.UnimplementedConsumerServer
}

func (s *TerrariumGrpcGateway) Configure(ctx context.Context, request *pb.ModuleConfigurationRequest) (*pb.TransactionStatusResponse, error) {


	//TODO implement me
	panic("implement me")
}

func (s *TerrariumGrpcGateway) BeginVersion(ctx context.Context, request *pb.BeginVersionRequest) (*pb.BeginVersionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *TerrariumGrpcGateway) RegisterModuleDependencies(ctx context.Context, request *pb.RegisterModuleDependenciesRequest) (*pb.TransactionStatusResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *TerrariumGrpcGateway) RegisterContainerDependencies(ctx context.Context, request *pb.RegisterContainerDependenciesRequest) (*pb.TransactionStatusResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *TerrariumGrpcGateway) UploadSourceZip(server pb.Publisher_UploadSourceZipServer) error {
	//TODO implement me
	panic("implement me")
}

func (s *TerrariumGrpcGateway) EndVersion(ctx context.Context, request *pb.EndVersionRequest) (*pb.TransactionStatusResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *TerrariumGrpcGateway) DownloadSourceZip(request *pb.DownloadSourceZipRequest, server pb.Consumer_DownloadSourceZipServer) error {
	//TODO implement me
	panic("implement me")
}

func (s *TerrariumGrpcGateway) RetrieveContainerDependencies(request *pb.RetrieveContainerDependenciesRequest, server pb.Consumer_RetrieveContainerDependenciesServer) error {
	//TODO implement me
	panic("implement me")
}

func (s *TerrariumGrpcGateway) RetrieveModuleDependencies(request *pb.RetrieveModuleDependenciesRequest, server pb.Consumer_RetrieveModuleDependenciesServer) error {
	//TODO implement me
	panic("implement me")
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
