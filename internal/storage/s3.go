package storage

import (
	"log"

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
