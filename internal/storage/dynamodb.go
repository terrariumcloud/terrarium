package storage

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBTableCreator interface {
	DescribeTable(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error)
	CreateTable(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error)
	Scan(ctx context.Context, in *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	GetItem(ctx context.Context, in *dynamodb.GetItemInput, opsFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, in *dynamodb.PutItemInput, opsFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	UpdateItem(ctx context.Context, in *dynamodb.UpdateItemInput, opsFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
	DeleteItem(ctx context.Context, in *dynamodb.DeleteItemInput, opsFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	Query(ctx context.Context, in *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}

// Create new DynamoDB client
func NewDynamoDbClient(key string, secret string, region string) *dynamodb.Client {
	cfg, err := NewAwsSession(key, secret, region)
	if err != nil {
		log.Fatalf("Unable to create AWS Session: %s", err.Error())
	}
	return dynamodb.NewFromConfig(*cfg)
}

// InitializeDynamoDb - checks if table exists, in case it doesn't it creates it
func InitializeDynamoDb(tableName string, schema *dynamodb.CreateTableInput, db DynamoDBTableCreator) error {
	if _, err := db.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{TableName: aws.String(tableName)}); err != nil {
		if err != nil {
			log.Printf("Creating DynamoDB Table: %s", tableName)
			_, err = db.CreateTable(context.TODO(), schema)
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
