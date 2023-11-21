package release

import (
	"context"
	"fmt"
	"log"
	"math"
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
	ReleaseNotFound     = status.Error(codes.NotFound, "Release not found.")

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

	return nil
}

// Publishes a new release.
func (s *ReleaseService) Publish(ctx context.Context, request *release.PublishRequest) (*release.PublishResponse, error) {

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("release.name", request.GetName()),
		attribute.String("release.version", request.GetVersion()),
		attribute.String("release.type", request.GetType()),
		attribute.String("release.organization", request.GetOrganization()),
	)
	defer span.End()

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
		return nil, MarshalReleaseError
	}

	in := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(ReleaseTableName),
	}

	if _, err = s.Db.PutItem(ctx, in); err != nil {
		span.RecordError(err)
		return nil, PublishReleaseError
	}

	return ReleasePublished, nil
}

// ListReleases Retrieves all releases.
// Only releases that have been published should be reported.
func (s *ReleaseService) ListReleases(ctx context.Context, request *releaseSvc.ListReleasesRequest) (*releaseSvc.ListReleasesResponse, error) {

	span := trace.SpanFromContext(ctx)

	// attribute does not support Uint64 so converting to int64
	MaxAgeSeconds := convertUint64ToInt64(request.GetMaxAgeSeconds())

	span.SetAttributes(
		attribute.StringSlice("release.organizations", request.GetOrganizations()),
		attribute.Int64("release.maxAge", MaxAgeSeconds),
		attribute.StringSlice("release.types", request.GetTypes()),
		attribute.String("release.page", request.GetPage().String()),
	)
	defer span.End()

	if request.MaxAgeSeconds == nil {
		request.MaxAgeSeconds = &DefaultMaxAgeSeconds
	}

	// Generate UTC string (Now - MaxAgeSeconds)
	splitTimeISO := time.Unix(time.Now().UTC().Unix()-int64(*request.MaxAgeSeconds), 0).UTC().String()

	// Construct the filter builder with a name and value.
	filter := expression.Name("createdAt").GreaterThanEqual(expression.Value(splitTimeISO))
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		span.RecordError(err)
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
		span.RecordError(err)
		return nil, err
	}
	span.SetAttributes(
		attribute.Int("release.count", len(response.Items)),
	)
	grpcResponse := releaseSvc.ListReleasesResponse{}
	if response.Items != nil {
		for _, item := range response.Items {

			Release := &releaseSvc.Release{}
			if err3 := attributevalue.UnmarshalMap(item, &Release); err3 != nil {
				span.RecordError(err3)
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

	span := trace.SpanFromContext(ctx)

	// attribute does not support Uint64 so converting to int64
	MaxAgeSeconds := convertUint64ToInt64(request.GetMaxAgeSeconds())

	span.SetAttributes(
		attribute.StringSlice("release.organizations", request.GetOrganizations()),
		attribute.Int64("release.maxAge", MaxAgeSeconds),
		attribute.StringSlice("release.types", request.GetTypes()),
		attribute.String("release.page", request.GetPage().String()),
	)
	defer span.End()
	scanQueryInputs := &dynamodb.ScanInput{
		TableName: aws.String(ReleaseTableName),
	}

	response, err := s.Db.Scan(ctx, scanQueryInputs)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	log.Println(response)

	if response == nil {
		span.RecordError(ReleaseNotFound)
		return &releaseSvc.ListReleasesResponse{}, nil
	}

	var releases []*releaseSvc.Release
	for _, item := range response.Items {
		releaseInfo := new(releaseSvc.Release)
		if err := attributevalue.UnmarshalMap(item, &releaseInfo); err != nil {
			span.RecordError(err)
			return nil, err
		}

		releases = append(releases, releaseInfo)
	}

	// Sort releases based on the "createdAt" attribute in descending order
	sort.SliceStable(releases, func(i, j int) bool {
		return releases[i].CreatedAt > releases[j].CreatedAt
	})

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
		log.Printf("ScanInput failed: %v", err)
		return nil, err
	}

	if response.Items == nil || len(response.Items) < 1 {
		return nil, fmt.Errorf("nothing to retrieve '%v'", request.GetPage())
	}

	typeValues := make([]string, 0, len(response.Items))
	typeStr := ""

	if response.Items != nil {
		for _, item := range response.Items {
			typeAttr, found := item["type"]
			if !found {
				log.Println("type attribute not found")
				continue
			}

			if err := attributevalue.Unmarshal(typeAttr, &typeStr); err != nil {
				span.RecordError(err)
				log.Printf("Marshal error: %v", err)
				continue
			}

			typeValues = append(typeValues, typeStr)
		}
	}
	grpcResponse := &releaseSvc.ListReleaseTypesResponse{Types: GetDistinctValues(typeValues)}
	return grpcResponse, nil

}

// ListOrganization is used to retrieve all distinct organizations.
func (s *ReleaseService) ListOrganization(ctx context.Context, request *releaseSvc.ListOrganizationRequest) (*releaseSvc.ListOrganizationResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("release.page", request.GetPage().String()),
	)

	scanQueryInputs := &dynamodb.ScanInput{
		ProjectionExpression: aws.String("organization"),
		TableName:            aws.String(ReleaseTableName),
	}

	response, err := s.Db.Scan(ctx, scanQueryInputs)
	if err != nil {
		span.RecordError(err)
		log.Printf("ScanInput failed: %v", err)
		return nil, err
	}

	if response.Items == nil || len(response.Items) < 1 {
		return nil, fmt.Errorf("nothing to retrieve '%v'", request.GetPage())
	}

	orgValues := make([]string, 0, len(response.Items))
	orgStr := ""

	if response.Items != nil {
		for _, item := range response.Items {
			typeAttr, found := item["organization"]
			if !found {
				log.Println("organizatiob attribte not found")
				continue
			}

			if err := attributevalue.Unmarshal(typeAttr, &orgStr); err != nil {
				span.RecordError(err)
				log.Printf("Marshal error: %v", err)
				continue
			}
			orgValues = append(orgValues, orgStr)
		}
	}
	grpcResponse := &releaseSvc.ListOrganizationResponse{Organizations: GetDistinctValues(orgValues)}
	return grpcResponse, nil
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

// Converts Uint64 to int64
func convertUint64ToInt64(uint64Value uint64) int64 {

	// Returning the max Int64 if the value is greater than max Int64 for spans.
	if uint64Value > math.MaxInt64 {
		return math.MaxInt64
	}

	return int64(uint64Value)
}
