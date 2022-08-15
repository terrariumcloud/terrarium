package storage

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// Create new DynamoDB client
func NewDynamoDbClient(key string, secret string, region string) dynamodbiface.DynamoDBAPI {
	sess, err := NewAwsSession(key, secret, region)
	if err != nil {
		log.Fatalf("Unable to create AWS Session: %s", err.Error())
	}
	return dynamodb.New(sess)
}

// InitializeDynamoDb - checks if table exists, in case it doesn't it creates it
func InitializeDynamoDb(tableName string, schema *dynamodb.CreateTableInput, db dynamodbiface.DynamoDBAPI) error {
	if _, err := db.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}); err != nil {
		var notFoundErr *dynamodb.ResourceNotFoundException
		if errors.As(err, &notFoundErr) {
			log.Printf("Creating DynamoDB Table: %s", tableName)
			_, err = db.CreateTable(schema)
			if err != nil {
				return err
			}
			log.Println("Table created.")
			return nil
		}
		return err
	}
	log.Printf("DynamoDB Table %s already exists.", tableName)
	return nil
}
