package mocks

import (
	"context"
	"io"

	"github.com/terrariumcloud/terrarium/internal/module/services"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/module"
)

type MockUploadSourceZipServer struct {
	services.Storage_UploadSourceZipServer
	SendAndCloseInvocations int
	SendAndCloseResponse    *module.Response
	SendAndCloseError       error
	RecvInvocations         int
	RecvMaxInvocations      int
	RecvRequest             *module.UploadSourceZipRequest
	RecvError               error
}

func (mus *MockUploadSourceZipServer) Context() context.Context {
	return context.TODO()
}

func (mus *MockUploadSourceZipServer) SendAndClose(response *module.Response) error {
	mus.SendAndCloseInvocations++
	mus.SendAndCloseResponse = response
	return mus.SendAndCloseError
}

func (mus *MockUploadSourceZipServer) Recv() (*module.UploadSourceZipRequest, error) {
	mus.RecvInvocations++

	if mus.RecvError != nil {
		return nil, mus.RecvError
	}

	if mus.RecvInvocations == mus.RecvMaxInvocations {
		return nil, io.EOF
	}

	return mus.RecvRequest, mus.RecvError
}

type MockDownloadSourceZipServer struct {
	services.Storage_DownloadSourceZipServer
	SendInvocations int
	SendResponse    *module.SourceZipResponse
	SendError       error
}

func (mds *MockDownloadSourceZipServer) Context() context.Context {
	return context.TODO()
}

func (mds *MockDownloadSourceZipServer) Send(res *module.SourceZipResponse) error {
	mds.SendInvocations++
	mds.SendResponse = res
	return mds.SendError
}

type MockConsumer_RetrieveContainerDependenciesServer struct {
	module.Consumer_RetrieveContainerDependenciesServer
	SendInvocations int
	SendResponse    *module.ContainerDependenciesResponse
	SendError       error
}

func (m *MockConsumer_RetrieveContainerDependenciesServer) Context() context.Context {
	return context.TODO()
}

func (m *MockConsumer_RetrieveContainerDependenciesServer) Send(res *module.ContainerDependenciesResponse) error {
	m.SendInvocations++
	m.SendResponse = res
	return m.SendError
}

type MockConsumer_RetrieveModuleDependenciesServer struct {
	module.Consumer_RetrieveModuleDependenciesServer
	SendInvocations int
	SendResponse    *module.ModuleDependenciesResponse
	SendError       error
}

func (m *MockConsumer_RetrieveModuleDependenciesServer) Context() context.Context {
	return context.TODO()
}

func (m *MockConsumer_RetrieveModuleDependenciesServer) Send(res *module.ModuleDependenciesResponse) error {
	m.SendInvocations++
	m.SendResponse = res
	return m.SendError
}
