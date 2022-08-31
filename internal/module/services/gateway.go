package services

import (
	"context"
	"io"
	"log"

	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	ConnectToRegistrarError           = status.Error(codes.Unavailable, "Failed to connect to Registrar service.")
	ConnectToVersionManagerError      = status.Error(codes.Unavailable, "Failed to connect to Version manager service.")
	ConnectToStorageError             = status.Error(codes.Unavailable, "Failed to connect to Storage service.")
	ConnectToDependencyManagerError   = status.Error(codes.Unavailable, "Failed to connect to Dependency manager service.")
	UnknownVersionManagerActionError  = status.Error(codes.InvalidArgument, "Unknown Version manager action requested.")
	ForwardModuleDependenciesError    = status.Error(codes.Unknown, "Failed to send module dependencies.")
	ForwardContainerDependenciesError = status.Error(codes.Unknown, "Failed to send module dependencies.")
)

type TerrariumGrpcGateway struct {
	terrarium.UnimplementedPublisherServer
	terrarium.UnimplementedConsumerServer
}

// Registers TerrariumGrpcGateway with grpc server
func (s *TerrariumGrpcGateway) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	terrarium.RegisterPublisherServer(grpcServer, s)
	terrarium.RegisterConsumerServer(grpcServer, s)
	return nil
}

// Register new module with Registrar service
func (s *TerrariumGrpcGateway) Register(ctx context.Context, request *terrarium.RegisterModuleRequest) (*terrarium.Response, error) {
	log.Println("Register => Registrar")
	conn, err := grpc.Dial(RegistrarServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
		return nil, ConnectToRegistrarError
	}

	defer conn.Close()

	client := NewRegistrarClient(conn)

	if res, delegateError := client.Register(ctx, request); delegateError != nil {
		log.Printf("Failed: %v", delegateError)
		return nil, delegateError
	} else {
		log.Println("Done <= Registrar")
		return res, nil
	}
}

// Begin (Create) Version with Version Manager service
func (s *TerrariumGrpcGateway) BeginVersion(ctx context.Context, request *terrarium.BeginVersionRequest) (*terrarium.Response, error) {
	log.Println("Begin version => Version Manager")
	conn, err := grpc.Dial(VersionManagerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
		return nil, ConnectToVersionManagerError
	}

	defer conn.Close()

	client := NewVersionManagerClient(conn)

	if res, delegateError := client.BeginVersion(ctx, request); delegateError != nil {
		log.Printf("Failed: %v", delegateError)
		return nil, delegateError
	} else {
		log.Println("Done <= Version Manager")
		return res, nil
	}
}

// End Version with Version Manger service
// This can mean either abort (remove) or publish Version
func (s *TerrariumGrpcGateway) EndVersion(ctx context.Context, request *terrarium.EndVersionRequest) (*terrarium.Response, error) {
	log.Println("End version => Version Manager")
	conn, err := grpc.Dial(VersionManagerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
		return nil, ConnectToVersionManagerError
	}

	defer conn.Close()

	client := NewVersionManagerClient(conn)

	terminateRequest := TerminateVersionRequest{
		Module: request.GetModule(),
	}

	if request.GetAction() == terrarium.EndVersionRequest_DISCARD {
		log.Println("Abort version => Version Manager")
		if res, delegateError := client.AbortVersion(ctx, &terminateRequest); delegateError != nil {
			log.Printf("Failed: %v", delegateError)
			return nil, delegateError
		} else {
			log.Println("Done <= Version Manager")
			return res, nil
		}
	} else if request.GetAction() == terrarium.EndVersionRequest_PUBLISH {
		log.Println("Publish version => Version Manager")
		if res, delegateError := client.PublishVersion(ctx, &terminateRequest); delegateError != nil {
			log.Printf("Failed: %v", delegateError)
			return nil, delegateError
		} else {
			log.Println("Done <= Version Manager")
			return res, nil
		}
	} else {
		log.Printf("Unknown Version manager action requested: %v", request.GetAction())
		return nil, UnknownVersionManagerActionError
	}
}

// Upload source zip to Storage service
func (s *TerrariumGrpcGateway) UploadSourceZip(server terrarium.Publisher_UploadSourceZipServer) error {
	log.Println("Upload source zip => Storage")
	conn, err := grpc.Dial(StorageServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
		return ConnectToStorageError
	}

	defer conn.Close()

	client := NewStorageClient(conn)

	upstream, upErr := client.UploadSourceZip(server.Context())

	if upErr != nil {
		log.Println(upErr)
		return upErr
	}

	for {
		req, err := server.Recv()

		if err == io.EOF {
			res, upErr := upstream.CloseAndRecv()

			if upErr != nil {
				return upErr
			}
			log.Println("Done <= Store")
			return server.SendAndClose(res)
		}

		if err != nil {
			log.Printf("Failed to recieve: %v", err)
			return RecieveSourceZipError
		}

		upErr = upstream.Send(req)

		if upErr == io.EOF {
			log.Println("Done <= Store")
			upstream.CloseSend()
			return server.SendAndClose(SourceZipUploaded)
		}

		if upErr != nil {
			log.Printf("Failed to send: %v", upErr)
			return upErr
		}
	}
}

// Download source zip from Storage service
func (s *TerrariumGrpcGateway) DownloadSourceZip(request *terrarium.DownloadSourceZipRequest, server terrarium.Consumer_DownloadSourceZipServer) error {
	log.Println("Download source zip => Storage")
	conn, err := grpc.Dial(StorageServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
		return ConnectToStorageError
	}

	defer conn.Close()

	client := NewStorageClient(conn)

	downstream, downErr := client.DownloadSourceZip(server.Context(), request)

	if downErr != nil {
		log.Println(downErr)
		return downErr
	}

	for {
		res, downErr := downstream.Recv()

		if downErr == io.EOF {
			log.Println("Done <= Storage")
			return nil
		}

		if downErr != nil {
			log.Printf("Failed to recieve: %v", downErr)
			return downErr
		}

		err = server.Send(res)

		if err != nil {
			log.Printf("Failed to send: %v", err)
			downstream.CloseSend()
			return SendSourceZipError
		}
	}
}

// Register Module dependencies with Dependency Manager service
func (s *TerrariumGrpcGateway) RegisterModuleDependencies(ctx context.Context, request *terrarium.RegisterModuleDependenciesRequest) (*terrarium.Response, error) {
	log.Println("Register module dependencies => Dependency Manager")
	conn, err := grpc.Dial(DependencyManagerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
		return nil, ConnectToDependencyManagerError
	}

	defer conn.Close()

	client := NewDependencyManagerClient(conn)

	if res, err := client.RegisterModuleDependencies(ctx, request); err != nil {
		log.Println(err)
		return nil, err
	} else {
		log.Println("Done <= Dependency Manager")
		return res, nil
	}
}

// Register Container dependencies with Dependency Manager service
func (s *TerrariumGrpcGateway) RegisterContainerDependencies(ctx context.Context, request *terrarium.RegisterContainerDependenciesRequest) (*terrarium.Response, error) {
	log.Println("Register container dependencies => Dependency Manager")
	conn, err := grpc.Dial(DependencyManagerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
		return nil, ConnectToDependencyManagerError
	}

	defer conn.Close()

	client := NewDependencyManagerClient(conn)

	if res, err := client.RegisterContainerDependencies(ctx, request); err != nil {
		log.Println(err)
		return nil, err
	} else {
		log.Println("Done <= Dependency Manager")
		return res, nil
	}
}

// Retrieve Container dependencies from Dependency Manager service
func (s *TerrariumGrpcGateway) RetrieveContainerDependencies(request *terrarium.RetrieveContainerDependenciesRequest, server terrarium.Consumer_RetrieveContainerDependenciesServer) error {
	log.Println("Retrieve container dependencies => Dependency Manager")
	conn, err := grpc.Dial(DependencyManagerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
		return ConnectToDependencyManagerError
	}

	defer conn.Close()

	client := NewDependencyManagerClient(conn)

	downStream, downErr := client.RetrieveContainerDependencies(server.Context(), request)

	if downErr != nil {
		return downErr
	}

	for {
		res, downErr := downStream.Recv()

		if downErr == io.EOF {
			log.Println("Done <= Dependency Manager")
			return nil
		}

		if downErr != nil {
			log.Printf("Failed to recieve: %v", downErr)
			return downErr
		}

		err = server.Send(res)

		if err != nil {
			log.Printf("Failed to send: %v", err)
			downStream.CloseSend()
			return ForwardModuleDependenciesError
		}
	}
}

// Retrieve Module dependences from Dependency Manager service
func (s *TerrariumGrpcGateway) RetrieveModuleDependencies(request *terrarium.RetrieveModuleDependenciesRequest, server terrarium.Consumer_RetrieveModuleDependenciesServer) error {
	log.Println("Retrieve module dependencies => Dependency Manager")
	conn, err := grpc.Dial(DependencyManagerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
		return ConnectToDependencyManagerError
	}

	defer conn.Close()

	client := NewDependencyManagerClient(conn)

	downStream, downErr := client.RetrieveModuleDependencies(server.Context(), request)

	if downErr != nil {
		return downErr
	}

	for {
		res, downErr := downStream.Recv()

		if downErr == io.EOF {
			log.Println("Done <= Dependency Manager")
			return nil
		}

		if downErr != nil {
			log.Printf("Failed to recieve: %v", err)
			return downErr
		}

		err = server.Send(res)

		if err != nil {
			log.Printf("Failed to send: %v", err)
			downStream.CloseSend()
			return ForwardContainerDependenciesError
		}
	}
}
