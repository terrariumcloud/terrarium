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

func (m *MockRegistrarClient) Register(ctx context.Context, in *module.RegisterModuleRequest, opts ...grpc.CallOption) (*module.Response, error) {
	m.RegisterInvocations++
	return m.RegisterResponse, m.RegisterError
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

func (m *MockVersionManagerClient) BeginVersion(ctx context.Context, in *module.BeginVersionRequest, opts ...grpc.CallOption) (*module.Response, error) {
	m.BeginVersionInvocations++
	return m.BeginVersionResponse, m.BeginVersionError
}

func (m *MockVersionManagerClient) PublishVersion(ctx context.Context, in *services.TerminateVersionRequest, opts ...grpc.CallOption) (*module.Response, error) {
	m.PublishVersionInvocations++
	return m.PublishVersionResponse, m.PublishVersionError
}

func (m *MockVersionManagerClient) AbortVersion(ctx context.Context, in *services.TerminateVersionRequest, opts ...grpc.CallOption) (*module.Response, error) {
	m.AbortVersionInvocations++
	return m.AbortVersionResponse, m.AbortVersionError
}

type MockStorageClient struct {
	services.StorageClient
	UploadSourceZipInvocations   int
	UploadSourceZipClient        services.Storage_UploadSourceZipClient
	UploadSourceZipError         error
	DownloadSourceZipInvocations int
	DownloadSourceZipClient      services.Storage_DownloadSourceZipClient
	DownloadSourceZipError       error
}

func (m *MockStorageClient) UploadSourceZip(ctx context.Context, opts ...grpc.CallOption) (services.Storage_UploadSourceZipClient, error) {
	m.UploadSourceZipInvocations++
	return m.UploadSourceZipClient, m.UploadSourceZipError
}

func (m *MockStorageClient) DownloadSourceZip(ctx context.Context, in *module.DownloadSourceZipRequest, opts ...grpc.CallOption) (services.Storage_DownloadSourceZipClient, error) {
	m.DownloadSourceZipInvocations++
	return m.DownloadSourceZipClient, m.DownloadSourceZipError
}

type MockStorage_UploadSourceZipClient struct {
	services.Storage_UploadSourceZipClient
	CloseAndRecvInvocations int
	CloseAndRecvResponse    *module.Response
	CloseAndRecvError       error
	SendInvocations         int
	SendError               error
	CloseSendInvocations    int
	CloseSendError          error
}

func (m *MockStorage_UploadSourceZipClient) CloseAndRecv() (*module.Response, error) {
	m.CloseAndRecvInvocations++
	return m.CloseAndRecvResponse, m.CloseAndRecvError
}

func (m *MockStorage_UploadSourceZipClient) Send(*module.UploadSourceZipRequest) error {
	m.SendInvocations++
	return m.SendError
}

func (m *MockStorage_UploadSourceZipClient) CloseSend() error {
	m.CloseSendInvocations++
	return m.CloseSendError
}

type MockStorage_DownloadSourceZipClient struct {
	services.Storage_DownloadSourceZipClient
	RecvInvocations      int
	RecvResponse         *module.SourceZipResponse
	RecvError            error
	CloseSendInvocations int
	CloseSendError       error
}

func (m *MockStorage_DownloadSourceZipClient) Recv() (*module.SourceZipResponse, error) {
	m.RecvInvocations++
	return m.RecvResponse, m.RecvError
}

func (m *MockStorage_DownloadSourceZipClient) CloseSend() error {
	m.CloseSendInvocations++
	return m.CloseSendError
}

type MockDependencyManagerClient struct {
	services.DependencyManagerClient
	RegisterModuleDependenciesInvocations    int
	RegisterModuleDependenciesResponse       *module.Response
	RegisterModuleDependenciesError          error
	RegisterContainerDependenciesInvocations int
	RegisterContainerDependenciesResponse    *module.Response
	RegisterContainerDependenciesError       error
	RetrieveContainerDependenciesInvocations int
	RetrieveContainerDependenciesClient      services.DependencyManager_RetrieveContainerDependenciesClient
	RetrieveContainerDependenciesError       error
	RetrieveModuleDependenciesInvocations    int
	RetrieveModuleDependenciesClient         services.DependencyManager_RetrieveModuleDependenciesClient
	RetrieveModuleDependenciesError          error
}

func (m *MockDependencyManagerClient) RegisterModuleDependencies(ctx context.Context, in *module.RegisterModuleDependenciesRequest, opts ...grpc.CallOption) (*module.Response, error) {
	m.RegisterModuleDependenciesInvocations++
	return m.RegisterModuleDependenciesResponse, m.RegisterModuleDependenciesError
}

func (m *MockDependencyManagerClient) RegisterContainerDependencies(ctx context.Context, in *module.RegisterContainerDependenciesRequest, opts ...grpc.CallOption) (*module.Response, error) {
	m.RegisterContainerDependenciesInvocations++
	return m.RegisterContainerDependenciesResponse, m.RegisterContainerDependenciesError
}

func (m *MockDependencyManagerClient) RetrieveContainerDependencies(ctx context.Context, in *module.RetrieveContainerDependenciesRequestV2, opts ...grpc.CallOption) (services.DependencyManager_RetrieveContainerDependenciesClient, error) {
	m.RetrieveContainerDependenciesInvocations++
	return m.RetrieveContainerDependenciesClient, m.RetrieveContainerDependenciesError
}
func (m *MockDependencyManagerClient) RetrieveModuleDependencies(ctx context.Context, in *module.RetrieveModuleDependenciesRequest, opts ...grpc.CallOption) (services.DependencyManager_RetrieveModuleDependenciesClient, error) {
	m.RetrieveModuleDependenciesInvocations++
	return m.RetrieveModuleDependenciesClient, m.RetrieveModuleDependenciesError
}

type MockDependencyManager_RetrieveContainerDependenciesClient struct {
	services.DependencyManager_RetrieveContainerDependenciesClient
	RecvInvocations      int
	RecvResponse         *module.ContainerDependenciesResponseV2
	RecvError            error
	CloseSendInvocations int
	CloseSendError       error
}

func (m *MockDependencyManager_RetrieveContainerDependenciesClient) Recv() (*module.ContainerDependenciesResponseV2, error) {
	m.RecvInvocations++
	return m.RecvResponse, m.RecvError
}

func (m *MockDependencyManager_RetrieveContainerDependenciesClient) CloseSend() error {
	m.CloseSendInvocations++
	return m.CloseSendError
}

type MockDependencyManager_RetrieveModuleDependenciesClient struct {
	services.DependencyManager_RetrieveModuleDependenciesClient
	RecvInvocations      int
	RecvResponse         *module.ModuleDependenciesResponse
	RecvError            error
	CloseSendInvocations int
	CloseSendError       error
}

func (m *MockDependencyManager_RetrieveModuleDependenciesClient) Recv() (*module.ModuleDependenciesResponse, error) {
	m.RecvInvocations++
	return m.RecvResponse, m.RecvError
}

func (m *MockDependencyManager_RetrieveModuleDependenciesClient) CloseSend() error {
	m.CloseSendInvocations++
	return m.CloseSendError
}
