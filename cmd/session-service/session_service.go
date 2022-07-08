package main

import (
	"context"
	"terrarium-grpc-gateway/internal/services"
	"time"

	"github.com/aws/aws-sdk-go/aws"
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

type PublishedModule struct {
	ID          string `json:"id" bson:"_id" dynamodbav:"_id"`
	Name        string `json:"name" bson:"name" dynamodbav:"name"`
	Version     string `json:"version" bson:"version" dynamodbav:"version"`
	PublishedOn string `json:"published_on" bson:"published_on" dynamodbav:"published_on"`
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
		TableName: aws.String(tableName),
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
	err := s.removeSessionKey(request.GetSessionKey())
	if err != nil {
		return &terrarium.TransactionStatusResponse{
			Status:        terrarium.Status_UNKNOWN_ERROR,
			StatusMessage: "Failed to remove session key.",
		}, err
	}

	return &terrarium.TransactionStatusResponse{
		Status:        terrarium.Status_OK,
		StatusMessage: "Version aborted.",
	}, nil
}

func (s *SessionService) PublishVersion(ctx context.Context, request *services.TerminateVersionRequest) (*terrarium.TransactionStatusResponse, error) {

	tableName := "terrarium-module-session"
	getInput := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"_id": {
				S: aws.String(request.GetSessionKey()),
			},
		},
		TableName: aws.String(tableName),
	}
	output, err := s.db.GetItem(getInput)
	if output.Item == nil {
		return &terrarium.TransactionStatusResponse{
			Status:        terrarium.Status_UNKNOWN_ERROR,
			StatusMessage: "Could not find '" + request.GetSessionKey() + "'",
		}, err
	}

	item := ModuleSession{}

	err = dynamodbattribute.UnmarshalMap(output.Item, &item)
	if err != nil {
		return &terrarium.TransactionStatusResponse{
			Status:        terrarium.Status_UNKNOWN_ERROR,
			StatusMessage: "Failed to unmarshal record.",
		}, err
	}

	pm := PublishedModule{
		ID:          item.ID,
		Name:        item.Name,
		Version:     item.Version,
		PublishedOn: time.Now().UTC().String(),
	}
	av, err := dynamodbattribute.MarshalMap(pm)
	if err != nil {
		return &terrarium.TransactionStatusResponse{
			Status:        terrarium.Status_UNKNOWN_ERROR,
			StatusMessage: "Failed to marshal record.",
		}, err
	}

	tableName = "terrarium-published-modules"

	putInput := &dynamodb.PutItemInput{
		Item:      av,
		TableName: &tableName,
	}

	_, err = s.db.PutItem(putInput)

	if err != nil {
		return nil, err
	}

	err = s.removeSessionKey(request.GetSessionKey())
	if err != nil {
		return &terrarium.TransactionStatusResponse{
			Status:        terrarium.Status_UNKNOWN_ERROR,
			StatusMessage: "Failed to remove session key.",
		}, err
	}

	return &terrarium.TransactionStatusResponse{
		Status:        terrarium.Status_OK,
		StatusMessage: "Version published",
	}, nil
}

func (s *SessionService) removeSessionKey(sessionKey string) error {
	tableName := "terrarium-module-session"
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"_id": {
				S: aws.String(sessionKey),
			},
		},
		TableName: aws.String(tableName),
	}
	_, err := s.db.DeleteItem(input)

	if err != nil {
		return err
	}
	return nil
}
