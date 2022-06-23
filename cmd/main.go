package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"terrarium-grpc-gateway/internal/services"

	pb "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium"
	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"
)

const (
	moduleCreationServiceEndpoint          = "module_creation:3000"
	moduleSessionManagerServiceEndpoint          = "module_session_manager:3001"
	moduleStorageServiceEndpoint          = "module_storage:3002"
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

func (s *TerrariumGrpcGateway) BeginVersion(ctx context.Context, request *pb.BeginVersionRequest) (*pb.BeginVersionResponse, error) {
	// TODO add support for authentication
	conn, err := grpc.Dial(moduleSessionManagerServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return nil, err
	}
	defer conn.Close()
	client := services.NewSessionManagerClient(conn)

	delegatedRequest := services.BeginVersionRequest{
		Module: request.GetModule(),
	}
	if response, delegateError := client.BeginVersion(ctx, &delegatedRequest); delegateError != nil {
		log.Printf("BeginVersion remote call failed: %v", delegateError)
		return nil, err
	} else {
		return response, nil
	}
}

func (s *TerrariumGrpcGateway) EndVersion(ctx context.Context, request *pb.EndVersionRequest) (*pb.TransactionStatusResponse, error) {
	// TODO add support for authentication
	conn, err := grpc.Dial(moduleSessionManagerServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return &pb.TransactionStatusResponse{
			Status:        pb.Status_UNKNOWN_ERROR,
			StatusMessage: connectionToModuleCreationServiceError,
		}, nil
	}
	defer conn.Close()
	client := services.NewSessionManagerClient(conn)

	delegatedRequest := services.TerminateVersionRequest{
		SessionKey: request.GetSessionKey(),
	}
	if request.GetAction() == pb.EndVersionRequest_DISCARD {
		response, delegateError := client.AbortVersion(ctx, &delegatedRequest)
		if delegateError != nil {
			log.Printf("AbortVersion remote call failed: %v", delegateError)
			return nil, delegateError
		} else {
			return response, nil
		}
	}
	if response, delegateError := client.PublishVersion(ctx, &delegatedRequest); delegateError != nil {
		log.Printf("PublishVersion remote call failed: %v", delegateError)
		client.AbortVersion(ctx, &delegatedRequest)
		return nil, err
	} else {
		return response, nil
	}
}

func (s *TerrariumGrpcGateway) UploadSourceZip(server pb.Publisher_UploadSourceZipServer) error {
	conn, err := grpc.Dial(moduleStorageServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return err
	}
	defer conn.Close()
	client := services.NewStorageClient(conn)
	uploadStream, err := client.UploadSourceZip(context.TODO())
	if err != nil {
		return err
	}
	for {
		chunk, err := server.Recv()
		if err == io.EOF {
			uploadStream.CloseSend()
			return server.SendAndClose(&pb.TransactionStatusResponse{
				Status:        pb.Status_OK,
				StatusMessage: "All is good",
			})
		}
		if err != nil {
			return err
		}
		err = uploadStream.Send(chunk)
		if err != nil {
			server.SendAndClose(&pb.TransactionStatusResponse{
				Status:        pb.Status_UNKNOWN_ERROR,
				StatusMessage: "Upload failed",
			})
			return err
		}
	}
}

func (s *TerrariumGrpcGateway) DownloadSourceZip(request *pb.DownloadSourceZipRequest, server pb.Consumer_DownloadSourceZipServer) error {
	conn, err := grpc.Dial(moduleStorageServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return err
	}
	defer conn.Close()
	client := services.NewStorageClient(conn)
	downloadStream, err := client.DownloadSourceZip(context.TODO(), &pb.DownloadSourceZipRequest{
		ApiKey: request.GetApiKey(),
		Module: request.GetModule(),
	})
	if err != nil {
		return err
	}
	for {
		chunk, err := downloadStream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = server.Send(chunk)
		if err != nil {
			downloadStream.CloseSend()
			return err
		}
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
