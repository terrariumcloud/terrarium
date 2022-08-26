package services_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	services "github.com/terrariumcloud/terrarium-grpc-gateway/internal/module/services"
	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"
)

func TestBeginVersion(t *testing.T) {
	t.Parallel()

	fd := &fakeDynamoDB{}

	versionService := &services.VersionManagerService{
		Db: fd,
	}
	request := terrarium.BeginVersionRequest{
		Module: &terrarium.Module{
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
	}

	if fd.numberOfPutItemCalls != 1 {
		t.Errorf("Expected number of calls to PutItem to be %d, got %d", 1, fd.numberOfPutItemCalls)
	}

	if fd.tableName == nil {
		t.Errorf("Expected tableName, got nil.")
	} else {
		if *fd.tableName != services.DefaultVersionsTableName {
			t.Errorf("Expected tableName to be %s, got %s", services.DefaultVersionsTableName, *fd.tableName)
		}
	}
}

func IgnoreTestBeginVersionE2E(t *testing.T) {
	t.Parallel()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	versionService := &services.VersionManagerService{
		Db: svc,
	}
	request := terrarium.BeginVersionRequest{
		Module: &terrarium.Module{
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
	}
}
