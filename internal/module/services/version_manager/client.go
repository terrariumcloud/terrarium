package version_manager

import (
	"context"
	"github.com/terrariumcloud/terrarium/internal/common/grpc_service"
	"github.com/terrariumcloud/terrarium/internal/module/services"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"google.golang.org/grpc"
)

type versionManagerGrpcClient struct {
	endpoint string
}

func NewVersionManagerGrpcClient(endpoint string) services.VersionManagerClient {
	return &versionManagerGrpcClient{endpoint: endpoint}
}

func (v versionManagerGrpcClient) BeginVersion(ctx context.Context, in *module.BeginVersionRequest, opts ...grpc.CallOption) (*module.Response, error) {
	if conn, err := grpc_service.CreateGRPCConnection(v.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = conn.Close() }()

		client := services.NewVersionManagerClient(conn)
		return client.BeginVersion(ctx, in, opts...)
	}
}

func (v versionManagerGrpcClient) AbortVersion(ctx context.Context, in *services.TerminateVersionRequest, opts ...grpc.CallOption) (*module.Response, error) {
	if conn, err := grpc_service.CreateGRPCConnection(v.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = conn.Close() }()

		client := services.NewVersionManagerClient(conn)
		return client.AbortVersion(ctx, in, opts...)
	}
}

func (v versionManagerGrpcClient) PublishVersion(ctx context.Context, in *services.TerminateVersionRequest, opts ...grpc.CallOption) (*module.Response, error) {
	if conn, err := grpc_service.CreateGRPCConnection(v.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = conn.Close() }()

		client := services.NewVersionManagerClient(conn)
		return client.PublishVersion(ctx, in, opts...)
	}
}

func (v versionManagerGrpcClient) ListModuleVersions(ctx context.Context, in *services.ListModuleVersionsRequest, opts ...grpc.CallOption) (*services.ListModuleVersionsResponse, error) {
	if conn, err := grpc_service.CreateGRPCConnection(v.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = conn.Close() }()

		client := services.NewVersionManagerClient(conn)
		return client.ListModuleVersions(ctx, in, opts...)
	}
}
