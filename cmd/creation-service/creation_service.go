package main

import (
	"context"
	"terrarium-grpc-gateway/internal/services"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
	"github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium"
)

type CreationService struct {
	services.UnimplementedCreatorServer
	db dynamodbiface.DynamoDBAPI
}

type ModuleStream struct {
	ID          interface{} `json:"id" bson:"_id" dynamodbav:"_id"`
	Name        string      `json:"name" bson:"name" dynamodbav:"name"`
	Description string      `json:"description" bson:"description" dynamodbav:"description"`
	SourceUrl   string      `json:"source_url" bson:"source_url" dynamodbav:"source_url"`
	CreatedOn   string      `json:"created_on" bson:"created_on" dynamodbav:"created_on"`
}

func (s *CreationService) SetupModule(ctx context.Context, request *services.SetupModuleRequest) (*services.SetupModuleResponse, error) {

	ms := ModuleStream{
		ID:          uuid.NewString(),
		Name:        request.GetName(),
		Description: request.GetDescription(),
		SourceUrl:   request.GetSourceUrl(),
		CreatedOn:   time.Now().UTC().String(),
	}
	av, err := dynamodbattribute.MarshalMap(ms)
	if err != nil {
		response := services.SetupModuleResponse{
			Status:        terrarium.Status_UNKNOWN_ERROR,
			StatusMessage: "Something went wrong.",
		}
		return &response, nil
	}

	tableName := "terrarium-module-stream"
	condition := "attribute_not_exists(SourceUrl)"

	input := &dynamodb.PutItemInput{
		Item:                av,
		TableName:           &tableName,
		ConditionExpression: &condition,
	}

	_, err = s.db.PutItem(input)

	if err != nil {
		response := services.SetupModuleResponse{
			Status:        terrarium.Status_UNKNOWN_ERROR,
			StatusMessage: "Something went wrong.",
		}
		return &response, nil
	}

	response := services.SetupModuleResponse{
		Status:        terrarium.Status_OK,
		StatusMessage: "All is good",
	}
	return &response, nil
}
