package authorization

import (
	"context"

	"github.com/terrariumcloud/terrarium/internal/oauth"
	"github.com/terrariumcloud/terrarium/internal/oauth/services"
	"google.golang.org/grpc"
)

type AuthorizationServer struct {
	services.UnimplementedAuthorizationServer
}

// RegisterWithServer Registers AuthorizationServer with grpc server
func (a *AuthorizationServer) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	services.RegisterAuthorizationServer(grpcServer, a)
	return nil
}

func (a *AuthorizationServer) CreateApplication(ctx context.Context, req *oauth.CreateApplicationRequest) (*oauth.ApplicationResponse, error) {
	return nil, nil
}
