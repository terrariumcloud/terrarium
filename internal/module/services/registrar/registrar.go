package registrar

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/terrariumcloud/terrarium/internal/module/services"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/terrariumcloud/terrarium/internal/storage"
	terrarium "github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	DefaultRegistrarTableName       = "terrarium-modules"
	DefaultRegistrarServiceEndpoint = "registrar:3001"
	DefaultApiKeyTableName          = "terrarium-organization-api-keys"
)

var (
	RegistrarTableName       = DefaultRegistrarTableName
	RegistrarServiceEndpoint = DefaultRegistrarServiceEndpoint
	ApiKeyTableName          = DefaultApiKeyTableName

	ModuleRegistered = &terrarium.Response{Message: "Module registered successfully."}

	ModuleTableInitializationError = status.Error(codes.Unknown, "Failed to initialize table for modules.")
	ApiKeyTableInitializationError = status.Error(codes.Unknown, "Failed to initialize table for organization API keys.")
	ModuleGetError                 = status.Error(codes.Unknown, "Failed to check if module already exists.")
	ModuleUpdateError              = status.Error(codes.Unknown, "Failed to update module.")
	ModuleRegisterError            = status.Error(codes.Unknown, "Failed to register module.")
	ExpressionBuildError           = status.Error(codes.Unknown, "Failed to build update expression.")
	MarshalModuleError             = status.Error(codes.Unknown, "Failed to marshal module.")
	ApiKeyValidationError          = status.Error(codes.Unauthenticated, "Invalid API key for organization.")
	ApiKeyGetError                 = status.Error(codes.Unknown, "Failed to retrieve organization API key.")
)

type RegistrarService struct {
	services.UnimplementedRegistrarServer
	Db           storage.DynamoDBTableCreator
	Table        string
	Schema       *dynamodb.CreateTableInput
	ApiKeyTable  string
	ApiKeySchema *dynamodb.CreateTableInput
}

type Module struct {
	Name        string `json:"name" bson:"name" dynamodbav:"name"`
	Description string `json:"description" bson:"description" dynamodbav:"description"`
	Source      string `json:"source_url" bson:"source_url" dynamodbav:"source"`
	Maturity    string `json:"maturity" bson:"maturity" dynamodbav:"maturity"`
	CreatedOn   string `json:"created_on" bson:"created_on" dynamodbav:"created_on"`
	ModifiedOn  string `json:"modified_on" bson:"modified_on" dynamodbav:"modified_on"`
}

type OrganizationApiKey struct {
	OrganizationName string `json:"organization_name" bson:"organization_name" dynamodbav:"organization_name"`
	ApiKey           string `json:"api_key" bson:"api_key" dynamodbav:"api_key"`
	CreatedOn        string `json:"created_on" bson:"created_on" dynamodbav:"created_on"`
	ModifiedOn       string `json:"modified_on" bson:"modified_on" dynamodbav:"modified_on"`
}

// Registers RegistrarService with grpc server
func (s *RegistrarService) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	if err := storage.InitializeDynamoDb(s.Table, s.Schema, s.Db); err != nil {
		return ModuleTableInitializationError
	}

	if err := storage.InitializeDynamoDb(s.ApiKeyTable, s.ApiKeySchema, s.Db); err != nil {
		return ApiKeyTableInitializationError
	}

	services.RegisterRegistrarServer(grpcServer, s)

	return nil
}

// Register new Module in Terrarium
func (s *RegistrarService) Register(ctx context.Context, request *terrarium.RegisterModuleRequest) (*terrarium.Response, error) {
	log.Println("Registering new module.")

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("module.name", request.GetName()),
	)

	// Validate API key for the organization
	if err := s.validateApiKey(ctx, request.GetName(), request.GetApiKey()); err != nil {
		span.RecordError(err)
		log.Printf("API key validation failed for module %s: %v", request.GetName(), err)
		return nil, err
	}

	name, err := attributevalue.Marshal(request.GetName())
	if err != nil {
		span.RecordError(err)
		log.Println(err)
		return nil, ModuleGetError
	}
	res, err := s.Db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(RegistrarTableName),
		Key: map[string]types.AttributeValue{
			"name": name,
		},
	})

	if err != nil {
		span.RecordError(err)
		log.Println(err)
		return nil, ModuleGetError
	}

	if res.Item == nil {
		ms := Module{
			Name:        request.GetName(),
			Description: request.GetDescription(),
			Source:      request.GetSource(),
			Maturity:    request.GetMaturity().String(),
			CreatedOn:   time.Now().UTC().String(),
			ModifiedOn:  time.Now().UTC().String(),
		}

		av, err := attributevalue.MarshalMap(ms)

		if err != nil {
			span.RecordError(err)
			log.Println(err)
			return nil, MarshalModuleError
		}

		in := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(RegistrarTableName),
		}

		if _, err = s.Db.PutItem(ctx, in); err != nil {
			span.RecordError(err)
			log.Println(err)
			return nil, ModuleRegisterError
		}
	} else {
		update := expression.Set(expression.Name("description"), expression.Value(request.GetDescription()))
		update.Set(expression.Name("source"), expression.Value(request.GetSource()))
		update.Set(expression.Name("maturity"), expression.Value(request.GetMaturity().String()))
		update.Set(expression.Name("modified_on"), expression.Value(time.Now().UTC().String()))
		expr, err := expression.NewBuilder().WithUpdate(update).Build()

		if err != nil {
			span.RecordError(err)
			log.Println(err)
			return nil, ExpressionBuildError
		}

		in := &dynamodb.UpdateItemInput{
			TableName: aws.String(RegistrarTableName),
			Key: map[string]types.AttributeValue{
				"name": name},
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
		}

		_, err = s.Db.UpdateItem(ctx, in)

		if err != nil {
			span.RecordError(err)
			log.Println(err)
			return nil, ModuleUpdateError
		}
	}

	log.Println("New module registered.")
	return ModuleRegistered, nil
}

func unmarshalModule(item map[string]types.AttributeValue) (*services.ModuleMetadata, error) {
	module := Module{}
	if err := attributevalue.UnmarshalMap(item, &module); err != nil {
		log.Printf("UnmarshalMap failed: %v", err)
		return nil, err
	}
	moduleAddress := strings.Split(module.Name, "/")

	result := services.ModuleMetadata{
		Organization: moduleAddress[0],
		Name:         moduleAddress[1],
		Provider:     moduleAddress[2],
		Description:  module.Description,
		SourceUrl:    module.Source,
		Maturity:     terrarium.Maturity(terrarium.Maturity_value[module.Maturity]),
	}
	return &result, nil
}

// validateApiKey validates the API key for the organization extracted from the module name
func (s *RegistrarService) validateApiKey(ctx context.Context, moduleName, apiKey string) error {
	// Extract organization name from module name
	moduleNameParts := strings.Split(moduleName, "/")
	if len(moduleNameParts) == 0 {
		return ApiKeyValidationError
	}
	organizationName := moduleNameParts[0]

	// Check if organization has API keys in the database
	orgName, err := attributevalue.Marshal(organizationName)
	if err != nil {
		log.Printf("Failed to marshal organization name: %v", err)
		return ApiKeyGetError
	}

	res, err := s.Db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(ApiKeyTableName),
		Key: map[string]types.AttributeValue{
			"organization_name": orgName,
		},
	})

	if err != nil {
		log.Printf("Failed to get organization API key: %v", err)
		return ApiKeyGetError
	}

	// If organization not found in API keys table, no validation required
	if res.Item == nil {
		log.Printf("No API key configured for organization: %s", organizationName)
		return nil
	}

	// Unmarshal the stored API key
	var storedApiKey OrganizationApiKey
	if err := attributevalue.UnmarshalMap(res.Item, &storedApiKey); err != nil {
		log.Printf("Failed to unmarshal organization API key: %v", err)
		return ApiKeyGetError
	}

	// Validate the provided API key against the stored one
	if storedApiKey.ApiKey != apiKey {
		log.Printf("Invalid API key for organization: %s", organizationName)
		return ApiKeyValidationError
	}

	log.Printf("API key validated successfully for organization: %s", organizationName)
	return nil
}

// GetModule Retrieve module metadata
func (s *RegistrarService) GetModule(ctx context.Context, request *services.GetModuleRequest) (*services.GetModuleResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("module.name", request.GetName()),
	)

	filter := expression.Name("name").Equal(expression.Value(request.Name))
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		span.RecordError(err)
		log.Printf("Expression Builder failed creation: %v", err)
		return nil, err
	}

	scanQueryInputs := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(RegistrarTableName),
	}
	response, err := s.Db.Scan(ctx, scanQueryInputs)
	if err != nil {
		span.RecordError(err)
		log.Printf("ScanInput failed: %v", err)
		return nil, err
	}

	if response.Items == nil || len(response.Items) < 1 {
		return nil, fmt.Errorf("module not found '%v'", request.GetName())
	}
	grpcResponse := services.GetModuleResponse{}
	if moduleMetadata, err := unmarshalModule(response.Items[0]); err != nil {
		span.RecordError(err)
		return nil, err
	} else {
		grpcResponse.Module = moduleMetadata
		return &grpcResponse, nil
	}
}

// ListModules Retrieve all published modules
func (s *RegistrarService) ListModules(ctx context.Context, request *services.ListModulesRequest) (*services.ListModulesResponse, error) {

	scanQueryInputs := &dynamodb.ScanInput{
		TableName: aws.String(RegistrarTableName),
	}

	response, err := s.Db.Scan(ctx, scanQueryInputs)
	if err != nil {
		log.Printf("ScanInput failed: %v", err)
		return nil, err
	}

	grpcResponse := services.ListModulesResponse{}
	if response.Items != nil {
		for _, item := range response.Items {
			if moduleMetadata, err := unmarshalModule(item); err != nil {
				return nil, err
			} else {
				grpcResponse.Modules = append(grpcResponse.Modules, moduleMetadata)
			}
		}
	}

	return &grpcResponse, nil
}

// GetModulesSchema returns CreateTableInput
// that can be used to create table if it does not exist
func GetModulesSchema(table string) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("name"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName:   aws.String(table),
		BillingMode: types.BillingModePayPerRequest,
	}
}

// GetOrganizationApiKeysSchema returns CreateTableInput
// that can be used to create organization API keys table if it does not exist
func GetOrganizationApiKeysSchema(table string) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("organization_name"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("organization_name"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName:   aws.String(table),
		BillingMode: types.BillingModePayPerRequest,
	}
}
