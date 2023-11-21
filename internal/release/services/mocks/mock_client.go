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
