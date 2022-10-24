package storage_test

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/terrariumcloud/terrarium/internal/mocks"
	"github.com/terrariumcloud/terrarium/internal/storage"
)

// This test checks if table is not recreated when it already exists
func TestInitializeDynamoDbWhenTableExists(t *testing.T) {
	t.Parallel()

	table := "Test"
	schema := &dynamodb.CreateTableInput{}
	db := &mocks.MockDynamoDB{}

	err := storage.InitializeDynamoDb(table, schema, db)

	if db.DescribeTableInvocations != 1 {
		t.Errorf("Expected 1 call to DescribeTable, got %v.", db.DescribeTableInvocations)
	}

	if db.TableName != table {
		t.Errorf("Expected %v, got %v.", table, db.TableName)
	}

	if err != nil {
		t.Errorf("Expected no error, got %v.", err)
	}
}

// This test checks if table is created when it does not exist
func TestInitializeDynamoDbWhenTableDoesNotExists(t *testing.T) {
	t.Parallel()

	table := "Test"
	schema := &dynamodb.CreateTableInput{}
	db := &mocks.MockDynamoDB{
		DescribeTableError: &dynamodb.ResourceNotFoundException{},
	}

	err := storage.InitializeDynamoDb(table, schema, db)

	if db.DescribeTableInvocations != 1 {
		t.Errorf("Expected 1 call to DescribeTable, got %v.", db.DescribeTableInvocations)
	}

	if db.TableName != table {
		t.Errorf("Expected %v, got %v.", table, db.TableName)
	}

	if db.CreateTableInvocations != 1 {
		t.Errorf("Expected 1 call to CreateTable, got %v.", db.CreateTableInvocations)
	}

	if db.Schema != schema {
		t.Errorf("Expected %v, got %v.", schema, db.Schema)
	}

	if err != nil {
		t.Errorf("Expected no error, got %v.", err)
	}
}

// This test checks if error is returned when checking for table existence fails
func TestInitializeDynamoDbWhenTableExistsErrors(t *testing.T) {
	t.Parallel()

	table := "Test"
	schema := &dynamodb.CreateTableInput{}
	someError := errors.New("some error")
	db := &mocks.MockDynamoDB{
		DescribeTableError: someError,
	}

	err := storage.InitializeDynamoDb(table, schema, db)

	if db.DescribeTableInvocations != 1 {
		t.Errorf("Expected 1 call to DescribeTable, got %v.", db.DescribeTableInvocations)
	}

	if db.TableName != table {
		t.Errorf("Expected %v, got %v.", table, db.TableName)
	}

	if err == nil {
		t.Error("Expected error, got nil.")
	}

	if err.Error() != someError.Error() {
		t.Errorf("Expected %v, got %v.", someError.Error(), err.Error())
	}
}

// This test checks if error is returned when create table fails
func TestInitializeDynamoDbWhenCreateTableErrors(t *testing.T) {
	t.Parallel()

	table := "Test"
	schema := &dynamodb.CreateTableInput{}
	someError := errors.New("some error")
	db := &mocks.MockDynamoDB{
		DescribeTableError: &dynamodb.ResourceNotFoundException{},
		CreateTableError:   someError,
	}

	err := storage.InitializeDynamoDb(table, schema, db)

	if db.DescribeTableInvocations != 1 {
		t.Errorf("Expected 1 call to DescribeTable, got %v.", db.DescribeTableInvocations)
	}

	if db.TableName != table {
		t.Errorf("Expected %v, got %v.", table, db.TableName)
	}

	if db.CreateTableInvocations != 1 {
		t.Errorf("Expected 1 call to CreateTable, got %v.", db.CreateTableInvocations)
	}

	if db.Schema != schema {
		t.Errorf("Expected %v, got %v.", schema, db.Schema)
	}

	if err == nil {
		t.Error("Expected error, got nil.")
	}

	if err.Error() != someError.Error() {
		t.Errorf("Expected %v, got %v.", someError.Error(), err.Error())
	}
}
