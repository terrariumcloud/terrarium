package main

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium"
)

type fakeDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	err                  error
	numberOfPutItemCalls int
	tableName            *string
}

func (fd *fakeDynamoDB) PutItem(item *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	fd.tableName = item.TableName
	output := new(dynamodb.PutItemOutput)
	output.Attributes = make(map[string]*dynamodb.AttributeValue)
	fd.numberOfPutItemCalls++
	return output, fd.err
}

func TestRegisterModuleDependencies(t *testing.T) {
	t.Run("It creates entry in DynamoDB", func(t *testing.T) {
		fd := &fakeDynamoDB{}

		dependencyService := &DependencyService{
			db: fd,
		}
		modules := []*terrarium.VersionedModule{
			&terrarium.VersionedModule{
				Name:    "test",
				Version: "v1.0.0",
			},
		}
		request := terrarium.RegisterModuleDependenciesRequest{
			SessionKey: "123",
			Modules:    modules,
		}
		_, err := dependencyService.RegisterModuleDependencies(context.TODO(), &request)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// if response == nil {
		// 	t.Errorf("Expected response, got nil.")
		// } else {
		// 	if response.Status != terrarium.Status_OK {
		// 		t.Errorf("Expected response status %v, got %v", terrarium.Status_OK, response.Status)
		// 	}
		// }

		// if fd.numberOfPutItemCalls != 1 {
		// 	t.Errorf("Expected number of calls to PutItem to be %d, got %d", 1, fd.numberOfPutItemCalls)
		// }

		// if fd.tableName == nil {
		// 	t.Errorf("Expected tableName, got nil.")
		// } else {
		// 	if *fd.tableName != "terrarium-module-stream" {
		// 		t.Errorf("Expected tableName to be %s, got %s", "terrarium-module-stream", *fd.tableName)
		// 	}
		// }
	})
}
