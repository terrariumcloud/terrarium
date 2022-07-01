package main

import (
	"log"
	"net"
	"terrarium-grpc-gateway/internal/services"

	"google.golang.org/grpc"
)

func main() {
	log.Println("Starting Terrarium GRPC Session Service")

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
