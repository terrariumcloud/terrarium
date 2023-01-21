package storage_test

import (
	"errors"
	"github.com/terrariumcloud/terrarium/internal/storage/mocks"
	"testing"

	"github.com/terrariumcloud/terrarium/internal/storage"
)

// Test_InitializeS3Bucket checks:
// - if bucket is not recreated when it already exists
// - if bucket is created when it does not exist
// - if error is returned when create bucket fails
func Test_InitializeS3Bucket(t *testing.T) {
	t.Parallel()

	t.Run("when bucket exists", func(t *testing.T) {
		bucket := "Test"
		region := "test"
		s3 := &mocks.S3{}

		err := storage.InitializeS3Bucket(bucket, region, s3)

		if s3.HeadBucketInvocations != 1 {
			t.Errorf("Expected 1 call to HeadBucket, got %v.", s3.HeadBucketInvocations)
		}

		if s3.BucketName != bucket {
			t.Errorf("Expected %v, got %v.", bucket, s3.BucketName)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when bucket does not exists", func(t *testing.T) {
		bucket := "Test"
		region := "test"
		s3Client := &mocks.S3{HeadBucketError: errors.New("some error")}

		err := storage.InitializeS3Bucket(bucket, region, s3Client)

		if s3Client.HeadBucketInvocations != 1 {
			t.Errorf("Expected 1 call to HeadBucket, got %v.", s3Client.HeadBucketInvocations)
		}

		if s3Client.BucketName != bucket {
			t.Errorf("Expected %v, got %v.", bucket, s3Client.BucketName)
		}

		if s3Client.CreateBucketInvocations != 1 {
			t.Errorf("Expected 1 call to CreateTable, got %v.", s3Client.CreateBucketInvocations)
		}

		if s3Client.Region != region {
			t.Errorf("Expected %v, got %v.", region, s3Client.Region)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}
	})

	t.Run("when create bucket fails", func(t *testing.T) {
		bucket := "Test"
		region := "test"
		someError := errors.New("some error")
		s3Client := &mocks.S3{HeadBucketError: someError, CreateBucketError: someError}

		err := storage.InitializeS3Bucket(bucket, region, s3Client)

		if s3Client.HeadBucketInvocations != 1 {
			t.Errorf("Expected 1 call to HeadBucket, got %v.", s3Client.HeadBucketInvocations)
		}

		if s3Client.BucketName != bucket {
			t.Errorf("Expected %v, got %v.", bucket, s3Client.BucketName)
		}

		if s3Client.CreateBucketInvocations != 1 {
			t.Errorf("Expected 1 call to CreateTable, got %v.", s3Client.CreateBucketInvocations)
		}

		if s3Client.Region != region {
			t.Errorf("Expected %v, got %v.", region, s3Client.Region)
		}

		if err == nil {
			t.Error("Expected error, got nil.")
		}

		if err.Error() != someError.Error() {
			t.Errorf("Expected %v, got %v.", someError.Error(), err.Error())
		}
	})
}
