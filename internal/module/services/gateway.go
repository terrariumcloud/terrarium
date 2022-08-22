package services

import (
	"context"
	"io"
	"log"

	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type TerrariumGrpcGateway struct {
	terrarium.UnimplementedPublisherServer
	terrarium.UnimplementedConsumerServer
}

func (s *TerrariumGrpcGateway) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	terrarium.RegisterPublisherServer(grpcServer, s)
	terrarium.RegisterConsumerServer(grpcServer, s)
	return nil
}

// Register new module with Registrar service
func (s *TerrariumGrpcGateway) Register(ctx context.Context, request *terrarium.RegisterModuleRequest) (*terrarium.TransactionStatusResponse, error) {
	conn, err := grpc.Dial(RegistrarServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return nil, err
	}

	defer conn.Close()

	client := NewRegistrarClient(conn)

	log.Println("Register module => Registrar")
	if res, delegateError := client.Register(ctx, request); delegateError != nil {
		log.Printf("Register module failed: %v", delegateError)
		return nil, delegateError
	} else {
		return res, nil
	}
}

// Begin (Create) Version with Version Manager service
func (s *TerrariumGrpcGateway) BeginVersion(ctx context.Context, request *terrarium.BeginVersionRequest) (*terrarium.BeginVersionResponse, error) {
	conn, err := grpc.Dial(VersionManagerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return nil, err
	}

	defer conn.Close()

	client := NewVersionManagerClient(conn)

	delegatedRequest := BeginVersionRequest{
		Module: request.GetModule(),
	}

	log.Println("Create new version => Version Manager")
	if res, delegateError := client.BeginVersion(ctx, &delegatedRequest); delegateError != nil {
		log.Printf("BeginVersion remote call failed: %v", delegateError)
		return nil, delegateError
	} else {
		return res, nil
	}
}

// End Version with Version Manger service
// This can mean either abort (remove) or publish Version
func (s *TerrariumGrpcGateway) EndVersion(ctx context.Context, request *terrarium.EndVersionRequest) (*terrarium.TransactionStatusResponse, error) {
	conn, err := grpc.Dial(VersionManagerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return nil, err
	}

	defer conn.Close()

	client := NewVersionManagerClient(conn)

	delegatedRequest := TerminateVersionRequest{
		Module: request.GetModule(),
	}

	if request.GetAction() == terrarium.EndVersionRequest_DISCARD {
		log.Println("Abort version => Version Manager")
		if res, delegateError := client.AbortVersion(ctx, &delegatedRequest); delegateError != nil {
			log.Printf("AbortVersion remote call failed: %v", delegateError)
			return nil, delegateError
		} else {
			return res, nil
		}
	} else if request.GetAction() == terrarium.EndVersionRequest_PUBLISH {
		log.Println("Publish version => Version Manager")
		if res, delegateError := client.PublishVersion(ctx, &delegatedRequest); delegateError != nil {
			log.Printf("PublishVersion remote call failed: %v", delegateError)
			return nil, delegateError
		} else {
			return res, nil
		}
	} else {
		log.Printf("Unknown Version manager action requested: %v", request.GetAction())
		return UnknownVersionManagerAction, nil
	}
}

// Upload source zip to Storage service
func (s *TerrariumGrpcGateway) UploadSourceZip(server terrarium.Publisher_UploadSourceZipServer) error {
	conn, err := grpc.Dial(StorageServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return err
	}

	defer conn.Close()

	ctx := server.Context()
	md, _ := metadata.FromIncomingContext(ctx)

	client := NewStorageClient(conn)

	log.Println("Upload source zip => Storage")
	ctx = metadata.NewOutgoingContext(ctx, md)
	upstream, upErr := client.UploadSourceZip(ctx)

	if upErr != nil {
		return upErr
	}

	for {
		req, err := server.Recv()

		if err == io.EOF {
			res, upErr := upstream.CloseAndRecv()
			if upErr != nil {
				return upErr
			}
			return server.SendAndClose(res)
		}

		if err != nil {
			return err
		}

		upErr = upstream.Send(req)

		if upErr == io.EOF {
			if upErr := upstream.CloseSend(); upErr != nil {
				return upErr
			}
			return server.SendAndClose(ArchiveUploaded)
		}
		if upErr != nil {
			server.SendAndClose(ArchiveUploadFailed)
			return upErr
		}
	}
}

// Download source zip from Storage service
func (s *TerrariumGrpcGateway) DownloadSourceZip(request *terrarium.DownloadSourceZipRequest, server terrarium.Consumer_DownloadSourceZipServer) error {
	conn, err := grpc.Dial(StorageServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return err
	}

	defer conn.Close()

	client := NewStorageClient(conn)

	downloadStream, err := client.DownloadSourceZip(server.Context(), &terrarium.DownloadSourceZipRequest{
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

// Register Module dependencies with Dependency Resolver service
func (s *TerrariumGrpcGateway) RegisterModuleDependencies(ctx context.Context, request *terrarium.RegisterModuleDependenciesRequest) (*terrarium.TransactionStatusResponse, error) {
	conn, err := grpc.Dial(DependencyManagerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return nil, err
	}

	defer conn.Close()

	client := NewDependencyManagerClient(conn)

	log.Println("Register module dependencies => Dependency Manager")
	if res, err := client.RegisterModuleDependencies(ctx, request); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

// Register Container dependencies with Dependency Resolver service
func (s *TerrariumGrpcGateway) RegisterContainerDependencies(ctx context.Context, request *terrarium.RegisterContainerDependenciesRequest) (*terrarium.TransactionStatusResponse, error) {
	conn, err := grpc.Dial(DependencyManagerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return nil, err
	}

	defer conn.Close()

	client := NewDependencyManagerClient(conn)

	log.Println("Register container dependencies => Dependency Manager")
	if res, err := client.RegisterContainerDependencies(ctx, request); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

// Retrieve Container dependencies from Dependency Resolver service
func (s *TerrariumGrpcGateway) RetrieveContainerDependencies(request *terrarium.RetrieveContainerDependenciesRequest, server terrarium.Consumer_RetrieveContainerDependenciesServer) error {
	conn, err := grpc.Dial(DependencyManagerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return err
	}

	defer conn.Close()

	client := NewDependencyManagerClient(conn)

	dependencyStream, err := client.RetrieveContainerDependencies(server.Context(), request)

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

// Retrieve Module dependences from Dependency Resolver service
func (s *TerrariumGrpcGateway) RetrieveModuleDependencies(request *terrarium.RetrieveModuleDependenciesRequest, server terrarium.Consumer_RetrieveModuleDependenciesServer) error {
	conn, err := grpc.Dial(DependencyManagerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return err
	}

	defer conn.Close()

	client := NewDependencyManagerClient(conn)

	dependencyStream, err := client.RetrieveModuleDependencies(server.Context(), request)

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
