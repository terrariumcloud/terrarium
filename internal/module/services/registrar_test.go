package services_test

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"testing"

	"github.com/terrariumcloud/terrarium/internal/mocks"
	"github.com/terrariumcloud/terrarium/internal/module/services"
	terrarium "github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"google.golang.org/grpc"
)

// Test_RegisterModule checks:
// - if correct response is returned when module is registered
// - if there was no error when version already exists
// - if error was returned when GetItem fails
// - if error is returned when marshal fails
// - if error is returned when PutItem fails
func Test_RegisterModule(t *testing.T) {
	t.Parallel()

	t.Run("when new version is created", func(t *testing.T) {
		db := &mocks.MockDynamoDB{
			GetItemOut: &dynamodb.GetItemOutput{},
		}

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

		if db.GetItemInvocations != 1 {
			t.Errorf("Expected 1 call to GetItem, got %d", db.GetItemInvocations)
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
	})

	t.Run("when version already exists", func(t *testing.T) {
		name := "test"
		emptyString := ""
		db := &mocks.MockDynamoDB{
			GetItemOut: &dynamodb.GetItemOutput{
				Item: map[string]*dynamodb.AttributeValue{
					"name":        {S: &name},
					"description": {S: &emptyString},
					"source_url":  {S: &emptyString},
					"maturity":    {S: &emptyString},
					"created_on":  {S: &emptyString},
					"modified_on": {S: &emptyString},
				},
			},
			UpdateItemOut: &dynamodb.UpdateItemOutput{},
		}

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

		if db.GetItemInvocations != 1 {
			t.Errorf("Expected 1 call to GetItem, got %d", db.GetItemInvocations)
		}

		if db.PutItemInvocations != 0 {
			t.Errorf("Expected 0 calls to PutItem, got %d", db.PutItemInvocations)
		}

		if db.UpdateItemInvocations != 1 {
			t.Errorf("Expected 1 call to UpdateItem, got %d", db.UpdateItemInvocations)
		}

		if db.TableName != services.RegistrarTableName {
			t.Errorf("Expected tableName to be %s, got %s", services.RegistrarTableName, db.TableName)
		}

		if res != services.ModuleRegistered {
			t.Errorf("Expected %v, got %v.", services.ModuleRegistered, res)
		}
	})

	t.Run("when GetItem fails", func(t *testing.T) {
		db := &mocks.MockDynamoDB{
			GetItemErrors: []error{errors.New("some error")},
		}

		svc := &services.RegistrarService{Db: db}

		req := terrarium.RegisterModuleRequest{
			Name:        "test",
			Description: "test desc",
			Source:      "http://test.com",
			Maturity:    terrarium.Maturity_ALPHA,
		}

		res, err := svc.Register(context.TODO(), &req)

		if err == nil {
			t.Errorf("Expected an error")
		}

		if db.GetItemInvocations != 1 {
			t.Errorf("Expected 1 call to GetItem, got %d", db.GetItemInvocations)
		}

		if db.PutItemInvocations != 0 {
			t.Errorf("Expected 0 calls to PutItem, got %d", db.PutItemInvocations)
		}

		if res != nil {
			t.Errorf("Expected no response, got %v.", res)
		}
	})

	// t.Run("when marshal fails", func(t *testing.T) {
	// 	db := &mocks.MockDynamoDB{}

	// 	svc := &services.RegistrarService{Db: db}

	// 	req := &terrarium.RegisterModuleRequest{} // TODO: need to make MarshalMap return error

	// 	res, err := svc.Register(context.TODO(), req)

	// 	if err != services.MarshalModuleError {
	// 		t.Errorf("Expected %v, got %v.", services.MarshalModuleError, err)
	// 	}

	// 	if db.PutItemInvocations == 0 {
	// 		t.Errorf("Expected 0 calls to PutItem, got %d", db.PutItemInvocations)
	// 	}

	// 	if res != nil {
	// 		t.Errorf("Expected no response, got %v.", res)
	// 	}
	// })

	t.Run("when PutItem fails", func(t *testing.T) {
		db := &mocks.MockDynamoDB{
			GetItemOuts:   []*dynamodb.GetItemOutput{{}},
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
	})
}

// Test_RegisterRegistrarWithServer checks:
// - if there was no error with table init
// - if error is returned when Table initialization fails
func Test_RegisterRegistrarWithServer(t *testing.T) {
	t.Parallel()

	t.Run("when there is no error with table init", func(t *testing.T) {
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
	})

	t.Run("when Table initialization fails", func(t *testing.T) {

		db := &mocks.MockDynamoDB{
			DescribeTableErrors: []error{errors.New("some error")},
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
	})
}
