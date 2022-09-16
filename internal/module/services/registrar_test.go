package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/mocks"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/module/services"
	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"
	"google.golang.org/grpc"
)

// This test checks if correct response is returned when module is registered
func TestRegisterModule(t *testing.T) {
	t.Parallel()

	db := &mocks.MockDynamoDB{}

	svc := &services.RegistrarService{Db: db}

	req := terrarium.RegisterModuleRequest{
		Name:        "test",
		Description: "test desc",
		Source:      "http://test.com",
		Maturity:    terrarium.Maturity_ALPHA,
	}

	res, err := svc.Register(context.TODO(), &req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if db.PutItemInvocations != 1 {
		t.Errorf("Expected 1 call to PutItem, got %d", db.PutItemInvocations)
	}

	if db.TableName != services.RegistrarTableName {
		t.Errorf("Expected tableName to be %s, got %s", services.RegistrarTableName, db.TableName)
	}

	if res != services.ModuleRegistered {
		t.Errorf("Expected %v, got %v.", services.ModuleRegistered, res)
	}
}

// This test checks if error is returned marshal error
func IgnoreTestRegisterModuleWhenMarshalModuleErrors(t *testing.T) {
	t.Parallel()

	db := &mocks.MockDynamoDB{}

	svc := &services.RegistrarService{Db: db}

	req := &terrarium.RegisterModuleRequest{} // TODO: need to make MarshalMap return error

	res, err := svc.Register(context.TODO(), req)

	if err != services.MarshalModuleError {
		t.Errorf("Expected %v, got %v.", services.MarshalModuleError, err)
	}

	if db.PutItemInvocations == 0 {
		t.Errorf("Expected 0 calls to PutItem, got %d", db.PutItemInvocations)
	}

	if res != nil {
		t.Errorf("Expected no response, got %v.", res)
	}
}

// This test checks if error is returned when PutItem fails
func TestRegisterModuleWhenPutItemErrors(t *testing.T) {
	t.Parallel()

	db := &mocks.MockDynamoDB{
		PutItemError: errors.New("some error"),
	}

	svc := &services.RegistrarService{Db: db}

	req := terrarium.RegisterModuleRequest{
		Name:        "test",
		Description: "test desc",
		Source:      "http://test.com",
		Maturity:    terrarium.Maturity_ALPHA,
	}

	res, err := svc.Register(context.TODO(), &req)

	if res != nil {
		t.Errorf("Expected no response, got %v", err)
	}

	if db.PutItemInvocations != 1 {
		t.Errorf("Expected 1 call to PutItem, got %d", db.PutItemInvocations)
	}

	if db.TableName != services.RegistrarTableName {
		t.Errorf("Expected tableName to be %s, got %s", services.RegistrarTableName, db.TableName)
	}

	if err != services.ModuleRegisterError {
		t.Errorf("Expected %v, got %v.", services.ModuleRegisterError, err)
	}
}

// This test checks if there was no error
func TestRegisterWithServer(t *testing.T) {
	t.Parallel()

	db := &mocks.MockDynamoDB{}

	rs := &services.RegistrarService{
		Db: db,
	}

	s := grpc.NewServer(*new([]grpc.ServerOption)...)

	err := rs.RegisterWithServer(s)

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
func TestRegisterWithServerWhenModuleTableInitializationErrors(t *testing.T) {
	t.Parallel()

	db := &mocks.MockDynamoDB{
		DescribeTableError: errors.New("some error"),
	}

	rs := &services.RegistrarService{
		Db: db,
	}

	s := grpc.NewServer(*new([]grpc.ServerOption)...)

	err := rs.RegisterWithServer(s)

	if err != services.ModuleTableInitializationError {
		t.Errorf("Expected %v, got %v.", services.ModuleTableInitializationError, err)
	}

	if db.DescribeTableInvocations != 1 {
		t.Errorf("Expected 1 call to DescribeTable, got %v.", db.DescribeTableInvocations)
	}

	if db.CreateTableInvocations != 0 {
		t.Errorf("Expected 0 calls to CreateTable, got %v.", db.CreateTableInvocations)
	}
}
