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
	GetItemInvocations       int
	GetItemOut               *dynamodb.GetItemOutput
	GetItemError             error
	PutItemInvocations       int
	PutItemOut               *dynamodb.PutItemOutput
	PutItemError             error
	UpdateItemInvocations    int
	UpdateItemOut            *dynamodb.UpdateItemOutput
	UpdateItemError          error
	DeleteItemInvocations    int
	DeleteItemOut            *dynamodb.DeleteItemOutput
	DeleteItemError          error
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

func (mdb *MockDynamoDB) GetItem(in *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	mdb.GetItemInvocations++
	mdb.TableName = *in.TableName
	return mdb.GetItemOut, mdb.GetItemError
}

func (mdb *MockDynamoDB) PutItem(in *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	mdb.PutItemInvocations++
	mdb.TableName = *in.TableName
	return mdb.PutItemOut, mdb.PutItemError
}

func (mdb *MockDynamoDB) UpdateItem(in *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	mdb.UpdateItemInvocations++
	mdb.TableName = *in.TableName
	return mdb.UpdateItemOut, mdb.UpdateItemError

}

func (mdb *MockDynamoDB) DeleteItem(in *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	mdb.DeleteItemInvocations++
	mdb.TableName = *in.TableName
	return mdb.DeleteItemOut, mdb.DeleteItemError
}
