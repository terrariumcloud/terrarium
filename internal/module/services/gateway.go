package services

import (
	"context"
	"io"
	"log"

	pb "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"
	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"
)

type TerrariumGrpcGateway struct {
	pb.UnimplementedPublisherServer
	pb.UnimplementedConsumerServer
}

func (s *TerrariumGrpcGateway) Register(ctx context.Context, request *pb.RegisterModuleRequest) (*pb.TransactionStatusResponse, error) {
	conn, err := grpc.Dial(RegistrarServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return FailedToConnectToRegistrar, nil // should return err?
	}
	defer conn.Close()
	client := NewRegistrarClient(conn)

	delegatedRequest := RegisterModuleRequest{
		Name:        request.GetName(),
		Description: request.GetDescription(),
		SourceUrl:   request.GetSourceUrl(),
		Maturity:    request.GetMaturity(),
	}
	if response, delegateError := client.Register(ctx, &delegatedRequest); delegateError != nil {
		log.Printf("Register module remote call failed: %v", delegateError)
		return FailedToExecuteRegister, nil // should return delegateError?
	} else {
		return &pb.TransactionStatusResponse{
			Status:        response.GetStatus(),
			StatusMessage: response.GetStatusMessage(),
		}, nil
	}
}

func (s *TerrariumGrpcGateway) BeginVersion(ctx context.Context, request *pb.BeginVersionRequest) (*pb.BeginVersionResponse, error) {
	conn, err := grpc.Dial(SessionServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return nil, err
	}
	defer conn.Close()
	client := NewVersionManagerClient(conn)

	delegatedRequest := BeginVersionRequest{
		Module: request.GetModule(),
	}
	if response, delegateError := client.BeginVersion(ctx, &delegatedRequest); delegateError != nil {
		log.Printf("BeginVersion remote call failed: %v", delegateError)
		return nil, delegateError
	} else {
		return response, nil
	}
}

func (s *TerrariumGrpcGateway) EndVersion(ctx context.Context, request *pb.EndVersionRequest) (*pb.TransactionStatusResponse, error) {
	conn, err := grpc.Dial(SessionServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return FailedToConnectToVersionManager, nil // should return err?
	}
	defer conn.Close()
	client := NewVersionManagerClient(conn)

	delegatedRequest := TerminateVersionRequest{
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
			client.AbortVersion(ctx, &delegatedRequest)
			return nil, delegateError
		} else {
			return response, nil
		}
	} else {
		log.Printf("Unknown Version manager action requested: %v", request.GetAction())
		return UnknownVersionManagerAction, nil
	}
}

func (s *TerrariumGrpcGateway) UploadSourceZip(server pb.Publisher_UploadSourceZipServer) error {
	conn, err := grpc.Dial(StorageServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return err
	}
	defer conn.Close()
	client := NewStorageClient(conn)

	uploadStream, err := client.UploadSourceZip(context.TODO())
	if err != nil {
		return err
	}
	for {
		chunk, err := server.Recv()
		if err == io.EOF {
			uploadStream.CloseSend()
			return server.SendAndClose(ArchiveUploaded)
		}
		if err != nil {
			return err
		}
		err = uploadStream.Send(chunk)
		if err != nil {
			server.SendAndClose(ArchiveUploadFailed)
			return err
		}
	}
}

func (s *TerrariumGrpcGateway) DownloadSourceZip(request *pb.DownloadSourceZipRequest, server pb.Consumer_DownloadSourceZipServer) error {
	conn, err := grpc.Dial(StorageServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return err
	}
	defer conn.Close()
	client := NewStorageClient(conn)

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
	conn, err := grpc.Dial(DependencyServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return FailedToConnectToDependencyResolver, nil // should return also err?
	}
	defer conn.Close()
	client := NewDependencyResolverClient(conn)

	if response, err := client.RegisterModuleDependencies(ctx, request); err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func (s *TerrariumGrpcGateway) RegisterContainerDependencies(ctx context.Context, request *pb.RegisterContainerDependenciesRequest) (*pb.TransactionStatusResponse, error) {
	conn, err := grpc.Dial(DependencyServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return FailedToConnectToDependencyResolver, nil // should return also err?
	}
	defer conn.Close()
	client := NewDependencyResolverClient(conn)

	if response, err := client.RegisterContainerDependencies(ctx, request); err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func (s *TerrariumGrpcGateway) RetrieveContainerDependencies(request *pb.RetrieveContainerDependenciesRequest, server pb.Consumer_RetrieveContainerDependenciesServer) error {
	conn, err := grpc.Dial(DependencyServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return err
	}
	defer conn.Close()
	client := NewDependencyResolverClient(conn)

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
	conn, err := grpc.Dial(DependencyServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return err
	}
	defer conn.Close()
	client := NewDependencyResolverClient(conn)

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
