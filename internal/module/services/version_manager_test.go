package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/terrariumcloud/terrarium/internal/mocks"
	"github.com/terrariumcloud/terrarium/internal/module/services"
	terrarium "github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"google.golang.org/grpc"
)

// This test checks if there was no error
func TestRegisterVersionManagerWithServer(t *testing.T) {
	t.Parallel()

	db := &mocks.MockDynamoDB{}

	vms := &services.VersionManagerService{Db: db}

	s := grpc.NewServer(*new([]grpc.ServerOption)...)

	err := vms.RegisterWithServer(s)

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
func TestRegisterWithServerWhenVersionsTableInitializationErrors(t *testing.T) {
	t.Parallel()

	db := &mocks.MockDynamoDB{DescribeTableError: errors.New("some error")}

	vms := &services.VersionManagerService{Db: db}

	s := grpc.NewServer(*new([]grpc.ServerOption)...)

	err := vms.RegisterWithServer(s)

	if err != services.ModuleVersionsTableInitializationError {
		t.Errorf("Expected %v, got %v.", services.ModuleVersionsTableInitializationError, err)
	}

	if db.DescribeTableInvocations != 1 {
		t.Errorf("Expected 1 call to DescribeTable, got %v.", db.DescribeTableInvocations)
	}

	if db.CreateTableInvocations != 0 {
		t.Errorf("Expected 0 calls to CreateTable, got %v.", db.CreateTableInvocations)
	}
}

// This test checks if correct response is returned when new version is created
func TestBeginVersion(t *testing.T) {
	t.Parallel()

	db := &mocks.MockDynamoDB{}

	svc := &services.VersionManagerService{Db: db}

	req := &terrarium.BeginVersionRequest{Module: &terrarium.Module{Name: "test", Version: "v1.0.0"}}

	res, err := svc.BeginVersion(context.TODO(), req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if db.PutItemInvocations != 1 {
		t.Errorf("Expected 1 call to PutItem, got %v", db.PutItemInvocations)
	}

	if db.TableName != services.VersionsTableName {
		t.Errorf("Expected tableName to be %v, got %v.", services.VersionsTableName, db.TableName)
	}

	if res != services.VersionCreated {
		t.Errorf("Expected %v, got %v.", services.VersionCreated, res)
	}
}

// This test checks if error is returned when PutItem fails
func TestBeginVersionWhenPutItemErrors(t *testing.T) {
	t.Parallel()

	db := &mocks.MockDynamoDB{PutItemError: errors.New("some error")}

	svc := &services.VersionManagerService{Db: db}

	req := &terrarium.BeginVersionRequest{Module: &terrarium.Module{Name: "test", Version: "v1.0.0"}}

	res, err := svc.BeginVersion(context.TODO(), req)

	if res != nil {
		t.Errorf("Expected no response, got %v", res)
	}

	if db.PutItemInvocations != 1 {
		t.Errorf("Expected 1 call to PutItem, got %v", db.PutItemInvocations)
	}

	if db.TableName != services.VersionsTableName {
		t.Errorf("Expected tableName to be %v, got %v.", services.VersionsTableName, db.TableName)
	}

	if err != services.CreateModuleVersionError {
		t.Errorf("Expected %v, got %v.", services.CreateModuleVersionError, err)
	}
}

// TODO: Test for MarshalModuleVersionError

// This test checks if correct response is returned when version is aborted
func TestAbortVersion(t *testing.T) {
	t.Parallel()

	db := &mocks.MockDynamoDB{}

	svc := &services.VersionManagerService{Db: db}

	req := &services.TerminateVersionRequest{Module: &terrarium.Module{Name: "test", Version: "v1.0.0"}}

	res, err := svc.AbortVersion(context.TODO(), req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if db.DeleteItemInvocations != 1 {
		t.Errorf("Expected 1 call to DeleteItem, got %v", db.DeleteItemInvocations)
	}

	if db.TableName != services.VersionsTableName {
		t.Errorf("Expected tableName to be %v, got %v.", services.VersionsTableName, db.TableName)
	}

	if res != services.VersionAborted {
		t.Errorf("Expected %v, got %v.", services.VersionAborted, res)
	}
}

// This test checks if error is returned when DeleteItem fails
func TestAbortVersionWhenDeleteItemErrors(t *testing.T) {
	t.Parallel()

	db := &mocks.MockDynamoDB{DeleteItemError: errors.New("some error")}

	svc := &services.VersionManagerService{Db: db}

	req := services.TerminateVersionRequest{Module: &terrarium.Module{Name: "test", Version: "v1.0.0"}}

	res, err := svc.AbortVersion(context.TODO(), &req)

	if res != nil {
		t.Errorf("Expected no response, got %v", res)
	}

	if db.DeleteItemInvocations != 1 {
		t.Errorf("Expected 1 call to DeleteItem, got %v", db.DeleteItemInvocations)
	}

	if db.TableName != services.VersionsTableName {
		t.Errorf("Expected tableName to be %v, got %v.", services.VersionsTableName, db.TableName)
	}

	if err != services.AbortModuleVersionError {
		t.Errorf("Expected %v, got %v.", services.AbortModuleVersionError, err)
	}
}

// This test checks if correct response is returned when version is published
func TestPublishVersion(t *testing.T) {
	t.Parallel()

	db := &mocks.MockDynamoDB{}

	svc := &services.VersionManagerService{Db: db}

	req := &services.TerminateVersionRequest{Module: &terrarium.Module{Name: "test", Version: "v1.0.0"}}

	res, err := svc.PublishVersion(context.TODO(), req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if db.UpdateItemInvocations != 1 {
		t.Errorf("Expected 1 call to UpdateItem, got %v", db.UpdateItemInvocations)
	}

	if db.TableName != services.VersionsTableName {
		t.Errorf("Expected tableName to be %v, got %v.", services.VersionsTableName, db.TableName)
	}

	if res != services.VersionPublished {
		t.Errorf("Expected %v, got %v.", services.VersionPublished, res)
	}
}

// This test checks if error is returned when UpdateItem fails
func TestPublishVersionWhenUpdateItemErrors(t *testing.T) {
	t.Parallel()

	db := &mocks.MockDynamoDB{UpdateItemError: errors.New("some error")}

	svc := &services.VersionManagerService{Db: db}

	req := services.TerminateVersionRequest{Module: &terrarium.Module{Name: "test", Version: "v1.0.0"}}

	res, err := svc.PublishVersion(context.TODO(), &req)

	if res != nil {
		t.Errorf("Expected no response, got %v", res)
	}

	if db.UpdateItemInvocations != 1 {
		t.Errorf("Expected 1 call to UpdateItem, got %v", db.UpdateItemInvocations)
	}

	if db.TableName != services.VersionsTableName {
		t.Errorf("Expected tableName to be %v, got %v.", services.VersionsTableName, db.TableName)
	}

	if err != services.PublishModuleVersionError {
		t.Errorf("Expected %v, got %v.", services.PublishModuleVersionError, err)
	}
}
