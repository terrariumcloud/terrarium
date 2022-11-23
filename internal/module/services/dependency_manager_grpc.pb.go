// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: pb/terrarium/module/services/dependency_manager.proto

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

// DependencyManagerClient is the client API for DependencyManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DependencyManagerClient interface {
	RegisterModuleDependencies(ctx context.Context, in *module.RegisterModuleDependenciesRequest, opts ...grpc.CallOption) (*module.Response, error)
	RegisterContainerDependencies(ctx context.Context, in *module.RegisterContainerDependenciesRequest, opts ...grpc.CallOption) (*module.Response, error)
	RetrieveContainerDependencies(ctx context.Context, in *module.RetrieveContainerDependenciesRequestV2, opts ...grpc.CallOption) (DependencyManager_RetrieveContainerDependenciesClient, error)
	RetrieveModuleDependencies(ctx context.Context, in *module.RetrieveModuleDependenciesRequest, opts ...grpc.CallOption) (DependencyManager_RetrieveModuleDependenciesClient, error)
}

type dependencyManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewDependencyManagerClient(cc grpc.ClientConnInterface) DependencyManagerClient {
	return &dependencyManagerClient{cc}
}

func (c *dependencyManagerClient) RegisterModuleDependencies(ctx context.Context, in *module.RegisterModuleDependenciesRequest, opts ...grpc.CallOption) (*module.Response, error) {
	out := new(module.Response)
	err := c.cc.Invoke(ctx, "/terrarium.module.services.DependencyManager/RegisterModuleDependencies", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dependencyManagerClient) RegisterContainerDependencies(ctx context.Context, in *module.RegisterContainerDependenciesRequest, opts ...grpc.CallOption) (*module.Response, error) {
	out := new(module.Response)
	err := c.cc.Invoke(ctx, "/terrarium.module.services.DependencyManager/RegisterContainerDependencies", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dependencyManagerClient) RetrieveContainerDependencies(ctx context.Context, in *module.RetrieveContainerDependenciesRequestV2, opts ...grpc.CallOption) (DependencyManager_RetrieveContainerDependenciesClient, error) {
	stream, err := c.cc.NewStream(ctx, &DependencyManager_ServiceDesc.Streams[0], "/terrarium.module.services.DependencyManager/RetrieveContainerDependencies", opts...)
	if err != nil {
		return nil, err
	}
	x := &dependencyManagerRetrieveContainerDependenciesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DependencyManager_RetrieveContainerDependenciesClient interface {
	Recv() (*module.ContainerDependenciesResponseV2, error)
	grpc.ClientStream
}

type dependencyManagerRetrieveContainerDependenciesClient struct {
	grpc.ClientStream
}

func (x *dependencyManagerRetrieveContainerDependenciesClient) Recv() (*module.ContainerDependenciesResponseV2, error) {
	m := new(module.ContainerDependenciesResponseV2)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *dependencyManagerClient) RetrieveModuleDependencies(ctx context.Context, in *module.RetrieveModuleDependenciesRequest, opts ...grpc.CallOption) (DependencyManager_RetrieveModuleDependenciesClient, error) {
	stream, err := c.cc.NewStream(ctx, &DependencyManager_ServiceDesc.Streams[1], "/terrarium.module.services.DependencyManager/RetrieveModuleDependencies", opts...)
	if err != nil {
		return nil, err
	}
	x := &dependencyManagerRetrieveModuleDependenciesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DependencyManager_RetrieveModuleDependenciesClient interface {
	Recv() (*module.ModuleDependenciesResponse, error)
	grpc.ClientStream
}

type dependencyManagerRetrieveModuleDependenciesClient struct {
	grpc.ClientStream
}

func (x *dependencyManagerRetrieveModuleDependenciesClient) Recv() (*module.ModuleDependenciesResponse, error) {
	m := new(module.ModuleDependenciesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// DependencyManagerServer is the server API for DependencyManager service.
// All implementations must embed UnimplementedDependencyManagerServer
// for forward compatibility
type DependencyManagerServer interface {
	RegisterModuleDependencies(context.Context, *module.RegisterModuleDependenciesRequest) (*module.Response, error)
	RegisterContainerDependencies(context.Context, *module.RegisterContainerDependenciesRequest) (*module.Response, error)
	RetrieveContainerDependencies(*module.RetrieveContainerDependenciesRequestV2, DependencyManager_RetrieveContainerDependenciesServer) error
	RetrieveModuleDependencies(*module.RetrieveModuleDependenciesRequest, DependencyManager_RetrieveModuleDependenciesServer) error
	mustEmbedUnimplementedDependencyManagerServer()
}

// UnimplementedDependencyManagerServer must be embedded to have forward compatible implementations.
type UnimplementedDependencyManagerServer struct {
}

func (UnimplementedDependencyManagerServer) RegisterModuleDependencies(context.Context, *module.RegisterModuleDependenciesRequest) (*module.Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterModuleDependencies not implemented")
}
func (UnimplementedDependencyManagerServer) RegisterContainerDependencies(context.Context, *module.RegisterContainerDependenciesRequest) (*module.Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterContainerDependencies not implemented")
}
func (UnimplementedDependencyManagerServer) RetrieveContainerDependencies(*module.RetrieveContainerDependenciesRequestV2, DependencyManager_RetrieveContainerDependenciesServer) error {
	return status.Errorf(codes.Unimplemented, "method RetrieveContainerDependencies not implemented")
}
func (UnimplementedDependencyManagerServer) RetrieveModuleDependencies(*module.RetrieveModuleDependenciesRequest, DependencyManager_RetrieveModuleDependenciesServer) error {
	return status.Errorf(codes.Unimplemented, "method RetrieveModuleDependencies not implemented")
}
func (UnimplementedDependencyManagerServer) mustEmbedUnimplementedDependencyManagerServer() {}

// UnsafeDependencyManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DependencyManagerServer will
// result in compilation errors.
type UnsafeDependencyManagerServer interface {
	mustEmbedUnimplementedDependencyManagerServer()
}

func RegisterDependencyManagerServer(s grpc.ServiceRegistrar, srv DependencyManagerServer) {
	s.RegisterService(&DependencyManager_ServiceDesc, srv)
}

func _DependencyManager_RegisterModuleDependencies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(module.RegisterModuleDependenciesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyManagerServer).RegisterModuleDependencies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terrarium.module.services.DependencyManager/RegisterModuleDependencies",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyManagerServer).RegisterModuleDependencies(ctx, req.(*module.RegisterModuleDependenciesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DependencyManager_RegisterContainerDependencies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(module.RegisterContainerDependenciesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DependencyManagerServer).RegisterContainerDependencies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terrarium.module.services.DependencyManager/RegisterContainerDependencies",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DependencyManagerServer).RegisterContainerDependencies(ctx, req.(*module.RegisterContainerDependenciesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DependencyManager_RetrieveContainerDependencies_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(module.RetrieveContainerDependenciesRequestV2)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DependencyManagerServer).RetrieveContainerDependencies(m, &dependencyManagerRetrieveContainerDependenciesServer{stream})
}

type DependencyManager_RetrieveContainerDependenciesServer interface {
	Send(*module.ContainerDependenciesResponseV2) error
	grpc.ServerStream
}

type dependencyManagerRetrieveContainerDependenciesServer struct {
	grpc.ServerStream
}

func (x *dependencyManagerRetrieveContainerDependenciesServer) Send(m *module.ContainerDependenciesResponseV2) error {
	return x.ServerStream.SendMsg(m)
}

func _DependencyManager_RetrieveModuleDependencies_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(module.RetrieveModuleDependenciesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DependencyManagerServer).RetrieveModuleDependencies(m, &dependencyManagerRetrieveModuleDependenciesServer{stream})
}

type DependencyManager_RetrieveModuleDependenciesServer interface {
	Send(*module.ModuleDependenciesResponse) error
	grpc.ServerStream
}

type dependencyManagerRetrieveModuleDependenciesServer struct {
	grpc.ServerStream
}

func (x *dependencyManagerRetrieveModuleDependenciesServer) Send(m *module.ModuleDependenciesResponse) error {
	return x.ServerStream.SendMsg(m)
}

// DependencyManager_ServiceDesc is the grpc.ServiceDesc for DependencyManager service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DependencyManager_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "terrarium.module.services.DependencyManager",
	HandlerType: (*DependencyManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterModuleDependencies",
			Handler:    _DependencyManager_RegisterModuleDependencies_Handler,
		},
		{
			MethodName: "RegisterContainerDependencies",
			Handler:    _DependencyManager_RegisterContainerDependencies_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RetrieveContainerDependencies",
			Handler:       _DependencyManager_RetrieveContainerDependencies_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "RetrieveModuleDependencies",
			Handler:       _DependencyManager_RetrieveModuleDependencies_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pb/terrarium/module/services/dependency_manager.proto",
}
