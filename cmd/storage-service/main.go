package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"terrarium-grpc-gateway/internal/services"
	"terrarium-grpc-gateway/pkg/terrarium"

	"google.golang.org/grpc"
)

type StorageService struct {
	services.UnimplementedStorageServer
}

func (s *StorageService) UploadSourceZip (server *terrarium.UploadSourceZipChunkRequest) error {
	return nil
}
func (s *StorageService) DownloadSourceZip(request *terrarium.DownloadSourceZipRequest, server *services.Storage_DownloadSourceZipServer) error {
	return nil
}

func main() {
	fmt.Println("Welcome to Terrarium GRPC API Gateway")
	flag.Parse()
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