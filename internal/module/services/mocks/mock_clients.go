package mocks

import (
	"context"
	moduleServices "github.com/terrariumcloud/terrarium/internal/module/services"
	providerServices "github.com/terrariumcloud/terrarium/internal/provider/services"

	terrariumModule "github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	terrariumProvider "github.com/terrariumcloud/terrarium/pkg/terrarium/provider"
	"google.golang.org/grpc"
)

type MockRegistrarClient struct {
	moduleServices.RegistrarClient
	RegisterInvocations int
	RegisterResponse    *terrariumModule.Response
	RegisterError       error
}

func (m *MockRegistrarClient) Register(ctx context.Context, in *terrariumModule.RegisterModuleRequest, opts ...grpc.CallOption) (*terrariumModule.Response, error) {
	m.RegisterInvocations++
	return m.RegisterResponse, m.RegisterError
}

type MockVersionManagerClient struct {
	moduleServices.VersionManagerClient
	BeginVersionInvocations   int
	BeginVersionResponse      *terrariumModule.Response
	BeginVersionError         error
	PublishVersionInvocations int
	PublishVersionResponse    *terrariumModule.Response
	PublishVersionError       error
	AbortVersionInvocations   int
	AbortVersionResponse      *terrariumModule.Response
	AbortVersionError         error
}

func (m *MockVersionManagerClient) BeginVersion(ctx context.Context, in *terrariumModule.BeginVersionRequest, opts ...grpc.CallOption) (*terrariumModule.Response, error) {
	m.BeginVersionInvocations++
	return m.BeginVersionResponse, m.BeginVersionError
}

func (m *MockVersionManagerClient) PublishVersion(ctx context.Context, in *moduleServices.TerminateVersionRequest, opts ...grpc.CallOption) (*terrariumModule.Response, error) {
	m.PublishVersionInvocations++
	return m.PublishVersionResponse, m.PublishVersionError
}

func (m *MockVersionManagerClient) AbortVersion(ctx context.Context, in *moduleServices.TerminateVersionRequest, opts ...grpc.CallOption) (*terrariumModule.Response, error) {
	m.AbortVersionInvocations++
	return m.AbortVersionResponse, m.AbortVersionError
}

type MockProviderVersionManagerClient struct {
	providerServices.VersionManagerClient
	RegisterInvocations 	  int
	RegisterResponse 		  *terrariumProvider.Response
	RegisterError       	  error
	PublishVersionInvocations int
	PublishVersionResponse    *terrariumProvider.Response
	PublishVersionError       error
	AbortVersionInvocations   int
	AbortVersionResponse      *terrariumProvider.Response
	AbortVersionError         error
}

func (m *MockProviderVersionManagerClient) Register(ctx context.Context, in *terrariumProvider.RegisterProviderRequest, opts ...grpc.CallOption) (*terrariumProvider.Response, error) {
	m.RegisterInvocations++
	return m.RegisterResponse, m.RegisterError
}

func (m *MockProviderVersionManagerClient) PublishVersion(ctx context.Context, in *providerServices.TerminateVersionRequest, opts ...grpc.CallOption) (*terrariumProvider.Response, error) {
	m.PublishVersionInvocations++
	return m.PublishVersionResponse, m.PublishVersionError
}

func (m *MockProviderVersionManagerClient) AbortProviderVersion(ctx context.Context, in *providerServices.TerminateVersionRequest, opts ...grpc.CallOption) (*terrariumProvider.Response, error) {
	m.AbortVersionInvocations++
	return m.AbortVersionResponse, m.AbortVersionError
}

type MockStorageClient struct {
	moduleServices.StorageClient
	UploadSourceZipInvocations   int
	UploadSourceZipClient        moduleServices.Storage_UploadSourceZipClient
	UploadSourceZipError         error
	DownloadSourceZipInvocations int
	DownloadSourceZipClient      moduleServices.Storage_DownloadSourceZipClient
	DownloadSourceZipError       error
}

func (m *MockStorageClient) UploadSourceZip(ctx context.Context, opts ...grpc.CallOption) (moduleServices.Storage_UploadSourceZipClient, error) {
	m.UploadSourceZipInvocations++
	return m.UploadSourceZipClient, m.UploadSourceZipError
}

func (m *MockStorageClient) DownloadSourceZip(ctx context.Context, in *terrariumModule.DownloadSourceZipRequest, opts ...grpc.CallOption) (moduleServices.Storage_DownloadSourceZipClient, error) {
	m.DownloadSourceZipInvocations++
	return m.DownloadSourceZipClient, m.DownloadSourceZipError
}

type MockStorage_UploadSourceZipClient struct {
	moduleServices.Storage_UploadSourceZipClient
	CloseAndRecvInvocations int
	CloseAndRecvResponse    *terrariumModule.Response
	CloseAndRecvError       error
	SendInvocations         int
	SendError               error
	CloseSendInvocations    int
	CloseSendError          error
}

func (m *MockStorage_UploadSourceZipClient) CloseAndRecv() (*terrariumModule.Response, error) {
	m.CloseAndRecvInvocations++
	return m.CloseAndRecvResponse, m.CloseAndRecvError
}

func (m *MockStorage_UploadSourceZipClient) Send(*terrariumModule.UploadSourceZipRequest) error {
	m.SendInvocations++
	return m.SendError
}

func (m *MockStorage_UploadSourceZipClient) CloseSend() error {
	m.CloseSendInvocations++
	return m.CloseSendError
}

type MockStorage_DownloadSourceZipClient struct {
	moduleServices.Storage_DownloadSourceZipClient
	RecvInvocations      int
	RecvResponse         *terrariumModule.SourceZipResponse
	RecvError            error
	CloseSendInvocations int
	CloseSendError       error
}

func (m *MockStorage_DownloadSourceZipClient) Recv() (*terrariumModule.SourceZipResponse, error) {
	m.RecvInvocations++
	return m.RecvResponse, m.RecvError
}

func (m *MockStorage_DownloadSourceZipClient) CloseSend() error {
	m.CloseSendInvocations++
	return m.CloseSendError
}

type MockDependencyManagerClient struct {
	moduleServices.DependencyManagerClient
	RegisterModuleDependenciesInvocations    int
	RegisterModuleDependenciesResponse       *terrariumModule.Response
	RegisterModuleDependenciesError          error
	RegisterContainerDependenciesInvocations int
	RegisterContainerDependenciesResponse    *terrariumModule.Response
	RegisterContainerDependenciesError       error
	RetrieveContainerDependenciesInvocations int
	RetrieveContainerDependenciesClient      moduleServices.DependencyManager_RetrieveContainerDependenciesClient
	RetrieveContainerDependenciesError       error
	RetrieveModuleDependenciesInvocations    int
	RetrieveModuleDependenciesClient         moduleServices.DependencyManager_RetrieveModuleDependenciesClient
	RetrieveModuleDependenciesError          error
}

func (m *MockDependencyManagerClient) RegisterModuleDependencies(ctx context.Context, in *terrariumModule.RegisterModuleDependenciesRequest, opts ...grpc.CallOption) (*terrariumModule.Response, error) {
	m.RegisterModuleDependenciesInvocations++
	return m.RegisterModuleDependenciesResponse, m.RegisterModuleDependenciesError
}

func (m *MockDependencyManagerClient) RegisterContainerDependencies(ctx context.Context, in *terrariumModule.RegisterContainerDependenciesRequest, opts ...grpc.CallOption) (*terrariumModule.Response, error) {
	m.RegisterContainerDependenciesInvocations++
	return m.RegisterContainerDependenciesResponse, m.RegisterContainerDependenciesError
}

func (m *MockDependencyManagerClient) RetrieveContainerDependencies(ctx context.Context, in *terrariumModule.RetrieveContainerDependenciesRequestV2, opts ...grpc.CallOption) (moduleServices.DependencyManager_RetrieveContainerDependenciesClient, error) {
	m.RetrieveContainerDependenciesInvocations++
	return m.RetrieveContainerDependenciesClient, m.RetrieveContainerDependenciesError
}
func (m *MockDependencyManagerClient) RetrieveModuleDependencies(ctx context.Context, in *terrariumModule.RetrieveModuleDependenciesRequest, opts ...grpc.CallOption) (moduleServices.DependencyManager_RetrieveModuleDependenciesClient, error) {
	m.RetrieveModuleDependenciesInvocations++
	return m.RetrieveModuleDependenciesClient, m.RetrieveModuleDependenciesError
}

type MockDependencyManager_RetrieveContainerDependenciesClient struct {
	moduleServices.DependencyManager_RetrieveContainerDependenciesClient
	RecvInvocations      int
	RecvResponse         *terrariumModule.ContainerDependenciesResponseV2
	RecvError            error
	CloseSendInvocations int
	CloseSendError       error
}

func (m *MockDependencyManager_RetrieveContainerDependenciesClient) Recv() (*terrariumModule.ContainerDependenciesResponseV2, error) {
	m.RecvInvocations++
	return m.RecvResponse, m.RecvError
}

func (m *MockDependencyManager_RetrieveContainerDependenciesClient) CloseSend() error {
	m.CloseSendInvocations++
	return m.CloseSendError
}

type MockDependencyManager_RetrieveModuleDependenciesClient struct {
	moduleServices.DependencyManager_RetrieveModuleDependenciesClient
	RecvInvocations      int
	RecvResponse         *terrariumModule.ModuleDependenciesResponse
	RecvError            error
	CloseSendInvocations int
	CloseSendError       error
}

func (m *MockDependencyManager_RetrieveModuleDependenciesClient) Recv() (*terrariumModule.ModuleDependenciesResponse, error) {
	m.RecvInvocations++
	return m.RecvResponse, m.RecvError
}

func (m *MockDependencyManager_RetrieveModuleDependenciesClient) CloseSend() error {
	m.CloseSendInvocations++
	return m.CloseSendError
}
