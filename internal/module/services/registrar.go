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
	"github.com/google/uuid"
)

const (
	DefaultRegistrarTableName              = "terrarium-modules"
	DefaultRegistrarServiceDefaultEndpoint = "registrar:3001"
)

var (
	RegistrarTableName       string = DefaultRegistrarTableName
	RegistrarServiceEndpoint string = DefaultRegistrarServiceDefaultEndpoint

	ModuleRegistered = &terrarium.Response{Message: "Module registered successfully."}

	ModuleTableInitializationError = status.Error(codes.Unknown, "Failed to initialize table for modules.")
	ModuleRegisterError            = status.Error(codes.Unknown, "Failed to register module.")
	MarshalModuleError             = status.Error(codes.Unknown, "Failed to marshal module.")
)

type RegistrarService struct {
	UnimplementedRegistrarServer
	Db     dynamodbiface.DynamoDBAPI
	Table  string
	Schema *dynamodb.CreateTableInput
}

type Module struct {
	ID          interface{} `json:"id" bson:"_id" dynamodbav:"_id"`
	Name        string      `json:"name" bson:"name" dynamodbav:"name"`
	Description string      `json:"description" bson:"description" dynamodbav:"description"`
	Source      string      `json:"source_url" bson:"source_url" dynamodbav:"source"`
	Maturity    string      `json:"maturity" bson:"maturity" dynamodbav:"maturity"`
	CreatedOn   string      `json:"created_on" bson:"created_on" dynamodbav:"created_on"`
}

// Registers RegistrarService with grpc server
func (s *RegistrarService) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	if err := storage.InitializeDynamoDb(s.Table, s.Schema, s.Db); err != nil {
		return ModuleTableInitializationError
	}

	RegisterRegistrarServer(grpcServer, s)

	return nil
}

// Register new Module in Terrarium
func (s *RegistrarService) Register(ctx context.Context, request *terrarium.RegisterModuleRequest) (*terrarium.Response, error) {
	log.Println("Registering new module.")

	ms := Module{
		ID:          uuid.NewString(),
		Name:        request.GetName(),
		Description: request.GetDescription(),
		Source:      request.GetSource(),
		Maturity:    request.GetMaturity().String(),
		CreatedOn:   time.Now().UTC().String(),
	}

	av, err := dynamodbattribute.MarshalMap(ms)

	if err != nil {
		log.Println(err)
		return nil, MarshalModuleError
	}

	in := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(RegistrarTableName),
	}

	if _, err = s.Db.PutItem(in); err != nil {
		log.Println(err)
		return nil, ModuleRegisterError
	}

	log.Println("New module registered.")
	return ModuleRegistered, nil
}

//

// ListModules Retrieve all published modules

func (s *VersionManagerService) ListModules(_ context.Context, request *ListModulesRequest) (*ListModulesResponse, error) {

	// return all

	expr, err := expression.NewBuilder().Build()
	if err != nil {
		log.Printf("Expression Builder failed creation: %v", err)
		return nil, err
	}

	scanQueryInputs := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(RegistrarTableName),
	}

	response, err := s.Db.Scan(scanQueryInputs)
	if err != nil {
		log.Printf("ScanInput failed: %v", err)
		return nil, err
	}

	grpcResponse := ListModulesResponse{}
	if response.Items != nil {
		for _, item := range response.Items {
			modules := Module{}
			if err3 := dynamodbattribute.UnmarshalMap(item, &modules); err3 != nil {
				log.Printf("UnmarshalMap failed: %v", err3)
				return nil, err3
			}
			grpcResponse.Modules = append(grpcResponse.Modules, modules)
		}
	}

	return &grpcResponse, nil
}

//

// GetModulesSchema returns CreateTableInput
// that can be used to create table if it does not exist
func GetModulesSchema(table string) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("_id"),
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("_id"),
				KeyType:       aws.String("HASH"),
			},
		},
		TableName:   aws.String(table),
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
	}
}
