package gateway

import (
	"context"
	"io"
	"log"

	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/services"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/services/creation"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/services/dependency"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/services/session"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/services/storage"

	pb "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium"
	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"
)

type TerrariumGrpcGateway struct {
	pb.UnimplementedPublisherServer
	pb.UnimplementedConsumerServer
}

func (s *TerrariumGrpcGateway) Configure(ctx context.Context, request *pb.ModuleConfigurationRequest) (*pb.TransactionStatusResponse, error) {
	conn, err := grpc.Dial(creation.CreationServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return Error("Internal server error, unable to connect to the module creation service"), nil
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
		return Error("Failed to execute SetupModule"), nil
	} else {
		return &pb.TransactionStatusResponse{
			Status:        response.GetStatus(),
			StatusMessage: response.GetStatusMessage(),
		}, nil
	}
}

func (s *TerrariumGrpcGateway) BeginVersion(ctx context.Context, request *pb.BeginVersionRequest) (*pb.BeginVersionResponse, error) {
	conn, err := grpc.Dial(session.SessionServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	conn, err := grpc.Dial(session.SessionServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return Error("Internal server error, unable to connect to the module session service"), nil
	}
	defer conn.Close()
	client := services.NewSessionManagerClient(conn)

	delegatedRequest := services.TerminateVersionRequest{
		SessionKey: request.GetSessionKey(),
	}
	if request.GetAction() == pb.EndVersionRequest_DISCARD {
		if response, delegateError := client.AbortVersion(ctx, &delegatedRequest); delegateError != nil {
			log.Printf("AbortVersion remote call failed: %v", delegateError)
			return nil, delegateError
		} else {
			return response, nil
		}
	} else if request.GetAction() == pb.EndVersionRequest_PUBLISH {
		if response, delegateError := client.PublishVersion(ctx, &delegatedRequest); delegateError != nil {
			log.Printf("PublishVersion remote call failed: %v", delegateError)
			return nil, delegateError
		} else {
			if _, delegateError = client.AbortVersion(ctx, &delegatedRequest); delegateError != nil {
				log.Print("Failed to clean up session after version published.")
				return nil, delegateError
			}
			return response, nil
		}
	} else {
		return Error("Unknown action requested."), nil
	}
}

func (s *TerrariumGrpcGateway) UploadSourceZip(server pb.Publisher_UploadSourceZipServer) error {
	conn, err := grpc.Dial(storage.StorageServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
			return server.SendAndClose(Ok("Archive uploaded successfully."))
		}
		if err != nil {
			return err
		}
		err = uploadStream.Send(chunk)
		if err != nil {
			server.SendAndClose(Error("Upload failed."))
			return err
		}
	}
}

func (s *TerrariumGrpcGateway) DownloadSourceZip(request *pb.DownloadSourceZipRequest, server pb.Consumer_DownloadSourceZipServer) error {
	conn, err := grpc.Dial(storage.StorageServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	conn, err := grpc.Dial(dependency.DependencyServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return Error("Internal server error, unable to connect to the module dependency service"), nil
	}
	defer conn.Close()
	client := services.NewDependencyResolverClient(conn)
	response, err := client.RegisterModuleDependencies(ctx, request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *TerrariumGrpcGateway) RegisterContainerDependencies(ctx context.Context, request *pb.RegisterContainerDependenciesRequest) (*pb.TransactionStatusResponse, error) {
	conn, err := grpc.Dial(dependency.DependencyServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return Error("Internal server error, unable to connect to the module dependency service"), nil
	}
	defer conn.Close()
	client := services.NewDependencyResolverClient(conn)
	response, err := client.RegisterContainerDependencies(ctx, request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *TerrariumGrpcGateway) RetrieveContainerDependencies(request *pb.RetrieveContainerDependenciesRequest, server pb.Consumer_RetrieveContainerDependenciesServer) error {
	conn, err := grpc.Dial(dependency.DependencyServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return err
	}
	defer conn.Close()
	client := services.NewDependencyResolverClient(conn)
	dependencyStream, err := client.RetrieveContainerDependencies(context.TODO(), request)
	if err != nil {
		return err
	}
	for {
		chunk, err := dependencyStream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = server.Send(chunk)
		if err != nil {
			dependencyStream.CloseSend()
			return err
		}
	}
}

func (s *TerrariumGrpcGateway) RetrieveModuleDependencies(request *pb.RetrieveModuleDependenciesRequest, server pb.Consumer_RetrieveModuleDependenciesServer) error {
	conn, err := grpc.Dial(dependency.DependencyServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return err
	}
	defer conn.Close()
	client := services.NewDependencyResolverClient(conn)
	dependencyStream, err := client.RetrieveModuleDependencies(context.TODO(), request)
	if err != nil {
		return err
	}
	for {
		chunk, err := dependencyStream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = server.Send(chunk)
		if err != nil {
			dependencyStream.CloseSend()
			return err
		}
	}
}

func Ok(message string) *pb.TransactionStatusResponse {
	return &pb.TransactionStatusResponse{
		Status:        pb.Status_OK,
		StatusMessage: message,
	}
}

func Error(message string) *pb.TransactionStatusResponse {
	return &pb.TransactionStatusResponse{
		Status:        pb.Status_UNKNOWN_ERROR,
		StatusMessage: message,
	}
}
