package module

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"google.golang.org/grpc"
	"io"
	"strings"
	"testing"
)

func TestPublish(t *testing.T) {
	type args struct {
		client   *MockPublisherClient
		source   io.Reader
		metadata Metadata
	}
	type expected struct {
		registerCalls                      int
		beginVersionCalls                  int
		registerContainerDependenciesCalls int
		uploadSourceZipCalls               int
		uploadSourceZipSendCalls           int
		endVersionCalls                    int
		endVersionRequest                  *module.EndVersionRequest
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
					uploadSourceZipResponse: &MockPublisher_UploadSourceZipClient{
						sendError:         nil,
						closeAndRecvError: nil,
						closeAndRecvResponse: &module.Response{
							Message: "done",
						},
					},
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
			expected: expected{
				registerCalls:                      1,
				beginVersionCalls:                  1,
				registerContainerDependenciesCalls: 0,
				uploadSourceZipCalls:               1,
				uploadSourceZipSendCalls:           1,
				endVersionCalls:                    1,
				endVersionRequest: &module.EndVersionRequest{
					Module: &module.Module{
						Name:    "org/test/provider",
						Version: "1.2.4",
					},
					Action: module.EndVersionRequest_PUBLISH,
				},
			},
			wantErr: false,
		},
		{
			name: "Fail during register()",
			args: args{
				client: &MockPublisherClient{
					registerError: errors.New("failed"),
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
			expected: expected{
				registerCalls:                      1,
				beginVersionCalls:                  0,
				registerContainerDependenciesCalls: 0,
				uploadSourceZipCalls:               0,
				uploadSourceZipSendCalls:           0,
				endVersionCalls:                    0,
			},
			wantErr: true,
		},
		{
			name: "Fail during register()",
			args: args{
				client: &MockPublisherClient{
					beginVersionError: errors.New("failed"),
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
			expected: expected{
				registerCalls:                      1,
				beginVersionCalls:                  1,
				registerContainerDependenciesCalls: 0,
				uploadSourceZipCalls:               0,
				uploadSourceZipSendCalls:           0,
				endVersionCalls:                    0,
			},
			wantErr: true,
		},
		{
			name: "Fail at upload",
			args: args{
				client: &MockPublisherClient{
					uploadSourceZipError: errors.New("failed"),
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
			expected: expected{
				registerCalls:                      1,
				beginVersionCalls:                  1,
				registerContainerDependenciesCalls: 0,
				uploadSourceZipCalls:               1,
				uploadSourceZipSendCalls:           0,
				endVersionCalls:                    1,
				endVersionRequest: &module.EndVersionRequest{
					Module: &module.Module{
						Name:    "org/test/provider",
						Version: "1.2.4",
					},
					Action: module.EndVersionRequest_DISCARD,
				},
			},
			wantErr: true,
		},
		{
			name: "Fail during upload send()",
			args: args{
				client: &MockPublisherClient{
					uploadSourceZipResponse: &MockPublisher_UploadSourceZipClient{
						sendError:         errors.New("failed"),
						closeAndRecvError: nil,
						closeAndRecvResponse: &module.Response{
							Message: "done",
						},
					},
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
			expected: expected{
				registerCalls:                      1,
				beginVersionCalls:                  1,
				registerContainerDependenciesCalls: 0,
				uploadSourceZipCalls:               1,
				uploadSourceZipSendCalls:           1,
				endVersionCalls:                    1,
				endVersionRequest: &module.EndVersionRequest{
					Module: &module.Module{
						Name:    "org/test/provider",
						Version: "1.2.4",
					},
					Action: module.EndVersionRequest_DISCARD,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Publish(tt.args.client, tt.args.source, tt.args.metadata); (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
			require.Equal(t, tt.expected.registerCalls, tt.args.client.registerCalls)
			require.Equal(t, tt.expected.beginVersionCalls, tt.args.client.beginVersionCalls)
			require.Equal(t, tt.expected.uploadSourceZipCalls, tt.args.client.uploadSourceZipCalls)
			require.Equal(t, tt.expected.endVersionCalls, tt.args.client.endVersionCalls)
			if tt.expected.endVersionRequest != nil {
				require.Equal(t, tt.expected.endVersionRequest, tt.args.client.endVersionRequest)
			}
			if tt.args.client.uploadSourceZipResponse != nil {
				require.Equal(t, tt.expected.uploadSourceZipSendCalls, tt.args.client.uploadSourceZipResponse.sendCalls)
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
	uploadSourceZipResponse               *MockPublisher_UploadSourceZipClient
	endVersionCalls                       int
	endVersionRequest                     *module.EndVersionRequest
	endVersionError                       error
	endVersionResponse                    *module.Response
	publishTagCalls                       int
	publishTagError                       error
	publishTagResponse                    *module.Response
}

func (m *MockPublisherClient) Register(ctx context.Context, in *module.RegisterModuleRequest, opts ...grpc.CallOption) (*module.Response, error) {
	m.registerCalls++
	return m.registerResponse, m.registerError
}

func (m *MockPublisherClient) BeginVersion(ctx context.Context, in *module.BeginVersionRequest, opts ...grpc.CallOption) (*module.Response, error) {
	m.beginVersionCalls++
	return m.beginVersionResponse, m.beginVersionError
}

func (m *MockPublisherClient) RegisterModuleDependencies(ctx context.Context, in *module.RegisterModuleDependenciesRequest, opts ...grpc.CallOption) (*module.Response, error) {
	m.registerModuleDependenciesCalls++
	return m.registerModuleDependenciesResponse, m.registerModuleDependenciesError
}

func (m *MockPublisherClient) RegisterContainerDependencies(ctx context.Context, in *module.RegisterContainerDependenciesRequest, opts ...grpc.CallOption) (*module.Response, error) {
	m.registerContainerDependenciesCalls++
	return m.registerContainerDependenciesResponse, m.registerContainerDependenciesError
}

func (m *MockPublisherClient) UploadSourceZip(ctx context.Context, opts ...grpc.CallOption) (module.Publisher_UploadSourceZipClient, error) {
	m.uploadSourceZipCalls++
	return m.uploadSourceZipResponse, m.uploadSourceZipError
}

func (m *MockPublisherClient) EndVersion(ctx context.Context, in *module.EndVersionRequest, opts ...grpc.CallOption) (*module.Response, error) {
	m.endVersionCalls++
	m.endVersionRequest = in
	return m.endVersionResponse, m.endVersionError
}

func (m *MockPublisherClient) PublishTag(ctx context.Context, in *module.PublishTagRequest, opts ...grpc.CallOption) (*module.Response, error) {
	m.publishTagCalls++
	return m.publishTagResponse, m.publishTagError
}

type MockPublisher_UploadSourceZipClient struct {
	grpc.ClientStream
	sendCalls            int
	sendError            error
	closeAndRecvError    error
	closeAndRecvResponse *module.Response
}

func (m *MockPublisher_UploadSourceZipClient) Send(request *module.UploadSourceZipRequest) error {
	m.sendCalls++
	return m.sendError
}

func (m *MockPublisher_UploadSourceZipClient) CloseAndRecv() (*module.Response, error) {
	return m.closeAndRecvResponse, m.closeAndRecvError
}
