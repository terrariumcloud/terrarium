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

// Test_RegisterVersionManagerWithServer checks:
// - if there was no error with table init
// - if error is returned when Table initialization fails
func Test_RegisterVersionManagerWithServer(t *testing.T) {
	t.Parallel()

	t.Run("when there is no error with table init", func(t *testing.T) {
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
	})

	t.Run("when Table initialization fails", func(t *testing.T) {
		db := &mocks.MockDynamoDB{DescribeTableErrors: []error{errors.New("some error")}}

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
	})
}

// Test_BeginVersion checks:
// - if correct response is returned when new version is created
// - if error is returned when PutItem fails
func Test_BeginVersion(t *testing.T) {
	t.Parallel()

	t.Run("when new version is created", func(t *testing.T) {
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
	})

	t.Run("when PutItem fails", func(t *testing.T) {
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
	})

	// TODO: Test for MarshalModuleVersionError
}

// Test_AbortVersion checks:
// - if correct response is returned when version is aborted
// - if error is returned when DeleteItem fails
func Test_AbortVersion(t *testing.T) {
	t.Parallel()

	t.Run("when version is aborted", func(t *testing.T) {
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
	})

	t.Run("when DeleteItem fails", func(t *testing.T) {
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
	})
}

// Test_PublishVersion checks:
// - if correct response is returned when version is published
// - if error is returned when UpdateItem fails
func Test_PublishVersion(t *testing.T) {
	t.Parallel()

	t.Run("when version is published", func(t *testing.T) {
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
	})

	t.Run("when UpdateItem fails", func(t *testing.T) {
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
	})
}
