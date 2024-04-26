package mocks

import (
	"context"

	providerServices "github.com/terrariumcloud/terrarium/internal/provider/services"
	terrariumProvider "github.com/terrariumcloud/terrarium/pkg/terrarium/provider"
)

type MockDownloadProviderSourceZipServer struct {
	providerServices.Storage_DownloadProviderSourceZipServer
	SendInvocations int
	SendResponse    *terrariumProvider.SourceZipResponse
	SendError       error
}

func (mds *MockDownloadProviderSourceZipServer) Context() context.Context {
	return context.TODO()
}

func (mds *MockDownloadProviderSourceZipServer) Send(res *terrariumProvider.SourceZipResponse) error {
	mds.SendInvocations++
	mds.SendResponse = res
	return mds.SendError
}

type MockDownloadProviderShasumServer struct {
	providerServices.Storage_DownloadShasumServer
	SendInvocations int
	SendResponse    *terrariumProvider.DownloadShasumResponse
	SendError       error
}

func (mds *MockDownloadProviderShasumServer) Context() context.Context {
	return context.TODO()
}

func (mds *MockDownloadProviderShasumServer) Send(res *terrariumProvider.DownloadShasumResponse) error {
	mds.SendInvocations++
	mds.SendResponse = res
	return mds.SendError
}

type MockDownloadProviderShasumSignatureServer struct {
	providerServices.Storage_DownloadShasumSignatureServer
	SendInvocations int
	SendResponse    *terrariumProvider.DownloadShasumResponse
	SendError       error
}

func (mds *MockDownloadProviderShasumSignatureServer) Context() context.Context {
	return context.TODO()
}

func (mds *MockDownloadProviderShasumSignatureServer) Send(res *terrariumProvider.DownloadShasumResponse) error {
	mds.SendInvocations++
	mds.SendResponse = res
	return mds.SendError
}
