package services_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/terrariumcloud/terrarium/internal/mocks"
	"github.com/terrariumcloud/terrarium/internal/module/services"
	terrarium "github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"google.golang.org/grpc"
)

type MockRetrieveContainerDependenciesServer struct {
	grpc.ServerStream
	SendInvocations int
	Response        *terrarium.ContainerDependenciesResponse
	Err             error
}

func (srv *MockRetrieveContainerDependenciesServer) Send(res *terrarium.ContainerDependenciesResponse) error {
	srv.SendInvocations++
	srv.Response = res
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

// This test checks if there was no error
func TestRegisterDependencyManagerWithServer(t *testing.T) {
	t.Parallel()

	db := &mocks.MockDynamoDB{}

	dms := &services.DependencyManagerService{Db: db}

	s := grpc.NewServer(*new([]grpc.ServerOption)...)

	err := dms.RegisterWithServer(s)

	if err != nil {
		t.Errorf("Expected no error, got %v.", err)
	}

	if db.DescribeTableInvocations != 1 {
		t.Errorf("Expected 1 call to DescribeTable, got %v.", db.DescribeTableInvocations)
	}

	if db.CreateTableInvocations != 0 {
		t.Errorf("Expected no calls to CreateTable, got %v.", db.CreateTableInvocations)
	}
}

// This test checks if error is returned when Table initialization fails
func TestRegisterWithServerWhenModuleDependenciesTableInitializationErrors(t *testing.T) {
	t.Parallel()

	db := &mocks.MockDynamoDB{DescribeTableError: errors.New("some error")}

	dms := &services.DependencyManagerService{Db: db}

	s := grpc.NewServer(*new([]grpc.ServerOption)...)

	err := dms.RegisterWithServer(s)

	if err != services.ModuleDependenciesTableInitializationError {
		t.Errorf("Expected %v, got %v.", services.ModuleDependenciesTableInitializationError, err)
	}

	if db.DescribeTableInvocations != 1 {
		t.Errorf("Expected 1 call to DescribeTable, got %v.", db.DescribeTableInvocations)
	}

	if db.CreateTableInvocations != 0 {
		t.Errorf("Expected 0 calls to CreateTable, got %v.", db.CreateTableInvocations)
	}
}

// This test checks if correct response is returned when module dependencies are registered
func TestRegisterModuleDependencies(t *testing.T) {
	t.Parallel()

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
}

// This test checks if error is returned when PutItem fails
func TestRegisterModuleDependenciesWhenPutItemErrors(t *testing.T) {
	t.Parallel()

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

	if err != services.RegisterModuleDependenciesError {
		t.Errorf("Expected %v, got %v.", services.RegisterModuleDependenciesError, err)
	}
}

// TODO: Test for MarshalModuleDependenciesError

// This test checks if correct response is returned when container dependencies are registered
func TestRegisterContainerDependencies(t *testing.T) {
	t.Parallel()

	db := &mocks.MockDynamoDB{}

	svc := &services.DependencyManagerService{Db: db}

	req := &terrarium.RegisterContainerDependenciesRequest{
		Module:       &terrarium.Module{Name: "test", Version: "v1"},
		Dependencies: []string{"test", "test2"},
	}

	res, err := svc.RegisterContainerDependencies(context.TODO(), req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if db.UpdateItemInvocations != 1 {
		t.Errorf("Expected 1 call to UpdateItem, got %v", db.UpdateItemInvocations)
	}

	if db.TableName != services.ModuleDependenciesTableName {
		t.Errorf("Expected tableName to be %v, got %v.", services.ModuleDependenciesTableName, db.TableName)
	}

	if res != services.ContainerDependenciesRegistered {
		t.Errorf("Expected %v, got %v.", services.ContainerDependenciesRegistered, res)
	}
}

// This test checks if error is returned when UpdateItem fails
func TestRegisterContainerDependenciesWhenPutItemErrors(t *testing.T) {
	t.Parallel()

	db := &mocks.MockDynamoDB{UpdateItemError: errors.New("some error")}

	svc := &services.DependencyManagerService{Db: db}

	req := &terrarium.RegisterContainerDependenciesRequest{
		Module:       &terrarium.Module{Name: "test", Version: "v1"},
		Dependencies: []string{"test", "test2"},
	}

	res, err := svc.RegisterContainerDependencies(context.TODO(), req)

	if res != nil {
		t.Errorf("Expected no response, got %v", res)
	}

	if db.UpdateItemInvocations != 1 {
		t.Errorf("Expected 1 call to UpdateItem, got %v", db.UpdateItemInvocations)
	}

	if db.TableName != services.ModuleDependenciesTableName {
		t.Errorf("Expected tableName to be %v, got %v.", services.ModuleDependenciesTableName, db.TableName)
	}

	if err != services.RegisterContainerDependenciesError {
		t.Errorf("Expected %v, got %v.", services.RegisterContainerDependenciesError, err)
	}
}

// This test checks if correct response is returned when container dependencies are retrieved
func TestRetrieveContainerDependencies(t *testing.T) {
	t.Parallel()

	out := &dynamodb.GetItemOutput{}

	db := &mocks.MockDynamoDB{GetItemOut: out}

	dms := &services.DependencyManagerService{Db: db}

	srv := &MockRetrieveContainerDependenciesServer{}

	req := &terrarium.RetrieveContainerDependenciesRequest{}

	err := dms.RetrieveContainerDependencies(req, srv)

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
}

// This test checks if error is returned when GetItem fails
func TestRetrieveContainerDependenciesWhenGetItemErrors(t *testing.T) {
	t.Parallel()

	db := &mocks.MockDynamoDB{GetItemError: errors.New("some error")}

	dms := &services.DependencyManagerService{Db: db}

	srv := &MockRetrieveContainerDependenciesServer{}

	req := &terrarium.RetrieveContainerDependenciesRequest{}

	err := dms.RetrieveContainerDependencies(req, srv)

	if err != services.GetContainerDependenciesError {
		t.Errorf("Expected %v, got %v", services.GetContainerDependenciesError, err)
	}

	if db.GetItemInvocations != 1 {
		t.Errorf("Expected 1 call to GetItem, got %v.", db.GetItemInvocations)
	}

	if srv.SendInvocations != 0 {
		t.Errorf("Expected 0 calls to Send, got %v.", srv.SendInvocations)
	}
}

// TODO: Test for UnmarshalContainerDependenciesError

// This test checks if error is returned when Send fails
func TestRetrieveContainerDependenciesWhenSendErrors(t *testing.T) {
	t.Parallel()

	out := &dynamodb.GetItemOutput{}

	db := &mocks.MockDynamoDB{GetItemOut: out}

	dms := &services.DependencyManagerService{Db: db}

	srv := &MockRetrieveContainerDependenciesServer{Err: errors.New("some error")}

	req := &terrarium.RetrieveContainerDependenciesRequest{}

	err := dms.RetrieveContainerDependencies(req, srv)

	if err != services.SendContainerDependenciesError {
		t.Errorf("Expected %v, got %v", services.SendContainerDependenciesError, err)
	}

	if db.GetItemInvocations != 1 {
		t.Errorf("Expected 1 call to GetItem, got %v.", db.GetItemInvocations)
	}

	if srv.SendInvocations != 1 {
		t.Errorf("Expected 1 call to Send, got %v.", srv.SendInvocations)
	}
}

// This test checks if correct response is returned when module dependencies are retrieved
func TestRetrieveModuleDependencies(t *testing.T) {
	t.Parallel()

	out := &dynamodb.GetItemOutput{}

	db := &mocks.MockDynamoDB{GetItemOut: out}

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
}

// This test checks if error is returned when GetItem fails
func TestRetrieveModuleDependenciesWhenGetItemErrors(t *testing.T) {
	t.Parallel()

	db := &mocks.MockDynamoDB{GetItemError: errors.New("some error")}

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
}

// TODO: Test for UnmarshalModuleDependenciesError

// This test checks if error is returned when Send fails
func TestRetrieveModuleDependenciesWhenSendErrors(t *testing.T) {
	t.Parallel()

	out := &dynamodb.GetItemOutput{}

	db := &mocks.MockDynamoDB{GetItemOut: out}

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
}

func TestGetDependencies(t *testing.T) {
	t.Parallel()

	out := &dynamodb.GetItemOutput{Item: map[string]*dynamodb.AttributeValue{
		"name":    {S: aws.String("cietest/notify/aws")},
		"version": {S: aws.String("1.0.2")},
		"modules": {L: []*dynamodb.AttributeValue{{M: map[string]*dynamodb.AttributeValue{"name": {S: aws.String("cietest/lambda/aws")}, "version": {S: aws.String("1.0.1")}}}}},
		"images":  {SS: aws.StringSlice([]string{"slice", "slicee"})},
	}}

	db := &mocks.MockDynamoDB{GetItemOut: out}

	dms := &services.DependencyManagerService{Db: db}


	mockRequestModule := &terrarium.Module{Name: "cietest/notify/aws", Version: "1.0.2"}
	mockResponseModule := &MockGetDependenciesResponse{
		Dependencies: []*terrarium.Module{{Name: "cietest/lambda/aws", Version: "1.0.1"}},
		Err:          nil,
	}

	deps, err := dms.GetDependencies(mockRequestModule)

	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if db.GetItemInvocations != 1 {
		t.Errorf("Expected 1 call to GetItem, got %v.", db.GetItemInvocations)
	}

	if deps == nil {
		t.Errorf("Expected response, got nil.")
	}
	if !reflect.DeepEqual(deps, mockResponseModule.Dependencies) {
		t.Errorf("Expected response %v, got different. %v", mockResponseModule.Dependencies, deps)
	}
}
