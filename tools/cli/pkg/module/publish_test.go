package module

import (
	"context"
	"errors"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"google.golang.org/grpc"
	"io"
	"strings"
	"testing"
)

func TestPublish(t *testing.T) {
	type args struct {
		client   module.PublisherClient
		source   io.Reader
		metadata Metadata
	}
	type expected struct {
		registerCalls                      int
		beginVersionCalls                  int
		registerContainerDependenciesCalls int
		endVersionCalls                    int
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		expected expected
	}{
		{
			name: "Success",
			args: args{
				client: &MockPublisherClient{
					uploadSourceZipResponse: &MockPublisher_UploadSourceZipClient{sendError: errors.New("not implemented")},
				},
				source: strings.NewReader("test"),
				metadata: Metadata{
					Name:        "org/test/provider",
					Version:     "1.2.4",
					Description: "A test",
					Source:      "http://test.com/test",
					Maturity:    module.Maturity_PLANNING,
				},
			},
			expected: expected{},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Publish(tt.args.client, tt.args.source, tt.args.metadata); (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type MockPublisherClient struct {
	registerCalls                         int
	registerError                         error
	registerResponse                      *module.Response
	beginVersionCalls                     int
	beginVersionError                     error
	beginVersionResponse                  *module.Response
	registerModuleDependenciesCalls       int
	registerModuleDependenciesError       error
	registerModuleDependenciesResponse    *module.Response
	registerContainerDependenciesCalls    int
	registerContainerDependenciesError    error
	registerContainerDependenciesResponse *module.Response
	uploadSourceZipCalls                  int
	uploadSourceZipError                  error
	uploadSourceZipResponse               module.Publisher_UploadSourceZipClient
	endVersionCalls                       int
	endVersionError                       error
	endVersionResponse                    *module.Response
	publishTagCalls                       int
	publishTagError                       error
	publishTagResponse                    *module.Response
}

func (m MockPublisherClient) Register(ctx context.Context, in *module.RegisterModuleRequest, opts ...grpc.CallOption) (*module.Response, error) {
	m.registerCalls++
	return m.registerResponse, m.registerError
}

func (m MockPublisherClient) BeginVersion(ctx context.Context, in *module.BeginVersionRequest, opts ...grpc.CallOption) (*module.Response, error) {
	m.beginVersionCalls++
	return m.beginVersionResponse, m.beginVersionError
}

func (m MockPublisherClient) RegisterModuleDependencies(ctx context.Context, in *module.RegisterModuleDependenciesRequest, opts ...grpc.CallOption) (*module.Response, error) {
	m.registerModuleDependenciesCalls++
	return m.registerModuleDependenciesResponse, m.registerModuleDependenciesError
}

func (m MockPublisherClient) RegisterContainerDependencies(ctx context.Context, in *module.RegisterContainerDependenciesRequest, opts ...grpc.CallOption) (*module.Response, error) {
	m.registerContainerDependenciesCalls++
	return m.registerContainerDependenciesResponse, m.registerContainerDependenciesError
}

func (m MockPublisherClient) UploadSourceZip(ctx context.Context, opts ...grpc.CallOption) (module.Publisher_UploadSourceZipClient, error) {
	m.uploadSourceZipCalls++
	return m.uploadSourceZipResponse, m.uploadSourceZipError
}

func (m MockPublisherClient) EndVersion(ctx context.Context, in *module.EndVersionRequest, opts ...grpc.CallOption) (*module.Response, error) {
	m.endVersionCalls++
	return m.endVersionResponse, m.endVersionError
}

func (m MockPublisherClient) PublishTag(ctx context.Context, in *module.PublishTagRequest, opts ...grpc.CallOption) (*module.Response, error) {
	m.publishTagCalls++
	return m.publishTagResponse, m.publishTagError
}

type MockPublisher_UploadSourceZipClient struct {
	grpc.ClientStream
	sendError            error
	closeAndRecvError    error
	closeAndRecvResponse *module.Response
}

func (m MockPublisher_UploadSourceZipClient) Send(request *module.UploadSourceZipRequest) error {
	return m.sendError
}

func (m MockPublisher_UploadSourceZipClient) CloseAndRecv() (*module.Response, error) {
	return m.closeAndRecvResponse, m.closeAndRecvError
}
