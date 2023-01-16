package mocks

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDB struct {
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

func (mdb *DynamoDB) DescribeTable(_ context.Context, in *dynamodb.DescribeTableInput, _ ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
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
func (mdb *DynamoDB) CreateTable(_ context.Context, in *dynamodb.CreateTableInput, _ ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
	mdb.CreateTableInvocations++
	mdb.Schema = in
	return mdb.CreateTableOut, mdb.CreateTableError
}

func (mdb *DynamoDB) GetItem(_ context.Context, in *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
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

func (mdb *DynamoDB) PutItem(_ context.Context, in *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	mdb.PutItemInvocations++
	mdb.TableName = *in.TableName
	return mdb.PutItemOut, mdb.PutItemError
}

func (mdb *DynamoDB) UpdateItem(_ context.Context, in *dynamodb.UpdateItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	mdb.UpdateItemInvocations++
	mdb.TableName = *in.TableName
	return mdb.UpdateItemOut, mdb.UpdateItemError

}

func (mdb *DynamoDB) DeleteItem(_ context.Context, in *dynamodb.DeleteItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	mdb.DeleteItemInvocations++
	mdb.TableName = *in.TableName
	return mdb.DeleteItemOut, mdb.DeleteItemError
}

func (mdb *DynamoDB) Scan(_ context.Context, in *dynamodb.ScanInput, _ ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	panic("Not implemented")
}
