package releaseMocks

import (
	"context"

	"github.com/terrariumcloud/terrarium/internal/release/services"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/release"

	"google.golang.org/grpc"
)

type MockPublisherClient struct {
	services.PublisherClient
	PublishInvocations int
	PublishResponse    *release.PublishResponse
	PublishError       error
}

func (m *MockPublisherClient) Publish(ctx context.Context, in *release.PublishRequest, opts ...grpc.CallOption) (*release.PublishResponse, error) {
	m.PublishInvocations++
	return m.PublishResponse, m.PublishError
}

// ListReleases, ListReleaseTypes, ListOrganization
type MockBrowseClient struct {
	services.BrowseClient
	ListReleaseTypesInvocations int
	ListReleaseTypesResponse    *services.ListReleaseTypesResponse
	ListReleaseTypesError       error
	ListOrganizationInvocations int
	ListOrganizationResponse    *services.ListOrganizationResponse
	ListOrganizationError       error
}

func (m *MockBrowseClient) ListReleaseTypes(ctx context.Context, in *services.ListReleaseTypesRequest, opts ...grpc.CallOption) (*services.ListReleaseTypesResponse, error) {
	m.ListReleaseTypesInvocations++
	return m.ListReleaseTypesResponse, m.ListReleaseTypesError
}

func (m *MockBrowseClient) ListOrganization(ctx context.Context, in *services.ListOrganizationRequest, opts ...grpc.CallOption) (*services.ListOrganizationResponse, error) {
	m.ListOrganizationInvocations++
	return m.ListOrganizationResponse, m.ListOrganizationError
}
