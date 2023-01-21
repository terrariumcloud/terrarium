package storage_test

import (
	"context"
	"github.com/terrariumcloud/terrarium/internal/storage"
	"testing"
)

// Test_NewAwsSession checks:
// - if it returns AWS session with provided API key, Secret and Region
func Test_NewAwsSession(t *testing.T) {

	t.Run("returns AWS session`", func(t *testing.T) {
		var key, secret, region string
		key = "test_key"
		secret = "test_secret"
		region = "eu-west-1"

		cfg, err := storage.NewAwsSession(key, secret, region)

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}

		if cfg.Region != region {
			t.Errorf("Expected %v, got %v.", region, cfg.Region)
		}
		creds, _ := cfg.Credentials.Retrieve(context.TODO())

		if creds.AccessKeyID != key {
			t.Errorf("Expected %v, got %v.", key, creds.AccessKeyID)

			if creds.SecretAccessKey != secret {
				t.Errorf("Expected %v, got %v.", secret, creds.SecretAccessKey)
			}
		}
	})
}
