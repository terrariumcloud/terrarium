package creation

import (
	"context"
	"time"

	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/services"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
	"github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium"
)

var TableName string

const (
	DefaultTableName = "terrarium-module-stream"
)

type CreationService struct {
	services.UnimplementedCreatorServer
	Db dynamodbiface.DynamoDBAPI
}

type ModuleStream struct {
	ID          interface{} `json:"id" bson:"_id" dynamodbav:"_id"`
	Name        string      `json:"name" bson:"name" dynamodbav:"name"`
	Description string      `json:"description" bson:"description" dynamodbav:"description"`
	SourceUrl   string      `json:"source_url" bson:"source_url" dynamodbav:"source_url"`
	Maturity    string      `json:"maturity" bson:"maturity" dynamodbav:"maturity"`
	CreatedOn   string      `json:"created_on" bson:"created_on" dynamodbav:"created_on"`
}

func (s *CreationService) SetupModule(ctx context.Context, request *services.SetupModuleRequest) (*services.SetupModuleResponse, error) {

	ms := ModuleStream{
		ID:          uuid.NewString(),
		Name:        request.GetName(),
		Description: request.GetDescription(),
		SourceUrl:   request.GetSourceUrl(),
		Maturity:    request.GetMaturity().String(),
		CreatedOn:   time.Now().UTC().String(),
	}
	av, err := dynamodbattribute.MarshalMap(ms)
	if err != nil {
		return Error("Failed to marshal module stream."), err
	}

	input := &dynamodb.PutItemInput{
		Item:                av,
		TableName:           aws.String(TableName),
		ConditionExpression: aws.String("attribute_not_exists(source_url)"),
	}
	_, err = s.Db.PutItem(input)

	if err != nil {
		return Error("Failed to setup module."), err
	}

	return Ok("Module setup successfully."), nil
}

func Error(message string) *services.SetupModuleResponse {
	return &services.SetupModuleResponse{
		Status:        terrarium.Status_UNKNOWN_ERROR,
		StatusMessage: message,
	}
}

func Ok(message string) *services.SetupModuleResponse {
	return &services.SetupModuleResponse{
		Status:        terrarium.Status_OK,
		StatusMessage: message,
	}
}
