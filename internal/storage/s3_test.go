package storage_test

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/storage"
)

type fakeS3 struct {
	s3iface.S3API
	headBucketInvocations   int
	headBucketOut           *s3.HeadBucketOutput
	headBucketError         error
	bucketName              string
	createBucketInvocations int
	createBucketOut         *s3.CreateBucketOutput
	createBucketError       error
	region                  string
}

func (fs3 *fakeS3) HeadBucket(in *s3.HeadBucketInput) (*s3.HeadBucketOutput, error) {
	fs3.headBucketInvocations++
	fs3.bucketName = *in.Bucket
	return fs3.headBucketOut, fs3.headBucketError
}

func (fs3 *fakeS3) CreateBucket(in *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	fs3.createBucketInvocations++
	fs3.bucketName = *in.Bucket
	fs3.region = *in.CreateBucketConfiguration.LocationConstraint
	return fs3.createBucketOut, fs3.createBucketError
}

// This test checks if bucket is not recreated when it already exists
func TestInitializeS3BucketWhenBucketExists(t *testing.T) {
	t.Parallel()

	bucket := "Test"
	region := "test"
	s3 := &fakeS3{}

	err := storage.InitializeS3Bucket(bucket, region, s3)

	if s3.headBucketInvocations != 1 {
		t.Errorf("Expected 1 call to HeadBucket, got %v.", s3.headBucketInvocations)
	}

	if s3.bucketName != bucket {
		t.Errorf("Expected %v, got %v.", bucket, s3.bucketName)
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
	s3 := &fakeS3{
		headBucketError: errors.New("some error"),
	}

	err := storage.InitializeS3Bucket(bucket, region, s3)

	if s3.headBucketInvocations != 1 {
		t.Errorf("Expected 1 call to HeadBucket, got %v.", s3.headBucketInvocations)
	}

	if s3.bucketName != bucket {
		t.Errorf("Expected %v, got %v.", bucket, s3.bucketName)
	}

	if s3.createBucketInvocations != 1 {
		t.Errorf("Expected 1 call to CreateTable, got %v.", s3.createBucketInvocations)
	}

	if s3.region != region {
		t.Errorf("Expected %v, got %v.", region, s3.region)
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
	s3 := &fakeS3{
		headBucketError:   someError,
		createBucketError: someError,
	}

	err := storage.InitializeS3Bucket(bucket, region, s3)

	if s3.headBucketInvocations != 1 {
		t.Errorf("Expected 1 call to HeadBucket, got %v.", s3.headBucketInvocations)
	}

	if s3.bucketName != bucket {
		t.Errorf("Expected %v, got %v.", bucket, s3.bucketName)
	}

	if s3.createBucketInvocations != 1 {
		t.Errorf("Expected 1 call to CreateTable, got %v.", s3.createBucketInvocations)
	}

	if s3.region != region {
		t.Errorf("Expected %v, got %v.", region, s3.region)
	}

	if err == nil {
		t.Error("Expected error, got nil.")
	}

	if err.Error() != someError.Error() {
		t.Errorf("Expected %v, got %v.", someError.Error(), err.Error())
	}
}
