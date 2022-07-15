package dependency

import (
	"context"
	"encoding/json"
	"terrarium-grpc-gateway/internal/services"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium"
)

var ModuleDependenciesTableName string
var ContainerDependenciesTableName string

const (
	DefaultModuleDependenciesTableName    = "terrarium-module-dependencies"
	DefaultContainerDependenciesTableName = "terrarium-container-dependencies"
)

type DependencyService struct {
	services.UnimplementedDependencyResolverServer
	Db dynamodbiface.DynamoDBAPI
}

type Dependencies struct {
	ID      interface{} `json:"id" bson:"_id" dynamodbav:"_id"`
	Modules string      `json:"modules" bson:"modules" dynamodbav:"modules"`
}

func (s *DependencyService) RegisterModuleDependencies(ctx context.Context, request *terrarium.RegisterModuleDependenciesRequest) (*terrarium.TransactionStatusResponse, error) {
	d, err := json.Marshal(request.Modules)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(ModuleDependenciesTableName),
		Item: map[string]*dynamodb.AttributeValue{
			"_id": {
				S: aws.String(request.GetSessionKey()),
			},
			"modules": {
				B: d,
			},
		},
	}

	_, err = s.Db.PutItem(input)
	if err != nil {
		return Error("Failed to register module dependencies."), err
	}

	return Ok("Module dependencies successfully registered."), nil
}

func (s *DependencyService) RegisterContainerDependencies(ctx context.Context, request *terrarium.RegisterContainerDependenciesRequest) (*terrarium.TransactionStatusResponse, error) {
	img, err := json.Marshal(request.ContainerImageReferences)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(ContainerDependenciesTableName),
		Item: map[string]*dynamodb.AttributeValue{
			"_id": {
				S: aws.String(request.GetSessionKey()),
			},
			"images": {
				B: img,
			},
		},
	}

	_, err = s.Db.PutItem(input)

	if err != nil {
		return Error("Failed to register container dependencies."), err
	}

	return Ok("Container dependencies successfully registered."), nil
}

func (s *DependencyService) RetrieveContainerDependencies(request *terrarium.RetrieveContainerDependenciesRequest, server services.DependencyResolver_RetrieveContainerDependenciesServer) error {
	return nil
}

func (s *DependencyService) RetrieveModuleDependencies(request *terrarium.RetrieveModuleDependenciesRequest, server services.DependencyResolver_RetrieveModuleDependenciesServer) error {
	return nil
}

func Error(message string) *terrarium.TransactionStatusResponse {
	return &terrarium.TransactionStatusResponse{
		Status:        terrarium.Status_UNKNOWN_ERROR,
		StatusMessage: message,
	}
}

func Ok(message string) *terrarium.TransactionStatusResponse {
	return &terrarium.TransactionStatusResponse{
		Status:        terrarium.Status_OK,
		StatusMessage: message,
	}
}
