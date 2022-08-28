package mocks

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type MockS3 struct {
	s3iface.S3API
	HeadBucketInvocations   int
	HeadBucketOut           *s3.HeadBucketOutput
	HeadBucketError         error
	BucketName              string
	CreateBucketInvocations int
	CreateBucketOut         *s3.CreateBucketOutput
	CreateBucketError       error
	Region                  string
	PutObjectInvocations    int
	Filename                string
	PutObjectOut            *s3.PutObjectOutput
	PutObjectError          error
	GetObjectInvocations    int
	GetObjectOut            *s3.GetObjectOutput
	GetObjectError          error
}

func (ms3 *MockS3) HeadBucket(in *s3.HeadBucketInput) (*s3.HeadBucketOutput, error) {
	ms3.HeadBucketInvocations++
	ms3.BucketName = *in.Bucket
	return ms3.HeadBucketOut, ms3.HeadBucketError
}

func (ms3 *MockS3) CreateBucket(in *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	ms3.CreateBucketInvocations++
	ms3.BucketName = *in.Bucket
	ms3.Region = *in.CreateBucketConfiguration.LocationConstraint
	return ms3.CreateBucketOut, ms3.CreateBucketError
}

func (ms3 *MockS3) PutObject(in *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	ms3.PutObjectInvocations++
	ms3.BucketName = *in.Bucket
	ms3.Filename = *in.Key
	return ms3.PutObjectOut, ms3.PutObjectError
}

func (ms3 *MockS3) GetObject(in *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	ms3.GetObjectInvocations++
	ms3.BucketName = *in.Bucket
	ms3.Filename = *in.Key
	return ms3.GetObjectOut, ms3.GetObjectError
}
