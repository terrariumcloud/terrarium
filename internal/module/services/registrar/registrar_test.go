package registrar

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/terrariumcloud/terrarium/internal/module/services"
	"github.com/terrariumcloud/terrarium/internal/storage/mocks"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

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
		db := &mocks.DynamoDB{
			GetItemOuts: []*dynamodb.GetItemOutput{
				{}, // First call for API key validation (no API key configured)
				{}, // Second call for module check
			},
		}

		svc := &RegistrarService{Db: db}

		req := terrarium.RegisterModuleRequest{
			ApiKey:      "test-api-key",
			Name:        "test",
			Description: "test desc",
			Source:      "http://test.com",
			Maturity:    terrarium.Maturity_ALPHA,
		}

		res, err := svc.Register(context.TODO(), &req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if db.GetItemInvocations != 2 {
			t.Errorf("Expected 2 calls to GetItem (API key validation + module check), got %d", db.GetItemInvocations)
		}

		if db.PutItemInvocations != 1 {
			t.Errorf("Expected 1 call to PutItem, got %d", db.PutItemInvocations)
		}

		if db.TableName != RegistrarTableName {
			t.Errorf("Expected tableName to be %s, got %s", RegistrarTableName, db.TableName)
		}

		if res != ModuleRegistered {
			t.Errorf("Expected %v, got %v.", ModuleRegistered, res)
		}
	})

	t.Run("when version already exists", func(t *testing.T) {
		name := "test"
		emptyString := ""
		db := &mocks.DynamoDB{
			GetItemOuts: []*dynamodb.GetItemOutput{
				{}, // First call for API key validation (no API key configured)
				{ // Second call for module check - module exists
					Item: map[string]types.AttributeValue{
						"name":        services.MustMarshallString(name, t),
						"description": services.MustMarshallString(emptyString, t),
						"source_url":  services.MustMarshallString(emptyString, t),
						"maturity":    services.MustMarshallString(emptyString, t),
						"created_on":  services.MustMarshallString(emptyString, t),
						"modified_on": services.MustMarshallString(emptyString, t),
					},
				},
			},
			UpdateItemOut: &dynamodb.UpdateItemOutput{},
		}

		svc := &RegistrarService{Db: db}

		req := terrarium.RegisterModuleRequest{
			ApiKey:      "test-api-key",
			Name:        "test",
			Description: "test desc",
			Source:      "http://test.com",
			Maturity:    terrarium.Maturity_ALPHA,
		}

		res, err := svc.Register(context.TODO(), &req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if db.GetItemInvocations != 2 {
			t.Errorf("Expected 2 calls to GetItem (API key validation + module check), got %d", db.GetItemInvocations)
		}

		if db.PutItemInvocations != 0 {
			t.Errorf("Expected 0 calls to PutItem, got %d", db.PutItemInvocations)
		}

		if db.UpdateItemInvocations != 1 {
			t.Errorf("Expected 1 call to UpdateItem, got %d", db.UpdateItemInvocations)
		}

		if db.TableName != RegistrarTableName {
			t.Errorf("Expected tableName to be %s, got %s", RegistrarTableName, db.TableName)
		}

		if res != ModuleRegistered {
			t.Errorf("Expected %v, got %v.", ModuleRegistered, res)
		}
	})

	t.Run("when GetItem fails", func(t *testing.T) {
		db := &mocks.DynamoDB{
			GetItemErrors: []error{errors.New("some error")},
		}

		svc := &RegistrarService{Db: db}

		req := terrarium.RegisterModuleRequest{
			ApiKey:      "test-api-key",
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

	t.Run("when PutItem fails", func(t *testing.T) {
		db := &mocks.DynamoDB{
			GetItemOuts: []*dynamodb.GetItemOutput{
				{}, // First call for API key validation (no API key configured)
				{}, // Second call for module check
			},
			PutItemError: errors.New("some error"),
		}

		svc := &RegistrarService{Db: db}

		req := terrarium.RegisterModuleRequest{
			ApiKey:      "test-api-key",
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

		if db.TableName != RegistrarTableName {
			t.Errorf("Expected tableName to be %s, got %s", RegistrarTableName, db.TableName)
		}

		if err != ModuleRegisterError {
			t.Errorf("Expected %v, got %v.", ModuleRegisterError, err)
		}
	})

	t.Run("when API key validation succeeds", func(t *testing.T) {
		apiKey := "valid-api-key"
		orgName := "cie"
		db := &mocks.DynamoDB{
			GetItemOuts: []*dynamodb.GetItemOutput{
				{ // First call for API key validation - organization has API key configured
					Item: map[string]types.AttributeValue{
						"organization_name": services.MustMarshallString(orgName, t),
						"api_key":           services.MustMarshallString(apiKey, t),
						"created_on":        services.MustMarshallString("2023-01-01", t),
						"modified_on":       services.MustMarshallString("2023-01-01", t),
					},
				},
				{}, // Second call for module check (new module)
			},
		}

		svc := &RegistrarService{Db: db}

		req := terrarium.RegisterModuleRequest{
			ApiKey:      apiKey,
			Name:        "cie/mymodule/aws",
			Description: "test desc",
			Source:      "http://test.com",
			Maturity:    terrarium.Maturity_ALPHA,
		}

		res, err := svc.Register(context.TODO(), &req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if db.GetItemInvocations != 2 {
			t.Errorf("Expected 2 calls to GetItem (API key validation + module check), got %d", db.GetItemInvocations)
		}

		if db.PutItemInvocations != 1 {
			t.Errorf("Expected 1 call to PutItem, got %d", db.PutItemInvocations)
		}

		if res != ModuleRegistered {
			t.Errorf("Expected %v, got %v.", ModuleRegistered, res)
		}
	})

	t.Run("when no API key configured for organization", func(t *testing.T) {
		db := &mocks.DynamoDB{
			GetItemOuts: []*dynamodb.GetItemOutput{
				{}, // First call for API key validation (no API key configured)
				{}, // Second call for module check (new module)
			},
		}

		svc := &RegistrarService{Db: db}

		req := terrarium.RegisterModuleRequest{
			ApiKey:      "any-api-key",
			Name:        "cie/mymodule/aws",
			Description: "test desc",
			Source:      "http://test.com",
			Maturity:    terrarium.Maturity_ALPHA,
		}

		res, err := svc.Register(context.TODO(), &req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if db.GetItemInvocations != 2 {
			t.Errorf("Expected 2 calls to GetItem (API key validation + module check), got %d", db.GetItemInvocations)
		}

		if res != ModuleRegistered {
			t.Errorf("Expected %v, got %v.", ModuleRegistered, res)
		}
	})

	t.Run("when API key validation fails", func(t *testing.T) {
		correctApiKey := "correct-api-key"
		wrongApiKey := "wrong-api-key"
		orgName := "cie"
		db := &mocks.DynamoDB{
			GetItemOuts: []*dynamodb.GetItemOutput{
				{ // First call for API key validation - organization has API key configured
					Item: map[string]types.AttributeValue{
						"organization_name": services.MustMarshallString(orgName, t),
						"api_key":           services.MustMarshallString(correctApiKey, t),
						"created_on":        services.MustMarshallString("2023-01-01", t),
						"modified_on":       services.MustMarshallString("2023-01-01", t),
					},
				},
			},
		}

		svc := &RegistrarService{Db: db}

		req := terrarium.RegisterModuleRequest{
			ApiKey:      wrongApiKey,
			Name:        "cie/mymodule/aws",
			Description: "test desc",
			Source:      "http://test.com",
			Maturity:    terrarium.Maturity_ALPHA,
		}

		res, err := svc.Register(context.TODO(), &req)

		if err != ApiKeyValidationError {
			t.Errorf("Expected %v, got %v", ApiKeyValidationError, err)
		}

		if db.GetItemInvocations != 1 {
			t.Errorf("Expected 1 call to GetItem (only API key validation), got %d", db.GetItemInvocations)
		}

		if db.PutItemInvocations != 0 {
			t.Errorf("Expected 0 calls to PutItem, got %d", db.PutItemInvocations)
		}

		if res != nil {
			t.Errorf("Expected no response, got %v.", res)
		}
	})

	t.Run("when API key retrieval fails", func(t *testing.T) {
		db := &mocks.DynamoDB{
			GetItemErrors: []error{errors.New("database error")},
		}

		svc := &RegistrarService{Db: db}

		req := terrarium.RegisterModuleRequest{
			ApiKey:      "test-api-key",
			Name:        "cie/mymodule/aws",
			Description: "test desc",
			Source:      "http://test.com",
			Maturity:    terrarium.Maturity_ALPHA,
		}

		res, err := svc.Register(context.TODO(), &req)

		if err != ApiKeyGetError {
			t.Errorf("Expected %v, got %v", ApiKeyGetError, err)
		}

		if db.GetItemInvocations != 1 {
			t.Errorf("Expected 1 call to GetItem (API key validation), got %d", db.GetItemInvocations)
		}

		if res != nil {
			t.Errorf("Expected no response, got %v.", res)
		}
	})
}

// Test_RegisterRegistrarWithServer checks:
// - if there was no error with table init
// - if error is returned when Table initialization fails
func Test_RegisterRegistrarWithServer(t *testing.T) {
	t.Parallel()

	t.Run("when there is no error with table init", func(t *testing.T) {
		db := &mocks.DynamoDB{}

		rs := &RegistrarService{
			Db: db,
		}

		s := grpc.NewServer(*new([]grpc.ServerOption)...)

		err := rs.RegisterWithServer(s)

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}

		if db.DescribeTableInvocations != 2 {
			t.Errorf("Expected 2 calls to DescribeTable (modules + API keys), got %v.", db.DescribeTableInvocations)
		}

		if db.CreateTableInvocations != 0 {
			t.Errorf("Expected no calls to CreateTable, got %v.", db.CreateTableInvocations)
		}
	})

	t.Run("when Module table initialization fails", func(t *testing.T) {

		db := &mocks.DynamoDB{
			DescribeTableErrors: []error{errors.New("some error")},
			CreateTableError:    errors.New("some error"),
		}

		rs := &RegistrarService{
			Db: db,
		}

		s := grpc.NewServer(*new([]grpc.ServerOption)...)

		err := rs.RegisterWithServer(s)

		if err != ModuleTableInitializationError {
			t.Errorf("Expected %v, got %v.", ModuleTableInitializationError, err)
		}

		if db.DescribeTableInvocations != 1 {
			t.Errorf("Expected 1 call to DescribeTable, got %v.", db.DescribeTableInvocations)
		}

		if db.CreateTableInvocations != 1 {
			t.Errorf("Expected 1 calls to CreateTable, got %v.", db.CreateTableInvocations)
		}
	})

	t.Run("when API key table initialization fails", func(t *testing.T) {

		db := &mocks.DynamoDB{
			DescribeTableErrors: []error{nil, errors.New("some error")}, // First table succeeds, second fails
			CreateTableError:    errors.New("some error"),
		}

		rs := &RegistrarService{
			Db: db,
		}

		s := grpc.NewServer(*new([]grpc.ServerOption)...)

		err := rs.RegisterWithServer(s)

		if err != ApiKeyTableInitializationError {
			t.Errorf("Expected %v, got %v.", ApiKeyTableInitializationError, err)
		}

		if db.DescribeTableInvocations != 2 {
			t.Errorf("Expected 2 calls to DescribeTable, got %v.", db.DescribeTableInvocations)
		}

		if db.CreateTableInvocations != 1 {
			t.Errorf("Expected 1 call to CreateTable (for API key table), got %v.", db.CreateTableInvocations)
		}
	})
}
