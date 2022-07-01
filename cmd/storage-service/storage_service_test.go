package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type fakeDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	err                  error
	numberOfPutItemCalls int
	tableName            *string
}

func (fd *fakeDynamoDB) PutItem(item *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	fd.tableName = item.TableName
	output := new(dynamodb.PutItemOutput)
	output.Attributes = make(map[string]*dynamodb.AttributeValue)
	fd.numberOfPutItemCalls++
	return output, fd.err
}

func TestUploadSourceZip(t *testing.T) {
	t.Run("It creates entry in DynamoDB", func(t *testing.T) {

		// storageService := &StorageService{
		// }
	})
}
