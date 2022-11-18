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
	DefaultModuleDependenciesTableName    = "terrarium-module-dependencies"
	DefaultContainerDependenciesTableName = "terrarium-container-dependencies"
	DefaultDependencyManagerEndpoint      = "dependency_manager:3001"
)

var (
	ModuleDependenciesTableName    = DefaultModuleDependenciesTableName
	ContainerDependenciesTableName = DefaultContainerDependenciesTableName
	DependencyManagerEndpoint      = DefaultDependencyManagerEndpoint

	ModuleDependenciesRegistered    = &terrarium.Response{Message: "Module dependencies successfully registered."}
	ContainerDependenciesRegistered = &terrarium.Response{Message: "Container dependencies successfully registered."}

	ModuleDependenciesTableInitializationError    = status.Error(codes.Unavailable, "Failed to initialize table for module dependencies.")
	ContainerDependenciesTableInitializationError = status.Error(codes.Unavailable, "Failed to initialize table for container dependencies.")
	RegisterDependenciesError                     = status.Error(codes.Unknown, "Failed to register dependencies.")
	MarshalDependenciesError                      = status.Error(codes.Unknown, "Failed to marshal dependencies.")
	SendModuleDependenciesError                   = status.Error(codes.Unknown, "Failed to send module dependencies.")
	SendContainerDependenciesError                = status.Error(codes.Unknown, "Failed to send container dependencies.")
	UnmarshalModuleDependenciesError              = status.Error(codes.Unknown, "Failed to unmarshal module dependencies.")
	UnmarshalContainerDependenciesError           = status.Error(codes.Unknown, "Failed to unmarshal container dependencies.")
	GetModuleDependenciesError                    = status.Error(codes.Unknown, "Failed to get module dependencies.")
	GetContainerDependenciesError                 = status.Error(codes.Unknown, "Failed to get container dependencies.")
)

type DependencyManagerService struct {
	UnimplementedDependencyManagerServer
	Db              dynamodbiface.DynamoDBAPI
	ModuleTable     string
	ModuleSchema    *dynamodb.CreateTableInput
	ContainerTable  string
	ContainerSchema *dynamodb.CreateTableInput
}

type ModuleDependencies struct {
	Name    string              `json:"name" bson:"name" dynamodbav:"name"`
	Version string              `json:"version" bson:"version" dynamodbav:"version"`
	Modules []*terrarium.Module `json:"modules" bson:"modules" dynamodbav:"modules"`
}
type ContainerDependencies struct {
	Name    string                                      `json:"name" bson:"name" dynamodbav:"name"`
	Version string                                      `json:"version" bson:"version" dynamodbav:"version"`
	Images  map[string]*terrarium.ContainerImageDetails `json:"images" bson:"images" dynamodbav:"images"`
}

// Registers DependencyManagerService with grpc server
func (s *DependencyManagerService) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	if err := storage.InitializeDynamoDb(s.ModuleTable, s.ModuleSchema, s.Db); err != nil {
		log.Println(err)
		return ModuleDependenciesTableInitializationError
	}

	if err := storage.InitializeDynamoDb(s.ContainerTable, s.ContainerSchema, s.Db); err != nil {
		log.Println(err)
		return ContainerDependenciesTableInitializationError
	}

	RegisterDependencyManagerServer(grpcServer, s)

	return nil
}

func (s *DependencyManagerService) registerDependencies(tableName string, in interface{}) error {
	marshalledItem, err := dynamodbattribute.MarshalMap(in)
	if err != nil {
		log.Println(err)
		return MarshalDependenciesError
	}

	item := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      marshalledItem,
	}

	if _, err = s.Db.PutItem(item); err != nil {
		log.Println(err)
		return RegisterDependenciesError
	}
	return nil
}

// Registers Module dependencies in Terrarium
func (s *DependencyManagerService) RegisterModuleDependencies(ctx context.Context, request *terrarium.RegisterModuleDependenciesRequest) (*terrarium.Response, error) {
	log.Printf("Registering module dependencies for %s/%s.\n", request.Module.GetName(), request.Module.GetVersion())

	item := ModuleDependencies{
		Name:    request.Module.GetName(),
		Version: request.Module.GetVersion(),
		Modules: request.GetDependencies(),
	}
	err := s.registerDependencies(ModuleDependenciesTableName, item)
	if err != nil {
		return nil, err
	}
	log.Printf("Module dependencies registered for %s/%s.\n", request.Module.GetName(), request.Module.GetVersion())
	return ModuleDependenciesRegistered, err
}

// Registers Container dependencies in Terrarium
func (s *DependencyManagerService) RegisterContainerDependencies(ctx context.Context, request *terrarium.RegisterContainerDependenciesRequest) (*terrarium.Response, error) {
	log.Printf("Registering container dependencies for %s/%s.\n", request.Module.GetName(), request.Module.GetVersion())

	item := ContainerDependencies{
		Name:    request.Module.GetName(),
		Version: request.Module.GetVersion(),
		Images:  request.Images,
	}

	err := s.registerDependencies(ContainerDependenciesTableName, item)
	if err != nil {
		return nil, err
	}
	log.Printf("Container dependencies registered for %s/%s.\n", request.Module.GetName(), request.Module.GetVersion())
	return ContainerDependenciesRegistered, nil
}

// Retrieve Container dependencies from Terrarium
func (s *DependencyManagerService) RetrieveContainerDependencies(request *terrarium.RetrieveContainerDependenciesRequestV2, server DependencyManager_RetrieveContainerDependenciesServer) error {
	log.Println("Retrieving container dependencies.")
	controlCh := make(chan *terrarium.Module, 250)
	controlCh <- request.Module

	moreModulesToProcess := true

	for moreModulesToProcess {
		select {
		case moduleToProcess := <-controlCh:
			dep, err := s.GetModuleDependencies(moduleToProcess)
			if err != nil {
				return err
			}

			images, err := s.GetContainerDependencies(moduleToProcess)
			if err != nil {
				return err
			}

			res := &terrarium.ContainerDependenciesResponseV2{
				Module:       moduleToProcess,
				Dependencies: images,
			}
			if err := server.Send(res); err != nil {
				log.Println(err)
				return SendContainerDependenciesError
			}

			for _, dependency := range dep {
				controlCh <- dependency
			}
		default:
			moreModulesToProcess = false
			close(controlCh)
		}

	}

	log.Println("Container dependencies retrieved.")
	return nil
}

func (s *DependencyManagerService) GetModuleDependencies(module *terrarium.Module) ([]*terrarium.Module, error) {
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

func (s *DependencyManagerService) GetContainerDependencies(module *terrarium.Module) (map[string]*terrarium.ContainerImageDetails, error) {
	in := &dynamodb.GetItemInput{
		TableName: aws.String(ContainerDependenciesTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"name":    {S: aws.String(module.GetName())},
			"version": {S: aws.String(module.GetVersion())},
		},
	}

	out, err := s.Db.GetItem(in)
	if err != nil {
		log.Println(err)
		return nil, GetContainerDependenciesError
	}

	dependencies := ContainerDependencies{}
	if err := dynamodbattribute.UnmarshalMap(out.Item, &dependencies); err != nil {
		log.Println(err)
		return nil, UnmarshalContainerDependenciesError
	}

	return dependencies.Images, nil
}

// Retrieve Module dependencies from Terrarium
func (s *DependencyManagerService) RetrieveModuleDependencies(request *terrarium.RetrieveModuleDependenciesRequest, server DependencyManager_RetrieveModuleDependenciesServer) error {
	controlCh := make(chan *terrarium.Module, 250)
	controlCh <- request.Module

	moreModulesToProcess := true

	for moreModulesToProcess {
		select {
		case moduleToProcess := <-controlCh:
			dep, err := s.GetModuleDependencies(moduleToProcess)
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

// GetDependenciesSchema returns CreateTableInput
// that can be used to create table if it does not exist
func GetDependenciesSchema(table string) *dynamodb.CreateTableInput {
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
