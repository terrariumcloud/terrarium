// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: pb/terrarium/provider/services/storage.proto

package services

import (
	context "context"
	provider "github.com/terrariumcloud/terrarium/pkg/terrarium/provider"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// StorageClient is the client API for Storage service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StorageClient interface {
	DownloadProviderSourceZip(ctx context.Context, in *DownloadSourceZipRequest, opts ...grpc.CallOption) (Storage_DownloadProviderSourceZipClient, error)
	DownloadShasum(ctx context.Context, in *DownloadShasumRequest, opts ...grpc.CallOption) (Storage_DownloadShasumClient, error)
	DownloadShasumSignature(ctx context.Context, in *DownloadShasumRequest, opts ...grpc.CallOption) (Storage_DownloadShasumSignatureClient, error)
	UploadProviderBinaryZip(ctx context.Context, opts ...grpc.CallOption) (Storage_UploadProviderBinaryZipClient, error)
	UploadShasum(ctx context.Context, opts ...grpc.CallOption) (Storage_UploadShasumClient, error)
	UploadShasumSignature(ctx context.Context, opts ...grpc.CallOption) (Storage_UploadShasumSignatureClient, error)
}

type storageClient struct {
	cc grpc.ClientConnInterface
}

func NewStorageClient(cc grpc.ClientConnInterface) StorageClient {
	return &storageClient{cc}
}

func (c *storageClient) DownloadProviderSourceZip(ctx context.Context, in *DownloadSourceZipRequest, opts ...grpc.CallOption) (Storage_DownloadProviderSourceZipClient, error) {
	stream, err := c.cc.NewStream(ctx, &Storage_ServiceDesc.Streams[0], "/terrarium.provider.services.Storage/DownloadProviderSourceZip", opts...)
	if err != nil {
		return nil, err
	}
	x := &storageDownloadProviderSourceZipClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Storage_DownloadProviderSourceZipClient interface {
	Recv() (*SourceZipResponse, error)
	grpc.ClientStream
}

type storageDownloadProviderSourceZipClient struct {
	grpc.ClientStream
}

func (x *storageDownloadProviderSourceZipClient) Recv() (*SourceZipResponse, error) {
	m := new(SourceZipResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *storageClient) DownloadShasum(ctx context.Context, in *DownloadShasumRequest, opts ...grpc.CallOption) (Storage_DownloadShasumClient, error) {
	stream, err := c.cc.NewStream(ctx, &Storage_ServiceDesc.Streams[1], "/terrarium.provider.services.Storage/DownloadShasum", opts...)
	if err != nil {
		return nil, err
	}
	x := &storageDownloadShasumClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Storage_DownloadShasumClient interface {
	Recv() (*DownloadShasumResponse, error)
	grpc.ClientStream
}

type storageDownloadShasumClient struct {
	grpc.ClientStream
}

func (x *storageDownloadShasumClient) Recv() (*DownloadShasumResponse, error) {
	m := new(DownloadShasumResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *storageClient) DownloadShasumSignature(ctx context.Context, in *DownloadShasumRequest, opts ...grpc.CallOption) (Storage_DownloadShasumSignatureClient, error) {
	stream, err := c.cc.NewStream(ctx, &Storage_ServiceDesc.Streams[2], "/terrarium.provider.services.Storage/DownloadShasumSignature", opts...)
	if err != nil {
		return nil, err
	}
	x := &storageDownloadShasumSignatureClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Storage_DownloadShasumSignatureClient interface {
	Recv() (*DownloadShasumResponse, error)
	grpc.ClientStream
}

type storageDownloadShasumSignatureClient struct {
	grpc.ClientStream
}

func (x *storageDownloadShasumSignatureClient) Recv() (*DownloadShasumResponse, error) {
	m := new(DownloadShasumResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *storageClient) UploadProviderBinaryZip(ctx context.Context, opts ...grpc.CallOption) (Storage_UploadProviderBinaryZipClient, error) {
	stream, err := c.cc.NewStream(ctx, &Storage_ServiceDesc.Streams[3], "/terrarium.provider.services.Storage/UploadProviderBinaryZip", opts...)
	if err != nil {
		return nil, err
	}
	x := &storageUploadProviderBinaryZipClient{stream}
	return x, nil
}

type Storage_UploadProviderBinaryZipClient interface {
	Send(*provider.UploadProviderBinaryZipRequest) error
	CloseAndRecv() (*provider.Response, error)
	grpc.ClientStream
}

type storageUploadProviderBinaryZipClient struct {
	grpc.ClientStream
}

func (x *storageUploadProviderBinaryZipClient) Send(m *provider.UploadProviderBinaryZipRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *storageUploadProviderBinaryZipClient) CloseAndRecv() (*provider.Response, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(provider.Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *storageClient) UploadShasum(ctx context.Context, opts ...grpc.CallOption) (Storage_UploadShasumClient, error) {
	stream, err := c.cc.NewStream(ctx, &Storage_ServiceDesc.Streams[4], "/terrarium.provider.services.Storage/UploadShasum", opts...)
	if err != nil {
		return nil, err
	}
	x := &storageUploadShasumClient{stream}
	return x, nil
}

type Storage_UploadShasumClient interface {
	Send(*provider.UploadShasumRequest) error
	CloseAndRecv() (*provider.Response, error)
	grpc.ClientStream
}

type storageUploadShasumClient struct {
	grpc.ClientStream
}

func (x *storageUploadShasumClient) Send(m *provider.UploadShasumRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *storageUploadShasumClient) CloseAndRecv() (*provider.Response, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(provider.Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *storageClient) UploadShasumSignature(ctx context.Context, opts ...grpc.CallOption) (Storage_UploadShasumSignatureClient, error) {
	stream, err := c.cc.NewStream(ctx, &Storage_ServiceDesc.Streams[5], "/terrarium.provider.services.Storage/UploadShasumSignature", opts...)
	if err != nil {
		return nil, err
	}
	x := &storageUploadShasumSignatureClient{stream}
	return x, nil
}

type Storage_UploadShasumSignatureClient interface {
	Send(*provider.UploadShasumRequest) error
	CloseAndRecv() (*provider.Response, error)
	grpc.ClientStream
}

type storageUploadShasumSignatureClient struct {
	grpc.ClientStream
}

func (x *storageUploadShasumSignatureClient) Send(m *provider.UploadShasumRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *storageUploadShasumSignatureClient) CloseAndRecv() (*provider.Response, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(provider.Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StorageServer is the server API for Storage service.
// All implementations must embed UnimplementedStorageServer
// for forward compatibility
type StorageServer interface {
	DownloadProviderSourceZip(*DownloadSourceZipRequest, Storage_DownloadProviderSourceZipServer) error
	DownloadShasum(*DownloadShasumRequest, Storage_DownloadShasumServer) error
	DownloadShasumSignature(*DownloadShasumRequest, Storage_DownloadShasumSignatureServer) error
	UploadProviderBinaryZip(Storage_UploadProviderBinaryZipServer) error
	UploadShasum(Storage_UploadShasumServer) error
	UploadShasumSignature(Storage_UploadShasumSignatureServer) error
	mustEmbedUnimplementedStorageServer()
}

// UnimplementedStorageServer must be embedded to have forward compatible implementations.
type UnimplementedStorageServer struct {
}

func (UnimplementedStorageServer) DownloadProviderSourceZip(*DownloadSourceZipRequest, Storage_DownloadProviderSourceZipServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadProviderSourceZip not implemented")
}
func (UnimplementedStorageServer) DownloadShasum(*DownloadShasumRequest, Storage_DownloadShasumServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadShasum not implemented")
}
func (UnimplementedStorageServer) DownloadShasumSignature(*DownloadShasumRequest, Storage_DownloadShasumSignatureServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadShasumSignature not implemented")
}
func (UnimplementedStorageServer) UploadProviderBinaryZip(Storage_UploadProviderBinaryZipServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadProviderBinaryZip not implemented")
}
func (UnimplementedStorageServer) UploadShasum(Storage_UploadShasumServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadShasum not implemented")
}
func (UnimplementedStorageServer) UploadShasumSignature(Storage_UploadShasumSignatureServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadShasumSignature not implemented")
}
func (UnimplementedStorageServer) mustEmbedUnimplementedStorageServer() {}

// UnsafeStorageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StorageServer will
// result in compilation errors.
type UnsafeStorageServer interface {
	mustEmbedUnimplementedStorageServer()
}

func RegisterStorageServer(s grpc.ServiceRegistrar, srv StorageServer) {
	s.RegisterService(&Storage_ServiceDesc, srv)
}

func _Storage_DownloadProviderSourceZip_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadSourceZipRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StorageServer).DownloadProviderSourceZip(m, &storageDownloadProviderSourceZipServer{stream})
}

type Storage_DownloadProviderSourceZipServer interface {
	Send(*SourceZipResponse) error
	grpc.ServerStream
}

type storageDownloadProviderSourceZipServer struct {
	grpc.ServerStream
}

func (x *storageDownloadProviderSourceZipServer) Send(m *SourceZipResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Storage_DownloadShasum_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadShasumRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StorageServer).DownloadShasum(m, &storageDownloadShasumServer{stream})
}

type Storage_DownloadShasumServer interface {
	Send(*DownloadShasumResponse) error
	grpc.ServerStream
}

type storageDownloadShasumServer struct {
	grpc.ServerStream
}

func (x *storageDownloadShasumServer) Send(m *DownloadShasumResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Storage_DownloadShasumSignature_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadShasumRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StorageServer).DownloadShasumSignature(m, &storageDownloadShasumSignatureServer{stream})
}

type Storage_DownloadShasumSignatureServer interface {
	Send(*DownloadShasumResponse) error
	grpc.ServerStream
}

type storageDownloadShasumSignatureServer struct {
	grpc.ServerStream
}

func (x *storageDownloadShasumSignatureServer) Send(m *DownloadShasumResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Storage_UploadProviderBinaryZip_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StorageServer).UploadProviderBinaryZip(&storageUploadProviderBinaryZipServer{stream})
}

type Storage_UploadProviderBinaryZipServer interface {
	SendAndClose(*provider.Response) error
	Recv() (*provider.UploadProviderBinaryZipRequest, error)
	grpc.ServerStream
}

type storageUploadProviderBinaryZipServer struct {
	grpc.ServerStream
}

func (x *storageUploadProviderBinaryZipServer) SendAndClose(m *provider.Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *storageUploadProviderBinaryZipServer) Recv() (*provider.UploadProviderBinaryZipRequest, error) {
	m := new(provider.UploadProviderBinaryZipRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Storage_UploadShasum_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StorageServer).UploadShasum(&storageUploadShasumServer{stream})
}

type Storage_UploadShasumServer interface {
	SendAndClose(*provider.Response) error
	Recv() (*provider.UploadShasumRequest, error)
	grpc.ServerStream
}

type storageUploadShasumServer struct {
	grpc.ServerStream
}

func (x *storageUploadShasumServer) SendAndClose(m *provider.Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *storageUploadShasumServer) Recv() (*provider.UploadShasumRequest, error) {
	m := new(provider.UploadShasumRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Storage_UploadShasumSignature_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StorageServer).UploadShasumSignature(&storageUploadShasumSignatureServer{stream})
}

type Storage_UploadShasumSignatureServer interface {
	SendAndClose(*provider.Response) error
	Recv() (*provider.UploadShasumRequest, error)
	grpc.ServerStream
}

type storageUploadShasumSignatureServer struct {
	grpc.ServerStream
}

func (x *storageUploadShasumSignatureServer) SendAndClose(m *provider.Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *storageUploadShasumSignatureServer) Recv() (*provider.UploadShasumRequest, error) {
	m := new(provider.UploadShasumRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Storage_ServiceDesc is the grpc.ServiceDesc for Storage service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Storage_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "terrarium.provider.services.Storage",
	HandlerType: (*StorageServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "DownloadProviderSourceZip",
			Handler:       _Storage_DownloadProviderSourceZip_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "DownloadShasum",
			Handler:       _Storage_DownloadShasum_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "DownloadShasumSignature",
			Handler:       _Storage_DownloadShasumSignature_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "UploadProviderBinaryZip",
			Handler:       _Storage_UploadProviderBinaryZip_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "UploadShasum",
			Handler:       _Storage_UploadShasum_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "UploadShasumSignature",
			Handler:       _Storage_UploadShasumSignature_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "pb/terrarium/provider/services/storage.proto",
}
