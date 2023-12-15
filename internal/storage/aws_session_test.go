package storage

import (
	"context"
	"testing"
)

// Test_NewAwsSession checks:
// - if it returns AWS session with provided API key, Secret and Region
func Test_NewAwsSession(t *testing.T) {

	t.Run("returns AWS session`", func(t *testing.T) {
		sessionConfig := AWSSessionConfig{
			UseLocalStack: false,
			Region:        "eu-west-1",
			Key:           "test_key",
			Secret:        "test_secret",
		}

		cfg, err := NewAwsSession(sessionConfig)

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}

		if cfg.Region != sessionConfig.Region {
			t.Errorf("Expected %v, got %v.", sessionConfig.Region, cfg.Region)
		}
		creds, _ := cfg.Credentials.Retrieve(context.TODO())

		if creds.AccessKeyID != sessionConfig.Key {
			t.Errorf("Expected %v, got %v.", sessionConfig.Key, creds.AccessKeyID)

			if creds.SecretAccessKey != sessionConfig.Secret {
				t.Errorf("Expected %v, got %v.", sessionConfig.Secret, creds.SecretAccessKey)
			}
		}
	})

	t.Run("use localstack", func(t *testing.T) {
		sessionConfig := AWSSessionConfig{
			UseLocalStack: true,
			Region:        "eu-west-1",
			Key:           "super_secret_key",
			Secret:        "super_secret_secret_key",
		}

		cfg, err := NewAwsSession(sessionConfig)

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}

		if cfg.Region != sessionConfig.Region {
			t.Errorf("Expected %v, got %v.", sessionConfig.Region, cfg.Region)
		}
		creds, _ := cfg.Credentials.Retrieve(context.TODO())

		if creds.AccessKeyID != "test" {
			t.Errorf("Expected %v, got %v.", sessionConfig.Key, creds.AccessKeyID)

			if creds.SecretAccessKey != "test" {
				t.Errorf("Expected %v, got %v.", sessionConfig.Secret, creds.SecretAccessKey)
			}
		}
	})
}
