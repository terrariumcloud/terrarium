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

type SessionService struct {
	services.UnimplementedSessionManagerServer
	db dynamodbiface.DynamoDBAPI
}

type ModuleSession struct {
	ID        string `json:"id" bson:"_id" dynamodbav:"_id"`
	Name      string `json:"name" bson:"name" dynamodbav:"name"`
	Version   string `json:"version" bson:"version" dynamodbav:"version"`
	CreatedOn string `json:"created_on" bson:"created_on" dynamodbav:"created_on"`
}

func (s *SessionService) BeginVersion(ctx context.Context, request *services.BeginVersionRequest) (*terrarium.BeginVersionResponse, error) {

	ms := ModuleSession{
		ID:        uuid.NewString(),
		Name:      request.GetModule().GetName(),
		Version:   request.GetModule().GetVersion(),
		CreatedOn: time.Now().UTC().String(),
	}
	av, err := dynamodbattribute.MarshalMap(ms)
	if err != nil {
		return nil, err
	}

	tableName := "terrarium-module-session"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: &tableName,
	}

	_, err = s.db.PutItem(input)

	if err != nil {
		return nil, err
	}

	response := terrarium.BeginVersionResponse{
		SessionKey: ms.ID,
	}

	return &response, nil
}

func (s *SessionService) AbortVersion(ctx context.Context, request *services.TerminateVersionRequest) (*terrarium.TransactionStatusResponse, error) {
	// call abort on storage and dependency service with session key
	// remove session key
	return nil, nil
}

func (s *SessionService) PublishVersion(ctx context.Context, request *services.TerminateVersionRequest) (*terrarium.TransactionStatusResponse, error) {
	// put item module version table
	// remove session key
	return nil, nil
}
