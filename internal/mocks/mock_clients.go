package mocks

import (
	"context"

	"github.com/terrariumcloud/terrarium/internal/module/services"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"google.golang.org/grpc"
)

type MockRegistrarClient struct {
	services.RegistrarClient
	RegisterInvocations int
	RegisterResponse    *module.Response
	RegisterError       error
}

func (mrc *MockRegistrarClient) Register(ctx context.Context, in *module.RegisterModuleRequest, opts ...grpc.CallOption) (*module.Response, error) {
	mrc.RegisterInvocations++
	return mrc.RegisterResponse, mrc.RegisterError
}

type MockVersionManagerClient struct {
	services.VersionManagerClient
	BeginVersionInvocations   int
	BeginVersionResponse      *module.Response
	BeginVersionError         error
	PublishVersionInvocations int
	PublishVersionResponse    *module.Response
	PublishVersionError       error
	AbortVersionInvocations   int
	AbortVersionResponse      *module.Response
	AbortVersionError         error
}

func (mvmc *MockVersionManagerClient) BeginVersion(ctx context.Context, in *module.BeginVersionRequest, opts ...grpc.CallOption) (*module.Response, error) {
	mvmc.BeginVersionInvocations++
	return mvmc.BeginVersionResponse, mvmc.BeginVersionError
}

func (mvmc *MockVersionManagerClient) PublishVersion(ctx context.Context, in *services.TerminateVersionRequest, opts ...grpc.CallOption) (*module.Response, error) {
	mvmc.PublishVersionInvocations++
	return mvmc.PublishVersionResponse, mvmc.PublishVersionError
}

func (mvmc *MockVersionManagerClient) AbortVersion(ctx context.Context, in *services.TerminateVersionRequest, opts ...grpc.CallOption) (*module.Response, error) {
	mvmc.AbortVersionInvocations++
	return mvmc.AbortVersionResponse, mvmc.AbortVersionError
}

type MockStorageClient struct {
	services.StorageClient
	UploadSourceZipInvocations int
	UploadSourceZipClient      services.Storage_UploadSourceZipClient
	UploadSourceZipError       error
}

func (msc *MockStorageClient) UploadSourceZip(ctx context.Context, opts ...grpc.CallOption) (services.Storage_UploadSourceZipClient, error) {
	msc.UploadSourceZipInvocations++
	return msc.UploadSourceZipClient, msc.UploadSourceZipError
}

type MockStorage_UploadSourceZipClient struct {
	services.Storage_UploadSourceZipClient
	CloseAndRecvInvocations int
	CloseAndRecvResponse    *module.Response
	CloseAndRecvError       error
}

func (msuszc *MockStorage_UploadSourceZipClient) CloseAndRecv() (*module.Response, error) {
	msuszc.CloseAndRecvInvocations++
	return msuszc.CloseAndRecvResponse, msuszc.CloseAndRecvError
}
