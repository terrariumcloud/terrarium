package storage

import (
	"log"

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
