package services_test

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/terrariumcloud/terrarium/internal/mocks"
	"github.com/terrariumcloud/terrarium/internal/module/services"
	terrarium "github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"google.golang.org/grpc"
)

var registerContainerDependenciesTestData = terrarium.RegisterContainerDependenciesRequest{
	Module: &terrarium.Module{Name: "test", Version: "v1"},
	Images: map[string]*terrarium.ContainerImageDetails{
		"grafana": {
			Tag:       "0.1.1",
			Namespace: "cie",
			Images: []*terrarium.ContainerImageRef{
				{
					Arch:  "amd64",
					Image: "random.server.com/my-grafana-image-for-linux-amd64:tag23",
				},
			},
		},
		"kubescaler": {
			Tag:       "0.2.2",
			Namespace: "cie",
			Images: []*terrarium.ContainerImageRef{
				{
					Arch:  "amd64",
					Image: "random.server.com/my-kubescaler-image-for-linux-amd64:tag25",
				},
				{
					Arch:  "arm64",
					Image: "random.server.com/my-kubescaler-image-for-linux-arm64:graviton2",
				},
			},
		},
	},
}

type MockRetrieveContainerDependenciesServer struct {
	grpc.ServerStream
	SendInvocations int
	Responses       []*terrarium.ContainerDependenciesResponseV2
	Err             error
}

func (srv *MockRetrieveContainerDependenciesServer) Send(res *terrarium.ContainerDependenciesResponseV2) error {
	srv.SendInvocations++
	srv.Responses = append(srv.Responses, res)
	return srv.Err
}

type MockRetrieveModuleDependenciesServer struct {
	grpc.ServerStream
	SendInvocations int
	Response        *terrarium.ModuleDependenciesResponse
	Err             error
}

func (srv *MockRetrieveModuleDependenciesServer) Send(res *terrarium.ModuleDependenciesResponse) error {
	srv.SendInvocations++
	srv.Response = res
	return srv.Err
}

type MockGetDependenciesResponse struct {
	Dependencies []*terrarium.Module
	Err          error
}

func (m *MockGetDependenciesResponse) GetDependencies(request *terrarium.Module) (deps []*terrarium.Module, err error) {
	m.Dependencies = deps
	m.Err = err
	return m.Dependencies, err
}

// Test_RegisterDependencyManagerWithServer checks:
// - if there was no error with table init
// - if error is returned when Table initialization fails
func Test_RegisterDependencyManagerWithServer(t *testing.T) {
	t.Parallel()
	// The dependencies are to be stored in two tables at this stage:
	// - Module dependencies
	// - Container dependencies
	expectedDescribeTableInvocations := 2
	expectedCreateTableInvocations := 0

	t.Run("when there is no error with table init", func(t *testing.T) {
		db := &mocks.MockDynamoDB{}

		dms := &services.DependencyManagerService{Db: db}

		s := grpc.NewServer(*new([]grpc.ServerOption)...)

		err := dms.RegisterWithServer(s)

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}

		if db.DescribeTableInvocations != expectedDescribeTableInvocations {
			t.Errorf("Expected %d call to DescribeTable, got %v.", expectedDescribeTableInvocations, db.DescribeTableInvocations)
		}

		if db.CreateTableInvocations != expectedCreateTableInvocations {
			t.Errorf("Expected no calls to CreateTable, got %v.", db.CreateTableInvocations)
		}
	})

	t.Run("when Table initialization fails", func(t *testing.T) {
		// The dependencies are to be stored in two tables at this stage:
		// - Module dependencies will fail...
		// - Container dependencies
		expectedDescribeTableInvocations := 1
		expectedCreateTableInvocations := 0
		expectedError := services.ModuleDependenciesTableInitializationError

		db := &mocks.MockDynamoDB{DescribeTableErrors: []error{errors.New("some error")}}

		dms := &services.DependencyManagerService{Db: db}

		s := grpc.NewServer(*new([]grpc.ServerOption)...)

		err := dms.RegisterWithServer(s)

		if err != expectedError {
			t.Errorf("Expected '%s', got '%s'.", expectedError, err)
		}

		if db.DescribeTableInvocations != expectedDescribeTableInvocations {
			t.Errorf("Expected %d call to DescribeTable, got %v.", expectedDescribeTableInvocations, db.DescribeTableInvocations)
		}

		if db.CreateTableInvocations != expectedCreateTableInvocations {
			t.Errorf("Expected %d calls to CreateTable, got %v.", expectedCreateTableInvocations, db.CreateTableInvocations)
		}
	})
}

// This test checks if error is returned when ModuleTable initialization fails
func TestRegisterWithServerWhenContainerDependenciesTableInitializationErrors(t *testing.T) {
	t.Parallel()
	// The dependencies are to be stored in two tables at this stage:
	// - Module dependencies
	// - Container dependencies will fail...
	expectedError := services.ContainerDependenciesTableInitializationError
	expectedDescribeTableInvocations := 2
	expectedCreateTableInvocations := 0

	db := &mocks.MockDynamoDB{DescribeTableErrors: []error{nil, errors.New("some error")}}

	dms := &services.DependencyManagerService{Db: db}

	s := grpc.NewServer(*new([]grpc.ServerOption)...)

	err := dms.RegisterWithServer(s)

	if err != expectedError {
		t.Errorf("Expected '%s', got '%s'.", expectedError, err)
	}

	if db.DescribeTableInvocations != expectedDescribeTableInvocations {
		t.Errorf("Expected %d call to DescribeTable, got %d.", expectedDescribeTableInvocations, db.DescribeTableInvocations)
	}

	if db.CreateTableInvocations != expectedCreateTableInvocations {
		t.Errorf("Expected %d calls to CreateTable, got %v.", expectedCreateTableInvocations, db.CreateTableInvocations)
	}
}

// Test_RegisterModuleDependencies checks:
// - if correct response is returned when module dependencies are registered
// - if error is returned when PutItem fails
func Test_RegisterModuleDependencies(t *testing.T) {
	t.Parallel()

	t.Run("when module dependencies are registered", func(t *testing.T) {
		db := &mocks.MockDynamoDB{}

		svc := &services.DependencyManagerService{Db: db}

		req := &terrarium.RegisterModuleDependenciesRequest{
			Module: &terrarium.Module{Name: "test", Version: "v1"},
			Dependencies: []*terrarium.Module{
				{Name: "test", Version: "v1.0.0"},
				{Name: "test2", Version: "v1.1.0"},
			},
		}

		res, err := svc.RegisterModuleDependencies(context.TODO(), req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if db.PutItemInvocations != 1 {
			t.Errorf("Expected 1 call to PutItem, got %v", db.PutItemInvocations)
		}

		if db.TableName != services.ModuleDependenciesTableName {
			t.Errorf("Expected tableName to be %v, got %v.", services.ModuleDependenciesTableName, db.TableName)
		}

		if res != services.ModuleDependenciesRegistered {
			t.Errorf("Expected %v, got %v.", services.ModuleDependenciesRegistered, res)
		}
	})

	// TODO: Test for MarshalModuleDependenciesError

	t.Run("when PutItem fails", func(t *testing.T) {
		db := &mocks.MockDynamoDB{PutItemError: errors.New("some error")}

		svc := &services.DependencyManagerService{Db: db}

		req := &terrarium.RegisterModuleDependenciesRequest{
			Module: &terrarium.Module{Name: "test", Version: "v1"},
			Dependencies: []*terrarium.Module{
				{Name: "test", Version: "v1.0.0"},
				{Name: "test2", Version: "v1.1.0"},
			},
		}

		res, err := svc.RegisterModuleDependencies(context.TODO(), req)

		if res != nil {
			t.Errorf("Expected no response, got %v", err)
		}

		if db.PutItemInvocations != 1 {
			t.Errorf("Expected 1 call to PutItem, got %v", db.PutItemInvocations)
		}

		if db.TableName != services.ModuleDependenciesTableName {
			t.Errorf("Expected tableName to be %v, got %v.", services.ModuleDependenciesTableName, db.TableName)
		}

		if err != services.RegisterDependenciesError {
			t.Errorf("Expected %v, got %v.", services.RegisterDependenciesError, err)
		}
	})
}

// Test_RegisterContainerDependencies checks:
// - if correct response is returned when container dependencies are registered
// - if error is returned when UpdateItem fails
func Test_RegisterContainerDependencies(t *testing.T) {
	t.Parallel()

	var expectedPutItemInvocations = 1
	var expectedTableName = services.ContainerDependenciesTableName
	var expectedResponse = services.ContainerDependenciesRegistered
	var expectedError error = nil

	db := &mocks.MockDynamoDB{}
	svc := &services.DependencyManagerService{Db: db}
	req := &registerContainerDependenciesTestData
	res, err := svc.RegisterContainerDependencies(context.TODO(), req)

	if err != expectedError {
		t.Errorf("Expected %v, got %v.", expectedError, err)
	}

	if res != expectedResponse {
		t.Errorf("Expected %v, got %v.", expectedResponse, res)
	}

	if db.PutItemInvocations != expectedPutItemInvocations {
		t.Errorf("Expected %d call to PutItem, got %d", expectedPutItemInvocations, db.UpdateItemInvocations)
	}

	if db.TableName != expectedTableName {
		t.Errorf("Expected tableName to be %v, got %v.", expectedTableName, db.TableName)
	}
}

// This test checks if error is returned when UpdateItem fails
func TestRegisterContainerDependenciesWhenPutItemErrors(t *testing.T) {
	t.Parallel()

	var expectedPutItemInvocations = 1
	var expectedTableName = services.ContainerDependenciesTableName
	var expectedResponse *terrarium.Response = nil
	var expectedError = services.RegisterDependenciesError

	db := &mocks.MockDynamoDB{PutItemError: errors.New("some error")}
	svc := &services.DependencyManagerService{Db: db}
	req := &registerContainerDependenciesTestData
	res, err := svc.RegisterContainerDependencies(context.TODO(), req)

	if err != expectedError {
		t.Errorf("Expected %v, got %v.", expectedError, err)
	}

	if res != expectedResponse {
		t.Errorf("Expected %v, got %v.", expectedResponse, res)
	}

	if db.PutItemInvocations != expectedPutItemInvocations {
		t.Errorf("Expected %d call to PutItem, got %d", expectedPutItemInvocations, db.UpdateItemInvocations)
	}

	if db.TableName != expectedTableName {
		t.Errorf("Expected tableName to be %v, got %v.", expectedTableName, db.TableName)
	}
}

// Test_RetrieveContainerDependencies checks:
// - if correct response is returned when container dependencies are retrieved
// - if error is returned when GetItem fails
// - if error is returned when Send fails
func Test_RetrieveContainerDependencies(t *testing.T) {
	t.Parallel()

	items := []*dynamodb.GetItemOutput{
		makeGetItemOutput(
			services.ModuleDependencies{
				Name:    registerContainerDependenciesTestData.Module.Name,
				Version: registerContainerDependenciesTestData.Module.Version,
			}, t),
		makeGetItemOutput(
			services.ContainerDependencies{
				Name:    registerContainerDependenciesTestData.Module.Name,
				Version: registerContainerDependenciesTestData.Module.Version,
				Images:  registerContainerDependenciesTestData.Images,
			}, t),
	}
	db := &mocks.MockDynamoDB{GetItemOuts: items}
	dms := &services.DependencyManagerService{Db: db}
	srv := &MockRetrieveContainerDependenciesServer{}
	req := &terrarium.RetrieveContainerDependenciesRequestV2{
		Module: registerContainerDependenciesTestData.Module,
	}

	var expectedError error = nil
	var expectedServerResponse = &terrarium.ContainerDependenciesResponseV2{
		Module:       registerContainerDependenciesTestData.Module,
		Dependencies: registerContainerDependenciesTestData.Images,
	}
	var expectedGetItemInvocations = 2

	var expectedServerResponses = []*terrarium.ContainerDependenciesResponseV2{
		{
			Module:       registerContainerDependenciesTestData.Module,
			Dependencies: registerContainerDependenciesTestData.Images,
		},
	}
	err := dms.RetrieveContainerDependencies(req, srv)

	if err != expectedError {
		t.Errorf("Expected %v, got %v", expectedError, err)
	} else {
		if !reflect.DeepEqual(srv.Responses, expectedServerResponses) {
			t.Errorf("Expected server response %v, got %v.", expectedServerResponse, srv.Responses)
		}
	}

	if db.GetItemInvocations != expectedGetItemInvocations {
		t.Errorf("Expected %v call to GetItem, got %v.", expectedGetItemInvocations, db.GetItemInvocations)
	}
}

func makeGetItemOutput(in interface{}, t *testing.T) *dynamodb.GetItemOutput {
	t.Helper()
	marshalledItem, err := dynamodbattribute.MarshalMap(in)
	if err != nil {
		t.Errorf("Failed to marshal test data as a list %s", err)
	}
	return &dynamodb.GetItemOutput{Item: marshalledItem}

}

// This test checks if correct response is returned when container dependencies are retrieved
func TestRetrieveRecursiveContainerDependencies(t *testing.T) {
	t.Parallel()

	subModule := terrarium.Module{
		Name:    "test/test-submodule/all",
		Version: "2.0.2",
	}

	submoduleContainerDependencies := map[string]*terrarium.ContainerImageDetails{
		"lighstep-micro-satellite": {
			Tag:       "0.1.3",
			Namespace: "cie",
			Images: []*terrarium.ContainerImageRef{
				{
					Arch:  "amd64",
					Image: "random.server.com/my-satellite-image-for-linux-amd64:tag25",
				},
			},
		},
	}

	items := []*dynamodb.GetItemOutput{

		makeGetItemOutput(
			services.ModuleDependencies{
				Name:    registerContainerDependenciesTestData.Module.Name,
				Version: registerContainerDependenciesTestData.Module.Version,
				Modules: []*terrarium.Module{&subModule},
			}, t),
		makeGetItemOutput(
			services.ContainerDependencies{
				Name:    registerContainerDependenciesTestData.Module.Name,
				Version: registerContainerDependenciesTestData.Module.Version,
				Images:  registerContainerDependenciesTestData.Images,
			}, t),
		makeGetItemOutput(
			services.ModuleDependencies{
				Name:    subModule.Name,
				Version: subModule.Version,
			}, t),
		makeGetItemOutput(
			services.ContainerDependencies{
				Name:    subModule.Name,
				Version: subModule.Version,
				Images:  submoduleContainerDependencies,
			}, t),
	}

	db := &mocks.MockDynamoDB{GetItemOuts: items}
	dms := &services.DependencyManagerService{Db: db}
	srv := &MockRetrieveContainerDependenciesServer{}
	req := &terrarium.RetrieveContainerDependenciesRequestV2{
		Module: registerContainerDependenciesTestData.Module,
	}

	var expectedError error = nil
	var expectedServerResponses = []*terrarium.ContainerDependenciesResponseV2{
		{
			Module:       registerContainerDependenciesTestData.Module,
			Dependencies: registerContainerDependenciesTestData.Images,
		},
		{
			Module:       &subModule,
			Dependencies: submoduleContainerDependencies,
		},
	}
	err := dms.RetrieveContainerDependencies(req, srv)

	if err != expectedError {
		t.Errorf("Expected %v, got %v", expectedError, err)
	} else {
		if !reflect.DeepEqual(srv.Responses, expectedServerResponses) {
			t.Errorf("Expected server response %v, got %v.", expectedServerResponses, srv.Responses)
		}
	}
}

// This test checks if error is returned when GetItem fails
func TestRetrieveContainerDependenciesWhenModuleGetItemErrors(t *testing.T) {
	t.Parallel()
	var expectedError = services.GetModuleDependenciesError
	var expectedGetItemInvocations = 1
	var expectedServerSendInvocations = 0

	db := &mocks.MockDynamoDB{GetItemErrors: []error{errors.New("some error")}}
	dms := &services.DependencyManagerService{Db: db}
	srv := &MockRetrieveContainerDependenciesServer{}
	req := &terrarium.RetrieveContainerDependenciesRequestV2{}
	err := dms.RetrieveContainerDependencies(req, srv)

	if err != expectedError {
		t.Errorf("Expected %v, got %v", expectedError, err)
	}

	if db.GetItemInvocations != expectedGetItemInvocations {
		t.Errorf("Expected %v call to GetItem, got %v.", expectedGetItemInvocations, db.GetItemInvocations)
	}

	if srv.SendInvocations != expectedServerSendInvocations {
		t.Errorf("Expected %v calls to Send, got %v.", expectedServerSendInvocations, srv.SendInvocations)
	}
}

// This test checks if error is returned when GetItem fails
func TestRetrieveContainerDependenciesWhenContainerGetItemErrors(t *testing.T) {
	t.Parallel()

	var expectedError = services.GetContainerDependenciesError
	var expectedGetItemInvocations = 2
	var expectedServerSendInvocations = 0

	var moduleDependencies []*terrarium.Module
	var expectedModule = terrarium.Module{
		Name:    "test/test/aws",
		Version: "1.0.0",
	}
	dependencyModuleList, err := dynamodbattribute.MarshalList(moduleDependencies)
	if err != nil {
		t.Errorf("Failed to marshal test data as a list %s", err)
	}
	var moduleGetItemOutput = &dynamodb.GetItemOutput{
		ConsumedCapacity: nil,
		Item: map[string]*dynamodb.AttributeValue{
			"name":    {S: aws.String(expectedModule.Name)},
			"version": {S: aws.String(expectedModule.Version)},
			"modules": {L: dependencyModuleList},
		},
	}
	db := &mocks.MockDynamoDB{
		GetItemErrors: []error{nil, errors.New("some error")},
		GetItemOuts:   []*dynamodb.GetItemOutput{moduleGetItemOutput, nil},
	}
	dms := &services.DependencyManagerService{Db: db}
	srv := &MockRetrieveContainerDependenciesServer{}
	req := &terrarium.RetrieveContainerDependenciesRequestV2{}
	err = dms.RetrieveContainerDependencies(req, srv)

	if err != expectedError {
		t.Errorf("Expected %v, got %v", expectedError, err)
	}

	if db.GetItemInvocations != expectedGetItemInvocations {
		t.Errorf("Expected %v call to GetItem, got %v.", expectedGetItemInvocations, db.GetItemInvocations)
	}

	if srv.SendInvocations != expectedServerSendInvocations {
		t.Errorf("Expected %v calls to Send, got %v.", expectedServerSendInvocations, srv.SendInvocations)
	}
}

// TODO: Test for UnmarshalContainerDependenciesError

// This test checks if error is returned when Send fails
func TestRetrieveContainerDependenciesWhenSendErrors(t *testing.T) {
	t.Parallel()
	var expectedError = services.SendContainerDependenciesError
	var expectedModule = terrarium.Module{
		Name:    "test/test/aws",
		Version: "1.0.0",
	}
	var containerDependencies = map[string]*terrarium.ContainerImageDetails{
		"grafana": {
			Tag:       "0.1.1",
			Namespace: "cie",
			Images: []*terrarium.ContainerImageRef{
				{
					Arch:  "amd64",
					Image: "random.server.com/my-grafana-image-for-linux-amd64:tag23",
				},
			},
		},
		"kubescaler": {
			Tag:       "0.2.2",
			Namespace: "cie",
			Images: []*terrarium.ContainerImageRef{
				{
					Arch:  "amd64",
					Image: "random.server.com/my-kubescaler-image-for-linux-amd64:tag25",
				},
				{
					Arch:  "arm64",
					Image: "random.server.com/my-kubescaler-image-for-linux-arm64:graviton2",
				},
			},
		},
	}
	var expectedGetItemInvocations = 2
	var expectedServerSendInvocations = 1
	var moduleDependencies []*terrarium.Module

	dependencyModuleList, err := dynamodbattribute.MarshalList(moduleDependencies)
	if err != nil {
		t.Errorf("Failed to marshal test data as a list %s", err)
	}
	var moduleGetItemOutput = &dynamodb.GetItemOutput{
		ConsumedCapacity: nil,
		Item: map[string]*dynamodb.AttributeValue{
			"name":    {S: aws.String(expectedModule.Name)},
			"version": {S: aws.String(expectedModule.Version)},
			"modules": {L: dependencyModuleList},
		},
	}

	dependencyContainerMap, err := dynamodbattribute.MarshalMap(containerDependencies)
	if err != nil {
		t.Errorf("Failed to marshal test data as a list %s", err)
	}
	var containerGetItemOutput = &dynamodb.GetItemOutput{
		ConsumedCapacity: nil,
		Item: map[string]*dynamodb.AttributeValue{
			"name":    {S: aws.String(expectedModule.Name)},
			"version": {S: aws.String(expectedModule.Version)},
			"images":  {M: dependencyContainerMap},
		},
	}
	db := &mocks.MockDynamoDB{GetItemOuts: []*dynamodb.GetItemOutput{moduleGetItemOutput, containerGetItemOutput}}
	dms := &services.DependencyManagerService{Db: db}

	srv := &MockRetrieveContainerDependenciesServer{Err: errors.New("some error")}
	req := &terrarium.RetrieveContainerDependenciesRequestV2{
		Module: &expectedModule,
	}
	err = dms.RetrieveContainerDependencies(req, srv)

	if err != expectedError {
		t.Errorf("Expected %v, got %v", expectedError, err)
	}

	if db.GetItemInvocations != expectedGetItemInvocations {
		t.Errorf("Expected %v call to GetItem, got %v.", expectedGetItemInvocations, db.GetItemInvocations)
	}

	if srv.SendInvocations != expectedServerSendInvocations {
		t.Errorf("Expected %v calls to Send, got %v.", expectedServerSendInvocations, srv.SendInvocations)
	}
}

// Test_RetrieveModuleDependencies checks:
// - if correct response is returned when module dependencies are retrieved
// - if error is returned when GetItem fails
// - if error is returned when Send fails
func Test_RetrieveModuleDependencies(t *testing.T) {
	t.Parallel()

	t.Run("when module dependencies are retrieved", func(t *testing.T) {
		out := &dynamodb.GetItemOutput{}

		db := &mocks.MockDynamoDB{GetItemOuts: []*dynamodb.GetItemOutput{out}}

		dms := &services.DependencyManagerService{Db: db}

		srv := &MockRetrieveModuleDependenciesServer{}

		req := &terrarium.RetrieveModuleDependenciesRequest{}

		err := dms.RetrieveModuleDependencies(req, srv)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if db.GetItemInvocations != 1 {
			t.Errorf("Expected 1 call to GetItem, got %v.", db.GetItemInvocations)
		}

		if srv.SendInvocations != 1 {
			t.Errorf("Expected 1 call to Send, got %v.", srv.SendInvocations)
		}

		if srv.Response == nil {
			t.Error("Expected response, got nil.")
		}
	})

	t.Run("when GetItem fails", func(t *testing.T) {
		db := &mocks.MockDynamoDB{GetItemErrors: []error{errors.New("some error")}}

		dms := &services.DependencyManagerService{Db: db}

		srv := &MockRetrieveModuleDependenciesServer{}

		req := &terrarium.RetrieveModuleDependenciesRequest{}

		err := dms.RetrieveModuleDependencies(req, srv)

		if err != services.GetModuleDependenciesError {
			t.Errorf("Expected %v, got %v", services.GetModuleDependenciesError, err)
		}

		if db.GetItemInvocations != 1 {
			t.Errorf("Expected 1 call to GetItem, got %v.", db.GetItemInvocations)
		}

		if srv.SendInvocations != 0 {
			t.Errorf("Expected 0 calls to Send, got %v.", srv.SendInvocations)
		}
	})

	// TODO: Test for UnmarshalModuleDependenciesError

	t.Run("when Send fails", func(t *testing.T) {
		out := &dynamodb.GetItemOutput{}

		db := &mocks.MockDynamoDB{GetItemOuts: []*dynamodb.GetItemOutput{out}}

		dms := &services.DependencyManagerService{Db: db}

		srv := &MockRetrieveModuleDependenciesServer{Err: errors.New("some error")}

		req := &terrarium.RetrieveModuleDependenciesRequest{}

		err := dms.RetrieveModuleDependencies(req, srv)

		if err != services.SendModuleDependenciesError {
			t.Errorf("Expected %v, got %v", services.SendModuleDependenciesError, err)
		}

		if db.GetItemInvocations != 1 {
			t.Errorf("Expected 1 call to GetItem, got %v.", db.GetItemInvocations)
		}

		if srv.SendInvocations != 1 {
			t.Errorf("Expected 1 call to Send, got %v.", srv.SendInvocations)
		}
	})
}

// Test_GetDependencies checks:
// - if correct response is returned when dependencies are retrieved
// - if error is returnd when GetItem fails
func Test_GetDependencies(t *testing.T) {
	t.Parallel()

	t.Run("when dependencies are retrieved", func(t *testing.T) {
		out := &dynamodb.GetItemOutput{Item: map[string]*dynamodb.AttributeValue{
			"name":    {S: aws.String("cietest/notify/aws")},
			"version": {S: aws.String("1.0.2")},
			"modules": {L: []*dynamodb.AttributeValue{{M: map[string]*dynamodb.AttributeValue{"name": {S: aws.String("cietest/lambda/aws")}, "version": {S: aws.String("1.0.1")}}}}},
			"images":  {SS: aws.StringSlice([]string{"slice", "slicee"})},
		}}

		db := &mocks.MockDynamoDB{GetItemOuts: []*dynamodb.GetItemOutput{out}}

		dms := &services.DependencyManagerService{Db: db}

		m := &terrarium.Module{Name: "cietest/notify/aws", Version: "1.0.2"}
		res := &MockGetDependenciesResponse{
			Dependencies: []*terrarium.Module{{Name: "cietest/lambda/aws", Version: "1.0.1"}},
			Err:          nil,
		}

		deps, err := dms.GetModuleDependencies(m)

		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if db.GetItemInvocations != 1 {
			t.Errorf("Expected 1 call to GetItem, got %v.", db.GetItemInvocations)
		}

		if deps == nil {
			t.Errorf("Expected response, got nil.")
		}

		if !reflect.DeepEqual(deps, res.Dependencies) {
			t.Errorf("Expected response %v, got different. %v", res.Dependencies, deps)
		}
	})

	// Test for UnmarshalModuleDependenciesError

	t.Run("when GetItem fails", func(t *testing.T) {
		db := &mocks.MockDynamoDB{GetItemErrors: []error{errors.New("some error")}}

		dms := &services.DependencyManagerService{Db: db}

		m := &terrarium.Module{Name: "cietest/notify/aws", Version: "1.0.2"}

		deps, err := dms.GetModuleDependencies(m)

		if deps != nil {
			t.Errorf("Expected nil, got %v", deps)
		}

		if db.GetItemInvocations != 1 {
			t.Errorf("Expected 1 call to GetItem, got %v.", db.GetItemInvocations)
		}

		if err == nil {
			t.Errorf("Expected error, got nil.")
		}
	})
}
