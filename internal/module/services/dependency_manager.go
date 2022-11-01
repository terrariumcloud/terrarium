package services

import (
	"context"
	"log"

	"github.com/terrariumcloud/terrarium/internal/storage"
	terrarium "github.com/terrariumcloud/terrarium/pkg/terrarium/module"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DefaultModuleDependenciesTableName = "terrarium-module-dependencies"
	DefaultDependencyManagerEndpoint   = "dependency_manager:3001"
)

var (
	ModuleDependenciesTableName string = DefaultModuleDependenciesTableName
	DependencyManagerEndpoint   string = DefaultDependencyManagerEndpoint

	ModuleDependenciesRegistered    = &terrarium.Response{Message: "Module dependencies successfully registered."}
	ContainerDependenciesRegistered = &terrarium.Response{Message: "Container dependencies successfully registered."}

	ModuleDependenciesTableInitializationError = status.Error(codes.Unknown, "Failed to initialize table for module dependencies.")
	RegisterModuleDependenciesError            = status.Error(codes.Unknown, "Failed to register module dependencies.")
	RegisterContainerDependenciesError         = status.Error(codes.Unknown, "Failed to register container dependencies.")
	MarshalModuleDependenciesError             = status.Error(codes.Unknown, "Failed to marshal module dependencies.")
	MarshalContainerDependenciesError          = status.Error(codes.Unknown, "Failed to marshal container dependencies.")
	SendModuleDependenciesError                = status.Error(codes.Unknown, "Failed to send module dependencies.")
	SendContainerDependenciesError             = status.Error(codes.Unknown, "Failed to send container dependencies.")
	UnmarshalModuleDependenciesError           = status.Error(codes.Unknown, "Failed to unmarshal module dependencies.")
	UnmarshalContainerDependenciesError        = status.Error(codes.Unknown, "Failed to unmarshal container dependencies.")
	GetModuleDependenciesError                 = status.Error(codes.Unknown, "Failed to get module dependencies.")
	GetContainerDependenciesError              = status.Error(codes.Unknown, "Failed to get container dependencies.")
)

type DependencyManagerService struct {
	UnimplementedDependencyManagerServer
	Db     dynamodbiface.DynamoDBAPI
	Table  string
	Schema *dynamodb.CreateTableInput
}

type ModuleDependencies struct {
	Name    string              `json:"name" bson:"name" dynamodbav:"name"`
	Version string              `json:"version" bson:"version" dynamodbav:"version"`
	Modules []*terrarium.Module `json:"modules" bson:"modules" dynamodbav:"modules"`
	Images  []string            `json:"images" bson:"images" dynamodbav:"images"`
}

// Registers DependencyManagerService with grpc server
func (s *DependencyManagerService) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	if err := storage.InitializeDynamoDb(s.Table, s.Schema, s.Db); err != nil {
		log.Println(err)
		return ModuleDependenciesTableInitializationError
	}

	RegisterDependencyManagerServer(grpcServer, s)

	return nil
}

// Registers Module dependencies in Terrarium
func (s *DependencyManagerService) RegisterModuleDependencies(ctx context.Context, request *terrarium.RegisterModuleDependenciesRequest) (*terrarium.Response, error) {
	log.Println("Registering module dependencies.")
	dep := []*terrarium.Module{}
	for _, dependency := range request.Dependencies {
		dep = append(dep, &terrarium.Module{Name: dependency.Name, Version: dependency.Version})
	}

	depList, err := dynamodbattribute.Marshal(dep)

	if err != nil {
		log.Println(err)
		return nil, MarshalModuleDependenciesError
	}

	in := &dynamodb.PutItemInput{
		TableName: aws.String(ModuleDependenciesTableName),
		Item: map[string]*dynamodb.AttributeValue{
			"name":    {S: aws.String(request.Module.GetName())},
			"version": {S: aws.String(request.Module.GetVersion())},
			"modules": depList,
		},
	}

	if _, err = s.Db.PutItem(in); err != nil {
		log.Println(err)
		return nil, RegisterModuleDependenciesError
	}

	log.Println("Module dependencies registered.")
	return ModuleDependenciesRegistered, nil
}

// Registers Container dependencies in Terrarium
func (s *DependencyManagerService) RegisterContainerDependencies(ctx context.Context, request *terrarium.RegisterContainerDependenciesRequest) (*terrarium.Response, error) {
	log.Println("Registering container dependencies.")
	in := &dynamodb.UpdateItemInput{
		TableName:        aws.String(ModuleDependenciesTableName),
		UpdateExpression: aws.String("set images = :images"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":images": {SS: aws.StringSlice(request.Dependencies)}},
		Key: map[string]*dynamodb.AttributeValue{
			"name":    {S: aws.String(request.Module.GetName())},
			"version": {S: aws.String(request.Module.GetVersion())},
		},
	}

	_, err := s.Db.UpdateItem(in)

	if err != nil {
		log.Println(err)
		return nil, RegisterContainerDependenciesError
	}

	log.Println("Container dependencies registered.")
	return ContainerDependenciesRegistered, nil
}

// Retrieve Container dependencies from Terrarium
func (s *DependencyManagerService) RetrieveContainerDependencies(request *terrarium.RetrieveContainerDependenciesRequest, server DependencyManager_RetrieveContainerDependenciesServer) error {
	log.Println("Retrieving container dependencies.")
	in := &dynamodb.GetItemInput{
		TableName: aws.String(ModuleDependenciesTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"name":    {S: aws.String(request.Module.GetName())},
			"version": {S: aws.String(request.Module.GetVersion())},
		},
	}

	out, err := s.Db.GetItem(in)

	if err != nil {
		log.Println(err)
		return GetContainerDependenciesError
	}

	dependencies := ModuleDependencies{}

	if err := dynamodbattribute.UnmarshalMap(out.Item, &dependencies); err != nil {
		log.Println(err)
		return UnmarshalContainerDependenciesError
	}

	res := &terrarium.ContainerDependenciesResponse{
		Module:       request.Module,
		Dependencies: dependencies.Images,
	}

	if err := server.Send(res); err != nil {
		log.Println(err)
		return SendContainerDependenciesError
	}

	log.Println("Container dependencies retrieved.")
	return nil
}

func (s *DependencyManagerService) GetDependencies(module *terrarium.Module) ([]*terrarium.Module, error) {
	log.Println("Retrieving module dependencies.")

	in := &dynamodb.GetItemInput{
		TableName: aws.String(ModuleDependenciesTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"name":    {S: aws.String(module.GetName())},
			"version": {S: aws.String(module.GetVersion())},
		},
	}

	out, err := s.Db.GetItem(in)

	if err != nil {
		log.Println(err)
		return nil, GetModuleDependenciesError
	}

	dependencies := ModuleDependencies{}

	if err := dynamodbattribute.UnmarshalMap(out.Item, &dependencies); err != nil {
		log.Println(err)
		return nil, UnmarshalModuleDependenciesError
	}
	return dependencies.Modules, nil
}

// Retrieve Module dependencies from Terrarium
func (s *DependencyManagerService) RetrieveModuleDependencies(request *terrarium.RetrieveModuleDependenciesRequest, server DependencyManager_RetrieveModuleDependenciesServer) error {
	controlCh := make(chan *terrarium.Module, 250)
	controlCh <- request.Module

	moreModulesToProcess := true

	for moreModulesToProcess {
		select {
		case moduleToProcess := <-controlCh:
			dep, err := s.GetDependencies(moduleToProcess)
			if err != nil {
				return err
			}

			res := &terrarium.ModuleDependenciesResponse{
				Module:       moduleToProcess,
				Dependencies: dep,
			}
			if err := server.Send(res); err != nil {
				log.Println(err)
				return SendModuleDependenciesError
			}

			for _, dependency := range dep {
				controlCh <- dependency
			}
		default:
			moreModulesToProcess = false
			close(controlCh)
		}

	}

	log.Println("Module dependencies retrieved.")
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
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
	}
}
