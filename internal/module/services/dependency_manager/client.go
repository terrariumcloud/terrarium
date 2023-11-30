package dependency_manager

import (
	"context"
	"github.com/terrariumcloud/terrarium/internal/module/services"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"google.golang.org/grpc"
)

type dependencyManagerGrpcClient struct {
	endpoint string
}

func NewDependencyManagerGrpcClient(endpoint string) services.DependencyManagerClient {
	return &dependencyManagerGrpcClient{endpoint: endpoint}
}

func (d dependencyManagerGrpcClient) RegisterModuleDependencies(ctx context.Context, in *module.RegisterModuleDependenciesRequest, opts ...grpc.CallOption) (*module.Response, error) {
	if conn, err := services.CreateGRPCConnection(d.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = conn.Close() }()

		client := services.NewDependencyManagerClient(conn)
		return client.RegisterModuleDependencies(ctx, in, opts...)
	}
}

func (d dependencyManagerGrpcClient) RegisterContainerDependencies(ctx context.Context, in *module.RegisterContainerDependenciesRequest, opts ...grpc.CallOption) (*module.Response, error) {
	if conn, err := services.CreateGRPCConnection(d.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = conn.Close() }()

		client := services.NewDependencyManagerClient(conn)
		return client.RegisterContainerDependencies(ctx, in, opts...)
	}
}

func (d dependencyManagerGrpcClient) RetrieveContainerDependencies(ctx context.Context, in *module.RetrieveContainerDependenciesRequestV2, opts ...grpc.CallOption) (services.DependencyManager_RetrieveContainerDependenciesClient, error) {
	if conn, err := services.CreateGRPCConnection(d.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = conn.Close() }()

		client := services.NewDependencyManagerClient(conn)
		return client.RetrieveContainerDependencies(ctx, in, opts...)
	}
}

func (d dependencyManagerGrpcClient) RetrieveModuleDependencies(ctx context.Context, in *module.RetrieveModuleDependenciesRequest, opts ...grpc.CallOption) (services.DependencyManager_RetrieveModuleDependenciesClient, error) {
	if conn, err := services.CreateGRPCConnection(d.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = conn.Close() }()

		client := services.NewDependencyManagerClient(conn)
		return client.RetrieveModuleDependencies(ctx, in, opts...)
	}
}
