package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"terrarium-grpc-gateway/internal/services"
	"github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium"

	"google.golang.org/grpc"
)

type SessionService struct {
	services.UnimplementedSessionManagerServer
}

func (s *SessionService) BeginVersion(ctx context.Context, request *services.BeginVersionRequest) (*terrarium.BeginVersionResponse, error) {
	return nil, nil
}

func (s *SessionService) AbortVersion(ctx context.Context, request *services.TerminateVersionRequest) (*terrarium.TransactionStatusResponse, error) {
	return nil, nil
}

func (s *SessionService) PublishVersion(ctx context.Context, request *services.TerminateVersionRequest) (*terrarium.TransactionStatusResponse, error) {
	return nil, nil
}

func main() {
	fmt.Println("Welcome to Terrarium GRPC API Gateway")
	flag.Parse()
	lis, err := net.Listen("tcp", "0.0.0.0:9443")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	sessionServiceServer := &SessionService{}
	// Need TLS
	grpcServer := grpc.NewServer(opts...)
	services.RegisterSessionManagerServer(grpcServer, sessionServiceServer)
	grpcServer.Serve(lis)
}