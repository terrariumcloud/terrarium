package storage

import (
	"context"
	"github.com/terrariumcloud/terrarium/internal/common/grpc_service"
	"github.com/terrariumcloud/terrarium/internal/provider/services"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/provider"
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

func (s storageGrpcClient) DownloadProviderSourceZip(ctx context.Context, in *services.DownloadSourceZipRequest, opts ...grpc.CallOption) (services.Storage_DownloadProviderSourceZipClient, error) {
	if conn, err := grpc_service.CreateGRPCConnection(s.endpoint); err != nil {
		return nil, err
	} else {
		client := services.NewStorageClient(conn)
		if download, err := client.DownloadProviderSourceZip(ctx, in, opts...); err == nil {
			return &downloadSourceZipClient{conn: conn, client: download}, nil
		} else {
			_ = conn.Close()
			return nil, err
		}
	}
}

func (s storageGrpcClient) DownloadShasum(ctx context.Context, in *services.DownloadShasumRequest, opts ...grpc.CallOption) (services.Storage_DownloadShasumClient, error) {
	if conn, err := grpc_service.CreateGRPCConnection(s.endpoint); err != nil {
		return nil, err
	} else {
		client := services.NewStorageClient(conn)
		if download, err := client.DownloadShasum(ctx, in, opts...); err == nil {
			return &downloadShasumClient{conn: conn, client: download}, nil
		} else {
			_ = conn.Close()
			return nil, err
		}
	}
}

func (s storageGrpcClient) DownloadShasumSignature(ctx context.Context, in *services.DownloadShasumRequest, opts ...grpc.CallOption) (services.Storage_DownloadShasumSignatureClient, error) {
	if conn, err := grpc_service.CreateGRPCConnection(s.endpoint); err != nil {
		return nil, err
	} else {
		client := services.NewStorageClient(conn)
		if download, err := client.DownloadShasumSignature(ctx, in, opts...); err == nil {
			return &downloadShasumSignatureClient{conn: conn, client: download}, nil
		} else {
			_ = conn.Close()
			return nil, err
		}
	}
}

func (s storageGrpcClient) UploadProviderBinaryZip(ctx context.Context, opts ...grpc.CallOption) (services.Storage_UploadProviderBinaryZipClient, error) {
	if conn, err := grpc_service.CreateGRPCConnection(s.endpoint); err != nil {
		return nil, err
	} else {
		client := services.NewStorageClient(conn)
		if upload, err := client.UploadProviderBinaryZip(ctx, opts...); err == nil {
			return &uploadBinaryZipClient{
				Storage_UploadProviderBinaryZipClient: upload,
				conn:                                  conn}, nil
		} else {
			_ = conn.Close()
			return nil, err
		}
	}
}

func (s storageGrpcClient) UploadShasum(ctx context.Context, opts ...grpc.CallOption) (services.Storage_UploadShasumClient, error) {
	if conn, err := grpc_service.CreateGRPCConnection(s.endpoint); err != nil {
		return nil, err
	} else {
		client := services.NewStorageClient(conn)
		if upload, err := client.UploadShasum(ctx, opts...); err == nil {
			return &uploadShasumClient{
				Storage_UploadShasumClient: upload,
				conn:                       conn}, nil
		} else {
			_ = conn.Close()
			return nil, err
		}
	}
}

func (s storageGrpcClient) UploadShasumSignature(ctx context.Context, opts ...grpc.CallOption) (services.Storage_UploadShasumSignatureClient, error) {
	if conn, err := grpc_service.CreateGRPCConnection(s.endpoint); err != nil {
		return nil, err
	} else {
		client := services.NewStorageClient(conn)
		if upload, err := client.UploadShasumSignature(ctx, opts...); err == nil {
			return &uploadShasumSignatureClient{
				Storage_UploadShasumSignatureClient: upload,
				conn:                                conn}, nil
		} else {
			_ = conn.Close()
			return nil, err
		}
	}
}

type downloadSourceZipClient struct {
	conn   *grpc.ClientConn
	client services.Storage_DownloadProviderSourceZipClient
}

func (d downloadSourceZipClient) Recv() (*services.SourceZipResponse, error) {
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

type downloadShasumClient struct {
	conn   *grpc.ClientConn
	client services.Storage_DownloadShasumClient
}

func (d downloadShasumClient) Recv() (*services.DownloadShasumResponse, error) {
	result, err := d.client.Recv()
	if err == io.EOF {
		_ = d.conn.Close()
	}
	return result, err
}

func (d downloadShasumClient) Header() (metadata.MD, error) {
	return d.client.Header()
}

func (d downloadShasumClient) Trailer() metadata.MD {
	return d.client.Trailer()
}

func (d downloadShasumClient) CloseSend() error {
	return d.client.CloseSend()
}

func (d downloadShasumClient) Context() context.Context {
	return d.client.Context()
}

func (d downloadShasumClient) SendMsg(m any) error {
	return d.client.SendMsg(m)
}

func (d downloadShasumClient) RecvMsg(m any) error {
	return d.client.RecvMsg(m)
}

type downloadShasumSignatureClient struct {
	conn   *grpc.ClientConn
	client services.Storage_DownloadShasumSignatureClient
}

func (d downloadShasumSignatureClient) Recv() (*services.DownloadShasumResponse, error) {
	result, err := d.client.Recv()
	if err == io.EOF {
		_ = d.conn.Close()
	}
	return result, err
}

func (d downloadShasumSignatureClient) Header() (metadata.MD, error) {
	return d.client.Header()
}

func (d downloadShasumSignatureClient) Trailer() metadata.MD {
	return d.client.Trailer()
}

func (d downloadShasumSignatureClient) CloseSend() error {
	return d.client.CloseSend()
}

func (d downloadShasumSignatureClient) Context() context.Context {
	return d.client.Context()
}

func (d downloadShasumSignatureClient) SendMsg(m any) error {
	return d.client.SendMsg(m)
}

func (d downloadShasumSignatureClient) RecvMsg(m any) error {
	return d.client.RecvMsg(m)
}

type uploadBinaryZipClient struct {
	conn *grpc.ClientConn
	services.Storage_UploadProviderBinaryZipClient
}

func (u *uploadBinaryZipClient) CloseAndRecv() (*provider.Response, error) {
	defer func() { _ = u.conn.Close() }()
	return u.Storage_UploadProviderBinaryZipClient.CloseAndRecv()
}

type uploadShasumClient struct {
	conn *grpc.ClientConn
	services.Storage_UploadShasumClient
}

func (u *uploadShasumClient) CloseAndRecv() (*provider.Response, error) {
	defer func() { _ = u.conn.Close() }()
	return u.Storage_UploadShasumClient.CloseAndRecv()
}

type uploadShasumSignatureClient struct {
	conn *grpc.ClientConn
	services.Storage_UploadShasumSignatureClient
}

func (u *uploadShasumSignatureClient) CloseAndRecv() (*provider.Response, error) {
	defer func() { _ = u.conn.Close() }()
	return u.Storage_UploadShasumSignatureClient.CloseAndRecv()
}
