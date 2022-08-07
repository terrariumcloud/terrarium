package services

import (
	"context"
	"time"

	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
)

const (
	DefaultVersionsTableName      = "terrarium-module-versions"
	DefaultVersionManagerEndpoint = "version_manager:3001"
)

var VersionsTableName string = DefaultVersionsTableName
var VersionManagerEndpoint string = DefaultVersionManagerEndpoint

type VersionManagerService struct {
	UnimplementedVersionManagerServer
	Db dynamodbiface.DynamoDBAPI
}

type ModuleVersion struct {
	ID          string `json:"id" bson:"_id" dynamodbav:"_id"`
	Name        string `json:"name" bson:"name" dynamodbav:"name"`
	Version     string `json:"version" bson:"version" dynamodbav:"version"`
	CreatedOn   string `json:"created_on" bson:"created_on" dynamodbav:"created_on"`
	PublishedOn string `json:"published_on" bson:"published_on" dynamodbav:"published_on"`
}

func (s *VersionManagerService) BeginVersion(ctx context.Context, request *BeginVersionRequest) (*terrarium.BeginVersionResponse, error) {

	ms := ModuleVersion{
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
		TableName: aws.String(VersionsTableName),
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

func (s *VersionManagerService) AbortVersion(ctx context.Context, request *TerminateVersionRequest) (*terrarium.TransactionStatusResponse, error) {
	if err := s.removeSessionKey(request.GetSessionKey()); err != nil {
		return SessionKeyNotRemoved, err
	}
	return VersionAborted, nil
}

func (s *VersionManagerService) PublishVersion(ctx context.Context, request *TerminateVersionRequest) (*terrarium.TransactionStatusResponse, error) {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":p": {
				N: aws.String(time.Now().UTC().String()),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"_id": {
				S: aws.String(request.GetSessionKey()),
			},
		},
		TableName:        aws.String(VersionsTableName),
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set publised_on = :p"),
	}

	_, err := s.Db.UpdateItem(input)
	if err != nil {
		return nil, err
	}

	return VersionPublished, nil
}

func (s *VersionManagerService) removeSessionKey(sessionKey string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"_id": {
				S: aws.String(sessionKey),
			},
		},
		TableName: aws.String(VersionsTableName),
	}
	_, err := s.Db.DeleteItem(input)

	if err != nil {
		return err
	}
	return nil
}
