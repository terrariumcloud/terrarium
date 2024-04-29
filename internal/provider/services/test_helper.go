package services

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"testing"
)

func MustMarshallString(value string, t *testing.T) types.AttributeValue {
	t.Helper()
	if result, err := attributevalue.Marshal(value); err != nil {
		t.Fatal("Failed to serialized string to AWS SDK AttributeValue")
		return nil
	} else {
		return result
	}
}
