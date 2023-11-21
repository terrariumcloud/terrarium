package gateway

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/terrariumcloud/terrarium/internal/module/services/mocks"
	"github.com/terrariumcloud/terrarium/internal/module/services/storage"
	"github.com/terrariumcloud/terrarium/internal/release/services"
	releaseMocks "github.com/terrariumcloud/terrarium/internal/release/services/mocks"

	"github.com/terrariumcloud/terrarium/pkg/terrarium/module"
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

// Release

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

// Test_ListReleaseTypesWithClient checks:
// - if correct response is returned when client returns response
// - if error is returned when client returns error
func Test_ListReleaseTypesWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when client returns response", func(t *testing.T) {
		response := &services.ListReleaseTypesResponse{}
		client := &releaseMocks.MockBrowseClient{ListReleaseTypesResponse: response}
		gw := &TerrariumGrpcGateway{}

		res, err := gw.ListReleaseTypesWithClient(context.TODO(), &services.ListReleaseTypesRequest{}, client)

		if res != response {
			t.Errorf("Expected %v, got %v.", response, res)
		}

		if client.ListReleaseTypesInvocations != 1 {
			t.Errorf("Expected 1 call to Release, got %v", client.ListReleaseTypesInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when client returns error", func(t *testing.T) {
		expected := errors.New("Test")
		client := &releaseMocks.MockBrowseClient{ListReleaseTypesError: expected}
		gw := &TerrariumGrpcGateway{}

		_, actual := gw.ListReleaseTypesWithClient(context.TODO(), &services.ListReleaseTypesRequest{}, client)

		if actual != expected {
			t.Errorf("Expected %v, got %v.", expected, actual)
		}

		if client.ListReleaseTypesInvocations != 1 {
			t.Errorf("Expected 1 call to Register, got %v", client.ListReleaseTypesInvocations)
		}
	})
}

// Test_ListOrganizationWithClient checks:
// - if correct response is returned when client returns response
// - if error is returned when client returns error
func Test_ListOrganizationWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when client returns response", func(t *testing.T) {
		response := &services.ListOrganizationResponse{}
		client := &releaseMocks.MockBrowseClient{ListOrganizationResponse: response}
		gw := &TerrariumGrpcGateway{}

		res, err := gw.ListOrganizationWithClient(context.TODO(), &services.ListOrganizationRequest{}, client)

		if res != response {
			t.Errorf("Expected %v, got %v.", response, res)
		}

		if client.ListOrganizationInvocations != 1 {
			t.Errorf("Expected 1 call to Release, got %v", client.ListOrganizationInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when client returns error", func(t *testing.T) {
		expected := errors.New("Test")
		client := &releaseMocks.MockBrowseClient{ListOrganizationError: expected}
		gw := &TerrariumGrpcGateway{}

		_, actual := gw.ListOrganizationWithClient(context.TODO(), &services.ListOrganizationRequest{}, client)

		if actual != expected {
			t.Errorf("Expected %v, got %v.", expected, actual)
		}

		if client.ListOrganizationInvocations != 1 {
			t.Errorf("Expected 1 call to Register, got %v", client.ListOrganizationInvocations)
		}
	})
}

// Test_ListReleasesWithClient checks:
// - if correct response is returned when client returns response
// - if error is returned when client returns error
func Test_ListReleasesWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when client returns response", func(t *testing.T) {
		response := &services.ListReleasesResponse{}
		client := &releaseMocks.MockBrowseClient{ListReleasesResponse: response}
		gw := &TerrariumGrpcGateway{}

		res, err := gw.ListReleasesWithClient(context.TODO(), &services.ListReleasesRequest{}, client)

		if res != response {
			t.Errorf("Expected %v, got %v.", response, res)
		}

		if client.ListReleasesInvocations != 1 {
			t.Errorf("Expected 1 call to Release, got %v", client.ListReleasesInvocations)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when client returns error", func(t *testing.T) {
		expected := errors.New("Test")
		client := &releaseMocks.MockBrowseClient{ListReleasesError: expected}
		gw := &TerrariumGrpcGateway{}

		_, actual := gw.ListReleasesWithClient(context.TODO(), &services.ListReleasesRequest{}, client)

		if actual != expected {
			t.Errorf("Expected %v, got %v.", expected, actual)
		}

		if client.ListReleasesInvocations != 1 {
			t.Errorf("Expected 1 call to Register, got %v", client.ListReleasesInvocations)
		}
	})
}
