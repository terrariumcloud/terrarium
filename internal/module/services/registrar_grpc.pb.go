// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: pb/terrarium/module/services/registrar.proto

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

// RegistrarClient is the client API for Registrar service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RegistrarClient interface {
	Register(ctx context.Context, in *module.RegisterModuleRequest, opts ...grpc.CallOption) (*module.Response, error)
	ListModules(ctx context.Context, in *ListModulesRequest, opts ...grpc.CallOption) (*ListModulesResponse, error)
	GetModule(ctx context.Context, in *GetModuleRequest, opts ...grpc.CallOption) (*GetModuleResponse, error)
}

type registrarClient struct {
	cc grpc.ClientConnInterface
}

func NewRegistrarClient(cc grpc.ClientConnInterface) RegistrarClient {
	return &registrarClient{cc}
}

func (c *registrarClient) Register(ctx context.Context, in *module.RegisterModuleRequest, opts ...grpc.CallOption) (*module.Response, error) {
	out := new(module.Response)
	err := c.cc.Invoke(ctx, "/terrarium.module.services.Registrar/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registrarClient) ListModules(ctx context.Context, in *ListModulesRequest, opts ...grpc.CallOption) (*ListModulesResponse, error) {
	out := new(ListModulesResponse)
	err := c.cc.Invoke(ctx, "/terrarium.module.services.Registrar/ListModules", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registrarClient) GetModule(ctx context.Context, in *GetModuleRequest, opts ...grpc.CallOption) (*GetModuleResponse, error) {
	out := new(GetModuleResponse)
	err := c.cc.Invoke(ctx, "/terrarium.module.services.Registrar/GetModule", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RegistrarServer is the server API for Registrar service.
// All implementations must embed UnimplementedRegistrarServer
// for forward compatibility
type RegistrarServer interface {
	Register(context.Context, *module.RegisterModuleRequest) (*module.Response, error)
	ListModules(context.Context, *ListModulesRequest) (*ListModulesResponse, error)
	GetModule(context.Context, *GetModuleRequest) (*GetModuleResponse, error)
	mustEmbedUnimplementedRegistrarServer()
}

// UnimplementedRegistrarServer must be embedded to have forward compatible implementations.
type UnimplementedRegistrarServer struct {
}

func (UnimplementedRegistrarServer) Register(context.Context, *module.RegisterModuleRequest) (*module.Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedRegistrarServer) ListModules(context.Context, *ListModulesRequest) (*ListModulesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListModules not implemented")
}
func (UnimplementedRegistrarServer) GetModule(context.Context, *GetModuleRequest) (*GetModuleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetModule not implemented")
}
func (UnimplementedRegistrarServer) mustEmbedUnimplementedRegistrarServer() {}

// UnsafeRegistrarServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RegistrarServer will
// result in compilation errors.
type UnsafeRegistrarServer interface {
	mustEmbedUnimplementedRegistrarServer()
}

func RegisterRegistrarServer(s grpc.ServiceRegistrar, srv RegistrarServer) {
	s.RegisterService(&Registrar_ServiceDesc, srv)
}

func _Registrar_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(module.RegisterModuleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistrarServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terrarium.module.services.Registrar/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistrarServer).Register(ctx, req.(*module.RegisterModuleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Registrar_ListModules_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListModulesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistrarServer).ListModules(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terrarium.module.services.Registrar/ListModules",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistrarServer).ListModules(ctx, req.(*ListModulesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Registrar_GetModule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetModuleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistrarServer).GetModule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terrarium.module.services.Registrar/GetModule",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistrarServer).GetModule(ctx, req.(*GetModuleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Registrar_ServiceDesc is the grpc.ServiceDesc for Registrar service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Registrar_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "terrarium.module.services.Registrar",
	HandlerType: (*RegistrarServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _Registrar_Register_Handler,
		},
		{
			MethodName: "ListModules",
			Handler:    _Registrar_ListModules_Handler,
		},
		{
			MethodName: "GetModule",
			Handler:    _Registrar_GetModule_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/terrarium/module/services/registrar.proto",
}
