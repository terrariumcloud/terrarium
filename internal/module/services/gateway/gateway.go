package gateway

import (
	"context"
	"io"
	"log"

	"github.com/terrariumcloud/terrarium/internal/module/services"
	"github.com/terrariumcloud/terrarium/internal/module/services/dependency_manager"
	"github.com/terrariumcloud/terrarium/internal/module/services/registrar"
	"github.com/terrariumcloud/terrarium/internal/module/services/storage"
	"github.com/terrariumcloud/terrarium/internal/module/services/tag_manager"
	"github.com/terrariumcloud/terrarium/internal/module/services/version_manager"
	release "github.com/terrariumcloud/terrarium/internal/release/services"
	releaseSvc "github.com/terrariumcloud/terrarium/internal/release/services/release"

	terrarium "github.com/terrariumcloud/terrarium/pkg/terrarium/module"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
	release.UnimplementedBrowseServer
}

// RegisterWithServer registers TerrariumGrpcGateway with grpc server
func (gw *TerrariumGrpcGateway) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	terrarium.RegisterPublisherServer(grpcServer, gw)
	terrarium.RegisterConsumerServer(grpcServer, gw)
	release.RegisterBrowseServer(grpcServer, gw)
	return nil
}

// Register new module with Registrar service
func (gw *TerrariumGrpcGateway) Register(ctx context.Context, request *terrarium.RegisterModuleRequest) (*terrarium.Response, error) {
	log.Println("Register => Registrar")

	span := trace.SpanFromContext(ctx)
	span.AddEvent("gateway: registering new Module with Registrar service", trace.WithAttributes(attribute.String("Module Name", request.GetName())))
	span.SetAttributes(
		attribute.String("module.name", request.GetName()),
	)

	conn, err := services.CreateGRPCConnection(registrar.RegistrarServiceEndpoint)

	if err != nil {
		log.Println(err)
		span.RecordError(err)
		return nil, ConnectToRegistrarError
	}

	defer conn.Close()

	client := services.NewRegistrarClient(conn)

	return gw.RegisterWithClient(ctx, request, client)
}

// RegisterWithClient calls Register on Registrar client
func (gw *TerrariumGrpcGateway) RegisterWithClient(ctx context.Context, request *terrarium.RegisterModuleRequest, client services.RegistrarClient) (*terrarium.Response, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent("gateway: register with Client", trace.WithAttributes(attribute.String("Module Name", request.GetName())))
	span.SetAttributes(
		attribute.String("module.name", request.GetName()),
	)

	if res, delegateError := client.Register(ctx, request); delegateError != nil {
		log.Printf("Failed: %v", delegateError)
		span.RecordError(delegateError)
		return nil, delegateError
	} else {
		log.Println("Done <= Registrar")
		span.AddEvent("Successful call to Register on Registrar client.")
		return res, nil
	}
}

// Register PublishTag with Registrar service
func (gw *TerrariumGrpcGateway) PublishTag(ctx context.Context, request *terrarium.PublishTagRequest) (*terrarium.Response, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent("gateway: registering PublishTag with Registrar service", trace.WithAttributes(attribute.String("Module Name", request.GetName()), attribute.StringSlice("Module Tags", request.GetTags())))
	span.SetAttributes(
		attribute.String("module.name", request.GetName()),
		attribute.StringSlice("module.tags", request.GetTags()),
	)

	conn, err := services.CreateGRPCConnection(tag_manager.TagManagerEndpoint)

	if err != nil {
		log.Println(err)
		span.RecordError(err)
		return nil, tag_manager.ConnectToTagManagerError
	}

	defer conn.Close()

	client := services.NewTagManagerClient(conn)

	return gw.PublishTagWithClient(ctx, request, client)
}

func (gw *TerrariumGrpcGateway) PublishTagWithClient(ctx context.Context, request *terrarium.PublishTagRequest, client services.TagManagerClient) (*terrarium.Response, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent("gateway: publishing tags with Client", trace.WithAttributes(attribute.String("Module Name", request.GetName()), attribute.StringSlice("Module Tags", request.GetTags())))
	span.SetAttributes(
		attribute.String("module.name", request.GetName()),
		attribute.StringSlice("module.tags", request.GetTags()),
	)

	if res, delegateError := client.PublishTag(ctx, request); delegateError != nil {
		log.Printf("Failed: %v", delegateError)
		span.RecordError(delegateError)
		return nil, delegateError
	} else {
		log.Println("Done <= Registrar")
		span.AddEvent("Successfully publishing tags with Client.")
		return res, nil
	}
}

// BeginVersion creates new version with Version Manager service
func (gw *TerrariumGrpcGateway) BeginVersion(ctx context.Context, request *terrarium.BeginVersionRequest) (*terrarium.Response, error) {
	log.Println("Begin version => Version Manager")

	span := trace.SpanFromContext(ctx)
	span.AddEvent("gateway: creating new version with Version Manager service", trace.WithAttributes(attribute.String("Module Name", request.Module.GetName()), attribute.String("Module Version", request.Module.GetVersion())))
	span.SetAttributes(
		attribute.String("module.name", request.Module.GetName()),
		attribute.String("module.version", request.Module.GetVersion()),
	)

	conn, err := services.CreateGRPCConnection(version_manager.VersionManagerEndpoint)

	if err != nil {
		log.Println(err)
		span.RecordError(err)
		return nil, ConnectToVersionManagerError
	}

	defer conn.Close()

	client := services.NewVersionManagerClient(conn)

	return gw.BeginVersionWithClient(ctx, request, client)
}

// BeginVersionWithClient calls BeginVersion on Version Manager client
func (gw *TerrariumGrpcGateway) BeginVersionWithClient(ctx context.Context, request *terrarium.BeginVersionRequest, client services.VersionManagerClient) (*terrarium.Response, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent("gateway: call BeginVersion on Version Manager Client", trace.WithAttributes(attribute.String("Module Name", request.Module.GetName()), attribute.String("Module Version", request.Module.GetVersion())))
	span.SetAttributes(
		attribute.String("module.name", request.Module.GetName()),
		attribute.String("module.version", request.Module.GetVersion()),
	)
	if res, delegateError := client.BeginVersion(ctx, request); delegateError != nil {
		log.Printf("Failed: %v", delegateError)
		span.RecordError(delegateError)
		return nil, delegateError
	} else {
		log.Println("Done <= Version Manager")
		span.AddEvent("Sucessful call to BeginVersion.")
		return res, nil
	}
}

// EndVersion publishes/aborts with Version Manger service
func (gw *TerrariumGrpcGateway) EndVersion(ctx context.Context, request *terrarium.EndVersionRequest) (*terrarium.Response, error) {
	log.Println("End version => Version Manager")

	span := trace.SpanFromContext(ctx)
	span.AddEvent("gateway: aborting version with Version Manager service", trace.WithAttributes(attribute.String("Module Name", request.Module.GetName()), attribute.String("Module Version", request.Module.GetVersion())))
	span.SetAttributes(
		attribute.String("module.name", request.Module.GetName()),
		attribute.String("module.version", request.Module.GetVersion()),
	)

	conn, err := services.CreateGRPCConnection(version_manager.VersionManagerEndpoint)
	if err != nil {
		log.Println(err)
		span.RecordError(err)
		return nil, ConnectToVersionManagerError
	}

	defer conn.Close()

	client := services.NewVersionManagerClient(conn)

	return gw.EndVersionWithClient(ctx, request, client)
}

// EndVersionWithClient calls AbortVersion/PublishVersion on Version Manager client
func (gw *TerrariumGrpcGateway) EndVersionWithClient(ctx context.Context, request *terrarium.EndVersionRequest, client services.VersionManagerClient) (*terrarium.Response, error) {
	terminateRequest := services.TerminateVersionRequest{
		Module: request.GetModule(),
	}
	span := trace.SpanFromContext(ctx)

	if request.GetAction() == terrarium.EndVersionRequest_DISCARD {
		log.Println("Abort version => Version Manager")
		if res, delegateError := client.AbortVersion(ctx, &terminateRequest); delegateError != nil {
			log.Printf("Failed: %v", delegateError)
			span.RecordError(delegateError)
			return nil, delegateError
		} else {
			log.Println("Done <= Version Manager")
			span.AddEvent("Successfully aborted version", trace.WithAttributes(attribute.String("Module Name", request.Module.GetName()), attribute.String("Module Version", request.Module.GetVersion())))
			return res, nil
		}
	} else if request.GetAction() == terrarium.EndVersionRequest_PUBLISH {
		log.Println("Publish version => Version Manager")
		if res, delegateError := client.PublishVersion(ctx, &terminateRequest); delegateError != nil {
			log.Printf("Failed: %v", delegateError)
			span.RecordError(delegateError)
			return nil, delegateError
		} else {
			log.Println("Done <= Version Manager")
			span.AddEvent("Successfully published version", trace.WithAttributes(attribute.String("Module Name", request.Module.GetName()), attribute.String("Module Version", request.Module.GetVersion())))
			return res, nil
		}
	} else {
		log.Printf("Unknown Version manager action requested: %v", request.GetAction())
		span.AddEvent("Unknown action.")
		return nil, UnknownVersionManagerActionError
	}
}

// UploadSourceZip uploads source zip to Storage service
func (gw *TerrariumGrpcGateway) UploadSourceZip(server terrarium.Publisher_UploadSourceZipServer) error {
	log.Println("Upload source zip => Storage")
	ctx := server.Context()
	span := trace.SpanFromContext(ctx)

	conn, err := services.CreateGRPCConnection(storage.StorageServiceEndpoint)

	if err != nil {
		log.Println(err)
		span.RecordError(err)
		return ConnectToStorageError
	}

	defer conn.Close()

	client := services.NewStorageClient(conn)

	return gw.UploadSourceZipWithClient(server, client)
}

// UploadSourceZipWithClient calls UploadSourceZip on Storage client
func (gw *TerrariumGrpcGateway) UploadSourceZipWithClient(server terrarium.Publisher_UploadSourceZipServer, client services.StorageClient) error {
	upstream, upErr := client.UploadSourceZip(server.Context())
	ctx := server.Context()
	span := trace.SpanFromContext(ctx)
	if upErr != nil {
		log.Println(upErr)
		span.RecordError(upErr)
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
			span.AddEvent("Success. End of file. No more input is available.")
			return server.SendAndClose(res)
		}

		if err != nil {
			log.Printf("Failed to receive: %v", err)
			span.RecordError(err)
			return storage.RecieveSourceZipError
		}

		upErr = upstream.Send(req)

		if upErr == io.EOF {
			log.Println("Done <= Store")
			span.AddEvent("Success.")
			upstream.CloseSend()
			return server.SendAndClose(storage.SourceZipUploaded)
		}

		if upErr != nil {
			log.Printf("Failed to send: %v", upErr)
			span.RecordError(upErr)
			return upErr
		}
	}
}

// DownloadSourceZip downloads source zip from Storage service
func (gw *TerrariumGrpcGateway) DownloadSourceZip(request *terrarium.DownloadSourceZipRequest, server terrarium.Consumer_DownloadSourceZipServer) error {
	log.Println("Download source zip => Storage")

	ctx := server.Context()
	span := trace.SpanFromContext(ctx)
	span.AddEvent("gateway: download source zip from Storage", trace.WithAttributes(attribute.String("Module Name", request.Module.GetName()), attribute.String("Module Version", request.Module.GetVersion())))
	span.SetAttributes(
		attribute.String("module.name", request.Module.GetName()),
		attribute.String("module.version", request.Module.GetVersion()),
	)

	conn, err := services.CreateGRPCConnection(storage.StorageServiceEndpoint)

	if err != nil {
		log.Println(err)
		span.RecordError(err)
		return ConnectToStorageError
	}

	defer conn.Close()

	client := services.NewStorageClient(conn)

	return gw.DownloadSourceZipWithClient(request, server, client)
}

// DownloadSourceZipWithClient calls DownloadSourceZip on Storage client
func (gw *TerrariumGrpcGateway) DownloadSourceZipWithClient(request *terrarium.DownloadSourceZipRequest, server terrarium.Consumer_DownloadSourceZipServer, client services.StorageClient) error {
	downstream, downErr := client.DownloadSourceZip(server.Context(), request)
	ctx := server.Context()
	span := trace.SpanFromContext(ctx)

	if downErr != nil {
		log.Println(downErr)
		span.RecordError(downErr)
		return downErr
	}

	for {
		res, downErr := downstream.Recv()

		if downErr == io.EOF {
			log.Println("Done <= Storage")
			span.AddEvent("Success. End of file. No more input is available.")
			return nil
		}

		if downErr != nil {
			log.Printf("Failed to recieve: %v", downErr)
			span.RecordError(downErr)
			return downErr
		}

		err := server.Send(res)

		if err != nil {
			log.Printf("Failed to send: %v", err)
			span.RecordError(err)
			downstream.CloseSend()
			return storage.SendSourceZipError
		}
	}
}

// RegisterModuleDependencies registers Module dependencies with Dependency Manager service
func (gw *TerrariumGrpcGateway) RegisterModuleDependencies(ctx context.Context, request *terrarium.RegisterModuleDependenciesRequest) (*terrarium.Response, error) {
	log.Println("Register module dependencies => Dependency Manager")

	span := trace.SpanFromContext(ctx)
	span.AddEvent("gateway: registering module dependencies with Dependency Manager service", trace.WithAttributes(attribute.String("Module Name", request.Module.GetName()), attribute.String("Module Version", request.Module.GetVersion())))
	span.SetAttributes(
		attribute.String("module.name", request.Module.GetName()),
		attribute.String("module.version", request.Module.GetVersion()),
	)

	conn, err := services.CreateGRPCConnection(dependency_manager.DependencyManagerEndpoint)

	if err != nil {
		log.Println(err)
		span.RecordError(err)
		return nil, ConnectToDependencyManagerError
	}

	defer conn.Close()

	client := services.NewDependencyManagerClient(conn)

	return gw.RegisterModuleDependenciesWithClient(ctx, request, client)
}

// RegisterModuleDependenciesWithClient calls RegisterModuleDependencies on Dependency Manager client
func (gw *TerrariumGrpcGateway) RegisterModuleDependenciesWithClient(ctx context.Context, request *terrarium.RegisterModuleDependenciesRequest, client services.DependencyManagerClient) (*terrarium.Response, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent("gateway: registering module dependencies with Client", trace.WithAttributes(attribute.String("Module Name", request.Module.GetName()), attribute.String("Module Version", request.Module.GetVersion())))
	span.SetAttributes(
		attribute.String("module.name", request.Module.GetName()),
		attribute.String("module.version", request.Module.GetVersion()),
	)

	if res, err := client.RegisterModuleDependencies(ctx, request); err != nil {
		log.Println(err)
		span.RecordError(err)
		return nil, err
	} else {
		log.Println("Done <= Dependency Manager")
		span.AddEvent("Successfully registered module dependencies with Client.")
		return res, nil
	}
}

// RegisterContainerDependencies registers Container dependencies with Dependency Manager service
func (gw *TerrariumGrpcGateway) RegisterContainerDependencies(ctx context.Context, request *terrarium.RegisterContainerDependenciesRequest) (*terrarium.Response, error) {
	log.Println("Register container dependencies => Dependency Manager")

	span := trace.SpanFromContext(ctx)
	span.AddEvent("gateway: registering container dependencies with Dependency Manager service", trace.WithAttributes(attribute.String("Module Name", request.Module.GetName()), attribute.String("Module Version", request.Module.GetVersion())))
	span.SetAttributes(
		attribute.String("module.name", request.Module.GetName()),
		attribute.String("module.version", request.Module.GetVersion()),
	)

	conn, err := services.CreateGRPCConnection(dependency_manager.DependencyManagerEndpoint)

	if err != nil {
		log.Println(err)
		span.RecordError(err)
		return nil, ConnectToDependencyManagerError
	}

	defer conn.Close()

	client := services.NewDependencyManagerClient(conn)

	return gw.RegisterContainerDependenciesWithClient(ctx, request, client)
}

// RegisterContainerDependenciesWithClient calls RegisterContainerDependencies on Dependency Manager client
func (gw *TerrariumGrpcGateway) RegisterContainerDependenciesWithClient(ctx context.Context, request *terrarium.RegisterContainerDependenciesRequest, client services.DependencyManagerClient) (*terrarium.Response, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent("gateway: registering container dependencies with Client", trace.WithAttributes(attribute.String("Module Name", request.Module.GetName()), attribute.String("Module Version", request.Module.GetVersion())))
	span.SetAttributes(
		attribute.String("module.name", request.Module.GetName()),
		attribute.String("module.version", request.Module.GetVersion()),
	)
	if res, err := client.RegisterContainerDependencies(ctx, request); err != nil {
		log.Println(err)
		span.RecordError(err)
		return nil, err
	} else {
		log.Println("Done <= Dependency Manager")
		span.AddEvent("Successfully registered container dependencies with Client.")
		return res, nil
	}
}

// RetrieveContainerDependencies retrieves Container dependencies from Dependency Manager service
func (gw *TerrariumGrpcGateway) RetrieveContainerDependencies(request *terrarium.RetrieveContainerDependenciesRequest, server terrarium.Consumer_RetrieveContainerDependenciesServer) error {
	log.Println("Retrieve container dependencies => NOOP")
	return nil
}

// RetrieveContainerDependenciesV2 retrieves Container dependencies from Dependency Manager service
func (gw *TerrariumGrpcGateway) RetrieveContainerDependenciesV2(request *terrarium.RetrieveContainerDependenciesRequestV2, server terrarium.Consumer_RetrieveContainerDependenciesV2Server) error {
	log.Println("Retrieve container dependencies => Dependency Manager")

	ctx := server.Context()
	span := trace.SpanFromContext(ctx)
	span.AddEvent("gateway: retrieving container dependencies from Dependency Manager service", trace.WithAttributes(attribute.String("Module Name", request.Module.GetName()), attribute.String("Module Version", request.Module.GetVersion())))
	span.SetAttributes(
		attribute.String("module.name", request.Module.GetName()),
		attribute.String("module.version", request.Module.GetVersion()),
	)

	conn, err := services.CreateGRPCConnection(dependency_manager.DependencyManagerEndpoint)

	if err != nil {
		log.Println(err)
		span.RecordError(err)
		return ConnectToDependencyManagerError
	}

	defer conn.Close()

	client := services.NewDependencyManagerClient(conn)

	return gw.RetrieveContainerDependenciesV2WithClient(request, server, client)
}

// RetrieveContainerDependenciesWithClient calls RetrieveContainerDependencies on Dependency Manager client
func (gw *TerrariumGrpcGateway) RetrieveContainerDependenciesV2WithClient(request *terrarium.RetrieveContainerDependenciesRequestV2, server terrarium.Consumer_RetrieveContainerDependenciesV2Server, client services.DependencyManagerClient) error {
	//ctx := metadata.AppendToOutgoingContext(server.Context(), "k", "v")

	downStream, downErr := client.RetrieveContainerDependencies(server.Context(), request)

	ctx := server.Context()
	span := trace.SpanFromContext(ctx)
	span.AddEvent("gateway: retrieving container dependencies with Client", trace.WithAttributes(attribute.String("Module Name", request.Module.GetName()), attribute.String("Module Version", request.Module.GetVersion())))
	span.SetAttributes(
		attribute.String("module.name", request.Module.GetName()),
		attribute.String("module.version", request.Module.GetVersion()),
	)

	if downErr != nil {
		span.RecordError(downErr)
		return downErr
	}

	for {
		res, downErr := downStream.Recv()

		if downErr == io.EOF {
			log.Println("Done <= Dependency Manager")
			span.AddEvent("Success. Retrieved container dependencies.")
			return nil
		}

		if downErr != nil {
			log.Printf("Failed to recieve: %v", downErr)
			span.RecordError(downErr)
			return downErr
		}

		err := server.Send(res)

		if err != nil {
			log.Printf("Failed to send: %v", err)
			span.RecordError(err)
			downStream.CloseSend()
			return ForwardModuleDependenciesError
		}
	}
}

// Retrieve Module dependences from Dependency Manager service
func (gw *TerrariumGrpcGateway) RetrieveModuleDependencies(request *terrarium.RetrieveModuleDependenciesRequest, server terrarium.Consumer_RetrieveModuleDependenciesServer) error {
	log.Println("Retrieve module dependencies => Dependency Manager")
	conn, err := services.CreateGRPCConnection(dependency_manager.DependencyManagerEndpoint)

	ctx := server.Context()
	span := trace.SpanFromContext(ctx)
	span.AddEvent("gateway: retrieving module dependencies from Dependency Manager service", trace.WithAttributes(attribute.String("Module Name", request.Module.GetName()), attribute.String("Module Version", request.Module.GetVersion())))
	span.SetAttributes(
		attribute.String("module.name", request.Module.GetName()),
		attribute.String("module.version", request.Module.GetVersion()),
	)

	if err != nil {
		log.Println(err)
		span.RecordError(err)
		return ConnectToDependencyManagerError
	}

	defer conn.Close()

	client := services.NewDependencyManagerClient(conn)

	return gw.RetrieveModuleDependenciesWithClient(request, server, client)
}

func (gw *TerrariumGrpcGateway) RetrieveModuleDependenciesWithClient(request *terrarium.RetrieveModuleDependenciesRequest, server terrarium.Consumer_RetrieveModuleDependenciesServer, client services.DependencyManagerClient) error {
	downStream, downErr := client.RetrieveModuleDependencies(server.Context(), request)

	ctx := server.Context()
	span := trace.SpanFromContext(ctx)
	span.AddEvent("gateway: retrieving module dependencies with Client", trace.WithAttributes(attribute.String("Module Name", request.Module.GetName()), attribute.String("Module Version", request.Module.GetVersion())))
	span.SetAttributes(
		attribute.String("module.name", request.Module.GetName()),
		attribute.String("module.version", request.Module.GetVersion()),
	)

	if downErr != nil {
		span.RecordError(downErr)
		return downErr
	}

	for {
		res, downErr := downStream.Recv()

		if downErr == io.EOF {
			log.Println("Done <= Dependency Manager")
			span.AddEvent("Success. Received all module dependencies.")
			return nil
		}

		if downErr != nil {
			log.Printf("Failed to recieve: %v", downErr)
			span.RecordError(downErr)
			return downErr
		}

		err := server.Send(res)

		if err != nil {
			log.Printf("Failed to send: %v", err)
			span.RecordError(err)
			downStream.CloseSend()
			return ForwardModuleDependenciesError
		}
	}
}

// Register new module with Registrar service
func (gw *TerrariumGrpcGateway) ListReleases(ctx context.Context, request *release.ListReleasesRequest) (*release.ListReleasesResponse, error) {
	log.Println("ListReleases => ListReleases")

	span := trace.SpanFromContext(ctx)
	span.AddEvent("gateway: Retrieving list of releases from release service", trace.WithAttributes(attribute.String("", request.Get)))

	// attribute does not support Uint64 so converting to int64
	// MaxAgeSeconds := releaseSvc.convertUint64ToInt64(request.GetMaxAgeSeconds())

	span.SetAttributes(
		attribute.StringSlice("release.organizations", request.GetOrganizations()),
		// attribute.Int64("release.maxAge", MaxAgeSeconds),
		attribute.StringSlice("release.types", request.GetTypes()),
		attribute.String("release.page", request.GetPage().String()),
	)

	conn, err := services.CreateGRPCConnection(releaseSvc.ReleaseServiceEndpoint)

	if err != nil {
		log.Println(err)
		span.RecordError(err)
		return nil, ConnectToRegistrarError
	}

	defer conn.Close()

	client := release.NewBrowseClient(conn)

	return gw.ListReleasesWithClient(ctx, request, client)
}

// RegisterWithClient calls Register on Registrar client
func (gw *TerrariumGrpcGateway) ListReleasesWithClient(ctx context.Context, request *release.ListReleasesRequest, client release.BrowseClient) (*release.ListReleasesResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent("gateway: register with Client", trace.WithAttributes(attribute.String("Module Name", request.GetName())))
	span.SetAttributes(
		attribute.String("module.name", request.GetName()),
	)

	if res, delegateError := client.Register(ctx, request); delegateError != nil {
		log.Printf("Failed: %v", delegateError)
		span.RecordError(delegateError)
		return nil, delegateError
	} else {
		log.Println("Done <= Registrar")
		span.AddEvent("Successful call to Register on Registrar client.")
		return res, nil
	}
}
