package mocks

import (
	"io"
	"context"

	providerServices "github.com/terrariumcloud/terrarium/internal/provider/services"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/provider"
)

type MockDownloadProviderSourceZipServer struct {
	providerServices.Storage_DownloadProviderSourceZipServer
	SendInvocations int
	SendResponse    *providerServices.SourceZipResponse
	SendError       error
	TotalReceived   []byte
}

func (mds *MockDownloadProviderSourceZipServer) Context() context.Context {
	return context.TODO()
}

func (mds *MockDownloadProviderSourceZipServer) Send(res *providerServices.SourceZipResponse) error {
	mds.SendInvocations++
	mds.SendResponse = res
	mds.TotalReceived = append(mds.TotalReceived, mds.SendResponse.ZipDataChunk...)
	return mds.SendError
}

type MockDownloadProviderShasumServer struct {
	providerServices.Storage_DownloadShasumServer
	SendInvocations int
	SendResponse    *providerServices.DownloadShasumResponse
	SendError       error
	TotalReceived   []byte
}

func (mds *MockDownloadProviderShasumServer) Context() context.Context {
	return context.TODO()
}

func (mds *MockDownloadProviderShasumServer) Send(res *providerServices.DownloadShasumResponse) error {
	mds.SendInvocations++
	mds.SendResponse = res
	mds.TotalReceived = append(mds.TotalReceived, mds.SendResponse.ShasumDataChunk...)
	return mds.SendError
}

type MockDownloadProviderShasumSignatureServer struct {
	providerServices.Storage_DownloadShasumSignatureServer
	SendInvocations int
	SendResponse    *providerServices.DownloadShasumResponse
	SendError       error
	TotalReceived   []byte
}

func (mds *MockDownloadProviderShasumSignatureServer) Context() context.Context {
	return context.TODO()
}

func (mds *MockDownloadProviderShasumSignatureServer) Send(res *providerServices.DownloadShasumResponse) error {
	mds.SendInvocations++
	mds.SendResponse = res
	mds.TotalReceived = append(mds.TotalReceived, mds.SendResponse.ShasumDataChunk...)
	return mds.SendError
}

type MockUploadProviderBinaryZipServer struct {
	providerServices.Storage_UploadProviderBinaryZipServer
	SendAndCloseInvocations int
	SendAndCloseResponse    *provider.Response
	SendAndCloseError       error
	RecvInvocations         int
	RecvMaxInvocations      int
	RecvRequest             *provider.UploadProviderBinaryZipRequest
	RecvError               error
}

func (mus *MockUploadProviderBinaryZipServer) Context() context.Context {
	return context.TODO()
}

func (mus *MockUploadProviderBinaryZipServer) SendAndClose(response *provider.Response) error {
	mus.SendAndCloseInvocations++
	mus.SendAndCloseResponse = response
	return mus.SendAndCloseError
}

func (mus *MockUploadProviderBinaryZipServer) Recv() (*provider.UploadProviderBinaryZipRequest, error) {
	mus.RecvInvocations++

	if mus.RecvError != nil {
		return nil, mus.RecvError
	}
	if mus.RecvInvocations == mus.RecvMaxInvocations {
		return nil, io.EOF
	}
	return mus.RecvRequest, mus.RecvError
}

type MockUploadShasumServer struct {
	providerServices.Storage_UploadShasumServer
	SendAndCloseInvocations int
	SendAndCloseResponse    *provider.Response
	SendAndCloseError       error
	RecvInvocations         int
	RecvMaxInvocations      int
	RecvRequest             *provider.UploadShasumRequest
	RecvError               error
}

func (mus *MockUploadShasumServer) Context() context.Context {
	return context.TODO()
}

func (mus *MockUploadShasumServer) SendAndClose(response *provider.Response) error {
	mus.SendAndCloseInvocations++
	mus.SendAndCloseResponse = response
	return mus.SendAndCloseError
}

func (mus *MockUploadShasumServer) Recv() (*provider.UploadShasumRequest, error) {
	mus.RecvInvocations++

	if mus.RecvError != nil {
		return nil, mus.RecvError
	}
	if mus.RecvInvocations == mus.RecvMaxInvocations {
		return nil, io.EOF
	}
	return mus.RecvRequest, mus.RecvError
}

type MockUploadShasumSignatureServer struct {
	providerServices.Storage_UploadShasumSignatureServer
	SendAndCloseInvocations int
	SendAndCloseResponse    *provider.Response
	SendAndCloseError       error
	RecvInvocations         int
	RecvMaxInvocations      int
	RecvRequest             *provider.UploadShasumRequest
	RecvError               error
}

func (mus *MockUploadShasumSignatureServer) Context() context.Context {
	return context.TODO()
}

func (mus *MockUploadShasumSignatureServer) SendAndClose(response *provider.Response) error {
	mus.SendAndCloseInvocations++
	mus.SendAndCloseResponse = response
	return mus.SendAndCloseError
}

func (mus *MockUploadShasumSignatureServer) Recv() (*provider.UploadShasumRequest, error) {
	mus.RecvInvocations++

	if mus.RecvError != nil {
		return nil, mus.RecvError
	}
	if mus.RecvInvocations == mus.RecvMaxInvocations {
		return nil, io.EOF
	}
	return mus.RecvRequest, mus.RecvError
}
