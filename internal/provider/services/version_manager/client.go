package version_manager

import (
	"context"
	"github.com/terrariumcloud/terrarium/internal/common/grpcService"
	"github.com/terrariumcloud/terrarium/internal/provider/services"
	terrarium "github.com/terrariumcloud/terrarium/pkg/terrarium/provider"
	"google.golang.org/grpc"
)

type versionManagerGrpcClient struct {
	endpoint string
}

func NewVersionManagerGrpcClient(endpoint string) services.VersionManagerClient {
	return &versionManagerGrpcClient{endpoint: endpoint}
}

func (v versionManagerGrpcClient) ListProviderVersions(ctx context.Context, in *services.ProviderName, opts ...grpc.CallOption) (*services.ProviderVersionsResponse, error) {
	if conn, err := grpcService.CreateGRPCConnection(v.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = conn.Close() }()

		client := services.NewVersionManagerClient(conn)
		return client.ListProviderVersions(ctx, in, opts...)
	}
}

func (v versionManagerGrpcClient) GetVersionData(ctx context.Context, in *services.VersionDataRequest, opts ...grpc.CallOption) (*services.PlatformMetadataResponse, error) {
	if conn, err := grpcService.CreateGRPCConnection(v.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = conn.Close() }()

		client := services.NewVersionManagerClient(conn)
		return client.GetVersionData(ctx, in, opts...)
	}
}

func (v versionManagerGrpcClient) ListProviders(ctx context.Context, in *services.ListProvidersRequest, opts ...grpc.CallOption) (*services.ListProvidersResponse, error) {
	if conn, err := grpcService.CreateGRPCConnection(v.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = conn.Close() }()

		client := services.NewVersionManagerClient(conn)
		return client.ListProviders(ctx, in, opts...)
	}
}

func (v versionManagerGrpcClient) GetProvider(ctx context.Context, in *services.ProviderName, opts ...grpc.CallOption) (*services.GetProviderResponse, error) {
	if conn, err := grpcService.CreateGRPCConnection(v.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = conn.Close() }()

		client := services.NewVersionManagerClient(conn)
		return client.GetProvider(ctx, in, opts...)
	}
}

func (v versionManagerGrpcClient) Register(ctx context.Context, in *terrarium.RegisterProviderRequest, opts ...grpc.CallOption) (*terrarium.Response, error) {
	if conn, err := grpcService.CreateGRPCConnection(v.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = conn.Close() }()

		client := services.NewVersionManagerClient(conn)
		return client.Register(ctx, in, opts...)
	}
}

func (v versionManagerGrpcClient) PublishVersion(ctx context.Context, in *services.TerminateVersionRequest, opts ...grpc.CallOption) (*terrarium.Response, error) {
	if conn, err := grpcService.CreateGRPCConnection(v.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = conn.Close() }()

		client := services.NewVersionManagerClient(conn)
		return client.PublishVersion(ctx, in, opts...)
	}
}

func (v versionManagerGrpcClient) AbortProviderVersion(ctx context.Context, in *services.TerminateVersionRequest, opts ...grpc.CallOption) (*terrarium.Response, error) {
	if conn, err := grpcService.CreateGRPCConnection(v.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = conn.Close() }()

		client := services.NewVersionManagerClient(conn)
		return client.AbortProviderVersion(ctx, in, opts...)
	}
}
