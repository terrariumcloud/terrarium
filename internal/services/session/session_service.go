package session

import (
	"context"
	"fmt"
	"terrarium-grpc-gateway/internal/services"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
	"github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium"
)

var SessionTableName string
var PublishedModulesTableName string

const (
	DefaultSessionTableName          = "terrarium-module-session"
	DefaultPublishedModulesTableName = "terrarium-published-modules"
)

type SessionService struct {
	services.UnimplementedSessionManagerServer
	Db dynamodbiface.DynamoDBAPI
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

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(SessionTableName),
	}

	_, err = s.Db.PutItem(input)

	if err != nil {
		return nil, err
	}

	response := terrarium.BeginVersionResponse{
		SessionKey: ms.ID,
	}

	return &response, nil
}

func (s *SessionService) AbortVersion(ctx context.Context, request *services.TerminateVersionRequest) (*terrarium.TransactionStatusResponse, error) {
	if err := s.removeSessionKey(request.GetSessionKey()); err != nil {
		return Error("Failed to remove session key."), err
	}
	return Ok("Version aborted."), nil
}

func (s *SessionService) PublishVersion(ctx context.Context, request *services.TerminateVersionRequest) (*terrarium.TransactionStatusResponse, error) {
	getInput := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"_id": {
				S: aws.String(request.GetSessionKey()),
			},
		},
		TableName: aws.String(SessionTableName),
	}
	output, err := s.Db.GetItem(getInput)
	if output.Item == nil {
		return Error(fmt.Sprintf("Could not find '%s'", request.GetSessionKey())), err
	}

	item := ModuleSession{}

	if err := dynamodbattribute.UnmarshalMap(output.Item, &item); err != nil {
		return Error("Failed to unmarshal record."), err
	}

	pm := PublishedModule{
		ID:          item.ID,
		Name:        item.Name,
		Version:     item.Version,
		PublishedOn: time.Now().UTC().String(),
	}
	av, err := dynamodbattribute.MarshalMap(pm)
	if err != nil {
		return Error("Failed to marshal record."), err
	}

	putInput := &dynamodb.PutItemInput{
		Item:      av,
		TableName: &PublishedModulesTableName,
	}

	_, err = s.Db.PutItem(putInput)

	if err != nil {
		return nil, err
	}

	if err := s.removeSessionKey(request.GetSessionKey()); err != nil {
		return Error("Failed to remove session key."), err
	}

	return Ok("Version published"), nil
}

func (s *SessionService) removeSessionKey(sessionKey string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"_id": {
				S: aws.String(sessionKey),
			},
		},
		TableName: aws.String(SessionTableName),
	}
	_, err := s.Db.DeleteItem(input)

	if err != nil {
		return err
	}
	return nil
}

func Ok(message string) *terrarium.TransactionStatusResponse {
	return &terrarium.TransactionStatusResponse{
		Status:        terrarium.Status_OK,
		StatusMessage: message,
	}
}

func Error(message string) *terrarium.TransactionStatusResponse {
	return &terrarium.TransactionStatusResponse{
		Status:        terrarium.Status_UNKNOWN_ERROR,
		StatusMessage: message,
	}
}
