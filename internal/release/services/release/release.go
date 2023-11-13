package release

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	releaseSvc "github.com/terrariumcloud/terrarium/internal/release/services"
	"github.com/terrariumcloud/terrarium/internal/storage"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/release"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DefaultReleaseTableName = "terrarium-module-releases"
	//DefaultReleaseEndpoint  = "release:3001"
)

var (
	ReleaseTableName                = DefaultReleaseTableName
	ReleaseTableInitializationError = status.Error(codes.Unknown, "Failed to initialize table for releases.")
	//ReleaseEndpoint                 = DefaultReleaseEndpoint
)

type ReleaseService struct {
	releaseSvc.UnimplementedPublisherServer
	Db     storage.DynamoDBTableCreator
	Table  string
	Schema *dynamodb.CreateTableInput
}

type ReleaseInfo struct {
	Type             string       `json:"type" bson:"type" dynamodbav:"type"`
	Name             string       `json:"name" bson:"name" dynamodbav:"name"`
	Version          string       `json:"version" bson:"version" dynamodbav:"version"`
	Date             string       `json:"date" bson:"date" dynamodbav:"date"`
	Link             release.Link `json:"link" bson:"link" dynamodbav:"link"`
	ShortDescription string       `json:"short_description" bson:"short_description" dynamodbav:"short_description"`
	CreatedOn        string       `json:"created_on" bson:"created_on" dynamodbav:"created_on"`
	ModifiedOn       string       `json:"modified_on" bson:"modified_on" dynamodbav:"modified_on"`
}

// Registers ReleaseService with grpc server
func (s *ReleaseService) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	if err := storage.InitializeDynamoDb(s.Table, s.Schema, s.Db); err != nil {
		log.Println(err)
		return ReleaseTableInitializationError
	}

	releaseSvc.RegisterPublisherServer(grpcServer, s)

	return nil
}

// GetReleaseSchema returns CreateTableInput that can be used to create table if it does not exist
func GetReleaseSchema(table string) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("name"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("version"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("version"),
				KeyType:       types.KeyTypeRange,
			},
		},
		TableName:   aws.String(table),
		BillingMode: types.BillingModePayPerRequest,
	}
}
