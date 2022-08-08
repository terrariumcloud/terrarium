package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	services "github.com/terrariumcloud/terrarium-grpc-gateway/internal/module/services"
	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"
)

func TestRegisterModule(t *testing.T) {
	t.Parallel()

	fd := &fakeDynamoDB{}

	creationService := &services.RegistrarService{
		Db: fd,
	}
	request := services.RegisterModuleRequest{
		Name:        "test",
		Description: "test desc",
		SourceUrl:   "http://test.com",
		Maturity:    terrarium.Maturity_ALPHA,
	}
	response, err := creationService.Register(context.TODO(), &request)

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
}

func TestRegisterModuleWhenPutItemReturnsError(t *testing.T) {
	t.Parallel()

	fd := &fakeDynamoDB{
		err: errors.New("test"),
	}

	creationService := &services.RegistrarService{
		Db: fd,
	}
	request := services.RegisterModuleRequest{
		Name:        "test",
		Description: "test desc",
		SourceUrl:   "http://test.com",
		Maturity:    terrarium.Maturity_ALPHA,
	}
	response, err := creationService.Register(context.TODO(), &request)

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
}

func IgnoreTestRegisterModuleE2E(t *testing.T) {
	t.Parallel()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	creationService := &services.RegistrarService{
		Db: svc,
	}
	request := services.RegisterModuleRequest{
		Name:        "test",
		Description: "test desc",
		SourceUrl:   "http://test.com",
		Maturity:    terrarium.Maturity_ALPHA,
	}
	response, _ := creationService.Register(context.TODO(), &request)

	if response != nil {
		if response.Status == terrarium.Status_OK {
			t.Log("Created.")
		} else {
			t.Error("Failed.")
		}
	}
}
