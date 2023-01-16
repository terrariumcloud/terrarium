package services

import "google.golang.org/grpc"

type Service interface {
	RegisterWithServer(grpcServer grpc.ServiceRegistrar) error
}
