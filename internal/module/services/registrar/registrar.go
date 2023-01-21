package registrar

import (
	"context"
	"fmt"
	"github.com/terrariumcloud/terrarium/internal/module/services"
	"log"
	"strings"
	"time"

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
)

var (
	RegistrarTableName       = DefaultRegistrarTableName
	RegistrarServiceEndpoint = DefaultRegistrarServiceEndpoint

	ModuleRegistered = &terrarium.Response{Message: "Module registered successfully."}

	ModuleTableInitializationError = status.Error(codes.Unknown, "Failed to initialize table for modules.")
	ModuleGetError                 = status.Error(codes.Unknown, "Failed to check if module already exists.")
	ModuleUpdateError              = status.Error(codes.Unknown, "Failed to update module.")
	ModuleRegisterError            = status.Error(codes.Unknown, "Failed to register module.")
	ExpressionBuildError           = status.Error(codes.Unknown, "Failed to build update expression.")
	MarshalModuleError             = status.Error(codes.Unknown, "Failed to marshal module.")
)

type RegistrarService struct {
	services.UnimplementedRegistrarServer
	Db     storage.DynamoDBTableCreator
	Table  string
	Schema *dynamodb.CreateTableInput
}

type Module struct {
	Name        string `json:"name" bson:"name" dynamodbav:"name"`
	Description string `json:"description" bson:"description" dynamodbav:"description"`
	Source      string `json:"source_url" bson:"source_url" dynamodbav:"source"`
	Maturity    string `json:"maturity" bson:"maturity" dynamodbav:"maturity"`
	CreatedOn   string `json:"created_on" bson:"created_on" dynamodbav:"created_on"`
	ModifiedOn  string `json:"modified_on" bson:"modified_on" dynamodbav:"modified_on"`
}

// Registers RegistrarService with grpc server
func (s *RegistrarService) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	if err := storage.InitializeDynamoDb(s.Table, s.Schema, s.Db); err != nil {
		return ModuleTableInitializationError
	}

	services.RegisterRegistrarServer(grpcServer, s)

	return nil
}

// Register new Module in Terrarium
func (s *RegistrarService) Register(ctx context.Context, request *terrarium.RegisterModuleRequest) (*terrarium.Response, error) {
	log.Println("Registering new module.")
	name, err := attributevalue.Marshal(request.GetName())
	if err != nil {
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
			log.Println(err)
			return nil, MarshalModuleError
		}

		in := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(RegistrarTableName),
		}

		if _, err = s.Db.PutItem(ctx, in); err != nil {
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

// GetModule Retrieve module metadata
func (s *RegistrarService) GetModule(ctx context.Context, request *services.GetModuleRequest) (*services.GetModuleResponse, error) {
	filter := expression.Name("name").Equal(expression.Value(request.Name))
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
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
		log.Printf("ScanInput failed: %v", err)
		return nil, err
	}

	if response.Items == nil || len(response.Items) < 1 {
		return nil, fmt.Errorf("module not found '%v'", request.GetName())
	}
	grpcResponse := services.GetModuleResponse{}
	if moduleMetadata, err := unmarshalModule(response.Items[0]); err != nil {
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
