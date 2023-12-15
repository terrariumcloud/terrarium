package storage

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type AWSS3BucketClient interface {
	HeadBucket(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error)
	CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

// NewS3Client Create new S3 client
func NewS3Client(sessionConfig AWSSessionConfig) *s3.Client {
	cfg, err := NewAwsSession(sessionConfig)
	if err != nil {
		log.Fatalf("Unable to create AWS Session: %s", err.Error())
	}
	return s3.NewFromConfig(*cfg, func(options *s3.Options) {
		options.UsePathStyle = sessionConfig.UseLocalStack
	})
}

// InitializeS3Bucket - checks if bucket exists, in case it doesn't it creates it
func InitializeS3Bucket(bucketName, region string, svc AWSS3BucketClient) error {
	in := &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	}
	if _, err := svc.HeadBucket(context.TODO(), in); err != nil {
		log.Printf("Creating S3 bucket: %s", bucketName)
		if _, err := svc.CreateBucket(context.TODO(), &s3.CreateBucketInput{
			Bucket: &bucketName,
			ACL:    types.BucketCannedACLPrivate,
			CreateBucketConfiguration: &types.CreateBucketConfiguration{
				LocationConstraint: types.BucketLocationConstraint(region),
			},
		}); err != nil {
			return err
		} else {
			log.Println("Bucket created.")
			return nil
		}
	}
	log.Printf("S3 bucket %s already exists.", bucketName)
	return nil

}
