// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: pb/terrarium/provider/provider.proto

package provider

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

// ProviderPublisherClient is the client API for ProviderPublisher service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProviderPublisherClient interface {
	UploadProviderSourceZip(ctx context.Context, opts ...grpc.CallOption) (ProviderPublisher_UploadProviderSourceZipClient, error)
	UploadShasum(ctx context.Context, opts ...grpc.CallOption) (ProviderPublisher_UploadShasumClient, error)
	UploadShasumSignature(ctx context.Context, opts ...grpc.CallOption) (ProviderPublisher_UploadShasumSignatureClient, error)
	RegisterProvider(ctx context.Context, in *RegisterProviderRequest, opts ...grpc.CallOption) (*Response, error)
	EndProvider(ctx context.Context, in *EndProviderRequest, opts ...grpc.CallOption) (*Response, error)
}

type providerPublisherClient struct {
	cc grpc.ClientConnInterface
}

func NewProviderPublisherClient(cc grpc.ClientConnInterface) ProviderPublisherClient {
	return &providerPublisherClient{cc}
}

func (c *providerPublisherClient) UploadProviderSourceZip(ctx context.Context, opts ...grpc.CallOption) (ProviderPublisher_UploadProviderSourceZipClient, error) {
	stream, err := c.cc.NewStream(ctx, &ProviderPublisher_ServiceDesc.Streams[0], "/terrarium.provider.ProviderPublisher/UploadProviderSourceZip", opts...)
	if err != nil {
		return nil, err
	}
	x := &providerPublisherUploadProviderSourceZipClient{stream}
	return x, nil
}

type ProviderPublisher_UploadProviderSourceZipClient interface {
	Send(*UploadProviderSourceZipRequest) error
	CloseAndRecv() (*Response, error)
	grpc.ClientStream
}

type providerPublisherUploadProviderSourceZipClient struct {
	grpc.ClientStream
}

func (x *providerPublisherUploadProviderSourceZipClient) Send(m *UploadProviderSourceZipRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *providerPublisherUploadProviderSourceZipClient) CloseAndRecv() (*Response, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *providerPublisherClient) UploadShasum(ctx context.Context, opts ...grpc.CallOption) (ProviderPublisher_UploadShasumClient, error) {
	stream, err := c.cc.NewStream(ctx, &ProviderPublisher_ServiceDesc.Streams[1], "/terrarium.provider.ProviderPublisher/UploadShasum", opts...)
	if err != nil {
		return nil, err
	}
	x := &providerPublisherUploadShasumClient{stream}
	return x, nil
}

type ProviderPublisher_UploadShasumClient interface {
	Send(*UploadShasumRequest) error
	CloseAndRecv() (*Response, error)
	grpc.ClientStream
}

type providerPublisherUploadShasumClient struct {
	grpc.ClientStream
}

func (x *providerPublisherUploadShasumClient) Send(m *UploadShasumRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *providerPublisherUploadShasumClient) CloseAndRecv() (*Response, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *providerPublisherClient) UploadShasumSignature(ctx context.Context, opts ...grpc.CallOption) (ProviderPublisher_UploadShasumSignatureClient, error) {
	stream, err := c.cc.NewStream(ctx, &ProviderPublisher_ServiceDesc.Streams[2], "/terrarium.provider.ProviderPublisher/UploadShasumSignature", opts...)
	if err != nil {
		return nil, err
	}
	x := &providerPublisherUploadShasumSignatureClient{stream}
	return x, nil
}

type ProviderPublisher_UploadShasumSignatureClient interface {
	Send(*UploadShasumRequest) error
	CloseAndRecv() (*Response, error)
	grpc.ClientStream
}

type providerPublisherUploadShasumSignatureClient struct {
	grpc.ClientStream
}

func (x *providerPublisherUploadShasumSignatureClient) Send(m *UploadShasumRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *providerPublisherUploadShasumSignatureClient) CloseAndRecv() (*Response, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *providerPublisherClient) RegisterProvider(ctx context.Context, in *RegisterProviderRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/terrarium.provider.ProviderPublisher/RegisterProvider", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *providerPublisherClient) EndProvider(ctx context.Context, in *EndProviderRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/terrarium.provider.ProviderPublisher/EndProvider", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProviderPublisherServer is the server API for ProviderPublisher service.
// All implementations must embed UnimplementedProviderPublisherServer
// for forward compatibility
type ProviderPublisherServer interface {
	UploadProviderSourceZip(ProviderPublisher_UploadProviderSourceZipServer) error
	UploadShasum(ProviderPublisher_UploadShasumServer) error
	UploadShasumSignature(ProviderPublisher_UploadShasumSignatureServer) error
	RegisterProvider(context.Context, *RegisterProviderRequest) (*Response, error)
	EndProvider(context.Context, *EndProviderRequest) (*Response, error)
	mustEmbedUnimplementedProviderPublisherServer()
}

// UnimplementedProviderPublisherServer must be embedded to have forward compatible implementations.
type UnimplementedProviderPublisherServer struct {
}

func (UnimplementedProviderPublisherServer) UploadProviderSourceZip(ProviderPublisher_UploadProviderSourceZipServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadProviderSourceZip not implemented")
}
func (UnimplementedProviderPublisherServer) UploadShasum(ProviderPublisher_UploadShasumServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadShasum not implemented")
}
func (UnimplementedProviderPublisherServer) UploadShasumSignature(ProviderPublisher_UploadShasumSignatureServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadShasumSignature not implemented")
}
func (UnimplementedProviderPublisherServer) RegisterProvider(context.Context, *RegisterProviderRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterProvider not implemented")
}
func (UnimplementedProviderPublisherServer) EndProvider(context.Context, *EndProviderRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EndProvider not implemented")
}
func (UnimplementedProviderPublisherServer) mustEmbedUnimplementedProviderPublisherServer() {}

// UnsafeProviderPublisherServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProviderPublisherServer will
// result in compilation errors.
type UnsafeProviderPublisherServer interface {
	mustEmbedUnimplementedProviderPublisherServer()
}

func RegisterProviderPublisherServer(s grpc.ServiceRegistrar, srv ProviderPublisherServer) {
	s.RegisterService(&ProviderPublisher_ServiceDesc, srv)
}

func _ProviderPublisher_UploadProviderSourceZip_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ProviderPublisherServer).UploadProviderSourceZip(&providerPublisherUploadProviderSourceZipServer{stream})
}

type ProviderPublisher_UploadProviderSourceZipServer interface {
	SendAndClose(*Response) error
	Recv() (*UploadProviderSourceZipRequest, error)
	grpc.ServerStream
}

type providerPublisherUploadProviderSourceZipServer struct {
	grpc.ServerStream
}

func (x *providerPublisherUploadProviderSourceZipServer) SendAndClose(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *providerPublisherUploadProviderSourceZipServer) Recv() (*UploadProviderSourceZipRequest, error) {
	m := new(UploadProviderSourceZipRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ProviderPublisher_UploadShasum_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ProviderPublisherServer).UploadShasum(&providerPublisherUploadShasumServer{stream})
}

type ProviderPublisher_UploadShasumServer interface {
	SendAndClose(*Response) error
	Recv() (*UploadShasumRequest, error)
	grpc.ServerStream
}

type providerPublisherUploadShasumServer struct {
	grpc.ServerStream
}

func (x *providerPublisherUploadShasumServer) SendAndClose(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *providerPublisherUploadShasumServer) Recv() (*UploadShasumRequest, error) {
	m := new(UploadShasumRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ProviderPublisher_UploadShasumSignature_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ProviderPublisherServer).UploadShasumSignature(&providerPublisherUploadShasumSignatureServer{stream})
}

type ProviderPublisher_UploadShasumSignatureServer interface {
	SendAndClose(*Response) error
	Recv() (*UploadShasumRequest, error)
	grpc.ServerStream
}

type providerPublisherUploadShasumSignatureServer struct {
	grpc.ServerStream
}

func (x *providerPublisherUploadShasumSignatureServer) SendAndClose(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *providerPublisherUploadShasumSignatureServer) Recv() (*UploadShasumRequest, error) {
	m := new(UploadShasumRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ProviderPublisher_RegisterProvider_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterProviderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProviderPublisherServer).RegisterProvider(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terrarium.provider.ProviderPublisher/RegisterProvider",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProviderPublisherServer).RegisterProvider(ctx, req.(*RegisterProviderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProviderPublisher_EndProvider_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EndProviderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProviderPublisherServer).EndProvider(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terrarium.provider.ProviderPublisher/EndProvider",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProviderPublisherServer).EndProvider(ctx, req.(*EndProviderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ProviderPublisher_ServiceDesc is the grpc.ServiceDesc for ProviderPublisher service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProviderPublisher_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "terrarium.provider.ProviderPublisher",
	HandlerType: (*ProviderPublisherServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterProvider",
			Handler:    _ProviderPublisher_RegisterProvider_Handler,
		},
		{
			MethodName: "EndProvider",
			Handler:    _ProviderPublisher_EndProvider_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadProviderSourceZip",
			Handler:       _ProviderPublisher_UploadProviderSourceZip_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "UploadShasum",
			Handler:       _ProviderPublisher_UploadShasum_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "UploadShasumSignature",
			Handler:       _ProviderPublisher_UploadShasumSignature_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "pb/terrarium/provider/provider.proto",
}
