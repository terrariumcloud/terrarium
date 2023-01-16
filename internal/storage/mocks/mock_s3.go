package mocks

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
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

func (ms3 *S3) HeadBucket(_ context.Context, in *s3.HeadBucketInput, _ ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
	ms3.HeadBucketInvocations++
	ms3.BucketName = *in.Bucket
	return ms3.HeadBucketOut, ms3.HeadBucketError
}
func (ms3 *S3) CreateBucket(_ context.Context, in *s3.CreateBucketInput, _ ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	ms3.CreateBucketInvocations++
	ms3.BucketName = *in.Bucket
	ms3.Region = string(in.CreateBucketConfiguration.LocationConstraint)
	return ms3.CreateBucketOut, ms3.CreateBucketError
}

func (ms3 *S3) PutObject(_ context.Context, in *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	ms3.PutObjectInvocations++
	ms3.BucketName = *in.Bucket
	ms3.Filename = *in.Key
	return ms3.PutObjectOut, ms3.PutObjectError
}

func (ms3 *S3) GetObject(_ context.Context, in *s3.GetObjectInput, _ ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	ms3.GetObjectInvocations++
	ms3.BucketName = *in.Bucket
	ms3.Filename = *in.Key
	return ms3.GetObjectOut, ms3.GetObjectError
}
