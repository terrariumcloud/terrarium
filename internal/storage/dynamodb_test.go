package storage_test

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/terrariumcloud/terrarium/internal/mocks"
	"github.com/terrariumcloud/terrarium/internal/storage"
)

// TestInitializeDynamoDb checks:
// - if table is not recreated when it already exists
// - if table is created when it does not exist
// - if error is returned when checking for table existence fails
// - if error is returned when create table fails
func Test_InitializeDynamoDb(t *testing.T) {
	t.Parallel()

	t.Run("when table exists", func(t *testing.T) {
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
	})

	t.Run("when table does not exists", func(t *testing.T) {
		table := "Test"
		schema := &dynamodb.CreateTableInput{}
		db := &mocks.MockDynamoDB{
			DescribeTableErrors: []error{&dynamodb.ResourceNotFoundException{}},
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
	})

	t.Run("when checking for table existence fails", func(t *testing.T) {
		table := "Test"
		schema := &dynamodb.CreateTableInput{}
		someError := errors.New("some error")
		db := &mocks.MockDynamoDB{
			DescribeTableErrors: []error{errors.New("some error")},
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
	})

	t.Run("when create table fails", func(t *testing.T) {
		table := "Test"
		schema := &dynamodb.CreateTableInput{}
		someError := errors.New("some error")
		db := &mocks.MockDynamoDB{
			DescribeTableErrors: []error{&dynamodb.ResourceNotFoundException{}},
			CreateTableError:    someError,
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
	})
}
