package main

import (
	"log"
	"net"
	"terrarium-grpc-gateway/internal/services"

	"google.golang.org/grpc"
)

func main() {
	log.Println("Starting Terrarium GRPC Storage Service")

	lis, err := net.Listen("tcp", "0.0.0.0:9443")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	storageServiceServer := &StorageService{}
	// Need TLS
	grpcServer := grpc.NewServer(opts...)
	services.RegisterStorageServer(grpcServer, storageServiceServer)
	grpcServer.Serve(lis)
}
