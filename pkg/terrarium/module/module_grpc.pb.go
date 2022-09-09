// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.5
// source: pb/terrarium/module/module.proto

package module

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PublisherClient is the client API for Publisher service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PublisherClient interface {
	Register(ctx context.Context, in *RegisterModuleRequest, opts ...grpc.CallOption) (*TransactionStatusResponse, error)
	BeginVersion(ctx context.Context, in *BeginVersionRequest, opts ...grpc.CallOption) (*BeginVersionResponse, error)
	RegisterModuleDependencies(ctx context.Context, in *RegisterModuleDependenciesRequest, opts ...grpc.CallOption) (*TransactionStatusResponse, error)
	RegisterContainerDependencies(ctx context.Context, in *RegisterContainerDependenciesRequest, opts ...grpc.CallOption) (*TransactionStatusResponse, error)
	// Register Audit Trail
	UploadSourceZip(ctx context.Context, opts ...grpc.CallOption) (Publisher_UploadSourceZipClient, error)
	// Upload Documentation
	EndVersion(ctx context.Context, in *EndVersionRequest, opts ...grpc.CallOption) (*TransactionStatusResponse, error)
}

type publisherClient struct {
	cc grpc.ClientConnInterface
}

func NewPublisherClient(cc grpc.ClientConnInterface) PublisherClient {
	return &publisherClient{cc}
}

func (c *publisherClient) Register(ctx context.Context, in *RegisterModuleRequest, opts ...grpc.CallOption) (*TransactionStatusResponse, error) {
	out := new(TransactionStatusResponse)
	err := c.cc.Invoke(ctx, "/terrarium.module.Publisher/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publisherClient) BeginVersion(ctx context.Context, in *BeginVersionRequest, opts ...grpc.CallOption) (*BeginVersionResponse, error) {
	out := new(BeginVersionResponse)
	err := c.cc.Invoke(ctx, "/terrarium.module.Publisher/BeginVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publisherClient) RegisterModuleDependencies(ctx context.Context, in *RegisterModuleDependenciesRequest, opts ...grpc.CallOption) (*TransactionStatusResponse, error) {
	out := new(TransactionStatusResponse)
	err := c.cc.Invoke(ctx, "/terrarium.module.Publisher/RegisterModuleDependencies", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publisherClient) RegisterContainerDependencies(ctx context.Context, in *RegisterContainerDependenciesRequest, opts ...grpc.CallOption) (*TransactionStatusResponse, error) {
	out := new(TransactionStatusResponse)
	err := c.cc.Invoke(ctx, "/terrarium.module.Publisher/RegisterContainerDependencies", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publisherClient) UploadSourceZip(ctx context.Context, opts ...grpc.CallOption) (Publisher_UploadSourceZipClient, error) {
	stream, err := c.cc.NewStream(ctx, &Publisher_ServiceDesc.Streams[0], "/terrarium.module.Publisher/UploadSourceZip", opts...)
	if err != nil {
		return nil, err
	}
	x := &publisherUploadSourceZipClient{stream}
	return x, nil
}

type Publisher_UploadSourceZipClient interface {
	Send(*UploadSourceZipRequest) error
	CloseAndRecv() (*TransactionStatusResponse, error)
	grpc.ClientStream
}

type publisherUploadSourceZipClient struct {
	grpc.ClientStream
}

func (x *publisherUploadSourceZipClient) Send(m *UploadSourceZipRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *publisherUploadSourceZipClient) CloseAndRecv() (*TransactionStatusResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(TransactionStatusResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *publisherClient) EndVersion(ctx context.Context, in *EndVersionRequest, opts ...grpc.CallOption) (*TransactionStatusResponse, error) {
	out := new(TransactionStatusResponse)
	err := c.cc.Invoke(ctx, "/terrarium.module.Publisher/EndVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PublisherServer is the server API for Publisher service.
// All implementations must embed UnimplementedPublisherServer
// for forward compatibility
type PublisherServer interface {
	Register(context.Context, *RegisterModuleRequest) (*TransactionStatusResponse, error)
	BeginVersion(context.Context, *BeginVersionRequest) (*BeginVersionResponse, error)
	RegisterModuleDependencies(context.Context, *RegisterModuleDependenciesRequest) (*TransactionStatusResponse, error)
	RegisterContainerDependencies(context.Context, *RegisterContainerDependenciesRequest) (*TransactionStatusResponse, error)
	// Register Audit Trail
	UploadSourceZip(Publisher_UploadSourceZipServer) error
	// Upload Documentation
	EndVersion(context.Context, *EndVersionRequest) (*TransactionStatusResponse, error)
	mustEmbedUnimplementedPublisherServer()
}

// UnimplementedPublisherServer must be embedded to have forward compatible implementations.
type UnimplementedPublisherServer struct {
}

func (UnimplementedPublisherServer) Register(context.Context, *RegisterModuleRequest) (*TransactionStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedPublisherServer) BeginVersion(context.Context, *BeginVersionRequest) (*BeginVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BeginVersion not implemented")
}
func (UnimplementedPublisherServer) RegisterModuleDependencies(context.Context, *RegisterModuleDependenciesRequest) (*TransactionStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterModuleDependencies not implemented")
}
func (UnimplementedPublisherServer) RegisterContainerDependencies(context.Context, *RegisterContainerDependenciesRequest) (*TransactionStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterContainerDependencies not implemented")
}
func (UnimplementedPublisherServer) UploadSourceZip(Publisher_UploadSourceZipServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadSourceZip not implemented")
}
func (UnimplementedPublisherServer) EndVersion(context.Context, *EndVersionRequest) (*TransactionStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EndVersion not implemented")
}
func (UnimplementedPublisherServer) mustEmbedUnimplementedPublisherServer() {}

// UnsafePublisherServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PublisherServer will
// result in compilation errors.
type UnsafePublisherServer interface {
	mustEmbedUnimplementedPublisherServer()
}

func RegisterPublisherServer(s grpc.ServiceRegistrar, srv PublisherServer) {
	s.RegisterService(&Publisher_ServiceDesc, srv)
}

func _Publisher_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterModuleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PublisherServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terrarium.module.Publisher/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PublisherServer).Register(ctx, req.(*RegisterModuleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Publisher_BeginVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BeginVersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PublisherServer).BeginVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terrarium.module.Publisher/BeginVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PublisherServer).BeginVersion(ctx, req.(*BeginVersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Publisher_RegisterModuleDependencies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterModuleDependenciesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PublisherServer).RegisterModuleDependencies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terrarium.module.Publisher/RegisterModuleDependencies",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PublisherServer).RegisterModuleDependencies(ctx, req.(*RegisterModuleDependenciesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Publisher_RegisterContainerDependencies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterContainerDependenciesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PublisherServer).RegisterContainerDependencies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terrarium.module.Publisher/RegisterContainerDependencies",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PublisherServer).RegisterContainerDependencies(ctx, req.(*RegisterContainerDependenciesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Publisher_UploadSourceZip_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PublisherServer).UploadSourceZip(&publisherUploadSourceZipServer{stream})
}

type Publisher_UploadSourceZipServer interface {
	SendAndClose(*TransactionStatusResponse) error
	Recv() (*UploadSourceZipRequest, error)
	grpc.ServerStream
}

type publisherUploadSourceZipServer struct {
	grpc.ServerStream
}

func (x *publisherUploadSourceZipServer) SendAndClose(m *TransactionStatusResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *publisherUploadSourceZipServer) Recv() (*UploadSourceZipRequest, error) {
	m := new(UploadSourceZipRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Publisher_EndVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EndVersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PublisherServer).EndVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terrarium.module.Publisher/EndVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PublisherServer).EndVersion(ctx, req.(*EndVersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Publisher_ServiceDesc is the grpc.ServiceDesc for Publisher service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Publisher_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "terrarium.module.Publisher",
	HandlerType: (*PublisherServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _Publisher_Register_Handler,
		},
		{
			MethodName: "BeginVersion",
			Handler:    _Publisher_BeginVersion_Handler,
		},
		{
			MethodName: "RegisterModuleDependencies",
			Handler:    _Publisher_RegisterModuleDependencies_Handler,
		},
		{
			MethodName: "RegisterContainerDependencies",
			Handler:    _Publisher_RegisterContainerDependencies_Handler,
		},
		{
			MethodName: "EndVersion",
			Handler:    _Publisher_EndVersion_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadSourceZip",
			Handler:       _Publisher_UploadSourceZip_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "pb/terrarium/module/module.proto",
}

// ConsumerClient is the client API for Consumer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConsumerClient interface {
	DownloadSourceZip(ctx context.Context, in *DownloadSourceZipRequest, opts ...grpc.CallOption) (Consumer_DownloadSourceZipClient, error)
	RetrieveContainerDependencies(ctx context.Context, in *RetrieveContainerDependenciesRequest, opts ...grpc.CallOption) (Consumer_RetrieveContainerDependenciesClient, error)
	RetrieveModuleDependencies(ctx context.Context, in *RetrieveModuleDependenciesRequest, opts ...grpc.CallOption) (Consumer_RetrieveModuleDependenciesClient, error)
}

type consumerClient struct {
	cc grpc.ClientConnInterface
}

func NewConsumerClient(cc grpc.ClientConnInterface) ConsumerClient {
	return &consumerClient{cc}
}

func (c *consumerClient) DownloadSourceZip(ctx context.Context, in *DownloadSourceZipRequest, opts ...grpc.CallOption) (Consumer_DownloadSourceZipClient, error) {
	stream, err := c.cc.NewStream(ctx, &Consumer_ServiceDesc.Streams[0], "/terrarium.module.Consumer/DownloadSourceZip", opts...)
	if err != nil {
		return nil, err
	}
	x := &consumerDownloadSourceZipClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Consumer_DownloadSourceZipClient interface {
	Recv() (*SourceZipResponse, error)
	grpc.ClientStream
}

type consumerDownloadSourceZipClient struct {
	grpc.ClientStream
}

func (x *consumerDownloadSourceZipClient) Recv() (*SourceZipResponse, error) {
	m := new(SourceZipResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *consumerClient) RetrieveContainerDependencies(ctx context.Context, in *RetrieveContainerDependenciesRequest, opts ...grpc.CallOption) (Consumer_RetrieveContainerDependenciesClient, error) {
	stream, err := c.cc.NewStream(ctx, &Consumer_ServiceDesc.Streams[1], "/terrarium.module.Consumer/RetrieveContainerDependencies", opts...)
	if err != nil {
		return nil, err
	}
	x := &consumerRetrieveContainerDependenciesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Consumer_RetrieveContainerDependenciesClient interface {
	Recv() (*ContainerDependenciesResponse, error)
	grpc.ClientStream
}

type consumerRetrieveContainerDependenciesClient struct {
	grpc.ClientStream
}

func (x *consumerRetrieveContainerDependenciesClient) Recv() (*ContainerDependenciesResponse, error) {
	m := new(ContainerDependenciesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *consumerClient) RetrieveModuleDependencies(ctx context.Context, in *RetrieveModuleDependenciesRequest, opts ...grpc.CallOption) (Consumer_RetrieveModuleDependenciesClient, error) {
	stream, err := c.cc.NewStream(ctx, &Consumer_ServiceDesc.Streams[2], "/terrarium.module.Consumer/RetrieveModuleDependencies", opts...)
	if err != nil {
		return nil, err
	}
	x := &consumerRetrieveModuleDependenciesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Consumer_RetrieveModuleDependenciesClient interface {
	Recv() (*ModuleDependenciesResponse, error)
	grpc.ClientStream
}

type consumerRetrieveModuleDependenciesClient struct {
	grpc.ClientStream
}

func (x *consumerRetrieveModuleDependenciesClient) Recv() (*ModuleDependenciesResponse, error) {
	m := new(ModuleDependenciesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ConsumerServer is the server API for Consumer service.
// All implementations must embed UnimplementedConsumerServer
// for forward compatibility
type ConsumerServer interface {
	DownloadSourceZip(*DownloadSourceZipRequest, Consumer_DownloadSourceZipServer) error
	RetrieveContainerDependencies(*RetrieveContainerDependenciesRequest, Consumer_RetrieveContainerDependenciesServer) error
	RetrieveModuleDependencies(*RetrieveModuleDependenciesRequest, Consumer_RetrieveModuleDependenciesServer) error
	mustEmbedUnimplementedConsumerServer()
}

// UnimplementedConsumerServer must be embedded to have forward compatible implementations.
type UnimplementedConsumerServer struct {
}

func (UnimplementedConsumerServer) DownloadSourceZip(*DownloadSourceZipRequest, Consumer_DownloadSourceZipServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadSourceZip not implemented")
}
func (UnimplementedConsumerServer) RetrieveContainerDependencies(*RetrieveContainerDependenciesRequest, Consumer_RetrieveContainerDependenciesServer) error {
	return status.Errorf(codes.Unimplemented, "method RetrieveContainerDependencies not implemented")
}
func (UnimplementedConsumerServer) RetrieveModuleDependencies(*RetrieveModuleDependenciesRequest, Consumer_RetrieveModuleDependenciesServer) error {
	return status.Errorf(codes.Unimplemented, "method RetrieveModuleDependencies not implemented")
}
func (UnimplementedConsumerServer) mustEmbedUnimplementedConsumerServer() {}

// UnsafeConsumerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConsumerServer will
// result in compilation errors.
type UnsafeConsumerServer interface {
	mustEmbedUnimplementedConsumerServer()
}

func RegisterConsumerServer(s grpc.ServiceRegistrar, srv ConsumerServer) {
	s.RegisterService(&Consumer_ServiceDesc, srv)
}

func _Consumer_DownloadSourceZip_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadSourceZipRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ConsumerServer).DownloadSourceZip(m, &consumerDownloadSourceZipServer{stream})
}

type Consumer_DownloadSourceZipServer interface {
	Send(*SourceZipResponse) error
	grpc.ServerStream
}

type consumerDownloadSourceZipServer struct {
	grpc.ServerStream
}

func (x *consumerDownloadSourceZipServer) Send(m *SourceZipResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Consumer_RetrieveContainerDependencies_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RetrieveContainerDependenciesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ConsumerServer).RetrieveContainerDependencies(m, &consumerRetrieveContainerDependenciesServer{stream})
}

type Consumer_RetrieveContainerDependenciesServer interface {
	Send(*ContainerDependenciesResponse) error
	grpc.ServerStream
}

type consumerRetrieveContainerDependenciesServer struct {
	grpc.ServerStream
}

func (x *consumerRetrieveContainerDependenciesServer) Send(m *ContainerDependenciesResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Consumer_RetrieveModuleDependencies_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RetrieveModuleDependenciesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ConsumerServer).RetrieveModuleDependencies(m, &consumerRetrieveModuleDependenciesServer{stream})
}

type Consumer_RetrieveModuleDependenciesServer interface {
	Send(*ModuleDependenciesResponse) error
	grpc.ServerStream
}

type consumerRetrieveModuleDependenciesServer struct {
	grpc.ServerStream
}

func (x *consumerRetrieveModuleDependenciesServer) Send(m *ModuleDependenciesResponse) error {
	return x.ServerStream.SendMsg(m)
}

// Consumer_ServiceDesc is the grpc.ServiceDesc for Consumer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Consumer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "terrarium.module.Consumer",
	HandlerType: (*ConsumerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "DownloadSourceZip",
			Handler:       _Consumer_DownloadSourceZip_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "RetrieveContainerDependencies",
			Handler:       _Consumer_RetrieveContainerDependencies_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "RetrieveModuleDependencies",
			Handler:       _Consumer_RetrieveModuleDependencies_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pb/terrarium/module/module.proto",
}
