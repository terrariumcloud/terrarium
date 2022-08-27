package mocks

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type MockDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	DescribeTableInvocations int
	TableName                string
	DescribeTableOut         *dynamodb.DescribeTableOutput
	DescribeTableError       error
	CreateTableInvocations   int
	Schema                   *dynamodb.CreateTableInput
	CreateTableOut           *dynamodb.CreateTableOutput
	CreateTableError         error
	PutItemInvocations       int
	PutItemOut               *dynamodb.PutItemOutput
	PutItemError             error
}

func (mdb *MockDynamoDB) DescribeTable(in *dynamodb.DescribeTableInput) (*dynamodb.DescribeTableOutput, error) {
	mdb.DescribeTableInvocations++
	mdb.TableName = *in.TableName
	return mdb.DescribeTableOut, mdb.DescribeTableError
}

func (fd *MockDynamoDB) CreateTable(in *dynamodb.CreateTableInput) (*dynamodb.CreateTableOutput, error) {
	fd.CreateTableInvocations++
	fd.Schema = in
	return fd.CreateTableOut, fd.CreateTableError
}

func (mdb *MockDynamoDB) PutItem(in *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	mdb.PutItemInvocations++
	mdb.TableName = *in.TableName
	return mdb.PutItemOut, mdb.PutItemError
}
