package release

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
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
)

type ReleaseService struct {
	releaseSvc.UnimplementedPublisherServer
	Db     storage.DynamoDBTableCreator
	Table  string
	Schema *dynamodb.CreateTableInput
}

type ReleaseInfo struct {
	Type         string          `json:"type" bson:"type" dynamodbav:"type"`
	Organization string          `json:"organization" bson:"organization" dynamodbav:"organization"`
	Name         string          `json:"name" bson:"name" dynamodbav:"name"`
	Version      string          `json:"version" bson:"version" dynamodbav:"version"`
	Description  string          `json:"description" bson:"description" dynamodbav:"description"`
	Link         []*release.Link `json:"link" bson:"link" dynamodbav:"link"`
	CreatedOn    string          `json:"created_on" bson:"created_on" dynamodbav:"created_on"`
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

	log.Println("Creating new release.")
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("release.name", request.GetName()),
		attribute.String("release.version", request.GetVersion()),
	)

	mv := ReleaseInfo{
		Type:         request.GetType(),
		Organization: request.GetOrganization(),
		Name:         request.GetName(),
		Version:      request.GetVersion(),
		Description:  request.GetDescription(),
		Link:         request.GetLinks(),
		CreatedOn:    time.Now().UTC().String(),
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

	// filter := expression.And(
	// 	expression.Name("type").Equal(expression.Value(request.Types)),
	// 	expression.Name("organizations").AttributeExists())
	// expr, err := expression.NewBuilder().WithFilter(filter).Build()
	// if err != nil {
	// 	log.Printf("Expression Builder failed creation: %v", err)
	// 	return nil, err
	// }

	scanQueryInputs := &dynamodb.ScanInput{
		// ExpressionAttributeNames:  expr.Names(),
		// ExpressionAttributeValues: expr.Values(),
		// FilterExpression:          expr.Filter(),
		TableName: aws.String(ReleaseTableName),
	}

	response, err := s.Db.Scan(ctx, scanQueryInputs)
	if err != nil {
		log.Printf("ScanInput failed: %v", err)
		return nil, err
	}

	grpcResponse := releaseSvc.ListReleasesResponse{}
	if response.Items != nil {
		for _, item := range response.Items {
			releaseInfo := &releaseSvc.Release{}
			if err3 := attributevalue.UnmarshalMap(item, &releaseInfo); err3 != nil {
				log.Printf("UnmarshalMap failed: %v", err3)
				return nil, err3
			}
			grpcResponse.Releases = append(grpcResponse.Releases, releaseInfo)
		}
	}

	return &grpcResponse, nil
}

// func (s *ReleaseService) ListReleaseTypes(ctx context.Context, request *releaseSvc.ListReleaseTypesRequest) (*releaseSvc.ListReleaseTypesResponse, error) {

// 	return
// }

// func (s *ReleaseService) ListOrganization(ctx context.Context, request *releaseSvc.ListOrganizationRequest) (*releaseSvc.ListOrganizationResponse, error) {

// 	return
// }

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
