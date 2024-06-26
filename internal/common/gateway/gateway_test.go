package gateway

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/terrariumcloud/terrarium/internal/module/services/mocks"
	"github.com/terrariumcloud/terrarium/internal/module/services/storage"
	providerMocks "github.com/terrariumcloud/terrarium/internal/provider/services/mocks"
	providerStorage "github.com/terrariumcloud/terrarium/internal/provider/services/storage"
	releaseMocks "github.com/terrariumcloud/terrarium/internal/release/services/mocks"

	"github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	terrariumProvider "github.com/terrariumcloud/terrarium/pkg/terrarium/provider"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/release"
)

// Test_RegisterWithClient checks:
// - if correct response is returned when client returns response
// - if error is returned when client returns error
func Test_RegisterWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when client returns response", func(t *testing.T) {
		response := &module.Response{}
		client := &mocks.MockRegistrarClient{RegisterResponse: response}
		gw := &TerrariumGrpcGateway{}

		res, err := gw.RegisterWithClient(context.TODO(), &module.RegisterModuleRequest{}, client)

		if res != response {
			t.Errorf("Expected %v, got %v.", response, res)
		}

		if client.RegisterInvocations != 1 {
			t.Errorf("Expected 1 call to Register, got %v", client.RegisterInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when client returns error", func(t *testing.T) {
		expected := errors.New("Test")
		client := &mocks.MockRegistrarClient{RegisterError: expected}
		gw := &TerrariumGrpcGateway{}

		_, actual := gw.RegisterWithClient(context.TODO(), &module.RegisterModuleRequest{}, client)

		if actual != expected {
			t.Errorf("Expected %v, got %v.", expected, actual)
		}

		if client.RegisterInvocations != 1 {
			t.Errorf("Expected 1 call to Register, got %v", client.RegisterInvocations)
		}
	})
}

// Test_BeginVersionWithClient checks:
// - if correct response is returned when client returns response
// - if error is returned when client returns error
func Test_BeginVersionWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when client returns response", func(t *testing.T) {
		response := &module.Response{}
		client := &mocks.MockVersionManagerClient{BeginVersionResponse: response}
		gw := &TerrariumGrpcGateway{}

		res, err := gw.BeginVersionWithClient(context.TODO(), &module.BeginVersionRequest{}, client)

		if res != response {
			t.Errorf("Expected %v, got %v.", response, res)
		}

		if client.BeginVersionInvocations != 1 {
			t.Errorf("Expected 1 call to Register, got %v", client.BeginVersionInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when client returns error", func(t *testing.T) {
		expected := errors.New("Test")
		client := &mocks.MockVersionManagerClient{BeginVersionError: expected}
		gw := &TerrariumGrpcGateway{}

		_, actual := gw.BeginVersionWithClient(context.TODO(), &module.BeginVersionRequest{}, client)

		if actual != expected {
			t.Errorf("Expected %v, got %v.", expected, actual)
		}

		if client.BeginVersionInvocations != 1 {
			t.Errorf("Expected 1 call to Register, got %v", client.BeginVersionInvocations)
		}
	})
}

// Test_EndVersionWithClient checks:
// - if correct response is returned when client returns publish response
// - if error is returned when client returns publish error
// - if correct response is returned when client returns abort response
// - if error is returned when client returns abort error
// - if error is returned when unknown action is requested
func Test_EndVersionWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when client returns publish response", func(t *testing.T) {
		response := &module.Response{}
		client := &mocks.MockVersionManagerClient{PublishVersionResponse: response}
		gw := &TerrariumGrpcGateway{}
		req := &module.EndVersionRequest{
			Action: module.EndVersionRequest_PUBLISH,
			Module: &module.Module{},
		}

		res, err := gw.EndVersionWithClient(context.TODO(), req, client)

		if res != response {
			t.Errorf("Expected %v, got %v.", response, res)
		}

		if client.PublishVersionInvocations != 1 {
			t.Errorf("Expected 1 call to PublishVersion, got %v", client.PublishVersionInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when client returns publish error", func(t *testing.T) {
		expected := errors.New("Test")
		client := &mocks.MockVersionManagerClient{PublishVersionError: expected}
		gw := &TerrariumGrpcGateway{}
		req := &module.EndVersionRequest{
			Action: module.EndVersionRequest_PUBLISH,
			Module: &module.Module{},
		}

		_, actual := gw.EndVersionWithClient(context.TODO(), req, client)

		if actual != expected {
			t.Errorf("Expected %v, got %v.", expected, actual)
		}

		if client.PublishVersionInvocations != 1 {
			t.Errorf("Expected 1 call to PublishVersion, got %v", client.PublishVersionInvocations)
		}
	})

	t.Run("when client returns abort response", func(t *testing.T) {
		response := &module.Response{}
		client := &mocks.MockVersionManagerClient{AbortVersionResponse: response}
		gw := &TerrariumGrpcGateway{}
		req := &module.EndVersionRequest{
			Action: module.EndVersionRequest_DISCARD,
			Module: &module.Module{},
		}

		res, err := gw.EndVersionWithClient(context.TODO(), req, client)

		if res != response {
			t.Errorf("Expected %v, got %v.", response, res)
		}

		if client.AbortVersionInvocations != 1 {
			t.Errorf("Expected 1 call to AbortVersion, got %v", client.AbortVersionInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when client returns abort error", func(t *testing.T) {
		expected := errors.New("Test")

		client := &mocks.MockVersionManagerClient{AbortVersionError: expected}

		gw := &TerrariumGrpcGateway{}

		req := &module.EndVersionRequest{
			Action: module.EndVersionRequest_DISCARD,
			Module: &module.Module{},
		}

		_, actual := gw.EndVersionWithClient(context.TODO(), req, client)

		if actual != expected {
			t.Errorf("Expected %v, got %v.", expected, actual)
		}

		if client.AbortVersionInvocations != 1 {
			t.Errorf("Expected 1 call to AbortVersion, got %v", client.AbortVersionInvocations)
		}
	})

	t.Run("when unknown action is requested", func(t *testing.T) {
		client := &mocks.MockVersionManagerClient{}

		gw := &TerrariumGrpcGateway{}

		req := &module.EndVersionRequest{
			Action: 123,
			Module: &module.Module{},
		}

		_, actual := gw.EndVersionWithClient(context.TODO(), req, client)

		if actual != UnknownVersionManagerActionError {
			t.Errorf("Expected %v, got %v.", UnknownVersionManagerActionError, actual)
		}

		if client.PublishVersionInvocations != 0 {
			t.Errorf("Expected 0 calls to PublishVersion, got %v", client.PublishVersionInvocations)
		}

		if client.AbortVersionInvocations != 0 {
			t.Errorf("Expected 0 calls to AbortVersion, got %v", client.AbortVersionInvocations)
		}
	})
}

// Test_UploadSourceZipWithClient checks:
// - if error is returned when client UploadSourceZip fails
// - if error is returned when Recv returns EOF and client fails
// - if no error is returned when Recv returns EOF and both client and server close stream
// - if error is returned when Recv fails
// - if no error is returned when Send returns EOF
// - if error is returned when Send fails
func Test_UploadSourceZipWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when client UploadSourceZip fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &mocks.MockUploadSourceZipServer{}

		client := &mocks.MockStorageClient{UploadSourceZipError: errors.New("some error")}

		err := gw.UploadSourceZipWithClient(server, client)

		if client.UploadSourceZipInvocations != 1 {
			t.Errorf("Expected 1 call to UploadSourceZip, got %v", client.UploadSourceZipInvocations)
		}

		if err == nil {
			t.Errorf("Expected error, got nil.")
		}
	})

	t.Run("when Recv returns EOF and client fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &mocks.MockUploadSourceZipServer{RecvError: io.EOF}

		c := &mocks.MockStorage_UploadSourceZipClient{CloseAndRecvError: errors.New("some error")}

		client := &mocks.MockStorageClient{UploadSourceZipClient: c}

		err := gw.UploadSourceZipWithClient(server, client)

		if c.CloseAndRecvInvocations != 1 {
			t.Errorf("Expected 1 call to CloseAndRecv, got %v", c.CloseAndRecvInvocations)
		}

		if server.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", server.RecvInvocations)
		}

		if client.UploadSourceZipInvocations != 1 {
			t.Errorf("Expected 1 call to UploadSourceZip, got %v", client.UploadSourceZipInvocations)
		}

		if err == nil {
			t.Errorf("Expected error, got nil.")
		}
	})

	t.Run("when Recv returns EOF and both client and server close stream", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &mocks.MockUploadSourceZipServer{RecvError: io.EOF}

		c := &mocks.MockStorage_UploadSourceZipClient{}

		client := &mocks.MockStorageClient{UploadSourceZipClient: c}

		err := gw.UploadSourceZipWithClient(server, client)

		if server.SendAndCloseInvocations != 1 {
			t.Errorf("Expected 1 call to SendAndClose, got %v", server.SendAndCloseInvocations)
		}

		if c.CloseAndRecvInvocations != 1 {
			t.Errorf("Expected 1 call to CloseAndRecv, got %v", c.CloseAndRecvInvocations)
		}

		if server.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", server.RecvInvocations)
		}

		if client.UploadSourceZipInvocations != 1 {
			t.Errorf("Expected 1 call to UploadSourceZip, got %v", client.UploadSourceZipInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when Recv fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &mocks.MockUploadSourceZipServer{RecvError: errors.New("some error")}

		client := &mocks.MockStorageClient{}

		err := gw.UploadSourceZipWithClient(server, client)

		if server.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", server.RecvInvocations)
		}

		if client.UploadSourceZipInvocations != 1 {
			t.Errorf("Expected 1 call to UploadSourceZip, got %v", client.UploadSourceZipInvocations)
		}

		if err != storage.RecieveSourceZipError {
			t.Errorf("Expected %v, got %v.", storage.RecieveSourceZipError, err)
		}
	})

	t.Run("when Send returns EOF", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &mocks.MockUploadSourceZipServer{}

		c := &mocks.MockStorage_UploadSourceZipClient{SendError: io.EOF}

		client := &mocks.MockStorageClient{UploadSourceZipClient: c}

		err := gw.UploadSourceZipWithClient(server, client)

		if server.SendAndCloseInvocations != 1 {
			t.Errorf("Expected 1 call to SendAndClose, got %v", server.SendAndCloseInvocations)
		}

		if c.CloseSendInvocations != 1 {
			t.Errorf("Expected 1 call to CloseSend, got %v", c.CloseSendInvocations)
		}

		if c.SendInvocations != 1 {
			t.Errorf("Expected 1 call to Send, got %v", c.SendInvocations)
		}

		if server.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", server.RecvInvocations)
		}

		if client.UploadSourceZipInvocations != 1 {
			t.Errorf("Expected 1 call to UploadSourceZip, got %v", client.UploadSourceZipInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when Send fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &mocks.MockUploadSourceZipServer{}

		c := &mocks.MockStorage_UploadSourceZipClient{SendError: errors.New("some error")}

		client := &mocks.MockStorageClient{UploadSourceZipClient: c}

		err := gw.UploadSourceZipWithClient(server, client)

		if c.SendInvocations != 1 {
			t.Errorf("Expected 1 call to Send, got %v", c.SendInvocations)
		}

		if server.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", server.RecvInvocations)
		}

		if client.UploadSourceZipInvocations != 1 {
			t.Errorf("Expected 1 call to UploadSourceZip, got %v", client.UploadSourceZipInvocations)
		}

		if err == nil {
			t.Errorf("Expected error, got nil.")
		}
	})
}

// Test_DownloadSourceZipWithClient checks:
// - if error is returned when client DownloadSourceZip fails
// - if no error is returned when Recv returns EOF
// - if error is returned when Recv fails
// - if error is returned when Send fails
func Test_DownloadSourceZipWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when client DownloadSourceZip fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		request := &module.DownloadSourceZipRequest{}

		server := &mocks.MockDownloadSourceZipServer{}

		client := &mocks.MockStorageClient{DownloadSourceZipError: errors.New("some error")}

		err := gw.DownloadSourceZipWithClient(request, server, client)

		if client.DownloadSourceZipInvocations != 1 {
			t.Errorf("Expected 1 call to UploadSourceZip, got %v", client.DownloadSourceZipClient)
		}

		if err == nil {
			t.Errorf("Expected error, got nil.")
		}
	})

	t.Run("when client Recv returns EOF", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		request := &module.DownloadSourceZipRequest{}

		server := &mocks.MockDownloadSourceZipServer{}

		c := &mocks.MockStorage_DownloadSourceZipClient{RecvError: io.EOF}

		client := &mocks.MockStorageClient{DownloadSourceZipClient: c}

		err := gw.DownloadSourceZipWithClient(request, server, client)

		if c.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", c.RecvInvocations)
		}

		if client.DownloadSourceZipInvocations != 1 {
			t.Errorf("Expected 1 call to DownloadSourceZip, got %v", client.DownloadSourceZipInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when Recv fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		request := &module.DownloadSourceZipRequest{}

		server := &mocks.MockDownloadSourceZipServer{}

		c := &mocks.MockStorage_DownloadSourceZipClient{RecvError: errors.New("some error")}

		client := &mocks.MockStorageClient{DownloadSourceZipClient: c}

		err := gw.DownloadSourceZipWithClient(request, server, client)

		if c.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", c.RecvInvocations)
		}

		if client.DownloadSourceZipInvocations != 1 {
			t.Errorf("Expected 1 call to DownloadSourceZip, got %v", client.DownloadSourceZipInvocations)
		}

		if err == nil {
			t.Errorf("Expected error, got nil.")
		}
	})

	t.Run("when Send fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		request := &module.DownloadSourceZipRequest{}

		server := &mocks.MockDownloadSourceZipServer{SendError: errors.New("some error")}

		c := &mocks.MockStorage_DownloadSourceZipClient{}

		client := &mocks.MockStorageClient{DownloadSourceZipClient: c}

		err := gw.DownloadSourceZipWithClient(request, server, client)

		if client.DownloadSourceZipInvocations != 1 {
			t.Errorf("Expected 1 call to DownloadSourceZip, got %v", client.DownloadSourceZipInvocations)
		}

		if c.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", c.RecvInvocations)
		}

		if server.SendInvocations != 1 {
			t.Errorf("Expected 1 call to Send, got %v", server.SendInvocations)
		}

		if err != storage.SendSourceZipError {
			t.Errorf("Expected %v, got %v.", storage.SendSourceZipError, err)
		}
	})
}

// Test_RegisterModuleDependenciesWithClient checks:
// - if correct response is returned when client returns response
// - if error is returned when client returns error
func Test_RegisterModuleDependenciesWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when client returns response", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		request := &module.RegisterModuleDependenciesRequest{}

		response := &module.Response{}

		client := &mocks.MockDependencyManagerClient{RegisterModuleDependenciesResponse: response}

		actual, err := gw.RegisterModuleDependenciesWithClient(context.TODO(), request, client)

		if actual == nil {
			t.Errorf("Expected %v, got %v", response, actual)
		}

		if client.RegisterModuleDependenciesInvocations != 1 {
			t.Errorf("Expected 1 call to RegisterModuleDependencies, got %v", client.RegisterModuleDependenciesInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when client returns error", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		request := &module.RegisterModuleDependenciesRequest{}

		client := &mocks.MockDependencyManagerClient{RegisterModuleDependenciesError: errors.New("some error")}

		_, err := gw.RegisterModuleDependenciesWithClient(context.TODO(), request, client)

		if client.RegisterModuleDependenciesInvocations != 1 {
			t.Errorf("Expected 1 call to RegisterModuleDependencies, got %v", client.RegisterModuleDependenciesInvocations)
		}

		if err == nil {
			t.Errorf("Expected error, got nil.")
		}
	})
}

// Test_RegisterContainerDependenciesWithClient checks:
// - if correct response is returned when client returns response
// - if error is returned when client returns error
func Test_RegisterContainerDependenciesWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when client returns response", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		request := &module.RegisterContainerDependenciesRequest{}

		response := &module.Response{}

		client := &mocks.MockDependencyManagerClient{RegisterContainerDependenciesResponse: response}

		actual, err := gw.RegisterContainerDependenciesWithClient(context.TODO(), request, client)

		if actual == nil {
			t.Errorf("Expected %v, got %v", response, actual)
		}

		if client.RegisterContainerDependenciesInvocations != 1 {
			t.Errorf("Expected 1 call to RegisterContainerDependencies, got %v", client.RegisterContainerDependenciesInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when client returns error", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		request := &module.RegisterContainerDependenciesRequest{}

		client := &mocks.MockDependencyManagerClient{RegisterContainerDependenciesError: errors.New("some error")}

		_, err := gw.RegisterContainerDependenciesWithClient(context.TODO(), request, client)

		if client.RegisterContainerDependenciesInvocations != 1 {
			t.Errorf("Expected 1 call to RegisterContainerDependencies, got %v", client.RegisterContainerDependenciesInvocations)
		}

		if err == nil {
			t.Errorf("Expected error, got nil.")
		}
	})
}

// Test_RetrieveContainerDependenciesWithClient checks:
// - if error is returned when client RetrieveContainerDependencies fails
// - if no error is returned when Recv returns EOF
// - if error is returned when Recv fails
// - if error is returned when Send fails
func Test_RetrieveContainerDependenciesV2WithClient(t *testing.T) {
	t.Parallel()

	t.Run("when client RetrieveContainerDependencies fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		request := &module.RetrieveContainerDependenciesRequestV2{}

		server := &mocks.MockConsumer_RetrieveContainerDependenciesServer{}

		client := &mocks.MockDependencyManagerClient{RetrieveContainerDependenciesError: errors.New("some error")}

		err := gw.RetrieveContainerDependenciesV2WithClient(request, server, client)

		if client.RetrieveContainerDependenciesInvocations != 1 {
			t.Errorf("Expected 1 call to RetrieveContainerDependencies, got %v", client.RetrieveContainerDependenciesInvocations)
		}

		if err == nil {
			t.Errorf("Expected error, got nil.")
		}
	})

	t.Run("when client Recv returns EOF", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		request := &module.RetrieveContainerDependenciesRequestV2{}

		server := &mocks.MockConsumer_RetrieveContainerDependenciesServer{}

		c := &mocks.MockDependencyManager_RetrieveContainerDependenciesClient{RecvError: io.EOF}

		client := &mocks.MockDependencyManagerClient{RetrieveContainerDependenciesClient: c}

		err := gw.RetrieveContainerDependenciesV2WithClient(request, server, client)

		if c.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", c.RecvInvocations)
		}

		if client.RetrieveContainerDependenciesInvocations != 1 {
			t.Errorf("Expected 1 call to RetrieveContainerDependencies, got %v", client.RetrieveContainerDependenciesInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when Recv fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		request := &module.RetrieveContainerDependenciesRequestV2{}

		server := &mocks.MockConsumer_RetrieveContainerDependenciesServer{}

		c := &mocks.MockDependencyManager_RetrieveContainerDependenciesClient{RecvError: errors.New("some error")}

		client := &mocks.MockDependencyManagerClient{RetrieveContainerDependenciesClient: c}

		err := gw.RetrieveContainerDependenciesV2WithClient(request, server, client)

		if c.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", c.RecvInvocations)
		}

		if client.RetrieveContainerDependenciesInvocations != 1 {
			t.Errorf("Expected 1 call to RetrieveContainerDependencies, got %v", client.RetrieveContainerDependenciesInvocations)
		}

		if err == nil {
			t.Errorf("Expected error, got nil.")
		}
	})

	t.Run("when Send fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		request := &module.RetrieveContainerDependenciesRequestV2{}

		server := &mocks.MockConsumer_RetrieveContainerDependenciesServer{SendError: errors.New("some error")}

		c := &mocks.MockDependencyManager_RetrieveContainerDependenciesClient{}

		client := &mocks.MockDependencyManagerClient{RetrieveContainerDependenciesClient: c}

		err := gw.RetrieveContainerDependenciesV2WithClient(request, server, client)

		if client.RetrieveContainerDependenciesInvocations != 1 {
			t.Errorf("Expected 1 call to RetrieveContainerDependencies, got %v", client.RetrieveContainerDependenciesInvocations)
		}

		if c.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", c.RecvInvocations)
		}

		if server.SendInvocations != 1 {
			t.Errorf("Expected 1 call to Send, got %v", server.SendInvocations)
		}

		if err != ForwardModuleDependenciesError {
			t.Errorf("Expected %v, got %v.", ForwardModuleDependenciesError, err)
		}
	})
}

// Test_RetrieveModuleDependenciesWithClient checks:
// - if error is returned when client RetrieveModuleDependencies fails
// - if no error is returned when Recv returns EOF
// - if error is returned when Recv fails
// - if error is returned when Send fails
func Test_RetrieveModuleDependenciesWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when client RetrieveModuleDependencies fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		request := &module.RetrieveModuleDependenciesRequest{}

		server := &mocks.MockConsumer_RetrieveModuleDependenciesServer{}

		client := &mocks.MockDependencyManagerClient{RetrieveModuleDependenciesError: errors.New("some error")}

		err := gw.RetrieveModuleDependenciesWithClient(request, server, client)

		if client.RetrieveModuleDependenciesInvocations != 1 {
			t.Errorf("Expected 1 call to RetrieveModuleDependencies, got %v", client.RetrieveModuleDependenciesInvocations)
		}

		if err == nil {
			t.Errorf("Expected error, got nil.")
		}
	})

	t.Run("when client Recv returns EOF", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		request := &module.RetrieveModuleDependenciesRequest{}

		server := &mocks.MockConsumer_RetrieveModuleDependenciesServer{}

		c := &mocks.MockDependencyManager_RetrieveModuleDependenciesClient{RecvError: io.EOF}

		client := &mocks.MockDependencyManagerClient{RetrieveModuleDependenciesClient: c}

		err := gw.RetrieveModuleDependenciesWithClient(request, server, client)

		if c.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", c.RecvInvocations)
		}

		if client.RetrieveModuleDependenciesInvocations != 1 {
			t.Errorf("Expected 1 call to RetrieveModuleDependencies, got %v", client.RetrieveModuleDependenciesInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when Recv fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		request := &module.RetrieveModuleDependenciesRequest{}

		server := &mocks.MockConsumer_RetrieveModuleDependenciesServer{}

		c := &mocks.MockDependencyManager_RetrieveModuleDependenciesClient{RecvError: errors.New("some error")}

		client := &mocks.MockDependencyManagerClient{RetrieveModuleDependenciesClient: c}

		err := gw.RetrieveModuleDependenciesWithClient(request, server, client)

		if c.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", c.RecvInvocations)
		}

		if client.RetrieveModuleDependenciesInvocations != 1 {
			t.Errorf("Expected 1 call to RetrieveModuleDependencies, got %v", client.RetrieveModuleDependenciesInvocations)
		}

		if err == nil {
			t.Errorf("Expected error, got nil.")
		}
	})

	t.Run("when Send fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		request := &module.RetrieveModuleDependenciesRequest{}

		server := &mocks.MockConsumer_RetrieveModuleDependenciesServer{SendError: errors.New("some error")}

		c := &mocks.MockDependencyManager_RetrieveModuleDependenciesClient{}

		client := &mocks.MockDependencyManagerClient{RetrieveModuleDependenciesClient: c}

		err := gw.RetrieveModuleDependenciesWithClient(request, server, client)

		if client.RetrieveModuleDependenciesInvocations != 1 {
			t.Errorf("Expected 1 call to RetrieveModuleDependencies, got %v", client.RetrieveModuleDependenciesInvocations)
		}

		if c.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", c.RecvInvocations)
		}

		if server.SendInvocations != 1 {
			t.Errorf("Expected 1 call to Send, got %v", server.SendInvocations)
		}

		if err != ForwardModuleDependenciesError {
			t.Errorf("Expected %v, got %v.", ForwardModuleDependenciesError, err)
		}
	})
}

// Test_PublishWithClient checks:
// - if correct response is returned when client returns response
// - if error is returned when client returns error
func Test_PublishWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when client returns response", func(t *testing.T) {
		response := &release.PublishResponse{}
		client := &releaseMocks.MockPublisherClient{PublishResponse: response}
		gw := &TerrariumGrpcGateway{}

		res, err := gw.PublishWithClient(context.TODO(), &release.PublishRequest{}, client)

		if res != response {
			t.Errorf("Expected %v, got %v.", response, res)
		}

		if client.PublishInvocations != 1 {
			t.Errorf("Expected 1 call to Release, got %v", client.PublishInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when client returns error", func(t *testing.T) {
		expected := errors.New("Test")
		client := &releaseMocks.MockPublisherClient{PublishError: expected}
		gw := &TerrariumGrpcGateway{}

		_, actual := gw.PublishWithClient(context.TODO(), &release.PublishRequest{}, client)

		if actual != expected {
			t.Errorf("Expected %v, got %v.", expected, actual)
		}

		if client.PublishInvocations != 1 {
			t.Errorf("Expected 1 call to Register, got %v", client.PublishInvocations)
		}
	})
}

// Test_RegisterProviderWithClient checks:
// - if correct response is returned when client returns response
// - if error is returned when client returns error
func Test_RegisterProviderWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when client returns response", func(t *testing.T) {
		response := &terrariumProvider.Response{}
		client := &mocks.MockProviderVersionManagerClient{RegisterResponse: response}
		gw := &TerrariumGrpcGateway{}

		res, err := gw.RegisterProviderWithClient(context.TODO(), &terrariumProvider.RegisterProviderRequest{}, client)

		if res != response {
			t.Errorf("Expected %v, got %v.", response, res)
		}

		if client.RegisterInvocations != 1 {
			t.Errorf("Expected 1 call to Register, got %v", client.RegisterInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when client returns error", func(t *testing.T) {
		expected := errors.New("Test")
		client := &mocks.MockProviderVersionManagerClient{RegisterError: expected}
		gw := &TerrariumGrpcGateway{}

		_, actual := gw.RegisterProviderWithClient(context.TODO(), &terrariumProvider.RegisterProviderRequest{}, client)

		if actual != expected {
			t.Errorf("Expected %v, got %v.", expected, actual)
		}

		if client.RegisterInvocations != 1 {
			t.Errorf("Expected 1 call to Register, got %v", client.RegisterInvocations)
		}
	})
}

// Test_EndProviderWithClient checks:
// - if correct response is returned when client returns publish response
// - if error is returned when client returns publish error
// - if correct response is returned when client returns abort response
// - if error is returned when client returns abort error
// - if error is returned when unknown action is requested
func Test_EndProviderWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when client returns publish response", func(t *testing.T) {
		response := &terrariumProvider.Response{}
		client := &mocks.MockProviderVersionManagerClient{PublishVersionResponse: response}
		gw := &TerrariumGrpcGateway{}
		req := &terrariumProvider.EndProviderRequest{
			Action:   terrariumProvider.EndProviderRequest_PUBLISH,
			Provider: &terrariumProvider.Provider{},
		}

		res, err := gw.EndProviderWithClient(context.TODO(), req, client)

		if res != response {
			t.Errorf("Expected %v, got %v.", response, res)
		}

		if client.PublishVersionInvocations != 1 {
			t.Errorf("Expected 1 call to PublishVersion, got %v", client.PublishVersionInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when client returns publish error", func(t *testing.T) {
		expected := errors.New("Test")
		client := &mocks.MockProviderVersionManagerClient{PublishVersionError: expected}
		gw := &TerrariumGrpcGateway{}
		req := &terrariumProvider.EndProviderRequest{
			Action:   terrariumProvider.EndProviderRequest_PUBLISH,
			Provider: &terrariumProvider.Provider{},
		}

		_, actual := gw.EndProviderWithClient(context.TODO(), req, client)

		if actual != expected {
			t.Errorf("Expected %v, got %v.", expected, actual)
		}

		if client.PublishVersionInvocations != 1 {
			t.Errorf("Expected 1 call to PublishVersion, got %v", client.PublishVersionInvocations)
		}
	})

	t.Run("when client returns abort response", func(t *testing.T) {
		response := &terrariumProvider.Response{}
		client := &mocks.MockProviderVersionManagerClient{AbortVersionResponse: response}
		gw := &TerrariumGrpcGateway{}
		req := &terrariumProvider.EndProviderRequest{
			Action:   terrariumProvider.EndProviderRequest_DISCARD_VERSION,
			Provider: &terrariumProvider.Provider{},
		}

		res, err := gw.EndProviderWithClient(context.TODO(), req, client)

		if res != response {
			t.Errorf("Expected %v, got %v.", response, res)
		}

		if client.AbortVersionInvocations != 1 {
			t.Errorf("Expected 1 call to AbortVersion, got %v", client.AbortVersionInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when client returns abort error", func(t *testing.T) {
		expected := errors.New("Test")

		client := &mocks.MockProviderVersionManagerClient{AbortVersionError: expected}

		gw := &TerrariumGrpcGateway{}

		req := &terrariumProvider.EndProviderRequest{
			Action:   terrariumProvider.EndProviderRequest_DISCARD_VERSION,
			Provider: &terrariumProvider.Provider{},
		}

		_, actual := gw.EndProviderWithClient(context.TODO(), req, client)

		if actual != expected {
			t.Errorf("Expected %v, got %v.", expected, actual)
		}

		if client.AbortVersionInvocations != 1 {
			t.Errorf("Expected 1 call to AbortVersion, got %v", client.AbortVersionInvocations)
		}
	})

	t.Run("when unknown action is requested", func(t *testing.T) {
		client := &mocks.MockProviderVersionManagerClient{}

		gw := &TerrariumGrpcGateway{}

		req := &terrariumProvider.EndProviderRequest{
			Action:   123,
			Provider: &terrariumProvider.Provider{},
		}

		_, actual := gw.EndProviderWithClient(context.TODO(), req, client)

		if actual != UnknownVersionManagerActionError {
			t.Errorf("Expected %v, got %v.", UnknownVersionManagerActionError, actual)
		}

		if client.PublishVersionInvocations != 0 {
			t.Errorf("Expected 0 calls to PublishVersion, got %v", client.PublishVersionInvocations)
		}

		if client.AbortVersionInvocations != 0 {
			t.Errorf("Expected 0 calls to AbortVersion, got %v", client.AbortVersionInvocations)
		}
	})
}

// Test_UploadProviderBinaryZipWithClient checks:
// - if error is returned when client UploadProviderBinaryZip fails
// - if error is returned when Recv returns EOF and client fails
// - if no error is returned when Recv returns EOF and both client and server close stream
// - if error is returned when Recv fails
// - if no error is returned when Send returns EOF
// - if error is returned when Send fails
func Test_UploadProviderBinaryZipWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when client UploadProviderBinaryZip fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &providerMocks.MockUploadProviderBinaryZipServer{}

		client := &providerMocks.MockProviderStorageClient{UploadProviderBinaryZipError: errors.New("some error")}

		err := gw.UploadProviderBinaryZipWithClient(server, client)

		if client.UploadProviderBinaryZipInvocations != 1 {
			t.Errorf("Expected 1 call to UploadProviderBinaryZip, got %v", client.UploadProviderBinaryZipInvocations)
		}

		if err == nil {
			t.Errorf("Expected error, got nil.")
		}
	})

	t.Run("when Recv returns EOF and client fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &providerMocks.MockUploadProviderBinaryZipServer{RecvError: io.EOF}

		c := &providerMocks.MockStorage_UploadProviderBinaryZipClient{CloseAndRecvError: errors.New("some error")}

		client := &providerMocks.MockProviderStorageClient{UploadProviderBinaryZipClient: c}

		err := gw.UploadProviderBinaryZipWithClient(server, client)

		if c.CloseAndRecvInvocations != 1 {
			t.Errorf("Expected 1 call to CloseAndRecv, got %v", c.CloseAndRecvInvocations)
		}

		if server.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", server.RecvInvocations)
		}

		if client.UploadProviderBinaryZipInvocations != 1 {
			t.Errorf("Expected 1 call to UploadProviderBinaryZip, got %v", client.UploadProviderBinaryZipInvocations)
		}

		if err == nil {
			t.Errorf("Expected error, got nil.")
		}
	})

	t.Run("when Recv returns EOF and both client and server close stream", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &providerMocks.MockUploadProviderBinaryZipServer{RecvError: io.EOF}

		c := &providerMocks.MockStorage_UploadProviderBinaryZipClient{}

		client := &providerMocks.MockProviderStorageClient{UploadProviderBinaryZipClient: c}

		err := gw.UploadProviderBinaryZipWithClient(server, client)

		if server.SendAndCloseInvocations != 1 {
			t.Errorf("Expected 1 call to SendAndClose, got %v", server.SendAndCloseInvocations)
		}

		if c.CloseAndRecvInvocations != 1 {
			t.Errorf("Expected 1 call to CloseAndRecv, got %v", c.CloseAndRecvInvocations)
		}

		if server.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", server.RecvInvocations)
		}

		if client.UploadProviderBinaryZipInvocations != 1 {
			t.Errorf("Expected 1 call to UploadProviderBinaryZip, got %v", client.UploadProviderBinaryZipInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when Recv fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &providerMocks.MockUploadProviderBinaryZipServer{RecvError: errors.New("some error")}

		client := &providerMocks.MockProviderStorageClient{}

		err := gw.UploadProviderBinaryZipWithClient(server, client)

		if server.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", server.RecvInvocations)
		}

		if client.UploadProviderBinaryZipInvocations != 1 {
			t.Errorf("Expected 1 call to UploadProviderBinaryZip, got %v", client.UploadProviderBinaryZipInvocations)
		}

		if err != providerStorage.ReceiveBinaryZipError {
			t.Errorf("Expected %v, got %v.", providerStorage.ReceiveBinaryZipError, err)
		}
	})

	t.Run("when Send returns EOF", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &providerMocks.MockUploadProviderBinaryZipServer{}

		c := &providerMocks.MockStorage_UploadProviderBinaryZipClient{SendError: io.EOF}

		client := &providerMocks.MockProviderStorageClient{UploadProviderBinaryZipClient: c}

		err := gw.UploadProviderBinaryZipWithClient(server, client)

		if server.SendAndCloseInvocations != 1 {
			t.Errorf("Expected 1 call to SendAndClose, got %v", server.SendAndCloseInvocations)
		}

		if c.CloseSendInvocations != 1 {
			t.Errorf("Expected 1 call to CloseSend, got %v", c.CloseSendInvocations)
		}

		if c.SendInvocations != 1 {
			t.Errorf("Expected 1 call to Send, got %v", c.SendInvocations)
		}

		if server.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", server.RecvInvocations)
		}

		if client.UploadProviderBinaryZipInvocations != 1 {
			t.Errorf("Expected 1 call to UploadProviderBinaryZip, got %v", client.UploadProviderBinaryZipInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when Send fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &providerMocks.MockUploadProviderBinaryZipServer{}

		c := &providerMocks.MockStorage_UploadProviderBinaryZipClient{SendError: errors.New("some error")}

		client := &providerMocks.MockProviderStorageClient{UploadProviderBinaryZipClient: c}

		err := gw.UploadProviderBinaryZipWithClient(server, client)

		if c.SendInvocations != 1 {
			t.Errorf("Expected 1 call to Send, got %v", c.SendInvocations)
		}

		if server.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", server.RecvInvocations)
		}

		if client.UploadProviderBinaryZipInvocations != 1 {
			t.Errorf("Expected 1 call to UploadProviderBinaryZip, got %v", client.UploadProviderBinaryZipInvocations)
		}

		if err == nil {
			t.Errorf("Expected error, got nil.")
		}
	})
}

// Test_UploadShasumWithClient checks:
// - if error is returned when client UploadShasum fails
// - if error is returned when Recv returns EOF and client fails
// - if no error is returned when Recv returns EOF and both client and server close stream
// - if error is returned when Recv fails
// - if no error is returned when Send returns EOF
// - if error is returned when Send fails
func Test_UploadShasumWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when client UploadShasum fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &providerMocks.MockUploadShasumServer{}

		client := &providerMocks.MockProviderStorageClient{UploadShasumError: errors.New("some error")}

		err := gw.UploadShasumWithClient(server, client)

		if client.UploadShasumInvocations != 1 {
			t.Errorf("Expected 1 call to UploadShasum, got %v", client.UploadShasumInvocations)
		}

		if err == nil {
			t.Errorf("Expected error, got nil.")
		}
	})

	t.Run("when Recv returns EOF and client fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &providerMocks.MockUploadShasumServer{RecvError: io.EOF}

		c := &providerMocks.MockStorage_UploadShasumClient{CloseAndRecvError: errors.New("some error")}

		client := &providerMocks.MockProviderStorageClient{UploadShasumClient: c}

		err := gw.UploadShasumWithClient(server, client)

		if c.CloseAndRecvInvocations != 1 {
			t.Errorf("Expected 1 call to CloseAndRecv, got %v", c.CloseAndRecvInvocations)
		}

		if server.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", server.RecvInvocations)
		}

		if client.UploadShasumInvocations != 1 {
			t.Errorf("Expected 1 call to UploadShasum, got %v", client.UploadShasumInvocations)
		}

		if err == nil {
			t.Errorf("Expected error, got nil.")
		}
	})

	t.Run("when Recv returns EOF and both client and server close stream", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &providerMocks.MockUploadShasumServer{RecvError: io.EOF}

		c := &providerMocks.MockStorage_UploadShasumClient{}

		client := &providerMocks.MockProviderStorageClient{UploadShasumClient: c}

		err := gw.UploadShasumWithClient(server, client)

		if server.SendAndCloseInvocations != 1 {
			t.Errorf("Expected 1 call to SendAndClose, got %v", server.SendAndCloseInvocations)
		}

		if c.CloseAndRecvInvocations != 1 {
			t.Errorf("Expected 1 call to CloseAndRecv, got %v", c.CloseAndRecvInvocations)
		}

		if server.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", server.RecvInvocations)
		}

		if client.UploadShasumInvocations != 1 {
			t.Errorf("Expected 1 call to UploadShasum, got %v", client.UploadShasumInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when Recv fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &providerMocks.MockUploadShasumServer{RecvError: errors.New("some error")}

		client := &providerMocks.MockProviderStorageClient{}

		err := gw.UploadShasumWithClient(server, client)

		if server.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", server.RecvInvocations)
		}

		if client.UploadShasumInvocations != 1 {
			t.Errorf("Expected 1 call to UploadShasum, got %v", client.UploadShasumInvocations)
		}

		if err != providerStorage.ReceiveShasumError {
			t.Errorf("Expected %v, got %v.", providerStorage.ReceiveShasumError, err)
		}
	})

	t.Run("when Send returns EOF", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &providerMocks.MockUploadShasumServer{}

		c := &providerMocks.MockStorage_UploadShasumClient{SendError: io.EOF}

		client := &providerMocks.MockProviderStorageClient{UploadShasumClient: c}

		err := gw.UploadShasumWithClient(server, client)

		if server.SendAndCloseInvocations != 1 {
			t.Errorf("Expected 1 call to SendAndClose, got %v", server.SendAndCloseInvocations)
		}

		if c.CloseSendInvocations != 1 {
			t.Errorf("Expected 1 call to CloseSend, got %v", c.CloseSendInvocations)
		}

		if c.SendInvocations != 1 {
			t.Errorf("Expected 1 call to Send, got %v", c.SendInvocations)
		}

		if server.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", server.RecvInvocations)
		}

		if client.UploadShasumInvocations != 1 {
			t.Errorf("Expected 1 call to UploadShasum, got %v", client.UploadShasumInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when Send fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &providerMocks.MockUploadShasumServer{}

		c := &providerMocks.MockStorage_UploadShasumClient{SendError: errors.New("some error")}

		client := &providerMocks.MockProviderStorageClient{UploadShasumClient: c}

		err := gw.UploadShasumWithClient(server, client)

		if c.SendInvocations != 1 {
			t.Errorf("Expected 1 call to Send, got %v", c.SendInvocations)
		}

		if server.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", server.RecvInvocations)
		}

		if client.UploadShasumInvocations != 1 {
			t.Errorf("Expected 1 call to UploadShasum, got %v", client.UploadShasumInvocations)
		}

		if err == nil {
			t.Errorf("Expected error, got nil.")
		}
	})
}

// Test_UploadShasumSignatureWithClient checks:
// - if error is returned when client UploadShasumSignature fails
// - if error is returned when Recv returns EOF and client fails
// - if no error is returned when Recv returns EOF and both client and server close stream
// - if error is returned when Recv fails
// - if no error is returned when Send returns EOF
// - if error is returned when Send fails
func Test_UploadShasumSignatureWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when client UploadShasumSignature fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &providerMocks.MockUploadShasumSignatureServer{}

		client := &providerMocks.MockProviderStorageClient{UploadShasumSignatureError: errors.New("some error")}

		err := gw.UploadShasumSignatureWithClient(server, client)

		if client.UploadShasumSignatureInvocations != 1 {
			t.Errorf("Expected 1 call to UploadShasumSignature, got %v", client.UploadShasumSignatureInvocations)
		}

		if err == nil {
			t.Errorf("Expected error, got nil.")
		}
	})

	t.Run("when Recv returns EOF and client fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &providerMocks.MockUploadShasumSignatureServer{RecvError: io.EOF}

		c := &providerMocks.MockStorage_UploadShasumSignatureClient{CloseAndRecvError: errors.New("some error")}

		client := &providerMocks.MockProviderStorageClient{UploadShasumSignatureClient: c}

		err := gw.UploadShasumSignatureWithClient(server, client)

		if c.CloseAndRecvInvocations != 1 {
			t.Errorf("Expected 1 call to CloseAndRecv, got %v", c.CloseAndRecvInvocations)
		}

		if server.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", server.RecvInvocations)
		}

		if client.UploadShasumSignatureInvocations != 1 {
			t.Errorf("Expected 1 call to UploadShasumSignature, got %v", client.UploadShasumSignatureInvocations)
		}

		if err == nil {
			t.Errorf("Expected error, got nil.")
		}
	})

	t.Run("when Recv returns EOF and both client and server close stream", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &providerMocks.MockUploadShasumSignatureServer{RecvError: io.EOF}

		c := &providerMocks.MockStorage_UploadShasumSignatureClient{}

		client := &providerMocks.MockProviderStorageClient{UploadShasumSignatureClient: c}

		err := gw.UploadShasumSignatureWithClient(server, client)

		if server.SendAndCloseInvocations != 1 {
			t.Errorf("Expected 1 call to SendAndClose, got %v", server.SendAndCloseInvocations)
		}

		if c.CloseAndRecvInvocations != 1 {
			t.Errorf("Expected 1 call to CloseAndRecv, got %v", c.CloseAndRecvInvocations)
		}

		if server.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", server.RecvInvocations)
		}

		if client.UploadShasumSignatureInvocations != 1 {
			t.Errorf("Expected 1 call to UploadShasumSignature, got %v", client.UploadShasumSignatureInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when Recv fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &providerMocks.MockUploadShasumSignatureServer{RecvError: errors.New("some error")}

		client := &providerMocks.MockProviderStorageClient{}

		err := gw.UploadShasumSignatureWithClient(server, client)

		if server.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", server.RecvInvocations)
		}

		if client.UploadShasumSignatureInvocations != 1 {
			t.Errorf("Expected 1 call to UploadShasumSignature, got %v", client.UploadShasumSignatureInvocations)
		}

		if err != providerStorage.ReceiveShasumSigError {
			t.Errorf("Expected %v, got %v.", providerStorage.ReceiveShasumSigError, err)
		}
	})

	t.Run("when Send returns EOF", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &providerMocks.MockUploadShasumSignatureServer{}

		c := &providerMocks.MockStorage_UploadShasumSignatureClient{SendError: io.EOF}

		client := &providerMocks.MockProviderStorageClient{UploadShasumSignatureClient: c}

		err := gw.UploadShasumSignatureWithClient(server, client)

		if server.SendAndCloseInvocations != 1 {
			t.Errorf("Expected 1 call to SendAndClose, got %v", server.SendAndCloseInvocations)
		}

		if c.CloseSendInvocations != 1 {
			t.Errorf("Expected 1 call to CloseSend, got %v", c.CloseSendInvocations)
		}

		if c.SendInvocations != 1 {
			t.Errorf("Expected 1 call to Send, got %v", c.SendInvocations)
		}

		if server.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", server.RecvInvocations)
		}

		if client.UploadShasumSignatureInvocations != 1 {
			t.Errorf("Expected 1 call to UploadShasumSignature, got %v", client.UploadShasumSignatureInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when Send fails", func(t *testing.T) {
		gw := &TerrariumGrpcGateway{}

		server := &providerMocks.MockUploadShasumSignatureServer{}

		c := &providerMocks.MockStorage_UploadShasumSignatureClient{SendError: errors.New("some error")}

		client := &providerMocks.MockProviderStorageClient{UploadShasumSignatureClient: c}

		err := gw.UploadShasumSignatureWithClient(server, client)

		if c.SendInvocations != 1 {
			t.Errorf("Expected 1 call to Send, got %v", c.SendInvocations)
		}

		if server.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", server.RecvInvocations)
		}

		if client.UploadShasumSignatureInvocations != 1 {
			t.Errorf("Expected 1 call to UploadShasumSignature, got %v", client.UploadShasumSignatureInvocations)
		}

		if err == nil {
			t.Errorf("Expected error, got nil.")
		}
	})
}
