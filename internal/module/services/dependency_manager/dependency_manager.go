package dependency_manager

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/terrariumcloud/terrarium/internal/module/services"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"

	"github.com/terrariumcloud/terrarium/internal/storage"
	terrarium "github.com/terrariumcloud/terrarium/pkg/terrarium/module"

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
	services.UnimplementedDependencyManagerServer
	Db              storage.DynamoDBTableCreator
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

// RegisterWithServer Registers DependencyManagerService with grpc server
func (s *DependencyManagerService) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	if err := storage.InitializeDynamoDb(s.ModuleTable, s.ModuleSchema, s.Db); err != nil {
		log.Println(err)
		return ModuleDependenciesTableInitializationError
	}

	if err := storage.InitializeDynamoDb(s.ContainerTable, s.ContainerSchema, s.Db); err != nil {
		log.Println(err)
		return ContainerDependenciesTableInitializationError
	}

	services.RegisterDependencyManagerServer(grpcServer, s)

	return nil
}

func (s *DependencyManagerService) registerDependencies(ctx context.Context, tableName string, in interface{}) error {
	marshalledItem, err := attributevalue.MarshalMap(in)
	if err != nil {
		log.Println(err)
		return MarshalDependenciesError
	}

	item := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      marshalledItem,
	}

	if _, err = s.Db.PutItem(ctx, item); err != nil {

		return RegisterDependenciesError
	}
	return nil
}

// Registers Module dependencies in Terrarium
func (s *DependencyManagerService) RegisterModuleDependencies(ctx context.Context, request *terrarium.RegisterModuleDependenciesRequest) (*terrarium.Response, error) {
	log.Printf("Registering module dependencies for %s/%s.\n", request.Module.GetName(), request.Module.GetVersion())
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("module.name", request.Module.GetName()),
		attribute.String("module.version", request.Module.GetVersion()),
	)
	item := ModuleDependencies{
		Name:    request.Module.GetName(),
		Version: request.Module.GetVersion(),
		Modules: request.GetDependencies(),
	}

	if err := s.registerDependencies(ctx, s.ModuleTable, item); err != nil {
		span.RecordError(err)
		log.Println(err)
		return nil, err
	}
	log.Printf("Module dependencies registered for %s/%s.\n", request.Module.GetName(), request.Module.GetVersion())
	return ModuleDependenciesRegistered, nil
}

// RegisterContainerDependencies Registers Container dependencies in Terrarium
func (s *DependencyManagerService) RegisterContainerDependencies(ctx context.Context, request *terrarium.RegisterContainerDependenciesRequest) (*terrarium.Response, error) {
	log.Printf("Registering container dependencies for %s/%s.\n", request.Module.GetName(), request.Module.GetVersion())
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("module.name", request.Module.GetName()),
		attribute.String("module.version", request.Module.GetVersion()),
	)
	item := ContainerDependencies{
		Name:    request.Module.GetName(),
		Version: request.Module.GetVersion(),
		Images:  request.Images,
	}

	if err := s.registerDependencies(ctx, s.ContainerTable, item); err != nil {
		span.RecordError(err)
		log.Println(err)
		return nil, err
	}
	log.Printf("Container dependencies registered for %s/%s.\n", request.Module.GetName(), request.Module.GetVersion())
	return ContainerDependenciesRegistered, nil
}

// RetrieveContainerDependencies Retrieve Container dependencies from Terrarium
func (s *DependencyManagerService) RetrieveContainerDependencies(request *terrarium.RetrieveContainerDependenciesRequestV2, server services.DependencyManager_RetrieveContainerDependenciesServer) error {
	log.Println("Retrieving container dependencies.")
	controlCh := make(chan *terrarium.Module, 250)
	controlCh <- request.Module

	moreModulesToProcess := true
	ctx := server.Context()
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("module.name", request.Module.GetName()),
		attribute.String("module.version", request.Module.GetVersion()),
	)
	for moreModulesToProcess {
		select {
		case moduleToProcess := <-controlCh:
			dep, err := s.GetModuleDependencies(ctx, moduleToProcess)
			if err != nil {
				log.Println(err)
				span.RecordError(err)
				return err
			}

			images, err := s.GetContainerDependencies(ctx, moduleToProcess)
			if err != nil {
				log.Println(err)
				span.RecordError(err)
				return err
			}

			res := &terrarium.ContainerDependenciesResponseV2{
				Module:       moduleToProcess,
				Dependencies: images,
			}
			if err := server.Send(res); err != nil {
				log.Println(err)
				span.RecordError(err)
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

func (s *DependencyManagerService) GetModuleKey(module *terrarium.Module) (map[string]types.AttributeValue, error) {
	moduleName, err := attributevalue.Marshal(module.GetName())
	if err != nil {
		return map[string]types.AttributeValue{}, err
	}
	moduleVersion, err := attributevalue.Marshal(module.GetVersion())
	if err != nil {
		return map[string]types.AttributeValue{}, err
	}
	return map[string]types.AttributeValue{
		"name":    moduleName,
		"version": moduleVersion,
	}, nil
}

func (s *DependencyManagerService) GetModuleDependencies(ctx context.Context, module *terrarium.Module) ([]*terrarium.Module, error) {
	log.Printf("GetModuleDependencies for module: %s/%s", module.GetName(), module.GetVersion())
	span := trace.SpanFromContext(ctx)
	moduleKey, err := s.GetModuleKey(module)
	if err != nil {
		log.Println(err)
		span.RecordError(err)
		return nil, GetModuleDependenciesError
	}

	in := &dynamodb.GetItemInput{
		TableName: aws.String(s.ModuleTable),
		Key:       moduleKey,
	}

	out, err := s.Db.GetItem(ctx, in)
	if err != nil {
		log.Println(err)
		span.RecordError(err)
		return nil, GetModuleDependenciesError
	}

	dependencies := ModuleDependencies{}
	if err := attributevalue.UnmarshalMap(out.Item, &dependencies); err != nil {
		log.Println(err)
		span.RecordError(err)
		return nil, UnmarshalModuleDependenciesError
	}
	log.Printf("GetModuleDependencies returned %d entries\n", len(dependencies.Modules))
	return dependencies.Modules, nil
}

func (s *DependencyManagerService) GetContainerDependencies(ctx context.Context, module *terrarium.Module) (map[string]*terrarium.ContainerImageDetails, error) {
	log.Printf("GetContainerDependencies for module: %s/%s\n", module.GetName(), module.GetVersion())
	span := trace.SpanFromContext(ctx)
	moduleKey, err := s.GetModuleKey(module)
	if err != nil {
		log.Println(err)
		span.RecordError(err)
		return nil, GetModuleDependenciesError
	}
	in := &dynamodb.GetItemInput{
		TableName: aws.String(s.ContainerTable),
		Key:       moduleKey,
	}

	out, err := s.Db.GetItem(ctx, in)
	if err != nil {
		log.Println(err)
		span.RecordError(err)
		return nil, GetContainerDependenciesError
	}

	dependencies := ContainerDependencies{}
	if err := attributevalue.UnmarshalMap(out.Item, &dependencies); err != nil {
		log.Println(err)
		span.RecordError(err)
		return nil, UnmarshalContainerDependenciesError
	}
	log.Printf("GetContainerDependencies returned %d entries\n", len(dependencies.Images))
	return dependencies.Images, nil
}

// RetrieveModuleDependencies Retrieve Module dependencies from Terrarium
func (s *DependencyManagerService) RetrieveModuleDependencies(request *terrarium.RetrieveModuleDependenciesRequest, server services.DependencyManager_RetrieveModuleDependenciesServer) error {
	ctx := server.Context()
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("module.name", request.Module.GetName()),
		attribute.String("module.version", request.Module.GetVersion()),
	)
	controlCh := make(chan *terrarium.Module, 250)
	controlCh <- request.Module

	moreModulesToProcess := true

	for moreModulesToProcess {
		select {
		case moduleToProcess := <-controlCh:
			dep, err := s.GetModuleDependencies(context.TODO(), moduleToProcess)
			if err != nil {
				log.Println(err)
				span.RecordError(err)
				return err
			}

			res := &terrarium.ModuleDependenciesResponse{
				Module:       moduleToProcess,
				Dependencies: dep,
			}
			if err := server.Send(res); err != nil {
				log.Println(err)
				span.RecordError(err)
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
