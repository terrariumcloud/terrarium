package services_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/terrariumcloud/terrarium/internal/mocks"
	"github.com/terrariumcloud/terrarium/internal/module/services"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/module"
)

// Test_RegisterWithClient checks:
// - if correct response is returned when client returns response
// - if error is returned when client returns error
func Test_RegisterWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when client returns response", func(t *testing.T) {
		response := &module.Response{}
		client := &mocks.MockRegistrarClient{RegisterResponse: response}
		gw := &services.TerrariumGrpcGateway{}

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
		gw := &services.TerrariumGrpcGateway{}

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
		gw := &services.TerrariumGrpcGateway{}

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
		gw := &services.TerrariumGrpcGateway{}

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
		gw := &services.TerrariumGrpcGateway{}
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
		gw := &services.TerrariumGrpcGateway{}
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
		gw := &services.TerrariumGrpcGateway{}
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

		gw := &services.TerrariumGrpcGateway{}

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

		gw := &services.TerrariumGrpcGateway{}

		req := &module.EndVersionRequest{
			Action: 123,
			Module: &module.Module{},
		}

		_, actual := gw.EndVersionWithClient(context.TODO(), req, client)

		if actual != services.UnknownVersionManagerActionError {
			t.Errorf("Expected %v, got %v.", services.UnknownVersionManagerActionError, actual)
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
// - when Recv returns EOF and both client and server close stream
// - when Recv fails
// - when Send returns EOF
// - when Send fails
func Test_UploadSourceZipWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when client UploadSourceZip fails", func(t *testing.T) {
		gw := &services.TerrariumGrpcGateway{}

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
		gw := &services.TerrariumGrpcGateway{}

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
}

// Test_DownloadSourceZipWithClient checks:
// - when Recv completes (returns EOF)
// - when Recv fails
// - when Send fails
func Test_DownloadSourceZipWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when ...", func(t *testing.T) {
	})
}

// Test_RegisterModuleDependenciesWithClient checks:
// - if correct response is returned when client returns response
// - if error is returned when client returns error
func Test_RegisterModuleDependenciesWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when client returns response", func(t *testing.T) {
	})

	t.Run("when client returns error", func(t *testing.T) {
	})
}

// Test_RegisterContainerDependenciesWithClient checks:
// - if correct response is returned when client returns response
// - if error is returned when client returns error
func Test_RegisterContainerDependenciesWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when client returns response", func(t *testing.T) {
	})

	t.Run("when client returns error", func(t *testing.T) {
	})
}

// Test_RetrieveContainerDependenciesWithClient checks:
// - when Recv completes (returns EOF)
// - when Recv fails
// - when Send fails
func Test_RetrieveContainerDependenciesWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when ...", func(t *testing.T) {
	})
}

// Test_RetrieveModuleDependenciesWithClient checks:
// - when Recv completes (returns EOF)
// - when Recv fails
// - when Send fails
func Test_RetrieveModuleDependenciesWithClient(t *testing.T) {
	t.Parallel()

	t.Run("when ...", func(t *testing.T) {
	})
}
