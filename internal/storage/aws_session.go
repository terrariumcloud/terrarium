package storage

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"os"
)

type AWSSessionConfig struct {
	UseLocalStack bool
	Region        string
	Key           string
	Secret        string
}

// Create new AWS Session with provided API key, Secret and Region
func NewAwsSession(sessionConfig AWSSessionConfig) (*aws.Config, error) {
	awsRegion := sessionConfig.Region
	if envRegion := os.Getenv("AWS_REGION"); envRegion != "" {
		awsRegion = envRegion
	}

	if sessionConfig.UseLocalStack {
		sessionConfig.Key = "test"
		sessionConfig.Secret = "test"
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if sessionConfig.UseLocalStack {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           "http://localhost:4566",
				SigningRegion: awsRegion,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(sessionConfig.Key, sessionConfig.Secret, "")))
	if err != nil {
		return nil, err
	}
	otelaws.AppendMiddlewares(&cfg.APIOptions)
	return &cfg, nil
}
