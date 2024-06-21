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
	UploadProviderBinaryZipInvocations int
	UploadProviderBinaryZipClient      providerServices.Storage_UploadProviderBinaryZipClient
	UploadProviderBinaryZipError       error
	UploadShasumInvocations            int
	UploadShasumClient                 providerServices.Storage_UploadShasumClient
	UploadShasumError                  error
	UploadShasumSignatureInvocations   int
	UploadShasumSignatureClient        providerServices.Storage_UploadShasumSignatureClient
	UploadShasumSignatureError         error
}

func (m *MockProviderStorageClient) DownloadProviderSourceZip(ctx context.Context, in *providerServices.DownloadSourceZipRequest, opts ...grpc.CallOption) (providerServices.Storage_DownloadProviderSourceZipClient, error) {
	m.DownloadSourceZipInvocations++
	return m.DownloadSourceZipClient, m.DownloadSourceZipError
}

func (m *MockProviderStorageClient) DownloadShasum(ctx context.Context, in *providerServices.DownloadShasumRequest, opts ...grpc.CallOption) (providerServices.Storage_DownloadShasumClient, error) {
	m.DownloadShasumInvocations++
	return m.DownloadShasumClient, m.DownloadShasumError
}

func (m *MockProviderStorageClient) DownloadShasumSignature(ctx context.Context, in *providerServices.DownloadShasumRequest, opts ...grpc.CallOption) (providerServices.Storage_DownloadShasumSignatureClient, error) {
	m.DownloadShasumSignatureInvocations++
	return m.DownloadShasumSignatureClient, m.DownloadShasumSignatureError
}

func (m *MockProviderStorageClient) UploadProviderBinaryZip(ctx context.Context, opts ...grpc.CallOption) (providerServices.Storage_UploadProviderBinaryZipClient, error) {
	m.UploadProviderBinaryZipInvocations++
	return m.UploadProviderBinaryZipClient, m.UploadProviderBinaryZipError
}

func (m *MockProviderStorageClient) UploadShasum(ctx context.Context, opts ...grpc.CallOption) (providerServices.Storage_UploadShasumClient, error) {
	m.UploadShasumInvocations++
	return m.UploadShasumClient, m.UploadShasumError
}

func (m *MockProviderStorageClient) UploadShasumSignature(ctx context.Context, opts ...grpc.CallOption) (providerServices.Storage_UploadShasumSignatureClient, error) {
	m.UploadShasumSignatureInvocations++
	return m.UploadShasumSignatureClient, m.UploadShasumSignatureError
}

type MockStorage_DownloadProviderSourceZipClient struct {
	providerServices.Storage_DownloadProviderSourceZipClient
	RecvInvocations      int
	RecvResponse         *providerServices.SourceZipResponse
	RecvError            error
	CloseSendInvocations int
	CloseSendError       error
}

func (m *MockStorage_DownloadProviderSourceZipClient) Recv() (*providerServices.SourceZipResponse, error) {
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
	RecvResponse         *providerServices.DownloadShasumResponse
	RecvError            error
	CloseSendInvocations int
	CloseSendError       error
}

func (m *MockStorage_DownloadProviderShasumClient) Recv() (*providerServices.DownloadShasumResponse, error) {
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
	RecvResponse         *providerServices.DownloadShasumResponse
	RecvError            error
	CloseSendInvocations int
	CloseSendError       error
}

func (m *MockStorage_DownloadProviderShasumSignatureClient) Recv() (*providerServices.DownloadShasumResponse, error) {
	m.RecvInvocations++
	return m.RecvResponse, m.RecvError
}

func (m *MockStorage_DownloadProviderShasumSignatureClient) CloseSend() error {
	m.CloseSendInvocations++
	return m.CloseSendError
}

type MockStorage_UploadProviderBinaryZipClient struct {
	providerServices.Storage_UploadProviderBinaryZipClient
	CloseAndRecvInvocations int
	CloseAndRecvResponse    *terrariumProvider.Response
	CloseAndRecvError       error
	SendInvocations         int
	SendError               error
	CloseSendInvocations    int
	CloseSendError          error
}

func (m *MockStorage_UploadProviderBinaryZipClient) CloseAndRecv() (*terrariumProvider.Response, error) {
	m.CloseAndRecvInvocations++
	return m.CloseAndRecvResponse, m.CloseAndRecvError
}

func (m *MockStorage_UploadProviderBinaryZipClient) Send(*terrariumProvider.UploadProviderBinaryZipRequest) error {
	m.SendInvocations++
	return m.SendError
}

func (m *MockStorage_UploadProviderBinaryZipClient) CloseSend() error {
	m.CloseSendInvocations++
	return m.CloseSendError
}

type MockStorage_UploadShasumClient struct {
	providerServices.Storage_UploadShasumClient
	CloseAndRecvInvocations int
	CloseAndRecvResponse    *terrariumProvider.Response
	CloseAndRecvError       error
	SendInvocations         int
	SendError               error
	CloseSendInvocations    int
	CloseSendError          error
}

func (m *MockStorage_UploadShasumClient) CloseAndRecv() (*terrariumProvider.Response, error) {
	m.CloseAndRecvInvocations++
	return m.CloseAndRecvResponse, m.CloseAndRecvError
}

func (m *MockStorage_UploadShasumClient) Send(*terrariumProvider.UploadShasumRequest) error {
	m.SendInvocations++
	return m.SendError
}

func (m *MockStorage_UploadShasumClient) CloseSend() error {
	m.CloseSendInvocations++
	return m.CloseSendError
}

type MockStorage_UploadShasumSignatureClient struct {
	providerServices.Storage_UploadShasumSignatureClient
	CloseAndRecvInvocations int
	CloseAndRecvResponse    *terrariumProvider.Response
	CloseAndRecvError       error
	SendInvocations         int
	SendError               error
	CloseSendInvocations    int
	CloseSendError          error
}

func (m *MockStorage_UploadShasumSignatureClient) CloseAndRecv() (*terrariumProvider.Response, error) {
	m.CloseAndRecvInvocations++
	return m.CloseAndRecvResponse, m.CloseAndRecvError
}

func (m *MockStorage_UploadShasumSignatureClient) Send(*terrariumProvider.UploadShasumRequest) error {
	m.SendInvocations++
	return m.SendError
}

func (m *MockStorage_UploadShasumSignatureClient) CloseSend() error {
	m.CloseSendInvocations++
	return m.CloseSendError
}
