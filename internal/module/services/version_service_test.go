package services

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"
)

// type fakeDynamoDB struct {
// 	dynamodbiface.DynamoDBAPI
// 	err                  error
// 	numberOfPutItemCalls int
// 	tableName            *string
// }

// func (fd *fakeDynamoDB) PutItem(item *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
// 	fd.tableName = item.TableName
// 	output := new(dynamodb.PutItemOutput)
// 	output.Attributes = make(map[string]*dynamodb.AttributeValue)
// 	fd.numberOfPutItemCalls++
// 	return output, fd.err
// }

func TestBeginVersion(t *testing.T) {
	t.Run("It creates entry in DynamoDB", func(t *testing.T) {
		fd := &fakeDynamoDB{}

		versionService := &VersionService{
			Db: fd,
		}
		request := BeginVersionRequest{
			Module: &terrarium.VersionedModule{
				Name:    "test",
				Version: "v1.0.0",
			},
		}
		response, err := versionService.BeginVersion(context.TODO(), &request)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response == nil {
			t.Errorf("Expected response, got nil.")
		} else {
			if response.GetSessionKey() == "" {
				t.Errorf("Expected session key to not be an empty string, got empty string")
			}
		}

		if fd.numberOfPutItemCalls != 1 {
			t.Errorf("Expected number of calls to PutItem to be %d, got %d", 1, fd.numberOfPutItemCalls)
		}

		if fd.tableName == nil {
			t.Errorf("Expected tableName, got nil.")
		} else {
			if *fd.tableName != DefaultVersionTableName {
				t.Errorf("Expected tableName to be %s, got %s", DefaultVersionTableName, *fd.tableName)
			}
		}
	})
}

func IgnoreTestBeginVersionE2E(t *testing.T) {
	t.Run("It creates entry in DynamoDB", func(t *testing.T) {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		svc := dynamodb.New(sess)

		versionService := &VersionService{
			Db: svc,
		}
		request := BeginVersionRequest{
			Module: &terrarium.VersionedModule{
				Name:    "test",
				Version: "v1.0.0",
			},
		}
		response, err := versionService.BeginVersion(context.TODO(), &request)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if response == nil {
			t.Errorf("Expected response, got nil.")
		} else {
			if response.GetSessionKey() == "" {
				t.Errorf("Expected session key to not be an empty string, got empty string")
			}
		}
	})
}
