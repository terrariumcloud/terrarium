package registrar

import (
	"context"
	"github.com/terrariumcloud/terrarium/internal/module/services"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"google.golang.org/grpc"
)

type registrarGrpcClient struct {
	endpoint string
}

func NewRegistrarGrpcClient(endpoint string) services.RegistrarClient {
	return &registrarGrpcClient{endpoint: endpoint}
}

func (r registrarGrpcClient) Register(ctx context.Context, in *module.RegisterModuleRequest, opts ...grpc.CallOption) (*module.Response, error) {
	if connVersion, err := services.CreateGRPCConnection(r.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = connVersion.Close() }()

		client := services.NewRegistrarClient(connVersion)
		return client.Register(ctx, in, opts...)
	}
}

func (r registrarGrpcClient) ListModules(ctx context.Context, in *services.ListModulesRequest, opts ...grpc.CallOption) (*services.ListModulesResponse, error) {
	if connVersion, err := services.CreateGRPCConnection(r.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = connVersion.Close() }()

		client := services.NewRegistrarClient(connVersion)
		return client.ListModules(ctx, in, opts...)
	}
}

func (r registrarGrpcClient) GetModule(ctx context.Context, in *services.GetModuleRequest, opts ...grpc.CallOption) (*services.GetModuleResponse, error) {
	if connVersion, err := services.CreateGRPCConnection(r.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = connVersion.Close() }()

		client := services.NewRegistrarClient(connVersion)
		return client.GetModule(ctx, in, opts...)
	}
}
