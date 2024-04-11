package release

import (
	"context"
	"github.com/terrariumcloud/terrarium/internal/common/grpc_service"
	releaseSvc "github.com/terrariumcloud/terrarium/internal/release/services"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/release"
	"google.golang.org/grpc"
)

type browseGrpcClient struct {
	endpoint string
}

func NewBrowseGrpcClient(endpoint string) releaseSvc.BrowseClient {
	return &browseGrpcClient{endpoint: endpoint}
}

func (b browseGrpcClient) ListReleases(ctx context.Context, in *releaseSvc.ListReleasesRequest, opts ...grpc.CallOption) (*releaseSvc.ListReleasesResponse, error) {
	if conn, err := grpc_service.CreateGRPCConnection(b.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = conn.Close() }()

		client := releaseSvc.NewBrowseClient(conn)
		return client.ListReleases(ctx, in, opts...)
	}
}

func (b browseGrpcClient) ListReleaseTypes(ctx context.Context, in *releaseSvc.ListReleaseTypesRequest, opts ...grpc.CallOption) (*releaseSvc.ListReleaseTypesResponse, error) {
	if conn, err := grpc_service.CreateGRPCConnection(b.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = conn.Close() }()

		client := releaseSvc.NewBrowseClient(conn)
		return client.ListReleaseTypes(ctx, in, opts...)
	}
}

func (b browseGrpcClient) ListOrganization(ctx context.Context, in *releaseSvc.ListOrganizationRequest, opts ...grpc.CallOption) (*releaseSvc.ListOrganizationResponse, error) {
	if conn, err := grpc_service.CreateGRPCConnection(b.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = conn.Close() }()

		client := releaseSvc.NewBrowseClient(conn)
		return client.ListOrganization(ctx, in, opts...)
	}
}

type publisherGrpcClient struct {
	endpoint string
}

func NewPublisherGrpcClient(endpoint string) releaseSvc.PublisherClient {
	return &publisherGrpcClient{endpoint: endpoint}
}

func (r publisherGrpcClient) Publish(ctx context.Context, in *release.PublishRequest, opts ...grpc.CallOption) (*release.PublishResponse, error) {
	if conn, err := grpc_service.CreateGRPCConnection(r.endpoint); err != nil {
		return nil, err
	} else {
		defer func() { _ = conn.Close() }()

		client := releaseSvc.NewPublisherClient(conn)
		return client.Publish(ctx, in, opts...)
	}
}
