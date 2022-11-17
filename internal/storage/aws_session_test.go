package storage_test

import (
	"testing"

	"github.com/terrariumcloud/terrarium/internal/storage"
)

//TestNewAwsSession checks:
// - if it returns AWS session with provided API key, Secret and Region
func Test_NewAwsSession(t *testing.T) {

	t.Run("returns AWS session`", func(t *testing.T) {
		var key, secret, region string
		key = "test_key"
		secret = "test_secret"
		region = "eu-west-1"

		sess, err := storage.NewAwsSession(key, secret, region)

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}

		if *sess.Config.Region != region {
			t.Errorf("Expected %v, got %v.", region, sess.Config.Region)
		}

		creds, _ := sess.Config.Credentials.Get()

		if creds.AccessKeyID != key {
			t.Errorf("Expected %v, got %v.", key, creds.AccessKeyID)

			if creds.SecretAccessKey != secret {
				t.Errorf("Expected %v, got %v.", secret, creds.SecretAccessKey)
			}
		}
	})
}
