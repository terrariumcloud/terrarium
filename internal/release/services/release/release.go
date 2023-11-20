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

	MarshalReleaseError   = status.Error(codes.Unknown, "Failed to marshal publish release.")
	PublishReleaseError   = status.Error(codes.Unknown, "Failed to publish release.")
	ListReleaseTypesError = status.Error(codes.Unknown, "Failed to retrieve release types")

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

// TODO: OPTIMIZE!!!
// GetLatestRelease retrieves the latest release.
func (s *ReleaseService) GetLatestRelease(ctx context.Context, request *releaseSvc.ListReleasesRequest) (*releaseSvc.ListReleasesResponse, error) {
	log.Printf("Getting latest published release.")
	scanQueryInputs := &dynamodb.ScanInput{
		TableName: aws.String(ReleaseTableName),
	}

	response, err := s.Db.Scan(ctx, scanQueryInputs)
	if err != nil {
		log.Printf("ScanInput failed: %v", err)
		return nil, err
	}
	log.Println(response)

	if response == nil {
		log.Println("Release not found")
		return &releaseSvc.ListReleasesResponse{}, nil
	}
	// Convert DynamoDB items to a slice of custom struct
	var releases []*releaseSvc.Release
	for _, item := range response.Items {
		releaseInfo := new(releaseSvc.Release)
		if err := attributevalue.UnmarshalMap(item, &releaseInfo); err != nil {
			log.Printf("UnmarshalMap failed: %v", err)
			return nil, err
		}

		releases = append(releases, releaseInfo)
	}
	log.Println("Unmarshalled releases:", releases)

	log.Println("Sorting releases...")
	// Sort releases based on the "createdAt" attribute in descending order
	sort.SliceStable(releases, func(i, j int) bool {
		return releases[i].CreatedAt > releases[j].CreatedAt
	})

	// Return only the latest release
	grpcResponse := releaseSvc.ListReleasesResponse{}
	grpcResponse.Releases = append(grpcResponse.Releases, releases[0])

	return &grpcResponse, nil
}

// GetDistinctValues is a helper function to filter the response and return only distinct values.
func GetDistinctValues(resp []string) []string {
	temp := make(map[string]bool)

	for _, item := range resp {
		temp[item] = true
	}

	var distinctList []string
	for i := range temp {
		distinctList = append(distinctList, i)
	}

	return distinctList
}

// ListReleaseTypes is used to retrieve all distinct release types.
func (s *ReleaseService) ListReleaseTypes(ctx context.Context, request *releaseSvc.ListReleaseTypesRequest) (*releaseSvc.ListReleaseTypesResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("release.page", request.GetPage().String()),
	)

	scanQueryInputs := &dynamodb.ScanInput{
		ProjectionExpression: aws.String("type"),
		TableName:            aws.String(ReleaseTableName),
	}

	response, err := s.Db.Scan(ctx, scanQueryInputs)
	if err != nil {
		span.RecordError(err)
		log.Printf("Couldn't scan for release types: %v", err)
		return nil, err
	}

	typeValues := make([]string, 0, len(response.Items))
	for _, item := range response.Items {
		typeAttr, found := item["type"]
		if !found {
			log.Println("type attribute not found")
			continue
		}

		var typeStr string
		if err := attributevalue.Unmarshal(typeAttr, &typeStr); err != nil {
			span.RecordError(err)
			log.Printf("Error converting attribute value to string: %v", err)
			continue
		}

		typeValues = append(typeValues, typeStr)
	}
	grpcResponse := &releaseSvc.ListReleaseTypesResponse{Types: GetDistinctValues(typeValues)}
	return grpcResponse, nil
}

// API: ListOrganization

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
