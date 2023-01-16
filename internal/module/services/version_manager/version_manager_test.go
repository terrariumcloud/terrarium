package version_manager

import (
	"context"
	"errors"
	"github.com/terrariumcloud/terrarium/internal/module/services"
	"github.com/terrariumcloud/terrarium/internal/storage/mocks"
	"testing"

	terrarium "github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"google.golang.org/grpc"
)

// Test_RegisterVersionManagerWithServer checks:
// - if there was no error with table init
// - if error is returned when Table initialization fails
func Test_RegisterVersionManagerWithServer(t *testing.T) {
	t.Parallel()

	t.Run("when table init is successful", func(t *testing.T) {
		db := &mocks.DynamoDB{}

		vms := &VersionManagerService{Db: db}

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
		db := &mocks.DynamoDB{DescribeTableErrors: []error{errors.New("some error")}}

		vms := &VersionManagerService{Db: db}

		s := grpc.NewServer(*new([]grpc.ServerOption)...)

		err := vms.RegisterWithServer(s)

		if err != ModuleVersionsTableInitializationError {
			t.Errorf("Expected %v, got %v.", ModuleVersionsTableInitializationError, err)
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
		db := &mocks.DynamoDB{}

		svc := &VersionManagerService{Db: db}

		req := &terrarium.BeginVersionRequest{Module: &terrarium.Module{Name: "test", Version: "v1.0.0"}}

		res, err := svc.BeginVersion(context.TODO(), req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if db.PutItemInvocations != 1 {
			t.Errorf("Expected 1 call to PutItem, got %v", db.PutItemInvocations)
		}

		if db.TableName != VersionsTableName {
			t.Errorf("Expected tableName to be %v, got %v.", VersionsTableName, db.TableName)
		}

		if res != VersionCreated {
			t.Errorf("Expected %v, got %v.", VersionCreated, res)
		}
	})

	t.Run("when PutItem fails", func(t *testing.T) {
		db := &mocks.DynamoDB{PutItemError: errors.New("some error")}

		svc := &VersionManagerService{Db: db}

		req := &terrarium.BeginVersionRequest{Module: &terrarium.Module{Name: "test", Version: "v1.0.0"}}

		res, err := svc.BeginVersion(context.TODO(), req)

		if res != nil {
			t.Errorf("Expected no response, got %v", res)
		}

		if db.PutItemInvocations != 1 {
			t.Errorf("Expected 1 call to PutItem, got %v", db.PutItemInvocations)
		}

		if db.TableName != VersionsTableName {
			t.Errorf("Expected tableName to be %v, got %v.", VersionsTableName, db.TableName)
		}

		if err != CreateModuleVersionError {
			t.Errorf("Expected %v, got %v.", CreateModuleVersionError, err)
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
		db := &mocks.DynamoDB{}

		svc := &VersionManagerService{Db: db}

		req := &services.TerminateVersionRequest{Module: &terrarium.Module{Name: "test", Version: "v1.0.0"}}

		res, err := svc.AbortVersion(context.TODO(), req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if db.DeleteItemInvocations != 1 {
			t.Errorf("Expected 1 call to DeleteItem, got %v", db.DeleteItemInvocations)
		}

		if db.TableName != VersionsTableName {
			t.Errorf("Expected tableName to be %v, got %v.", VersionsTableName, db.TableName)
		}

		if res != VersionAborted {
			t.Errorf("Expected %v, got %v.", VersionAborted, res)
		}
	})

	t.Run("when DeleteItem fails", func(t *testing.T) {
		db := &mocks.DynamoDB{DeleteItemError: errors.New("some error")}

		svc := &VersionManagerService{Db: db}

		req := services.TerminateVersionRequest{Module: &terrarium.Module{Name: "test", Version: "v1.0.0"}}

		res, err := svc.AbortVersion(context.TODO(), &req)

		if res != nil {
			t.Errorf("Expected no response, got %v", res)
		}

		if db.DeleteItemInvocations != 1 {
			t.Errorf("Expected 1 call to DeleteItem, got %v", db.DeleteItemInvocations)
		}

		if db.TableName != VersionsTableName {
			t.Errorf("Expected tableName to be %v, got %v.", VersionsTableName, db.TableName)
		}

		if err != AbortModuleVersionError {
			t.Errorf("Expected %v, got %v.", AbortModuleVersionError, err)
		}
	})
}

// Test_PublishVersion checks:
// - if correct response is returned when version is published
// - if error is returned when UpdateItem fails
func Test_PublishVersion(t *testing.T) {
	t.Parallel()

	t.Run("when version is published", func(t *testing.T) {
		db := &mocks.DynamoDB{}

		svc := &VersionManagerService{Db: db}

		req := &services.TerminateVersionRequest{Module: &terrarium.Module{Name: "test", Version: "v1.0.0"}}

		res, err := svc.PublishVersion(context.TODO(), req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if db.UpdateItemInvocations != 1 {
			t.Errorf("Expected 1 call to UpdateItem, got %v", db.UpdateItemInvocations)
		}

		if db.TableName != VersionsTableName {
			t.Errorf("Expected tableName to be %v, got %v.", VersionsTableName, db.TableName)
		}

		if res != VersionPublished {
			t.Errorf("Expected %v, got %v.", VersionPublished, res)
		}
	})

	t.Run("when UpdateItem fails", func(t *testing.T) {
		db := &mocks.DynamoDB{UpdateItemError: errors.New("some error")}

		svc := &VersionManagerService{Db: db}

		req := services.TerminateVersionRequest{Module: &terrarium.Module{Name: "test", Version: "v1.0.0"}}

		res, err := svc.PublishVersion(context.TODO(), &req)

		if res != nil {
			t.Errorf("Expected no response, got %v", res)
		}

		if db.UpdateItemInvocations != 1 {
			t.Errorf("Expected 1 call to UpdateItem, got %v", db.UpdateItemInvocations)
		}

		if db.TableName != VersionsTableName {
			t.Errorf("Expected tableName to be %v, got %v.", VersionsTableName, db.TableName)
		}

		if err != PublishModuleVersionError {
			t.Errorf("Expected %v, got %v.", PublishModuleVersionError, err)
		}
	})
}
