package release

import (
	releaseSvc "github.com/terrariumcloud/terrarium/internal/release/services"
	"google.golang.org/grpc"
)

type ReleaseService struct {
	releaseSvc.UnimplementedPublisherServer
	// Db     storage.DynamoDBTableCreator
	// Table  string
	// Schema *dynamodb.CreateTableInput
}

// Registers ReleaseService with grpc server
func (s *ReleaseService) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	// if err := storage.InitializeDynamoDb(s.Table, s.Schema, s.Db); err != nil {
	// 	return ReleaseTableInitializationError
	// }

	releaseSvc.RegisterPublisherServer(grpcServer, s)

	return nil
}
