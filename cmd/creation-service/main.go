package main

import (
	"log"
	"net"
	"terrarium-grpc-gateway/internal/services"

	"google.golang.org/grpc"
)

func main() {
	log.Println("Starting Terrarium GRPC Creation Service")

	lis, err := net.Listen("tcp", "0.0.0.0:9443")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	creationServiceServer := &CreationService{}
	// Need TLS
	grpcServer := grpc.NewServer(opts...)
	services.RegisterCreatorServer(grpcServer, creationServiceServer)
	grpcServer.Serve(lis)
}
