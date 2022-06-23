package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"terrarium-grpc-gateway/internal/services"

	"google.golang.org/grpc"
)

const (
	moduleCreationServiceEndpoint          = "module_creation:3000"
	connectionToModuleCreationServiceError = "Internal server error, unable to connect to the module creation service"
)

type TerrariumGrpcGateway struct {
	pb.UnimplementedPublisherServer
	pb.UnimplementedConsumerServer
}

func (s *TerrariumGrpcGateway) Configure(ctx context.Context, request *pb.ModuleConfigurationRequest) (*pb.TransactionStatusResponse, error) {
	conn, err := grpc.Dial(moduleCreationServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return &pb.TransactionStatusResponse{
			Status:        pb.Status_UNKNOWN_ERROR,
			StatusMessage: connectionToModuleCreationServiceError,
		}, nil
	}
	defer conn.Close()
	client := services.NewCreatorClient(conn)

	delegatedRequest := services.SetupModuleRequest{
		Name:        request.GetName(),
		Description: request.GetDescription(),
		SourceUrl:   request.GetSourceUrl(),
		Maturity:    request.GetMaturity(),
	}
	if response, delegateError := client.SetupModule(ctx, &delegatedRequest); delegateError != nil {
		log.Printf("SetupModule remote call failed: %v", delegateError)
		return &pb.TransactionStatusResponse{
			Status:        pb.Status_UNKNOWN_ERROR,
			StatusMessage: "Failed to execute SetupModule",
		}, nil
	} else {
		return &pb.TransactionStatusResponse{
			Status:        response.GetStatus(),
			StatusMessage: response.GetStatusMessage(),
		}, nil
	}
}

func (s *TerrariumGrpcGateway) RegisterModuleDependencies(ctx context.Context, request *pb.RegisterModuleDependenciesRequest) (*pb.TransactionStatusResponse, error) {
	return &pb.TransactionStatusResponse{
		Status:        pb.Status_OK,
		StatusMessage: "All is good",
	}, nil
}

func (s *TerrariumGrpcGateway) RegisterContainerDependencies(ctx context.Context, request *pb.RegisterContainerDependenciesRequest) (*pb.TransactionStatusResponse, error) {
	return &pb.TransactionStatusResponse{
		Status:        pb.Status_OK,
		StatusMessage: "All is good",
	}, nil
}

func (s *TerrariumGrpcGateway) UploadSourceZip(server pb.Publisher_UploadSourceZipServer) error {
	return nil
}

func (s *TerrariumGrpcGateway) BeginVersion(ctx context.Context, request *pb.BeginVersionRequest) (*pb.BeginVersionResponse, error) {
	// TODO add support for authentication
	return &pb.BeginVersionResponse{
		SessionKey: "1234",
	}, nil
}

func (s *TerrariumGrpcGateway) EndVersion(ctx context.Context, request *pb.EndVersionRequest) (*pb.TransactionStatusResponse, error) {
	return &pb.TransactionStatusResponse{
		Status:        pb.Status_OK,
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
	// Need TLS
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterPublisherServer(grpcServer, gatewayServer)
	pb.RegisterConsumerServer(grpcServer, gatewayServer)
	grpcServer.Serve(lis)
}
