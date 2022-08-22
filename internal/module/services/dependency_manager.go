package services

import (
	"context"
	"encoding/json"

	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/storage"
	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"
	grpc "google.golang.org/grpc"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
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

	// filter := expression.Name("Name").Equal(expression.Value(request.Module.Name))
	// expr, err := expression.NewBuilder().WithFilter(filter).Build()
	// if err != nil {
	// 	return err
	// }

	// sin := &dynamodb.ScanInput{
	// 	ExpressionAttributeNames:  expr.Names(),
	// 	ExpressionAttributeValues: expr.Values(),
	// 	FilterExpression:          expr.Filter(),
	// 	TableName:                 aws.String(VersionsTableName),
	// }

	// sout, err := s.Db.Scan(sin)

	// if sout.Items == nil {
	// 	return err
	// }

	// moduleVersion := ModuleVersion{}

	// if *sout.Count > 1 {
	// 	return errors.New("unexpected number of results returned")
	// }

	// for _, i := range sout.Items {
	// 	if err := dynamodbattribute.UnmarshalMap(i, &moduleVersion); err != nil {
	// 		return err
	// 	}
	// }

	// in := &dynamodb.GetItemInput{
	// 	Key: map[string]*dynamodb.AttributeValue{
	// 		"_id": {
	// 			S: aws.String(moduleVersion.Version),
	// 		},
	// 	},
	// 	TableName: aws.String(ContainerDependenciesTableName),
	// }

	// out, err := s.Db.GetItem(in)

	// if out.Item == nil {
	// 	return err
	// }

	// dep := ContainerDependencies{}

	// if err := dynamodbattribute.UnmarshalMap(out.Item, &dep); err != nil {
	// 	return err
	// }

	dependencies := []string{}

	res := &terrarium.ContainerDependenciesResponse{
		Module:       request.Module,
		Dependencies: dependencies,
	}

	if err := server.Send(res); err != nil {
		return err
	}

	return nil
}

func (s *DependencyManagerService) RetrieveModuleDependencies(request *terrarium.RetrieveModuleDependenciesRequest, server DependencyManager_RetrieveModuleDependenciesServer) error {

	dependencies := []*terrarium.VersionedModule{}

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
				KeyType:       aws.String("HASH"),
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
