package services

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"
)

const (
	DefaultVersionTableName              = "terrarium-module-version"
	DefaultPublishedModulesTableName     = "terrarium-published-modules"
	DefaultVersionServiceDefaultEndpoint = "version_service:3001"
)

var SessionTableName string = DefaultVersionTableName
var PublishedModulesTableName string = DefaultPublishedModulesTableName
var SessionServiceEndpoint string = DefaultVersionServiceDefaultEndpoint

type VersionService struct {
	UnimplementedVersionManagerServer
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

func (s *VersionService) BeginVersion(ctx context.Context, request *BeginVersionRequest) (*terrarium.BeginVersionResponse, error) {

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

func (s *VersionService) AbortVersion(ctx context.Context, request *TerminateVersionRequest) (*terrarium.TransactionStatusResponse, error) {
	if err := s.removeSessionKey(request.GetSessionKey()); err != nil {
		return SessionKeyNotRemoved, err
	}
	return VersionAborted, nil
}

func (s *VersionService) PublishVersion(ctx context.Context, request *TerminateVersionRequest) (*terrarium.TransactionStatusResponse, error) {
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
		return NotFound, err
	}

	item := ModuleSession{}

	if err := dynamodbattribute.UnmarshalMap(output.Item, &item); err != nil {
		return FailedToUnmarshal, err
	}

	pm := PublishedModule{
		ID:          item.ID,
		Name:        item.Name,
		Version:     item.Version,
		PublishedOn: time.Now().UTC().String(),
	}
	av, err := dynamodbattribute.MarshalMap(pm)
	if err != nil {
		return FailedToMarshal, err
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
		return SessionKeyNotRemoved, err
	}

	return VersionPublished, nil
}

func (s *VersionService) removeSessionKey(sessionKey string) error {
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
