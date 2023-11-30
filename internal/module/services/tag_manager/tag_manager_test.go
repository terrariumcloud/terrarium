package tag_manager

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/terrariumcloud/terrarium/internal/module/services"
	"github.com/terrariumcloud/terrarium/internal/storage/mocks"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	terrarium "github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"google.golang.org/grpc"
)

func Test_RegisterTagManagerWithServer(t *testing.T) {
	t.Parallel()

	t.Run("when there is no error with table init", func(t *testing.T) {
		db := &mocks.DynamoDB{}

		tm := &TagManagerService{
			Db: db,
		}

		s := grpc.NewServer(*new([]grpc.ServerOption)...)

		err := tm.RegisterWithServer(s)

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

		db := &mocks.DynamoDB{
			DescribeTableErrors: []error{errors.New("some error")},
			CreateTableError:    errors.New("some error"),
		}

		tm := &TagManagerService{
			Db: db,
		}

		s := grpc.NewServer(*new([]grpc.ServerOption)...)

		err := tm.RegisterWithServer(s)

		if err != ModuleTagTableInitializationError {
			t.Errorf("Expected %v, got %v.", ModuleTagTableInitializationError, err)
		}

		if db.DescribeTableInvocations != 1 {
			t.Errorf("Expected 1 call to DescribeTable, got %v.", db.DescribeTableInvocations)
		}

		if db.CreateTableInvocations != 1 {
			t.Errorf("Expected 1 calls to CreateTable, got %v.", db.CreateTableInvocations)
		}
	})
}

func Test_PublishTag(t *testing.T) {
	t.Parallel()

	t.Run("when tag is published", func(t *testing.T) {
		db := &mocks.DynamoDB{
			GetItemOuts: []*dynamodb.GetItemOutput{{}},
		}

		svc := &TagManagerService{Db: db}

		listOfTags := []string{"eks"}
		req := terrarium.PublishTagRequest{
			Name:   "test",
			ApiKey: "test desc",
			Tags:   listOfTags,
		}

		res, err := svc.PublishTag(context.TODO(), &req)

		if err != nil {
			t.Errorf("Expected %v, got %v.", nil, err)
		}

		if res != TagPublished {
			t.Errorf("Expected %v, got %v.", TagPublished, res)
		}
	})

	t.Run("when UpdateItem is successful", func(t *testing.T) {
		name := "test"
		tagsList := []string{"eks"}
		marshaledTags, _ := attributevalue.Marshal(tagsList)

		db := &mocks.DynamoDB{
			GetItemOuts: []*dynamodb.GetItemOutput{{
				Item: map[string]types.AttributeValue{
					"name": services.MustMarshallString(name, t),
					"tags": marshaledTags,
				},
			},
			},
			UpdateItemOut: &dynamodb.UpdateItemOutput{},
		}

		svc := &TagManagerService{Db: db}

		listOfTags := []string{"eks", "eks1"}
		req := terrarium.PublishTagRequest{
			Name:   "test",
			ApiKey: "test desc",
			Tags:   listOfTags,
		}

		res, err := svc.PublishTag(context.TODO(), &req)

		if err != nil {
			t.Errorf("Expected %v, got %v.", nil, err)
		}

		if res != TagPublished {
			t.Errorf("Expected %v, got %v.", TagPublished, res)
		}
	})

	t.Run("when UpdateItem is not successful", func(t *testing.T) {
		name := "test"
		tagsList := []string{"eks"}
		marshaledTags, _ := attributevalue.Marshal(tagsList)
		db := &mocks.DynamoDB{
			GetItemOuts: []*dynamodb.GetItemOutput{{
				Item: map[string]types.AttributeValue{
					"name": services.MustMarshallString(name, t),
					"tags": marshaledTags,
				},
			},
			},
			UpdateItemOut:   &dynamodb.UpdateItemOutput{},
			UpdateItemError: errors.New("Failed to update module tag."),
		}

		svc := &TagManagerService{Db: db}

		listOfTags := []string{"eks"}
		req := terrarium.PublishTagRequest{
			Name:   "test",
			ApiKey: "test desc",
			Tags:   listOfTags,
		}

		res, err := svc.PublishTag(context.TODO(), &req)

		if res != nil {
			t.Errorf("Expected no response, got %v", err)
		}

		if db.UpdateItemError == nil {
			t.Errorf("Expected failed update, got %d", db.UpdateItemError)
		}

		if err != UpdateModuleTagError {
			t.Errorf("Expected %v, got %v.", UpdateModuleTagError, err)
		}
	})
}
