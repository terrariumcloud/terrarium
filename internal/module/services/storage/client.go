package storage

import (
	"context"
	"github.com/terrariumcloud/terrarium/internal/module/services"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"google.golang.org/grpc"
)

type storageGrpcClient struct {
	endpoint string
}

func NewStorageGrpcClient(endpoint string) services.StorageClient {
	return &storageGrpcClient{endpoint: endpoint}
}

func (s storageGrpcClient) UploadSourceZip(ctx context.Context, opts ...grpc.CallOption) (services.Storage_UploadSourceZipClient, error) {
	if connVersion, err := services.CreateGRPCConnection(s.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = connVersion.Close() }()

		client := services.NewStorageClient(connVersion)
		return client.UploadSourceZip(ctx, opts...)
	}
}

func (s storageGrpcClient) DownloadSourceZip(ctx context.Context, in *module.DownloadSourceZipRequest, opts ...grpc.CallOption) (services.Storage_DownloadSourceZipClient, error) {
	if connVersion, err := services.CreateGRPCConnection(s.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = connVersion.Close() }()

		client := services.NewStorageClient(connVersion)
		return client.DownloadSourceZip(ctx, in, opts...)
	}
}
