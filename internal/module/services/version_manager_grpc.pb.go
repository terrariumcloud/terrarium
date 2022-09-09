// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.5
// source: pb/terrarium/module/services/version_manager.proto

package services

import (
	context "context"
	module "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// VersionManagerClient is the client API for VersionManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VersionManagerClient interface {
	BeginVersion(ctx context.Context, in *BeginVersionRequest, opts ...grpc.CallOption) (*module.BeginVersionResponse, error)
	AbortVersion(ctx context.Context, in *TerminateVersionRequest, opts ...grpc.CallOption) (*module.TransactionStatusResponse, error)
	PublishVersion(ctx context.Context, in *TerminateVersionRequest, opts ...grpc.CallOption) (*module.TransactionStatusResponse, error)
	ListModuleVersions(ctx context.Context, in *ListModuleVersionsRequest, opts ...grpc.CallOption) (*ListModuleVersionsResponse, error)
}

type versionManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewVersionManagerClient(cc grpc.ClientConnInterface) VersionManagerClient {
	return &versionManagerClient{cc}
}

func (c *versionManagerClient) BeginVersion(ctx context.Context, in *BeginVersionRequest, opts ...grpc.CallOption) (*module.BeginVersionResponse, error) {
	out := new(module.BeginVersionResponse)
	err := c.cc.Invoke(ctx, "/terrarium.module.services.VersionManager/BeginVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *versionManagerClient) AbortVersion(ctx context.Context, in *TerminateVersionRequest, opts ...grpc.CallOption) (*module.TransactionStatusResponse, error) {
	out := new(module.TransactionStatusResponse)
	err := c.cc.Invoke(ctx, "/terrarium.module.services.VersionManager/AbortVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *versionManagerClient) PublishVersion(ctx context.Context, in *TerminateVersionRequest, opts ...grpc.CallOption) (*module.TransactionStatusResponse, error) {
	out := new(module.TransactionStatusResponse)
	err := c.cc.Invoke(ctx, "/terrarium.module.services.VersionManager/PublishVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *versionManagerClient) ListModuleVersions(ctx context.Context, in *ListModuleVersionsRequest, opts ...grpc.CallOption) (*ListModuleVersionsResponse, error) {
	out := new(ListModuleVersionsResponse)
	err := c.cc.Invoke(ctx, "/terrarium.module.services.VersionManager/ListModuleVersions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VersionManagerServer is the server API for VersionManager service.
// All implementations must embed UnimplementedVersionManagerServer
// for forward compatibility
type VersionManagerServer interface {
	BeginVersion(context.Context, *BeginVersionRequest) (*module.BeginVersionResponse, error)
	AbortVersion(context.Context, *TerminateVersionRequest) (*module.TransactionStatusResponse, error)
	PublishVersion(context.Context, *TerminateVersionRequest) (*module.TransactionStatusResponse, error)
	ListModuleVersions(context.Context, *ListModuleVersionsRequest) (*ListModuleVersionsResponse, error)
	mustEmbedUnimplementedVersionManagerServer()
}

// UnimplementedVersionManagerServer must be embedded to have forward compatible implementations.
type UnimplementedVersionManagerServer struct {
}

func (UnimplementedVersionManagerServer) BeginVersion(context.Context, *BeginVersionRequest) (*module.BeginVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BeginVersion not implemented")
}
func (UnimplementedVersionManagerServer) AbortVersion(context.Context, *TerminateVersionRequest) (*module.TransactionStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AbortVersion not implemented")
}
func (UnimplementedVersionManagerServer) PublishVersion(context.Context, *TerminateVersionRequest) (*module.TransactionStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishVersion not implemented")
}
func (UnimplementedVersionManagerServer) ListModuleVersions(context.Context, *ListModuleVersionsRequest) (*ListModuleVersionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListModuleVersions not implemented")
}
func (UnimplementedVersionManagerServer) mustEmbedUnimplementedVersionManagerServer() {}

// UnsafeVersionManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VersionManagerServer will
// result in compilation errors.
type UnsafeVersionManagerServer interface {
	mustEmbedUnimplementedVersionManagerServer()
}

func RegisterVersionManagerServer(s grpc.ServiceRegistrar, srv VersionManagerServer) {
	s.RegisterService(&VersionManager_ServiceDesc, srv)
}

func _VersionManager_BeginVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BeginVersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VersionManagerServer).BeginVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terrarium.module.services.VersionManager/BeginVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VersionManagerServer).BeginVersion(ctx, req.(*BeginVersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VersionManager_AbortVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TerminateVersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VersionManagerServer).AbortVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terrarium.module.services.VersionManager/AbortVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VersionManagerServer).AbortVersion(ctx, req.(*TerminateVersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VersionManager_PublishVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TerminateVersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VersionManagerServer).PublishVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terrarium.module.services.VersionManager/PublishVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VersionManagerServer).PublishVersion(ctx, req.(*TerminateVersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VersionManager_ListModuleVersions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListModuleVersionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VersionManagerServer).ListModuleVersions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terrarium.module.services.VersionManager/ListModuleVersions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VersionManagerServer).ListModuleVersions(ctx, req.(*ListModuleVersionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// VersionManager_ServiceDesc is the grpc.ServiceDesc for VersionManager service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VersionManager_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "terrarium.module.services.VersionManager",
	HandlerType: (*VersionManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "BeginVersion",
			Handler:    _VersionManager_BeginVersion_Handler,
		},
		{
			MethodName: "AbortVersion",
			Handler:    _VersionManager_AbortVersion_Handler,
		},
		{
			MethodName: "PublishVersion",
			Handler:    _VersionManager_PublishVersion_Handler,
		},
		{
			MethodName: "ListModuleVersions",
			Handler:    _VersionManager_ListModuleVersions_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/terrarium/module/services/version_manager.proto",
}
