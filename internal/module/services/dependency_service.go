package services

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"
)

const (
	DefaultModuleDependenciesTableName      = "terrarium-module-dependencies"
	DefaultContainerDependenciesTableName   = "terrarium-container-dependencies"
	DefaultDependencyServiceDefaultEndpoint = "dependency_service:3001"
)

var ModuleDependenciesTableName string = DefaultModuleDependenciesTableName
var ContainerDependenciesTableName string = DefaultContainerDependenciesTableName
var DependencyServiceEndpoint string = DefaultDependencyServiceDefaultEndpoint

type DependencyService struct {
	UnimplementedDependencyResolverServer
	Db dynamodbiface.DynamoDBAPI
}

type ModuleDependencies struct {
	ID      interface{} `json:"id" bson:"_id" dynamodbav:"_id"`
	Modules string      `json:"modules" bson:"modules" dynamodbav:"modules"`
}

type ContainerDependencies struct {
	ID     interface{} `json:"id" bson:"_id" dynamodbav:"_id"`
	Images string      `json:"images" bson:"images" dynamodbav:"images"`
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

func (s *DependencyService) RetrieveContainerDependencies(request *terrarium.RetrieveContainerDependenciesRequest, server DependencyResolver_RetrieveContainerDependenciesServer) error {
	getInput := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"_id": {
				S: aws.String(request.GetApiKey()),
			},
		},
		TableName: aws.String(ContainerDependenciesTableName),
	}
	output, err := s.Db.GetItem(getInput)
	if output.Item == nil {
		return err
	}

	item := ContainerDependencies{}

	if err := dynamodbattribute.UnmarshalMap(output.Item, &item); err != nil {
		return err
	}

	getInput = &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"_id": {
				S: aws.String(request.GetApiKey()),
			},
		},
		TableName: aws.String(SessionTableName),
	}
	output, err = s.Db.GetItem(getInput)
	if output.Item == nil {
		return err
	}

	item2 := &ModuleSession{}

	if err := dynamodbattribute.UnmarshalMap(output.Item, &item2); err != nil {
		return err
	}
	origin := &terrarium.VersionedModule{
		Name:    item2.Name,
		Version: item2.Version,
	}

	dependencies := &terrarium.ContainerDependenciesResponse{
		Origin:                   origin,
		ContainerImageReferences: []string{item.Images},
	}

	if err = server.Send(dependencies); err != nil {
		return err
	}

	return nil
}

func (s *DependencyService) RetrieveModuleDependencies(request *terrarium.RetrieveModuleDependenciesRequest, server DependencyResolver_RetrieveModuleDependenciesServer) error {
	return nil
}
