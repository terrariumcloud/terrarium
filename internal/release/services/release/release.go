package release

import (
	"context"
	"log"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	releaseSvc "github.com/terrariumcloud/terrarium/internal/release/services"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/terrariumcloud/terrarium/internal/storage"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/release"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"
)

const (
	DefaultReleaseTableName       = "terrarium-releases"
	DefaultReleaseServiceEndpoint = "release:3001"
)

var (
	ReleaseTableName                = DefaultReleaseTableName
	ReleaseTableInitializationError = status.Error(codes.Unknown, "Failed to initialize table for releases.")
	ReleaseServiceEndpoint          = DefaultReleaseServiceEndpoint

	ReleasePublished = &release.PublishResponse{} // No return information at this stage

	MarshalReleaseError = status.Error(codes.Unknown, "Failed to marshal publish release.")
	PublishReleaseError = status.Error(codes.Unknown, "Failed to publish release.")

	TimeFormatLayout     = "2006-01-02 15:04:05.999999999 -0700 MST"
	DefaultMaxAgeSeconds = uint64(86400) // 1 day in seconds
)

type ReleaseService struct {
	releaseSvc.UnimplementedPublisherServer
	Db     storage.DynamoDBTableCreator
	Table  string
	Schema *dynamodb.CreateTableInput
}

type Release struct {
	Type         string          `json:"type" bson:"type" dynamodbav:"type"`
	Organization string          `json:"organization" bson:"organization" dynamodbav:"organization"`
	Name         string          `json:"name" bson:"name" dynamodbav:"name"`
	Version      string          `json:"version" bson:"version" dynamodbav:"version"`
	Description  string          `json:"description" bson:"description" dynamodbav:"description"`
	Links        []*release.Link `json:"links" bson:"links" dynamodbav:"links"`
	CreatedAt    string          `json:"createdAt" bson:"createdAt" dynamodbav:"createdAt"`
}

// Registers ReleaseService with grpc server
func (s *ReleaseService) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	if err := storage.InitializeDynamoDb(s.Table, s.Schema, s.Db); err != nil {
		log.Println(err)
		return ReleaseTableInitializationError
	}

	releaseSvc.RegisterPublisherServer(grpcServer, s)

	// Code below to be uncommented for testing ListReleases
	// var maxAge uint64 = 72000
	// var request = &releaseSvc.ListReleasesRequest{MaxAgeSeconds: &maxAge}

	// s.ListReleases(context.TODO(), request)

	return nil
}

// Publishes a new release.
func (s *ReleaseService) Publish(ctx context.Context, request *release.PublishRequest) (*release.PublishResponse, error) {

	log.Println("Creating new release.")
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("release.name", request.GetName()),
		attribute.String("release.version", request.GetVersion()),
	)

	mv := Release{
		Type:         request.GetType(),
		Organization: request.GetOrganization(),
		Name:         request.GetName(),
		Version:      request.GetVersion(),
		Description:  request.GetDescription(),
		Links:        request.GetLinks(),
		CreatedAt:    time.Now().UTC().String(),
	}

	av, err := attributevalue.MarshalMap(mv)

	if err != nil {
		span.RecordError(err)
		log.Println(err)
		return nil, MarshalReleaseError
	}

	in := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(ReleaseTableName),
	}

	if _, err = s.Db.PutItem(ctx, in); err != nil {
		span.RecordError(err)
		log.Println(err)
		return nil, PublishReleaseError
	}

	log.Println("New release created.")
	return ReleasePublished, nil
}

// to do: latest release to be returned first
// to do: return last N releases published in the last hours (3600 seconds => max_age_seconds)
// now := time.Now().Unix()
// max_age_seconds := int64(3600)
// filter := expression.Name("created_on").Between(expression.Value(now-max_age_seconds), expression.Value(max_age_seconds))

// ListReleases Retrieves all releases.
// Only releases that have been published should be reported.
func (s *ReleaseService) ListReleases(ctx context.Context, request *releaseSvc.ListReleasesRequest) (*releaseSvc.ListReleasesResponse, error) {

	if request.MaxAgeSeconds == nil {
		request.MaxAgeSeconds = &DefaultMaxAgeSeconds
	}

	// Generate UTC string (Now - MaxAgeSeconds)
	splitTimeISO := time.Unix(time.Now().UTC().Unix()-int64(*request.MaxAgeSeconds), 0).UTC().String()

	// Construct the filter builder with a name and value.
	filter := expression.Name("createdAt").GreaterThanEqual(expression.Value(splitTimeISO))
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		log.Printf("Expression Builder failed creation: %v", err)
		return nil, err
	}

	scanQueryInputs := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(ReleaseTableName),
	}

	response, err := s.Db.Scan(ctx, scanQueryInputs)
	if err != nil {
		log.Printf("ScanInput failed: %v", err)
		return nil, err
	}
	log.Println("Filtered Release Count: ", len(response.Items))

	grpcResponse := releaseSvc.ListReleasesResponse{}
	if response.Items != nil {
		for _, item := range response.Items {

			Release := &releaseSvc.Release{}
			if err3 := attributevalue.UnmarshalMap(item, &Release); err3 != nil {
				log.Printf("UnmarshalMap failed: %v", err3)
				return nil, err3
			}
			grpcResponse.Releases = append(grpcResponse.Releases, Release)
		}
	}

	// Sort list of releases based on createdAt field
	grpcResponse.Releases = sortReleaseList(grpcResponse.Releases)

	return &grpcResponse, nil
}

// GetReleaseSchema returns CreateTableInput that can be used to create table if it does not exist
func GetReleaseSchema(table string) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("name"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("version"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("version"),
				KeyType:       types.KeyTypeRange,
			},
		},
		TableName:   aws.String(table),
		BillingMode: types.BillingModePayPerRequest,
	}
}

// Sort a slice of releaseSvc.Release struct, using sort.SliceStable method
func sortReleaseList(releases []*releaseSvc.Release) []*releaseSvc.Release {
	sort.SliceStable(releases, func(a, b int) bool {

		TimeA, err := time.Parse(TimeFormatLayout, releases[a].CreatedAt)
		if err != nil {
			log.Printf("Failed to parse createdAt field: %v", err)
		}

		TimeB, err := time.Parse(TimeFormatLayout, releases[b].CreatedAt)
		if err != nil {
			log.Printf("Failed to parse createdAt field: %v", err)
		}
		return TimeA.Before(TimeB)
	})

	return releases
}
