// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: pb/terrarium/module/services/storage.proto

package services

import (
	context "context"
	module "github.com/terrariumcloud/terrarium/pkg/terrarium/module"
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
	UploadSourceZip(ctx context.Context, opts ...grpc.CallOption) (Storage_UploadSourceZipClient, error)
	DownloadSourceZip(ctx context.Context, in *module.DownloadSourceZipRequest, opts ...grpc.CallOption) (Storage_DownloadSourceZipClient, error)
}

type storageClient struct {
	cc grpc.ClientConnInterface
}

func NewStorageClient(cc grpc.ClientConnInterface) StorageClient {
	return &storageClient{cc}
}

func (c *storageClient) UploadSourceZip(ctx context.Context, opts ...grpc.CallOption) (Storage_UploadSourceZipClient, error) {
	stream, err := c.cc.NewStream(ctx, &Storage_ServiceDesc.Streams[0], "/terrarium.module.services.Storage/UploadSourceZip", opts...)
	if err != nil {
		return nil, err
	}
	x := &storageUploadSourceZipClient{stream}
	return x, nil
}

type Storage_UploadSourceZipClient interface {
	Send(*module.UploadSourceZipRequest) error
	CloseAndRecv() (*module.Response, error)
	grpc.ClientStream
}

type storageUploadSourceZipClient struct {
	grpc.ClientStream
}

func (x *storageUploadSourceZipClient) Send(m *module.UploadSourceZipRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *storageUploadSourceZipClient) CloseAndRecv() (*module.Response, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(module.Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *storageClient) DownloadSourceZip(ctx context.Context, in *module.DownloadSourceZipRequest, opts ...grpc.CallOption) (Storage_DownloadSourceZipClient, error) {
	stream, err := c.cc.NewStream(ctx, &Storage_ServiceDesc.Streams[1], "/terrarium.module.services.Storage/DownloadSourceZip", opts...)
	if err != nil {
		return nil, err
	}
	x := &storageDownloadSourceZipClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Storage_DownloadSourceZipClient interface {
	Recv() (*module.SourceZipResponse, error)
	grpc.ClientStream
}

type storageDownloadSourceZipClient struct {
	grpc.ClientStream
}

func (x *storageDownloadSourceZipClient) Recv() (*module.SourceZipResponse, error) {
	m := new(module.SourceZipResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StorageServer is the server API for Storage service.
// All implementations must embed UnimplementedStorageServer
// for forward compatibility
type StorageServer interface {
	UploadSourceZip(Storage_UploadSourceZipServer) error
	DownloadSourceZip(*module.DownloadSourceZipRequest, Storage_DownloadSourceZipServer) error
	mustEmbedUnimplementedStorageServer()
}

// UnimplementedStorageServer must be embedded to have forward compatible implementations.
type UnimplementedStorageServer struct {
}

func (UnimplementedStorageServer) UploadSourceZip(Storage_UploadSourceZipServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadSourceZip not implemented")
}
func (UnimplementedStorageServer) DownloadSourceZip(*module.DownloadSourceZipRequest, Storage_DownloadSourceZipServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadSourceZip not implemented")
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

func _Storage_UploadSourceZip_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StorageServer).UploadSourceZip(&storageUploadSourceZipServer{stream})
}

type Storage_UploadSourceZipServer interface {
	SendAndClose(*module.Response) error
	Recv() (*module.UploadSourceZipRequest, error)
	grpc.ServerStream
}

type storageUploadSourceZipServer struct {
	grpc.ServerStream
}

func (x *storageUploadSourceZipServer) SendAndClose(m *module.Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *storageUploadSourceZipServer) Recv() (*module.UploadSourceZipRequest, error) {
	m := new(module.UploadSourceZipRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Storage_DownloadSourceZip_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(module.DownloadSourceZipRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StorageServer).DownloadSourceZip(m, &storageDownloadSourceZipServer{stream})
}

type Storage_DownloadSourceZipServer interface {
	Send(*module.SourceZipResponse) error
	grpc.ServerStream
}

type storageDownloadSourceZipServer struct {
	grpc.ServerStream
}

func (x *storageDownloadSourceZipServer) Send(m *module.SourceZipResponse) error {
	return x.ServerStream.SendMsg(m)
}

// Storage_ServiceDesc is the grpc.ServiceDesc for Storage service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Storage_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "terrarium.module.services.Storage",
	HandlerType: (*StorageServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadSourceZip",
			Handler:       _Storage_UploadSourceZip_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "DownloadSourceZip",
			Handler:       _Storage_DownloadSourceZip_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pb/terrarium/module/services/storage.proto",
}
