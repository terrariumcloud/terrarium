package services

import (
	"context"
	"io"
	"log"

	pb "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type TerrariumGrpcGateway struct {
	pb.UnimplementedPublisherServer
	pb.UnimplementedConsumerServer
}

// Register new module with Registrar service
func (s *TerrariumGrpcGateway) Register(ctx context.Context, request *pb.RegisterModuleRequest) (*pb.TransactionStatusResponse, error) {
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
func (s *TerrariumGrpcGateway) BeginVersion(ctx context.Context, request *pb.BeginVersionRequest) (*pb.BeginVersionResponse, error) {
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
func (s *TerrariumGrpcGateway) EndVersion(ctx context.Context, request *pb.EndVersionRequest) (*pb.TransactionStatusResponse, error) {
	conn, err := grpc.Dial(VersionManagerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return nil, err
	}

	defer conn.Close()

	client := NewVersionManagerClient(conn)

	delegatedRequest := TerminateVersionRequest{
		SessionKey: request.GetSessionKey(),
	}

	if request.GetAction() == pb.EndVersionRequest_DISCARD {
		log.Println("Abort version => Version Manager")
		if res, delegateError := client.AbortVersion(ctx, &delegatedRequest); delegateError != nil {
			log.Printf("AbortVersion remote call failed: %v", delegateError)
			return nil, delegateError
		} else {
			return res, nil
		}
	} else if request.GetAction() == pb.EndVersionRequest_PUBLISH {
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
func (s *TerrariumGrpcGateway) UploadSourceZip(server pb.Publisher_UploadSourceZipServer) error {
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
	uploadStream, upErr := client.UploadSourceZip(ctx)

	if upErr != nil {
		return upErr
	}

	for {
		req, err := server.Recv()

		if err == io.EOF {
			res, upErr := uploadStream.CloseAndRecv()
			if upErr != nil {
				return upErr
			}
			return server.SendAndClose(res)
		}

		if err != nil {
			return err
		}

		upErr = uploadStream.Send(req)

		if upErr == io.EOF {
			if upErr := uploadStream.CloseSend(); upErr != nil {
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
func (s *TerrariumGrpcGateway) DownloadSourceZip(request *pb.DownloadSourceZipRequest, server pb.Consumer_DownloadSourceZipServer) error {
	conn, err := grpc.Dial(StorageServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return err
	}

	defer conn.Close()

	client := NewStorageClient(conn)

	downloadStream, err := client.DownloadSourceZip(server.Context(), &pb.DownloadSourceZipRequest{
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
func (s *TerrariumGrpcGateway) RegisterModuleDependencies(ctx context.Context, request *pb.RegisterModuleDependenciesRequest) (*pb.TransactionStatusResponse, error) {
	conn, err := grpc.Dial(DependencyServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return nil, err
	}

	defer conn.Close()

	client := NewDependencyResolverClient(conn)

	log.Println("Register module dependencies => Dependency resolver")
	if res, err := client.RegisterModuleDependencies(ctx, request); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

// Register Container dependencies with Dependency Resolver service
func (s *TerrariumGrpcGateway) RegisterContainerDependencies(ctx context.Context, request *pb.RegisterContainerDependenciesRequest) (*pb.TransactionStatusResponse, error) {
	conn, err := grpc.Dial(DependencyServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return nil, err
	}

	defer conn.Close()

	client := NewDependencyResolverClient(conn)

	log.Println("Register container dependencies => Dependency resolver")
	if res, err := client.RegisterContainerDependencies(ctx, request); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

// Retrieve Container dependencies from Dependency Resolver service
func (s *TerrariumGrpcGateway) RetrieveContainerDependencies(request *pb.RetrieveContainerDependenciesRequest, server pb.Consumer_RetrieveContainerDependenciesServer) error {
	conn, err := grpc.Dial(DependencyServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return err
	}

	defer conn.Close()

	client := NewDependencyResolverClient(conn)

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
func (s *TerrariumGrpcGateway) RetrieveModuleDependencies(request *pb.RetrieveModuleDependenciesRequest, server pb.Consumer_RetrieveModuleDependenciesServer) error {
	conn, err := grpc.Dial(DependencyServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return err
	}

	defer conn.Close()

	client := NewDependencyResolverClient(conn)

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
