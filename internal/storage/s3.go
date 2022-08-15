package storage

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// Create new S3 client
func NewS3Client(key string, secret string, region string) s3iface.S3API {
	sess, err := NewAwsSession(key, secret, region)
	if err != nil {
		log.Fatalf("Unable to create AWS Session: %s", err.Error())
	}
	return s3.New(sess)
}

// InitializeS3Bucket - checks if bucket exists, in case it doesn't it creates it
func InitializeS3Bucket(bucketName, region string, svc s3iface.S3API) error {
	in := &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	}
	if _, err := svc.HeadBucket(in); err != nil {
		log.Printf("Creating S3 bucket: %s", bucketName)
		if _, err := svc.CreateBucket(&s3.CreateBucketInput{
			ACL:    aws.String(s3.BucketCannedACLPrivate),
			Bucket: &bucketName,
			CreateBucketConfiguration: &s3.CreateBucketConfiguration{
				LocationConstraint: aws.String(region),
			},
		}); err != nil {
			return err
		} else {
			log.Println("Bucket created.")
			return nil
		}
	} else {
		log.Printf("S3 bucket %s already exists.", bucketName)
		return nil
	}

}
