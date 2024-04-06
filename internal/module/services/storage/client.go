package storage

import (
	"context"
	"github.com/terrariumcloud/terrarium/internal/module/services"
	"github.com/terrariumcloud/terrarium/internal/common/grpcService"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
)

type storageGrpcClient struct {
	endpoint string
}

func NewStorageGrpcClient(endpoint string) services.StorageClient {
	return &storageGrpcClient{endpoint: endpoint}
}

func (s storageGrpcClient) UploadSourceZip(ctx context.Context, opts ...grpc.CallOption) (services.Storage_UploadSourceZipClient, error) {
	if conn, err := grpcService.CreateGRPCConnection(s.endpoint); err != nil {
		return nil, err
	} else {
		client := services.NewStorageClient(conn)
		if upload, err := client.UploadSourceZip(ctx, opts...); err == nil {
			return &uploadSourceZipClient{client: upload, conn: conn}, nil
		} else {
			_ = conn.Close()
			return nil, err
		}
	}
}

func (s storageGrpcClient) DownloadSourceZip(ctx context.Context, in *module.DownloadSourceZipRequest, opts ...grpc.CallOption) (services.Storage_DownloadSourceZipClient, error) {
	if conn, err := grpcService.CreateGRPCConnection(s.endpoint); err != nil {
		return nil, err
	} else {
		client := services.NewStorageClient(conn)
		if download, err := client.DownloadSourceZip(ctx, in, opts...); err == nil {
			return &downloadSourceZipClient{conn: conn, client: download}, nil
		} else {
			_ = conn.Close()
			return nil, err
		}
	}
}

type uploadSourceZipClient struct {
	conn   *grpc.ClientConn
	client services.Storage_UploadSourceZipClient
}

func (u *uploadSourceZipClient) Send(request *module.UploadSourceZipRequest) error {
	return u.client.Send(request)
}

func (u *uploadSourceZipClient) CloseAndRecv() (*module.Response, error) {
	defer func() { _ = u.conn.Close() }()
	return u.client.CloseAndRecv()
}

func (u *uploadSourceZipClient) Header() (metadata.MD, error) {
	return u.client.Header()
}

func (u *uploadSourceZipClient) Trailer() metadata.MD {
	return u.client.Trailer()
}

func (u *uploadSourceZipClient) CloseSend() error {
	return u.client.CloseSend()
}

func (u *uploadSourceZipClient) Context() context.Context {
	return u.client.Context()
}

func (u *uploadSourceZipClient) SendMsg(m any) error {
	return u.client.SendMsg(m)
}

func (u *uploadSourceZipClient) RecvMsg(m any) error {
	return u.client.RecvMsg(m)
}

type downloadSourceZipClient struct {
	conn   *grpc.ClientConn
	client services.Storage_DownloadSourceZipClient
}

func (d downloadSourceZipClient) Recv() (*module.SourceZipResponse, error) {
	result, err := d.client.Recv()
	if err == io.EOF {
		_ = d.conn.Close()
	}
	return result, err
}

func (d downloadSourceZipClient) Header() (metadata.MD, error) {
	return d.client.Header()
}

func (d downloadSourceZipClient) Trailer() metadata.MD {
	return d.client.Trailer()
}

func (d downloadSourceZipClient) CloseSend() error {
	return d.client.CloseSend()
}

func (d downloadSourceZipClient) Context() context.Context {
	return d.client.Context()
}

func (d downloadSourceZipClient) SendMsg(m any) error {
	return d.client.SendMsg(m)
}

func (d downloadSourceZipClient) RecvMsg(m any) error {
	return d.client.RecvMsg(m)
}
