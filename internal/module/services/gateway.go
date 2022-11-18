package services

import (
	"context"
	"io"
	"log"

	terrarium "github.com/terrariumcloud/terrarium/pkg/terrarium/module"

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

// RegisterWithServer registers TerrariumGrpcGateway with grpc server
func (gw *TerrariumGrpcGateway) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	terrarium.RegisterPublisherServer(grpcServer, gw)
	terrarium.RegisterConsumerServer(grpcServer, gw)
	return nil
}

// Register new module with Registrar service
func (gw *TerrariumGrpcGateway) Register(ctx context.Context, request *terrarium.RegisterModuleRequest) (*terrarium.Response, error) {
	log.Println("Register => Registrar")
	conn, err := grpc.Dial(RegistrarServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
		return nil, ConnectToRegistrarError
	}

	defer conn.Close()

	client := NewRegistrarClient(conn)

	return gw.RegisterWithClient(ctx, request, client)
}

// RegisterWithClient calls Register on Registrar client
func (gw *TerrariumGrpcGateway) RegisterWithClient(ctx context.Context, request *terrarium.RegisterModuleRequest, client RegistrarClient) (*terrarium.Response, error) {
	if res, delegateError := client.Register(ctx, request); delegateError != nil {
		log.Printf("Failed: %v", delegateError)
		return nil, delegateError
	} else {
		log.Println("Done <= Registrar")
		return res, nil
	}
}

// BeginVersion creates new version with Version Manager service
func (gw *TerrariumGrpcGateway) BeginVersion(ctx context.Context, request *terrarium.BeginVersionRequest) (*terrarium.Response, error) {
	log.Println("Begin version => Version Manager")
	conn, err := grpc.Dial(VersionManagerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
		return nil, ConnectToVersionManagerError
	}

	defer conn.Close()

	client := NewVersionManagerClient(conn)

	return gw.BeginVersionWithClient(ctx, request, client)
}

// BeginVersionWithClient calls BeginVersion on Version Manager client
func (gw *TerrariumGrpcGateway) BeginVersionWithClient(ctx context.Context, request *terrarium.BeginVersionRequest, client VersionManagerClient) (*terrarium.Response, error) {
	if res, delegateError := client.BeginVersion(ctx, request); delegateError != nil {
		log.Printf("Failed: %v", delegateError)
		return nil, delegateError
	} else {
		log.Println("Done <= Version Manager")
		return res, nil
	}
}

// EndVersion publishes/aborts with Version Manger service
func (gw *TerrariumGrpcGateway) EndVersion(ctx context.Context, request *terrarium.EndVersionRequest) (*terrarium.Response, error) {
	log.Println("End version => Version Manager")
	conn, err := grpc.Dial(VersionManagerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
		return nil, ConnectToVersionManagerError
	}

	defer conn.Close()

	client := NewVersionManagerClient(conn)

	return gw.EndVersionWithClient(ctx, request, client)
}

// EndVersionWithClient calls AbortVersion/PublishVersion on Version Manager client
func (gw *TerrariumGrpcGateway) EndVersionWithClient(ctx context.Context, request *terrarium.EndVersionRequest, client VersionManagerClient) (*terrarium.Response, error) {
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

// UploadSourceZip uploads source zip to Storage service
func (gw *TerrariumGrpcGateway) UploadSourceZip(server terrarium.Publisher_UploadSourceZipServer) error {
	log.Println("Upload source zip => Storage")
	conn, err := grpc.Dial(StorageServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
		return ConnectToStorageError
	}

	defer conn.Close()

	client := NewStorageClient(conn)

	return gw.UploadSourceZipWithClient(server, client)
}

// UploadSourceZipWithClient calls UploadSourceZip on Storage client
func (gw *TerrariumGrpcGateway) UploadSourceZipWithClient(server terrarium.Publisher_UploadSourceZipServer, client StorageClient) error {
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

// DownloadSourceZip downloads source zip from Storage service
func (gw *TerrariumGrpcGateway) DownloadSourceZip(request *terrarium.DownloadSourceZipRequest, server terrarium.Consumer_DownloadSourceZipServer) error {
	log.Println("Download source zip => Storage")
	conn, err := grpc.Dial(StorageServiceEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
		return ConnectToStorageError
	}

	defer conn.Close()

	client := NewStorageClient(conn)

	return gw.DownloadSourceZipWithClient(request, server, client)
}

// DownloadSourceZipWithClient calls DownloadSourceZip on Storage client
func (gw *TerrariumGrpcGateway) DownloadSourceZipWithClient(request *terrarium.DownloadSourceZipRequest, server terrarium.Consumer_DownloadSourceZipServer, client StorageClient) error {
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

		err := server.Send(res)

		if err != nil {
			log.Printf("Failed to send: %v", err)
			downstream.CloseSend()
			return SendSourceZipError
		}
	}
}

// RegisterModuleDependencies registers Module dependencies with Dependency Manager service
func (gw *TerrariumGrpcGateway) RegisterModuleDependencies(ctx context.Context, request *terrarium.RegisterModuleDependenciesRequest) (*terrarium.Response, error) {
	log.Println("Register module dependencies => Dependency Manager")
	conn, err := grpc.Dial(DependencyManagerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
		return nil, ConnectToDependencyManagerError
	}

	defer conn.Close()

	client := NewDependencyManagerClient(conn)

	return gw.RegisterModuleDependenciesWithClient(ctx, request, client)
}

// RegisterModuleDependenciesWithClient calls RegisterModuleDependencies on Dependency Manager client
func (gw *TerrariumGrpcGateway) RegisterModuleDependenciesWithClient(ctx context.Context, request *terrarium.RegisterModuleDependenciesRequest, client DependencyManagerClient) (*terrarium.Response, error) {
	if res, err := client.RegisterModuleDependencies(ctx, request); err != nil {
		log.Println(err)
		return nil, err
	} else {
		log.Println("Done <= Dependency Manager")
		return res, nil
	}
}

// RegisterContainerDependencies registers Container dependencies with Dependency Manager service
func (gw *TerrariumGrpcGateway) RegisterContainerDependencies(ctx context.Context, request *terrarium.RegisterContainerDependenciesRequest) (*terrarium.Response, error) {
	log.Println("Register container dependencies => Dependency Manager")
	conn, err := grpc.Dial(DependencyManagerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
		return nil, ConnectToDependencyManagerError
	}

	defer conn.Close()

	client := NewDependencyManagerClient(conn)

	return gw.RegisterContainerDependenciesWithClient(ctx, request, client)
}

// RegisterContainerDependenciesWithClient calls RegisterContainerDependencies on Dependency Manager client
func (gw *TerrariumGrpcGateway) RegisterContainerDependenciesWithClient(ctx context.Context, request *terrarium.RegisterContainerDependenciesRequest, client DependencyManagerClient) (*terrarium.Response, error) {
	if res, err := client.RegisterContainerDependencies(ctx, request); err != nil {
		log.Println(err)
		return nil, err
	} else {
		log.Println("Done <= Dependency Manager")
		return res, nil
	}
}

// RetrieveContainerDependencies retrieves Container dependencies from Dependency Manager service
func (gw *TerrariumGrpcGateway) RetrieveContainerDependencies(request *terrarium.RetrieveContainerDependenciesRequest, server terrarium.Consumer_RetrieveContainerDependenciesServer) error {
	log.Println("Retrieve container dependencies => Dependency Manager")
	conn, err := grpc.Dial(DependencyManagerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
		return ConnectToDependencyManagerError
	}

	defer conn.Close()

	client := NewDependencyManagerClient(conn)

	return gw.RetrieveContainerDependenciesWithClient(request, server, client)
}

// RetrieveContainerDependenciesWithClient calls RetrieveContainerDependencies on Dependency Manager client
func (gw *TerrariumGrpcGateway) RetrieveContainerDependenciesWithClient(request *terrarium.RetrieveContainerDependenciesRequest, server terrarium.Consumer_RetrieveContainerDependenciesServer, client DependencyManagerClient) error {
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

		err := server.Send(res)

		if err != nil {
			log.Printf("Failed to send: %v", err)
			downStream.CloseSend()
			return ForwardModuleDependenciesError
		}
	}
}

// Retrieve Module dependences from Dependency Manager service
func (gw *TerrariumGrpcGateway) RetrieveModuleDependencies(request *terrarium.RetrieveModuleDependenciesRequest, server terrarium.Consumer_RetrieveModuleDependenciesServer) error {
	log.Println("Retrieve module dependencies => Dependency Manager")
	conn, err := grpc.Dial(DependencyManagerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
		return ConnectToDependencyManagerError
	}

	defer conn.Close()

	client := NewDependencyManagerClient(conn)

	return gw.RetrieveModuleDependenciesWithClient(request, server, client)
}

func (gw *TerrariumGrpcGateway) RetrieveModuleDependenciesWithClient(request *terrarium.RetrieveModuleDependenciesRequest, server terrarium.Consumer_RetrieveModuleDependenciesServer, client DependencyManagerClient) error {
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
			log.Printf("Failed to recieve: %v", downErr)
			return downErr
		}

		err := server.Send(res)

		if err != nil {
			log.Printf("Failed to send: %v", err)
			downStream.CloseSend()
			return ForwardModuleDependenciesError
		}
	}
}
