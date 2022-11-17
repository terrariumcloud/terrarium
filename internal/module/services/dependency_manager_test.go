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

// Test_RegisterDependencyManagerWithServer checks:
// - if there was no error with table init
// - if error is returned when Table initialization fails
func Test_RegisterDependencyManagerWithServer(t *testing.T) {
	t.Parallel()

	t.Run("when there is no error with table init", func(t *testing.T) {
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
	})

	t.Run("when Table initialization fails", func(t *testing.T) {
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
	})
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

		if err != services.RegisterModuleDependenciesError {
			t.Errorf("Expected %v, got %v.", services.RegisterModuleDependenciesError, err)
		}
	})
}

// Test_RegisterContainerDependencies checks:
// - if correct response is returned when container dependencies are registered
// - if error is returned when UpdateItem fails
func Test_RegisterContainerDependencies(t *testing.T) {
	t.Parallel()

	t.Run("when container dependencies are registered", func(t *testing.T) {
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
	})

	t.Run("when UpdateItem fails", func(t *testing.T) {
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
	})
}

// Test_RetrieveContainerDependencies checks:
// - if correct response is returned when container dependencies are retrieved
// - if error is returned when GetItem fails
// - if error is returned when Send fails
func Test_RetrieveContainerDependencies(t *testing.T) {
	t.Parallel()

	t.Run("when container dependencies are retrieved", func(t *testing.T) {
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
	})

	t.Run("when GetItem fails", func(t *testing.T) {
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
	})

	// TODO: Test for UnmarshalContainerDependenciesError

	t.Run("when Send fails", func(t *testing.T) {
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
	})
}

// Test_RetrieveModuleDependencies checks:
// - if correct response is returned when module dependencies are retrieved
// - if error is returned when GetItem fails
// - if error is returned when Send fails
func Test_RetrieveModuleDependencies(t *testing.T) {
	t.Parallel()

	t.Run("when module dependencies are retrieved", func(t *testing.T) {
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
	})

	t.Run("when GetItem fails", func(t *testing.T) {
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
	})

	// TODO: Test for UnmarshalModuleDependenciesError

	t.Run("when Send fails", func(t *testing.T) {
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

		db := &mocks.MockDynamoDB{GetItemOut: out}

		dms := &services.DependencyManagerService{Db: db}

		m := &terrarium.Module{Name: "cietest/notify/aws", Version: "1.0.2"}
		res := &MockGetDependenciesResponse{
			Dependencies: []*terrarium.Module{{Name: "cietest/lambda/aws", Version: "1.0.1"}},
			Err:          nil,
		}

		deps, err := dms.GetDependencies(m)

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
		db := &mocks.MockDynamoDB{GetItemError: errors.New("some error")}

		dms := &services.DependencyManagerService{Db: db}

		m := &terrarium.Module{Name: "cietest/notify/aws", Version: "1.0.2"}

		deps, err := dms.GetDependencies(m)

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
