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

func Test_RegisterTagManagerWithServer(t *testing.T) {
	t.Parallel()

	t.Run("when table init is successful", func(t *testing.T) {
		db := &mocks.MockDynamoDB{}

		tms := &services.TagManagerService{Db: db}

		s := grpc.NewServer(*new([]grpc.ServerOption)...)

		err := tms.RegisterWithServer(s)

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

		tms := &services.TagManagerService{Db: db}

		s := grpc.NewServer(*new([]grpc.ServerOption)...)

		err := tms.RegisterWithServer(s)

		if err != services.ModuleTagTableInitializationError {
			t.Errorf("Expected %v, got %v.", services.ModuleTagTableInitializationError, err)
		}

		if db.DescribeTableInvocations != 1 {
			t.Errorf("Expected 1 call to DescribeTable, got %v.", db.DescribeTableInvocations)
		}

		if db.CreateTableInvocations != 0 {
			t.Errorf("Expected 0 calls to CreateTable, got %v.", db.CreateTableInvocations)
		}
	})
}

func Test_PublishTag(t *testing.T) {
	t.Parallel()

	t.Run("when tag is published", func(t *testing.T) {
		db := &mocks.MockDynamoDB{}

		svc := &services.TagManagerService{Db: db}

		req := &services.PublishTagRequest{
		Module: &terrarium.Module{Name: "test", Version: "v1.0.0"}, 
		Tags: map[string]*services.TagList{
			"EKS": {
			  Tag: "eks-bundle-tag",
			},
		},
		}

		res, err := svc.PublishTag(context.TODO(), req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if db.UpdateItemInvocations != 1 {
			t.Errorf("Expected 1 call to UpdateItem, got %v", db.UpdateItemInvocations)
		}

		if db.TableName != services.TagTableName {
			t.Errorf("Expected tableName to be %v, got %v.", services.TagTableName, db.TableName)
		}

		if res != services.TagPublished {
			t.Errorf("Expected %v, got %v.", services.TagPublished, res)
		}
	})

	t.Run("when UpdateItem fails", func(t *testing.T) {
		db := &mocks.MockDynamoDB{UpdateItemError: errors.New("some error")}

		svc := &services.TagManagerService{Db: db}

		req := &services.PublishTagRequest{
			Module: &terrarium.Module{Name: "test", Version: "v1.0.0"}, 
			Tags: map[string]*services.TagList{
				"EKS": {
				  Tag: "eks-bundle-tag",
				},
			},
			}
		res, err := svc.PublishTag(context.TODO(), req)

		if res != nil {
			t.Errorf("Expected no response, got %v", res)
		}

		if db.UpdateItemInvocations != 1 {
			t.Errorf("Expected 1 call to UpdateItem, got %v", db.UpdateItemInvocations)
		}

		if db.TableName != services.TagTableName {
			t.Errorf("Expected tableName to be %v, got %v.", services.TagTableName, db.TableName)
		}

		if err != services.PublishModuleTagError {
			t.Errorf("Expected %v, got %v.", services.PublishModuleVersionError, err)
		}
	})
}