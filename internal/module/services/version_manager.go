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

// Creates new Module Version with Version Manager service
func (s *VersionManagerService) BeginVersion(ctx context.Context, request *BeginVersionRequest) (*terrarium.BeginVersionResponse, error) {
	mv := ModuleVersion{
		ID:        uuid.NewString(),
		Name:      request.GetModule().GetName(),
		Version:   request.GetModule().GetVersion(),
		CreatedOn: time.Now().UTC().String(),
	}

	av, err := dynamodbattribute.MarshalMap(mv)

	if err != nil {
		return nil, err
	}

	in := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(VersionsTableName),
	}

	if _, err = s.Db.PutItem(in); err != nil {
		return nil, err
	}

	response := &terrarium.BeginVersionResponse{
		SessionKey: mv.ID,
	}

	return response, nil
}

// Removes Module Version with Version Manager service
func (s *VersionManagerService) AbortVersion(ctx context.Context, request *TerminateVersionRequest) (*terrarium.TransactionStatusResponse, error) {
	if err := s.removeSessionKey(request.GetSessionKey()); err != nil {
		return SessionKeyNotRemoved, err
	}

	return VersionAborted, nil
}

// Updates Module Version to published with Verison Manager service
func (s *VersionManagerService) PublishVersion(ctx context.Context, request *TerminateVersionRequest) (*terrarium.TransactionStatusResponse, error) {
	in := &dynamodb.UpdateItemInput{
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

	if _, err := s.Db.UpdateItem(in); err != nil {
		return nil, err
	}

	return VersionPublished, nil
}

func (s *VersionManagerService) removeSessionKey(sessionKey string) error {
	in := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"_id": {
				S: aws.String(sessionKey),
			},
		},
		TableName: aws.String(VersionsTableName),
	}

	if _, err := s.Db.DeleteItem(in); err != nil {
		return err
	}

	return nil
}
