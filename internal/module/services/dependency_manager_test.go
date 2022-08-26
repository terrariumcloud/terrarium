package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	services "github.com/terrariumcloud/terrarium-grpc-gateway/internal/module/services"
	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"
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

func (fd *fakeDynamoDB) UpdateItem(item *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	fd.tableName = item.TableName
	output := new(dynamodb.UpdateItemOutput)
	output.Attributes = make(map[string]*dynamodb.AttributeValue)
	fd.numberOfPutItemCalls++
	return output, fd.err

}

func TestRegisterModuleDependencies(t *testing.T) {
	t.Parallel()

	fd := &fakeDynamoDB{}

	dependencyService := &services.DependencyManagerService{
		Db: fd,
	}
	dep := []*terrarium.Module{
		{
			Name:    "test",
			Version: "v1.0.0",
		},
		{
			Name:    "test2",
			Version: "v1.1.0",
		},
	}
	request := terrarium.RegisterModuleDependenciesRequest{
		Module:       &terrarium.Module{Name: "test", Version: "v1"},
		Dependencies: dep,
	}
	response, err := dependencyService.RegisterModuleDependencies(context.TODO(), &request)

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
		if *fd.tableName != services.DefaultModuleDependenciesTableName {
			t.Errorf("Expected tableName to be %s, got %s", services.DefaultModuleDependenciesTableName, *fd.tableName)
		}
	}
}

func TestRegisterModuleDependenciesWhenPutItemReturnsError(t *testing.T) {
	t.Parallel()

	fd := &fakeDynamoDB{
		err: errors.New("test"),
	}

	dependencyService := &services.DependencyManagerService{
		Db: fd,
	}
	dep := []*terrarium.Module{
		{
			Name:    "test",
			Version: "v1.0.0",
		},
		{
			Name:    "test2",
			Version: "v1.1.0",
		},
	}
	request := terrarium.RegisterModuleDependenciesRequest{
		Module:       &terrarium.Module{Name: "test", Version: "v1"},
		Dependencies: dep,
	}
	_, err := dependencyService.RegisterModuleDependencies(context.TODO(), &request)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if fd.numberOfPutItemCalls != 1 {
		t.Errorf("Expected number of calls to PutItem to be %d, got %d", 1, fd.numberOfPutItemCalls)
	}

	if fd.tableName == nil {
		t.Errorf("Expected tableName, got nil.")
	} else {
		if *fd.tableName != services.DefaultModuleDependenciesTableName {
			t.Errorf("Expected tableName to be %s, got %s", services.DefaultModuleDependenciesTableName, *fd.tableName)
		}
	}
}

func IgnoreTestRegisterModuleDependenciesE2E(t *testing.T) {
	t.Parallel()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	dependencyService := &services.DependencyManagerService{
		Db: svc,
	}
	dep := []*terrarium.Module{
		{
			Name:    "test",
			Version: "v1.0.0",
		},
		{
			Name:    "test2",
			Version: "v1.1.0",
		},
	}
	request := terrarium.RegisterModuleDependenciesRequest{
		Module:       &terrarium.Module{Name: "test", Version: "v1"},
		Dependencies: dep,
	}
	response, err := dependencyService.RegisterModuleDependencies(context.TODO(), &request)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response == nil {
		t.Errorf("Expected response, got nil.")
	}
}

func TestRegisterContainerDependencies(t *testing.T) {
	t.Parallel()

	fd := &fakeDynamoDB{}

	dependencyService := &services.DependencyManagerService{
		Db: fd,
	}
	request := terrarium.RegisterContainerDependenciesRequest{
		Module:       &terrarium.Module{Name: "test", Version: "v1"},
		Dependencies: []string{"test", "test2"},
	}
	response, err := dependencyService.RegisterContainerDependencies(context.TODO(), &request)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response == nil {
		t.Errorf("Expected response, got nil.")
	}

	if fd.numberOfPutItemCalls != 1 {
		t.Errorf("Expected number of calls to PutItem to be %d, got %d", 1, fd.numberOfPutItemCalls)
	}

	// if fd.tableName == nil {
	// 	t.Errorf("Expected tableName, got nil.")
	// } else {
	// 	if *fd.tableName != services.DefaultContainerDependenciesTableName {
	// 		t.Errorf("Expected tableName to be %s, got %s", services.DefaultContainerDependenciesTableName, *fd.tableName)
	// 	}
	// }
}

func TestRegisterContainerDependenciesWhenPutItemReturnsError(t *testing.T) {
	t.Parallel()

	fd := &fakeDynamoDB{
		err: errors.New("test"),
	}

	dependencyService := &services.DependencyManagerService{
		Db: fd,
	}
	request := terrarium.RegisterContainerDependenciesRequest{
		Module:       &terrarium.Module{Name: "test", Version: "v1"},
		Dependencies: []string{"test", "test2"},
	}
	_, err := dependencyService.RegisterContainerDependencies(context.TODO(), &request)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if fd.numberOfPutItemCalls != 1 {
		t.Errorf("Expected number of calls to PutItem to be %d, got %d", 1, fd.numberOfPutItemCalls)
	}

	// if fd.tableName == nil {
	// 	t.Errorf("Expected tableName, got nil.")
	// } else {
	// 	if *fd.tableName != services.DefaultContainerDependenciesTableName {
	// 		t.Errorf("Expected tableName to be %s, got %s", services.DefaultContainerDependenciesTableName, *fd.tableName)
	// 	}
	// }
}

func IgnoreTestRegisterContainerDependenciesE2E(t *testing.T) {
	t.Parallel()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	dependencyService := &services.DependencyManagerService{
		Db: svc,
	}

	request := terrarium.RegisterContainerDependenciesRequest{
		Module:       &terrarium.Module{Name: "test", Version: "v1"},
		Dependencies: []string{"test", "test2"},
	}
	response, err := dependencyService.RegisterContainerDependencies(context.TODO(), &request)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if response == nil {
		t.Errorf("Expected response, got nil.")
	}
}

// func TestRetrieveContainerDependencies(t *testing.T) {
// 	t.Run("It retrieves container dependencies details", func(t *testing.T) {

// 		fd := &fakeDynamoDB{}

// 		dependencyService := &DependencyService{
// 			Db: fd,
// 		}

// 		request := &terrarium.RetrieveContainerDependenciesRequest{}
// 		fds := &FakeDependencyService{}

// 		if err := dependencyService.RetrieveContainerDependencies(request, fds); err != nil {
// 			t.Errorf("Expected no error, got %v", err)
// 		}

// 		// if fd.numberOfPutItemCalls != 1 {
// 		// 	t.Errorf("Expected number of calls to PutItem to be %d, got %d", 1, fd.numberOfPutItemCalls)
// 		// }

// 		// if fd.tableName == nil {
// 		// 	t.Errorf("Expected tableName, got nil.")
// 		// } else {
// 		// 	if *fd.tableName != DefaultModuleDependenciesTableName {
// 		// 		t.Errorf("Expected tableName to be %s, got %s", DefaultModuleDependenciesTableName, *fd.tableName)
// 		// 	}
// 		// }

// 	})
// }
