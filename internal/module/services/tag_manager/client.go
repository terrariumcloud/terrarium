package tag_manager

import (
	"context"
	"github.com/terrariumcloud/terrarium/internal/module/services"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"google.golang.org/grpc"
)

type tagManagerGrpcClient struct {
	endpoint string
}

func NewTagManagerGrpcClient(endpoint string) services.TagManagerClient {
	return &tagManagerGrpcClient{endpoint: endpoint}
}

func (t tagManagerGrpcClient) PublishTag(ctx context.Context, in *module.PublishTagRequest, opts ...grpc.CallOption) (*module.Response, error) {
	if conn, err := services.CreateGRPCConnection(t.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = conn.Close() }()

		client := services.NewTagManagerClient(conn)
		return client.PublishTag(ctx, in, opts...)
	}
}
