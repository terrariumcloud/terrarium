package dependency_manager

import (
	"context"
	"github.com/terrariumcloud/terrarium/internal/common/grpc_service"
	"github.com/terrariumcloud/terrarium/internal/module/services"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"google.golang.org/grpc"
	"io"
)

type dependencyManagerGrpcClient struct {
	endpoint string
}

func NewDependencyManagerGrpcClient(endpoint string) services.DependencyManagerClient {
	return &dependencyManagerGrpcClient{endpoint: endpoint}
}

func (d dependencyManagerGrpcClient) RegisterModuleDependencies(ctx context.Context, in *module.RegisterModuleDependenciesRequest, opts ...grpc.CallOption) (*module.Response, error) {
	if conn, err := grpc_service.CreateGRPCConnection(d.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = conn.Close() }()

		client := services.NewDependencyManagerClient(conn)
		return client.RegisterModuleDependencies(ctx, in, opts...)
	}
}

func (d dependencyManagerGrpcClient) RegisterContainerDependencies(ctx context.Context, in *module.RegisterContainerDependenciesRequest, opts ...grpc.CallOption) (*module.Response, error) {
	if conn, err := grpc_service.CreateGRPCConnection(d.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = conn.Close() }()

		client := services.NewDependencyManagerClient(conn)
		return client.RegisterContainerDependencies(ctx, in, opts...)
	}
}

func (d dependencyManagerGrpcClient) RetrieveContainerDependencies(ctx context.Context, in *module.RetrieveContainerDependenciesRequestV2, opts ...grpc.CallOption) (services.DependencyManager_RetrieveContainerDependenciesClient, error) {
	if conn, err := grpc_service.CreateGRPCConnection(d.endpoint); err != nil {
		return nil, err
	} else {
		client := services.NewDependencyManagerClient(conn)
		if deps, err := client.RetrieveContainerDependencies(ctx, in, opts...); err == nil {
			return &dependencyManager_RetrieveContainerDependenciesClient{
				conn:   conn,
				client: deps,
			}, nil
		} else {
			_ = conn.Close()
			return nil, err
		}
	}
}

func (d dependencyManagerGrpcClient) RetrieveModuleDependencies(ctx context.Context, in *module.RetrieveModuleDependenciesRequest, opts ...grpc.CallOption) (services.DependencyManager_RetrieveModuleDependenciesClient, error) {
	if conn, err := grpc_service.CreateGRPCConnection(d.endpoint); err != nil {
		return nil, err
	} else {
		client := services.NewDependencyManagerClient(conn)
		if deps, err := client.RetrieveModuleDependencies(ctx, in, opts...); err == nil {
			return &dependencyManager_RetrieveModuleDependenciesClient{
				conn:   conn,
				client: deps,
			}, nil
		} else {
			_ = conn.Close()
			return nil, err
		}
	}
}

type dependencyManager_RetrieveContainerDependenciesClient struct {
	grpc.ClientStream
	conn   *grpc.ClientConn
	client services.DependencyManager_RetrieveContainerDependenciesClient
}

func (d *dependencyManager_RetrieveContainerDependenciesClient) Recv() (*module.ContainerDependenciesResponseV2, error) {
	result, err := d.client.Recv()
	if err == io.EOF {
		_ = d.conn.Close()
	}
	return result, err
}

type dependencyManager_RetrieveModuleDependenciesClient struct {
	grpc.ClientStream
	conn   *grpc.ClientConn
	client services.DependencyManager_RetrieveModuleDependenciesClient
}

func (d dependencyManager_RetrieveModuleDependenciesClient) Recv() (*module.ModuleDependenciesResponse, error) {
	result, err := d.client.Recv()
	if err == io.EOF {
		_ = d.conn.Close()
	}
	return result, err
}
