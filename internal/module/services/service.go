package services

import grpc "google.golang.org/grpc"

type Service interface {
	RegisterWithServer(grpcServer grpc.ServiceRegistrar) error
}
