package services_test

import (
	"context"
	"errors"
	"testing"
	
	services "github.com/terrariumcloud/terrarium-grpc-gateway/internal/module/services"
	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"
)

func TestRegisterModule(t *testing.T) {
	t.Parallel()

	fd := &fakeDynamoDB{}

	creationService := &services.RegistrarService{
		Db: fd,
	}
	request := terrarium.RegisterModuleRequest{
		Name:        "test",
		Description: "test desc",
		Source:      "http://test.com",
		Maturity:    terrarium.RegisterModuleRequest_ALPHA,
	}
	response, err := creationService.Register(context.TODO(), &request)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response == nil {
		t.Errorf("Expected response, got nil.")
	}

	if fd.numberOfPutItemCalls != 1 {
		t.Errorf("Expected number of calls to PutItem to be %d, got %d", 1, fd.numberOfPutItemCalls)
	}

	if fd.tableName == nil {
		t.Errorf("Expected tableName, got nil.")
	} else {
		if *fd.tableName != services.RegistrarTableName {
			t.Errorf("Expected tableName to be %s, got %s", services.RegistrarTableName, *fd.tableName)
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
	request := terrarium.RegisterModuleRequest{
		Name:        "test",
		Description: "test desc",
		Source:      "http://test.com",
		Maturity:    terrarium.RegisterModuleRequest_ALPHA,
	}
	_, err := creationService.Register(context.TODO(), &request)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if fd.numberOfPutItemCalls != 1 {
		t.Errorf("Expected number of calls to PutItem to be %d, got %d", 1, fd.numberOfPutItemCalls)
	}

	if fd.tableName == nil {
		t.Errorf("Expected tableName, got nil.")
	} else {
		if *fd.tableName != services.RegistrarTableName {
			t.Errorf("Expected tableName to be %s, got %s", services.RegistrarTableName, *fd.tableName)
		}
	}
}

// func IgnoreTestRegisterModuleE2E(t *testing.T) {
// 	t.Parallel()

// 	sess := session.Must(session.NewSessionWithOptions(session.Options{
// 		SharedConfigState: session.SharedConfigEnable,
// 	}))

// 	svc := dynamodb.New(sess)

// 	creationService := &services.RegistrarService{
// 		Db: svc,
// 	}
// 	request := terrarium.RegisterModuleRequest{
// 		Name:        "test",
// 		Description: "test desc",
// 		Source:      "http://test.com",
// 		Maturity:    terrarium.RegisterModuleRequest_ALPHA,
// 	}
// 	response, _ := creationService.Register(context.TODO(), &request)

// 	if response != nil {
// 	}
// }
