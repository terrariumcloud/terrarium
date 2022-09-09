package services

import (
	"context"
	"log"
	"time"

	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/storage"
	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	DefaultVersionsTableName      = "terrarium-module-versions"
	DefaultVersionManagerEndpoint = "version_manager:3001"
)

var (
	VersionsTableName      string = DefaultVersionsTableName
	VersionManagerEndpoint string = DefaultVersionManagerEndpoint

	VersionCreated   = &terrarium.Response{Message: "Version created."}
	VersionPublished = &terrarium.Response{Message: "Version published."}
	VersionAborted   = &terrarium.Response{Message: "Version aborted."}

	ModuleVersionsTableInitializationError = status.Error(codes.Unknown, "Failed to initialize table for module versions.")
	MarshalModuleVersionError              = status.Error(codes.Unknown, "Failed to marshal module version.")
	CreateModuleVersionError               = status.Error(codes.Unknown, "Failed to create module version.")
	AbortModuleVersionError                = status.Error(codes.Unknown, "Failed to abort module version.")
	PublishModuleVersionError              = status.Error(codes.Unknown, "Failed to publish module version.")
)

type VersionManagerService struct {
	UnimplementedVersionManagerServer
	Db     dynamodbiface.DynamoDBAPI
	Table  string
	Schema *dynamodb.CreateTableInput
}

type ModuleVersion struct {
	Name        string `json:"name" bson:"name" dynamodbav:"name"`
	Version     string `json:"version" bson:"version" dynamodbav:"version"`
	CreatedOn   string `json:"created_on" bson:"created_on" dynamodbav:"created_on"`
	PublishedOn string `json:"published_on" bson:"published_on" dynamodbav:"published_on"`
}

// RegisterWithServer Registers VersionManagerService with grpc server
func (s *VersionManagerService) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	if err := storage.InitializeDynamoDb(s.Table, s.Schema, s.Db); err != nil {
		log.Println(err)
		return ModuleVersionsTableInitializationError
	}

	RegisterVersionManagerServer(grpcServer, s)

	return nil
}

// BeginVersion Creates new Module Version with Version Manager service
func (s *VersionManagerService) BeginVersion(_ context.Context, request *terrarium.BeginVersionRequest) (*terrarium.Response, error) {
	log.Println("Creating new version.")
	mv := ModuleVersion{
		Name:      request.Module.GetName(),
		Version:   request.Module.GetVersion(),
		CreatedOn: time.Now().UTC().String(),
	}

	av, err := dynamodbattribute.MarshalMap(mv)

	if err != nil {
		log.Println(err)
		return nil, MarshalModuleVersionError
	}

	in := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(VersionsTableName),
	}

	if _, err = s.Db.PutItem(in); err != nil {
		log.Println(err)
		return nil, CreateModuleVersionError
	}

	log.Println("New version created.")
	return VersionCreated, nil
}

// AbortVersion Removes Module Version with Version Manager service
func (s *VersionManagerService) AbortVersion(_ context.Context, request *TerminateVersionRequest) (*terrarium.Response, error) {
	log.Println("Aborting module version.")
	in := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"name":    {S: aws.String(request.Module.GetName())},
			"version": {S: aws.String(request.Module.GetVersion())},
		},
		TableName: aws.String(VersionsTableName),
	}

	if _, err := s.Db.DeleteItem(in); err != nil {
		log.Println(err)
		return nil, AbortModuleVersionError
	}

	log.Println("Module version aborted.")
	return VersionAborted, nil
}

// PublishVersion Updates Module Version to published with Version Manager service
func (s *VersionManagerService) PublishVersion(_ context.Context, request *TerminateVersionRequest) (*terrarium.Response, error) {
	log.Println("Publishing module version.")
	in := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":published_on": {S: aws.String(time.Now().UTC().String())},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"name":    {S: aws.String(request.Module.GetName())},
			"version": {S: aws.String(request.Module.GetVersion())},
		},
		TableName:        aws.String(VersionsTableName),
		UpdateExpression: aws.String("set published_on = :published_on"),
	}

	if _, err := s.Db.UpdateItem(in); err != nil {
		log.Println(err)
		return nil, PublishModuleVersionError
	}

	log.Println("Module version published.")
	return VersionPublished, nil
}

// ListModuleVersions Retrieve all versions of a given module and return an array of versions.
// Only versions that have been published should be reported
func (s *VersionManagerService) ListModuleVersions(_ context.Context, request *ListModuleVersionsRequest) (*ListModuleVersionsResponse, error) {
	filter := expression.And(expression.Name("name").Equal(expression.Value(request.Module)), expression.Name("published_on").AttributeExists())
	builder := expression.NewBuilder()
	builder.WithProjection(expression.NamesList(expression.Name("version")))
	builder.WithFilter(filter)
	expr, err := builder.Build()
	if err != nil {
		return nil, err
	}

	scanQueryInputs := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(VersionsTableName),
	}

	response, err := s.Db.Scan(scanQueryInputs)
	if response.Items == nil {
		return nil, err
	}

	grpcResponse := ListModuleVersionsResponse{}
	for _, item := range response.Items {
		moduleVersion := ModuleVersion{}
		if err := dynamodbattribute.UnmarshalMap(item, &moduleVersion); err != nil {
			return nil, err
		}
		grpcResponse.Versions = append(grpcResponse.Versions, moduleVersion.Version)
	}

	return &grpcResponse, nil
}

// GetModuleVersionsSchema returns CreateTableInput that can be used to create table if it does not exist
func GetModuleVersionsSchema(table string) *dynamodb.CreateTableInput {
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
