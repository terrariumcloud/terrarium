package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/storage"
	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"
	grpc "google.golang.org/grpc"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	DefaultModuleDependenciesTableName = "terrarium-module-dependencies"
	DefaultDependencyManagerEndpoint   = "dependency_manager:3001"
)

var ModuleDependenciesTableName string = DefaultModuleDependenciesTableName
var DependencyManagerEndpoint string = DefaultDependencyManagerEndpoint

type DependencyManagerService struct {
	UnimplementedDependencyManagerServer
	Db     dynamodbiface.DynamoDBAPI
	Table  string
	Schema *dynamodb.CreateTableInput
}

type ModuleDependencies struct {
	Name    string                      `json:"name" bson:"name" dynamodbav:"name"`
	Version string                      `json:"version" bson:"version" dynamodbav:"version"`
	Modules []terrarium.VersionedModule `json:"modules" bson:"modules" dynamodbav:"modules"`
	Images  []string                    `json:"images" bson:"images" dynamodbav:"images"`
}

func (s *DependencyManagerService) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	RegisterDependencyManagerServer(grpcServer, s)
	if err := storage.InitializeDynamoDb(s.Table, s.Schema, s.Db); err != nil {
		return err
	}
	return nil
}

// Registers Module dependencies in Terrarium
func (s *DependencyManagerService) RegisterModuleDependencies(ctx context.Context, request *terrarium.RegisterModuleDependenciesRequest) (*terrarium.TransactionStatusResponse, error) {
	dep, err := json.Marshal(request.Dependencies)

	if err != nil {
		return nil, err
	}

	in := &dynamodb.PutItemInput{
		TableName: aws.String(ModuleDependenciesTableName),
		Item: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(request.Module.GetName()),
			},
			"version": {
				S: aws.String(request.Module.GetVersion()),
			},
			"modules": {
				S: aws.String(string(dep)),
			},
		},
	}

	if _, err = s.Db.PutItem(in); err != nil {
		return RegisterModuleDependenciesFailed, err
	}

	return ModuleDependenciesRegistered, nil
}

func (s *DependencyManagerService) RegisterContainerDependencies(ctx context.Context, request *terrarium.RegisterContainerDependenciesRequest) (*terrarium.TransactionStatusResponse, error) {
	img, err := json.Marshal(request.Dependencies)

	if err != nil {
		return nil, err
	}

	in := &dynamodb.UpdateItemInput{
		TableName:        aws.String(ModuleDependenciesTableName),
		UpdateExpression: aws.String("set images = :images"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":images": {
				S: aws.String(string(img)),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(request.Module.GetName()),
			},
			"version": {
				S: aws.String(request.Module.GetVersion()),
			},
		},
	}

	_, err = s.Db.UpdateItem(in)

	if err != nil {
		return RegisterContainerDependenciesFailed, err
	}

	return ContainerDependenciesRegistered, nil
}

func (s *DependencyManagerService) RetrieveContainerDependencies(request *terrarium.RetrieveContainerDependenciesRequest, server DependencyManager_RetrieveContainerDependenciesServer) error {
	in := &dynamodb.GetItemInput{
		TableName: aws.String(ModuleDependenciesTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"name":    {S: aws.String(request.Module.Name)},
			"version": {S: aws.String(request.Module.Version)},
		},
	}

	out, err := s.Db.GetItem(in)

	if err != nil {
		return err
	}

	if out.Item == nil {
		return err
	}

	dependencies := ModuleDependencies{}

	if err := dynamodbattribute.UnmarshalMap(out.Item, &dependencies); err != nil {
		return err
	}

	res := &terrarium.ContainerDependenciesResponse{
		Module:       request.Module,
		Dependencies: dependencies.Images,
	}

	if err := server.Send(res); err != nil {
		return err
	}

	return nil
}

func (s *DependencyManagerService) RetrieveModuleDependencies(request *terrarium.RetrieveModuleDependenciesRequest, server DependencyManager_RetrieveModuleDependenciesServer) error {

	projEx := expression.NamesList(expression.Name("modules"))
	expr, err := expression.NewBuilder().WithProjection(projEx).Build()
	if err != nil {
		log.Printf("Couldn't build expressions %v\n", err)
	}

	output, err := s.Db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(ModuleDependenciesTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"name":    {S: aws.String(request.Module.Name)},
			"version": {S: aws.String(request.Module.Version)},
		},
		ProjectionExpression: expr.Projection(),
	})

	if output.Item == nil {
		return err
	}

	dependencies := []*terrarium.VersionedModule{}

	if err := dynamodbattribute.UnmarshalMap(output.Item, &dependencies); err != nil {
		return err
	}

	res := &terrarium.ModuleDependenciesResponse{
		Module:       request.Module,
		Dependencies: dependencies,
	}

	if err := server.Send(res); err != nil {
		return err
	}

	return nil
}

// GetModuleDependenciesSchema returns CreateTableInput
// that can be used to create table if it does not exist
func GetModuleDependenciesSchema(table string) *dynamodb.CreateTableInput {
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
		BillingMode: aws.String(dynamodb.BillingModeProvisioned),
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	}
}
