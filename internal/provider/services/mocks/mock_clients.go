package mocks

import (
	"context"

	providerServices "github.com/terrariumcloud/terrarium/internal/provider/services"
	terrariumProvider "github.com/terrariumcloud/terrarium/pkg/terrarium/provider"

	"google.golang.org/grpc"
)

type MockProviderStorageClient struct {
	providerServices.StorageClient
	DownloadSourceZipInvocations       int
	DownloadSourceZipClient            providerServices.Storage_DownloadProviderSourceZipClient
	DownloadSourceZipError             error
	DownloadShasumInvocations          int
	DownloadShasumClient               providerServices.Storage_DownloadShasumClient
	DownloadShasumError                error
	DownloadShasumSignatureInvocations int
	DownloadShasumSignatureClient      providerServices.Storage_DownloadShasumSignatureClient
	DownloadShasumSignatureError       error
}

func (m *MockProviderStorageClient) DownloadProviderSourceZip(ctx context.Context, in *terrariumProvider.DownloadSourceZipRequest, opts ...grpc.CallOption) (providerServices.Storage_DownloadProviderSourceZipClient, error) {
	m.DownloadSourceZipInvocations++
	return m.DownloadSourceZipClient, m.DownloadSourceZipError
}

func (m *MockProviderStorageClient) DownloadShasum(ctx context.Context, in *terrariumProvider.DownloadShasumRequest, opts ...grpc.CallOption) (providerServices.Storage_DownloadShasumClient, error) {
	m.DownloadShasumInvocations++
	return m.DownloadShasumClient, m.DownloadShasumError
}

func (m *MockProviderStorageClient) DownloadShasumSignature(ctx context.Context, in *terrariumProvider.DownloadShasumRequest, opts ...grpc.CallOption) (providerServices.Storage_DownloadShasumSignatureClient, error) {
	m.DownloadShasumSignatureInvocations++
	return m.DownloadShasumSignatureClient, m.DownloadShasumSignatureError
}

type MockStorage_DownloadProviderSourceZipClient struct {
	providerServices.Storage_DownloadProviderSourceZipClient
	RecvInvocations      int
	RecvResponse         *terrariumProvider.SourceZipResponse
	RecvError            error
	CloseSendInvocations int
	CloseSendError       error
}

func (m *MockStorage_DownloadProviderSourceZipClient) Recv() (*terrariumProvider.SourceZipResponse, error) {
	m.RecvInvocations++
	return m.RecvResponse, m.RecvError
}

func (m *MockStorage_DownloadProviderSourceZipClient) CloseSend() error {
	m.CloseSendInvocations++
	return m.CloseSendError
}

type MockStorage_DownloadProviderShasumClient struct {
	providerServices.Storage_DownloadShasumClient
	RecvInvocations      int
	RecvResponse         *terrariumProvider.DownloadShasumResponse
	RecvError            error
	CloseSendInvocations int
	CloseSendError       error
}

func (m *MockStorage_DownloadProviderShasumClient) Recv() (*terrariumProvider.DownloadShasumResponse, error) {
	m.RecvInvocations++
	return m.RecvResponse, m.RecvError
}

func (m *MockStorage_DownloadProviderShasumClient) CloseSend() error {
	m.CloseSendInvocations++
	return m.CloseSendError
}

type MockStorage_DownloadProviderShasumSignatureClient struct {
	providerServices.Storage_DownloadShasumSignatureClient
	RecvInvocations      int
	RecvResponse         *terrariumProvider.DownloadShasumResponse
	RecvError            error
	CloseSendInvocations int
	CloseSendError       error
}

func (m *MockStorage_DownloadProviderShasumSignatureClient) Recv() (*terrariumProvider.DownloadShasumResponse, error) {
	m.RecvInvocations++
	return m.RecvResponse, m.RecvError
}

func (m *MockStorage_DownloadProviderShasumSignatureClient) CloseSend() error {
	m.CloseSendInvocations++
	return m.CloseSendError
}
