package mocks

import (
	"io"

	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/module/services"
	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"
)

type MockUploadSourceZipServer struct {
	services.Storage_UploadSourceZipServer
	SendAndCloseInvocations int
	SendAndCloseResponse    *terrarium.Response
	SendAndCloseError       error
	RecvInvocations         int
	RecvMaxInvocations      int
	RecvRequest             *terrarium.UploadSourceZipRequest
	RecvError               error
}

func (mus *MockUploadSourceZipServer) SendAndClose(response *terrarium.Response) error {
	mus.SendAndCloseInvocations++
	mus.SendAndCloseResponse = response
	return mus.SendAndCloseError
}

func (mus *MockUploadSourceZipServer) Recv() (*terrarium.UploadSourceZipRequest, error) {
	mus.RecvInvocations++

	if mus.RecvError != nil {
		return nil, mus.RecvError
	}

	if mus.RecvMaxInvocations >= mus.RecvInvocations {
		return nil, io.EOF
	}

	return mus.RecvRequest, mus.RecvError
}

type MockDownloadSourceZipServer struct {
	services.Storage_DownloadSourceZipServer
	SendInvocations int
	SendResponse    *terrarium.SourceZipResponse
	SendError       error
}

func (mds *MockDownloadSourceZipServer) Send(res *terrarium.SourceZipResponse) error {
	mds.SendInvocations++
	mds.SendResponse = res
	return mds.SendError
}
