package session

import (
	"context"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/services"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium"
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

func TestBeginVersion(t *testing.T) {
	t.Run("It creates entry in DynamoDB", func(t *testing.T) {
		fd := &fakeDynamoDB{}

		sessionService := &SessionService{
			Db: fd,
		}
		request := services.BeginVersionRequest{
			Module: &terrarium.VersionedModule{
				Name:    "test",
				Version: "v1.0.0",
			},
		}
		response, err := sessionService.BeginVersion(context.TODO(), &request)

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
			if *fd.tableName != DefaultSessionTableName {
				t.Errorf("Expected tableName to be %s, got %s", DefaultSessionTableName, *fd.tableName)
			}
		}
	})
}

func TestBeginVersionE2E(t *testing.T) {
	t.Run("It creates entry in DynamoDB", func(t *testing.T) {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		svc := dynamodb.New(sess)

		sessionService := &SessionService{
			Db: svc,
		}
		request := services.BeginVersionRequest{
			Module: &terrarium.VersionedModule{
				Name:    "test",
				Version: "v1.0.0",
			},
		}
		response, err := sessionService.BeginVersion(context.TODO(), &request)

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
