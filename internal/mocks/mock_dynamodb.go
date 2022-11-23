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
	DescribeTableErrors      []error
	CreateTableInvocations   int
	Schema                   *dynamodb.CreateTableInput
	CreateTableOut           *dynamodb.CreateTableOutput
	CreateTableError         error
	GetItemInvocations       int
	GetItemOuts              []*dynamodb.GetItemOutput
	GetItemErrors            []error
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
	var err error = nil
	if len(mdb.DescribeTableErrors) > mdb.DescribeTableInvocations {
		err = mdb.DescribeTableErrors[mdb.DescribeTableInvocations]
	}
	//} else {
	//	panic("Not enough errors for call to DescribeTable - Invalid Test")
	//}
	mdb.DescribeTableInvocations++
	mdb.TableName = *in.TableName
	return mdb.DescribeTableOut, err
}

func (mdb *MockDynamoDB) CreateTable(in *dynamodb.CreateTableInput) (*dynamodb.CreateTableOutput, error) {
	mdb.CreateTableInvocations++
	mdb.Schema = in
	return mdb.CreateTableOut, mdb.CreateTableError
}

func (mdb *MockDynamoDB) GetItem(in *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	var err error = nil
	if len(mdb.GetItemErrors) > mdb.GetItemInvocations {
		err = mdb.GetItemErrors[mdb.GetItemInvocations]
	}

	var out *dynamodb.GetItemOutput = nil
	if len(mdb.GetItemOuts) > mdb.GetItemInvocations {
		out = mdb.GetItemOuts[mdb.GetItemInvocations]
	}

	mdb.GetItemInvocations++
	mdb.TableName = *in.TableName
	return out, err
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
