package services

import (
	"context"
	"log"
	"time"

	"github.com/terrariumcloud/terrarium/internal/storage"
	terrarium "github.com/terrariumcloud/terrarium/pkg/terrarium/module"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	//"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	//"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	DefaultTagTableName      = "terrarium-module-tags"
	DefaultTagManagerEndpoint = "tag_manager:3001"
)

var (
	TagTableName      string = DefaultTagTableName
	TagManagerEndpoint string = DefaultTagManagerEndpoint
	TagPublished = &terrarium.Response{Message: "Tag published."}
	ModuleTagTableInitializationError = status.Error(codes.Unknown, "Failed to initialize table for tags.")
	//MarshalModuleVersionError              = status.Error(codes.Unknown, "Failed to marshal module version.")
	PublishModuleTagError              = status.Error(codes.Unknown, "Failed to publish module tag.")
)

type TagManagerService struct {
	UnimplementedTagManagerServer
	Db     dynamodbiface.DynamoDBAPI
	Table  string
	Schema *dynamodb.CreateTableInput
}

type ModuleTag struct {
	Name        string `json:"name" bson:"name" dynamodbav:"name"`
	Version     string `json:"version" bson:"version" dynamodbav:"version"`
	Tag         []string `json:"tags" bson:"tags" dynamodbav:"tags"`
	CreatedOn   string `json:"created_on" bson:"created_on" dynamodbav:"created_on"`
	PublishedOn string `json:"published_on" bson:"published_on" dynamodbav:"published_on"`
}




// RegisterWithServer registers TagManagerService with grpc server
func (s *TagManagerService) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	if err := storage.InitializeDynamoDb(s.Table, s.Schema, s.Db); err != nil {
		log.Println(err)
		return ModuleTagTableInitializationError
	}

	RegisterTagManagerServer(grpcServer, s)  

	return nil
}


// PublishTag Updates Module Tag to published with Tag Manager service
func (s *TagManagerService) PublishTag(_ context.Context, request *PublishTagRequest) (*terrarium.Response, error) {
	log.Println("Publishing module tag.")
	in := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":published_on": {S: aws.String(time.Now().UTC().String())},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"name":    {S: aws.String(request.Module.GetName())},
			"version": {S: aws.String(request.Module.GetVersion())},			
		},
		TableName:        aws.String(TagTableName),
		UpdateExpression: aws.String("set published_on = :published_on"),
	}

	if _, err := s.Db.UpdateItem(in); err != nil {
		log.Println(err)
		return nil, PublishModuleTagError
	}

	log.Println("Module version published.")
	return TagPublished, nil
}

// ListModuleTags retrieve all tags of a given module for specified version and return an array of tags.

// func (s *TagManagerService) ListModuleTags(_ context.Context, request *ListModuleVersionsRequest) (*ListModuleVersionsResponse, error) {
// 	projection := expression.NamesList(expression.Name("version"))
// 	filter := expression.And(
// 		expression.Name("name").Equal(expression.Value(request.Module)),
// 		expression.Name("published_on").AttributeExists())
// 	expr, err := expression.NewBuilder().WithProjection(projection).WithFilter(filter).Build()
// 	if err != nil {
// 		log.Printf("Expression Builder failed creation: %v", err)
// 		return nil, err
// 	}

// 	scanQueryInputs := &dynamodb.ScanInput{
// 		ExpressionAttributeNames:  expr.Names(),
// 		ExpressionAttributeValues: expr.Values(),
// 		FilterExpression:          expr.Filter(),
// 		ProjectionExpression:      expr.Projection(),
// 		TableName:                 aws.String(VersionsTableName),
// 	}

// 	response, err := s.Db.Scan(scanQueryInputs)
// 	if err != nil {
// 		log.Printf("ScanInput failed: %v", err)
// 		return nil, err
// 	}

// 	grpcResponse := ListModuleVersionsResponse{}
// 	if response.Items != nil {
// 		for _, item := range response.Items {
// 			moduleVersion := ModuleVersion{}
// 			if err3 := dynamodbattribute.UnmarshalMap(item, &moduleVersion); err3 != nil {
// 				log.Printf("UnmarshalMap failed: %v", err3)
// 				return nil, err3
// 			}
// 			grpcResponse.Versions = append(grpcResponse.Versions, moduleVersion.Version)
// 		}
// 	}

// 	return &grpcResponse, nil
// }






// GetTagsSchema returns CreateTableInput that can be used to create table if it does not exist
func GetTagsSchema(table string) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("name"),
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
			},
			{
				AttributeName: aws.String("version"),
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
			},
			{
				AttributeName: aws.String("tag"),
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("version"),
				KeyType:       aws.String("RANGE"),
			},
		},
		TableName:   aws.String(table),
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
	}
}