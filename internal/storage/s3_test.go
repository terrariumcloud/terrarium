package storage_test

import (
	"errors"
	"testing"

	"github.com/terrariumcloud/terrarium/internal/mocks"
	"github.com/terrariumcloud/terrarium/internal/storage"
)

// This test checks if bucket is not recreated when it already exists
func TestInitializeS3BucketWhenBucketExists(t *testing.T) {
	t.Parallel()

	bucket := "Test"
	region := "test"
	s3 := &mocks.MockS3{}

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
}

// This test checks if bucket is created when it does not exist
func TestInitializeS3BucketWhenBucketDoesNotExists(t *testing.T) {
	t.Parallel()

	bucket := "Test"
	region := "test"
	s3 := &mocks.MockS3{HeadBucketError: errors.New("some error")}

	err := storage.InitializeS3Bucket(bucket, region, s3)

	if s3.HeadBucketInvocations != 1 {
		t.Errorf("Expected 1 call to HeadBucket, got %v.", s3.HeadBucketInvocations)
	}

	if s3.BucketName != bucket {
		t.Errorf("Expected %v, got %v.", bucket, s3.BucketName)
	}

	if s3.CreateBucketInvocations != 1 {
		t.Errorf("Expected 1 call to CreateTable, got %v.", s3.CreateBucketInvocations)
	}

	if s3.Region != region {
		t.Errorf("Expected %v, got %v.", region, s3.Region)
	}

	if err != nil {
		t.Errorf("Expected no error, got %v.", err)
	}
}

// This test checks if error is returned when create bucket fails
func TestInitializeS3BucketWhenCreateBucketErrors(t *testing.T) {
	t.Parallel()

	bucket := "Test"
	region := "test"
	someError := errors.New("some error")
	s3 := &mocks.MockS3{HeadBucketError: someError, CreateBucketError: someError}

	err := storage.InitializeS3Bucket(bucket, region, s3)

	if s3.HeadBucketInvocations != 1 {
		t.Errorf("Expected 1 call to HeadBucket, got %v.", s3.HeadBucketInvocations)
	}

	if s3.BucketName != bucket {
		t.Errorf("Expected %v, got %v.", bucket, s3.BucketName)
	}

	if s3.CreateBucketInvocations != 1 {
		t.Errorf("Expected 1 call to CreateTable, got %v.", s3.CreateBucketInvocations)
	}

	if s3.Region != region {
		t.Errorf("Expected %v, got %v.", region, s3.Region)
	}

	if err == nil {
		t.Error("Expected error, got nil.")
	}

	if err.Error() != someError.Error() {
		t.Errorf("Expected %v, got %v.", someError.Error(), err.Error())
	}
}
