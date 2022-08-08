package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Create new AWS Session with provided API key, Secret and Region
func NewAwsSession(key string, secret string, region string) (*session.Session, error) {
	if sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
	}); err != nil {
		return nil, err
	} else {
		return sess, nil
	}
}

// Create new AWS Session using shared config (.aws/ dir and env)
func NewAwsSessionFromSharedConfig() (*session.Session, error) {
	if sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}); err != nil {
		return nil, err
	} else {
		return sess, nil
	}
}
