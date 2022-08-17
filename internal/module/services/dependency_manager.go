package services

import (
	"context"
	"encoding/json"
	"errors"

	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	DefaultModuleDependenciesTableName    = "terrarium-module-dependencies"
	DefaultContainerDependenciesTableName = "terrarium-container-dependencies"
	DefaultDependencyManagerEndpoint      = "dependency_manager:3001"
)

var ModuleDependenciesTableName string = DefaultModuleDependenciesTableName
var ContainerDependenciesTableName string = DefaultContainerDependenciesTableName
var DependencyManagerEndpoint string = DefaultDependencyManagerEndpoint

type DependencyManagerService struct {
	UnimplementedDependencyManagerServer
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

// Registers Module dependencies in Terrarium
func (s *DependencyManagerService) RegisterModuleDependencies(ctx context.Context, request *terrarium.RegisterModuleDependenciesRequest) (*terrarium.TransactionStatusResponse, error) {
	dep, err := json.Marshal(request.Modules)

	if err != nil {
		return nil, err
	}

	in := &dynamodb.PutItemInput{
		TableName: aws.String(ModuleDependenciesTableName),
		Item: map[string]*dynamodb.AttributeValue{
			"_id": {
				S: aws.String(request.GetSessionKey()),
			},
			"modules": {
				B: dep,
			},
		},
	}

	if _, err = s.Db.PutItem(in); err != nil {
		return RegisterModuleDependenciesFailed, err
	}

	return ModuleDependenciesRegistered, nil
}

func (s *DependencyManagerService) RegisterContainerDependencies(ctx context.Context, request *terrarium.RegisterContainerDependenciesRequest) (*terrarium.TransactionStatusResponse, error) {
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
		return RegisterContainerDependenciesFailed, err
	}

	return ContainerDependenciesRegistered, nil
}

func (s *DependencyManagerService) RetrieveContainerDependencies(request *terrarium.RetrieveContainerDependenciesRequest, server DependencyManager_RetrieveContainerDependenciesServer) error {

	filter := expression.Name("Name").Equal(expression.Value(request.Module.Name))
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return err
	}

	sin := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(VersionsTableName),
	}

	sout, err := s.Db.Scan(sin)

	if sout.Items == nil {
		return err
	}

	moduleVersion := ModuleVersion{}

	if *sout.Count > 1 {
		return errors.New("unexpected number of results returned")
	}

	for _, i := range sout.Items {
		if err := dynamodbattribute.UnmarshalMap(i, &moduleVersion); err != nil {
			return err
		}
	}

	in := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"_id": {
				S: aws.String(moduleVersion.Version),
			},
		},
		TableName: aws.String(ContainerDependenciesTableName),
	}

	out, err := s.Db.GetItem(in)

	if out.Item == nil {
		return err
	}

	dep := ContainerDependencies{}

	if err := dynamodbattribute.UnmarshalMap(out.Item, &dep); err != nil {
		return err
	}

	res := &terrarium.ContainerDependenciesResponse{
		Origin:                   request.Module,
		ContainerImageReferences: []string{dep.Images},
	}

	if err = server.Send(res); err != nil {
		return err
	}

	return nil
}

func (s *DependencyManagerService) RetrieveModuleDependencies(request *terrarium.RetrieveModuleDependenciesRequest, server DependencyManager_RetrieveModuleDependenciesServer) error {
	return nil
}
