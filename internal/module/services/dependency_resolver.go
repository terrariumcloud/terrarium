package services

import (
	"context"
	"encoding/json"

	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const (
	DefaultModuleDependenciesTableName      = "terrarium-module-dependencies"
	DefaultDependencyServiceDefaultEndpoint = "dependency_resolver:3001"
)

var ModuleDependenciesTableName string = DefaultModuleDependenciesTableName
var DependencyServiceEndpoint string = DefaultDependencyServiceDefaultEndpoint

type DependencyResolverService struct {
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

func (s *DependencyResolverService) RegisterModuleDependencies(ctx context.Context, request *terrarium.RegisterModuleDependenciesRequest) (*terrarium.TransactionStatusResponse, error) {
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
		return RegisterModuleDependenciesFailed, err
	}

	return ModuleDependenciesRegistered, nil
}

func (s *DependencyResolverService) RegisterContainerDependencies(ctx context.Context, request *terrarium.RegisterContainerDependenciesRequest) (*terrarium.TransactionStatusResponse, error) {
	// img, err := json.Marshal(request.ContainerImageReferences)
	// if err != nil {
	// 	return nil, err
	// }

	// input := &dynamodb.PutItemInput{
	// 	TableName: aws.String(ContainerDependenciesTableName),
	// 	Item: map[string]*dynamodb.AttributeValue{
	// 		"_id": {
	// 			S: aws.String(request.GetSessionKey()),
	// 		},
	// 		"images": {
	// 			B: img,
	// 		},
	// 	},
	// }

	// _, err = s.Db.PutItem(input)

	// if err != nil {
	// 	return RegisterContainerDependenciesFailed, err
	// }

	return ContainerDependenciesRegistered, nil
}

func (s *DependencyResolverService) RetrieveContainerDependencies(request *terrarium.RetrieveContainerDependenciesRequest, server DependencyResolver_RetrieveContainerDependenciesServer) error {
	// getInput := &dynamodb.GetItemInput{
	// 	Key: map[string]*dynamodb.AttributeValue{
	// 		"_id": {
	// 			S: aws.String(request.GetApiKey()),
	// 		},
	// 	},
	// 	TableName: aws.String(ContainerDependenciesTableName),
	// }
	// output, err := s.Db.GetItem(getInput)
	// if output.Item == nil {
	// 	return err
	// }

	// item := ContainerDependencies{}

	// if err := dynamodbattribute.UnmarshalMap(output.Item, &item); err != nil {
	// 	return err
	// }

	// getInput = &dynamodb.GetItemInput{
	// 	Key: map[string]*dynamodb.AttributeValue{
	// 		"_id": {
	// 			S: aws.String(request.GetApiKey()),
	// 		},
	// 	},
	// 	TableName: aws.String(VersionsTableName),
	// }
	// output, err = s.Db.GetItem(getInput)
	// if output.Item == nil {
	// 	return err
	// }

	// item2 := &ModuleSession{}

	// if err := dynamodbattribute.UnmarshalMap(output.Item, &item2); err != nil {
	// 	return err
	// }
	// origin := &terrarium.VersionedModule{
	// 	Name:    item2.Name,
	// 	Version: item2.Version,
	// }

	// dependencies := &terrarium.ContainerDependenciesResponse{
	// 	Origin:                   origin,
	// 	ContainerImageReferences: []string{item.Images},
	// }

	// if err = server.Send(dependencies); err != nil {
	// 	return err
	// }

	return nil
}

func (s *DependencyResolverService) RetrieveModuleDependencies(request *terrarium.RetrieveModuleDependenciesRequest, server DependencyResolver_RetrieveModuleDependenciesServer) error {
	return nil
}
