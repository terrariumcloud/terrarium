package services

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"
)

// type fakeDynamoDB struct {
// 	dynamodbiface.DynamoDBAPI
// 	err                  error
// 	numberOfPutItemCalls int
// 	tableName            *string
// }

// func (fd *fakeDynamoDB) PutItem(item *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
// 	fd.tableName = item.TableName
// 	output := new(dynamodb.PutItemOutput)
// 	output.Attributes = make(map[string]*dynamodb.AttributeValue)
// 	fd.numberOfPutItemCalls++
// 	return output, fd.err
// }

func TestSetupModule(t *testing.T) {
	t.Run("It creates entry in DynamoDB", func(t *testing.T) {
		fd := &fakeDynamoDB{}

		creationService := &RegistrarService{
			Db: fd,
		}
		request := RegisterModuleRequest{
			Name:        "test",
			Description: "test desc",
			SourceUrl:   "http://test.com",
			Maturity:    terrarium.Maturity_ALPHA,
		}
		response, err := creationService.SetupModule(context.TODO(), &request)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response == nil {
			t.Errorf("Expected response, got nil.")
		} else {
			if response.Status != terrarium.Status_OK {
				t.Errorf("Expected response status %v, got %v", terrarium.Status_OK, response.Status)
			}
		}

		if fd.numberOfPutItemCalls != 1 {
			t.Errorf("Expected number of calls to PutItem to be %d, got %d", 1, fd.numberOfPutItemCalls)
		}

		if fd.tableName == nil {
			t.Errorf("Expected tableName, got nil.")
		} else {
			if *fd.tableName != "terrarium-module-stream" {
				t.Errorf("Expected tableName to be %s, got %s", "terrarium-module-stream", *fd.tableName)
			}
		}
	})
}

func TestSetupModuleWhenPutItemReturnsError(t *testing.T) {
	t.Run("It returns an error", func(t *testing.T) {
		fd := &fakeDynamoDB{
			err: errors.New("test"),
		}

		creationService := &RegistrarService{
			Db: fd,
		}
		request := RegisterModuleRequest{
			Name:        "test",
			Description: "test desc",
			SourceUrl:   "http://test.com",
			Maturity:    terrarium.Maturity_ALPHA,
		}
		response, err := creationService.SetupModule(context.TODO(), &request)

		if err == nil {
			t.Error("Expected error, got nil")
		} else {
			if response.Status != terrarium.Status_UNKNOWN_ERROR {
				t.Errorf("Expected response status %v, got %v", terrarium.Status_UNKNOWN_ERROR, response.Status)
			}
		}

		if fd.numberOfPutItemCalls != 1 {
			t.Errorf("Expected number of calls to PutItem to be %d, got %d", 1, fd.numberOfPutItemCalls)
		}

		if fd.tableName == nil {
			t.Errorf("Expected tableName, got nil.")
		} else {
			if *fd.tableName != "terrarium-module-stream" {
				t.Errorf("Expected tableName to be %s, got %s", "terrarium-module-stream", *fd.tableName)
			}
		}
	})
}

func IgnoreTestSetupModuleE2E(t *testing.T) {
	t.Run("It creates entry in DynamoDB", func(t *testing.T) {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		svc := dynamodb.New(sess)

		creationService := &RegistrarService{
			Db: svc,
		}
		request := RegisterModuleRequest{
			Name:        "test",
			Description: "test desc",
			SourceUrl:   "http://test.com",
			Maturity:    terrarium.Maturity_ALPHA,
		}
		response, _ := creationService.SetupModule(context.TODO(), &request)

		if response != nil {
			if response.Status == terrarium.Status_OK {
				t.Log("Created.")
			} else {
				t.Error("Failed.")
			}
		}
	})
}
