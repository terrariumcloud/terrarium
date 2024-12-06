package version_manager

import (
	"context"
	"log"
	"strings"
	"time"

	releasePkg "github.com/terrariumcloud/terrarium/pkg/terrarium/release"

	"github.com/terrariumcloud/terrarium/internal/module/services"
	releaseSvc "github.com/terrariumcloud/terrarium/internal/release/services"
	"github.com/terrariumcloud/terrarium/internal/storage"
	terrarium "github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/apparentlymart/go-versions/versions"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	DefaultVersionsTableName      = "terrarium-module-versions-ciedev-4757"
	DefaultVersionManagerEndpoint = "version_manager:3001"
)

var (
	VersionsTableName      = DefaultVersionsTableName
	VersionManagerEndpoint = DefaultVersionManagerEndpoint

	VersionCreated   = &terrarium.Response{Message: "Version created."}
	VersionPublished = &terrarium.Response{Message: "Version published."}
	VersionAborted   = &terrarium.Response{Message: "Version aborted."}

	ModuleVersionsTableInitializationError = status.Error(codes.Unknown, "Failed to initialize table for module versions.")
	MarshalModuleVersionError              = status.Error(codes.Unknown, "Failed to marshal module version.")
	CreateModuleVersionError               = status.Error(codes.Unknown, "Failed to create module version.")
	AbortModuleVersionError                = status.Error(codes.Unknown, "Failed to abort module version.")
	PublishModuleVersionError              = status.Error(codes.Unknown, "Failed to publish module version.")
	DevelopmentVersion                     = versions.MustParseVersion("0.0.0")
)

type VersionManagerService struct {
	services.UnimplementedVersionManagerServer
	Db             storage.DynamoDBTableCreator
	Table          string
	Schema         *dynamodb.CreateTableInput
	ReleaseService releaseSvc.PublisherClient
}

type ModuleVersion struct {
	Name        string `json:"name" bson:"name" dynamodbav:"name"`
	Version     string `json:"version" bson:"version" dynamodbav:"version"`
	CreatedOn   string `json:"created_on" bson:"created_on" dynamodbav:"created_on"`
	PublishedOn string `json:"published_on" bson:"published_on" dynamodbav:"published_on"`
}

// RegisterWithServer Registers VersionManagerService with grpc server
func (s *VersionManagerService) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	if err := storage.InitializeDynamoDb(s.Table, s.Schema, s.Db); err != nil {
		log.Println(err)
		return ModuleVersionsTableInitializationError
	}
	services.RegisterVersionManagerServer(grpcServer, s)

	return nil
}

// BeginVersion Creates new Module Version with Version Manager service
func (s *VersionManagerService) BeginVersion(ctx context.Context, request *terrarium.BeginVersionRequest) (*terrarium.Response, error) {
	log.Println("Creating new version.")
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("module.name", request.Module.GetName()),
		attribute.String("module.version", request.Module.GetVersion()),
	)

	mv := ModuleVersion{
		Name:      request.Module.GetName(),
		Version:   request.Module.GetVersion(),
		CreatedOn: time.Now().UTC().String(),
	}

	av, err := attributevalue.MarshalMap(mv)

	if err != nil {
		span.RecordError(err)
		log.Println(err)
		return nil, MarshalModuleVersionError
	}

	in := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(VersionsTableName),
	}

	if _, err = s.Db.PutItem(ctx, in); err != nil {
		span.RecordError(err)
		log.Println(err)
		return nil, CreateModuleVersionError
	}

	log.Println("New version created.")
	return VersionCreated, nil
}

func (s *VersionManagerService) GetModuleKey(module *terrarium.Module) (map[string]types.AttributeValue, error) {
	moduleName, err := attributevalue.Marshal(module.GetName())
	if err != nil {
		return map[string]types.AttributeValue{}, err
	}
	moduleVersion, err := attributevalue.Marshal(module.GetVersion())
	if err != nil {
		return map[string]types.AttributeValue{}, err
	}
	return map[string]types.AttributeValue{
		"name":    moduleName,
		"version": moduleVersion,
	}, nil
}

// AbortVersion Removes Module Version with Version Manager service
func (s *VersionManagerService) AbortVersion(ctx context.Context, request *services.TerminateVersionRequest) (*terrarium.Response, error) {
	log.Println("Aborting module version.")

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("module.name", request.Module.GetName()),
		attribute.String("module.version", request.Module.GetVersion()),
	)

	moduleKey, err := s.GetModuleKey(request.Module)
	if err != nil {
		span.RecordError(err)
		log.Println(err)
		return nil, AbortModuleVersionError
	}

	in := &dynamodb.DeleteItemInput{
		Key:       moduleKey,
		TableName: aws.String(VersionsTableName),
	}

	if _, err := s.Db.DeleteItem(ctx, in); err != nil {
		span.RecordError(err)
		log.Println(err)
		return nil, AbortModuleVersionError
	}

	log.Println("Module version aborted.")
	return VersionAborted, nil
}

// PublishVersion Updates Module Version to published with Version Manager service
// And publishes a release.
func (s *VersionManagerService) PublishVersion(ctx context.Context, request *services.TerminateVersionRequest) (*terrarium.Response, error) {
	log.Println("Publishing module version.")

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("module.name", request.Module.GetName()),
		attribute.String("module.version", request.Module.GetVersion()),
	)
	moduleKey, err := s.GetModuleKey(request.Module)
	if err != nil {
		span.RecordError(err)
		log.Println(err)
		return nil, PublishModuleVersionError
	}

	publishOn, err := attributevalue.Marshal(time.Now().UTC().String())
	if err != nil {
		span.RecordError(err)
		log.Println(err)
		return nil, PublishModuleVersionError
	}

	in := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":published_on": publishOn,
		},
		Key:              moduleKey,
		TableName:        aws.String(VersionsTableName),
		UpdateExpression: aws.String("set published_on = :published_on"),
	}

	if _, err := s.Db.UpdateItem(ctx, in); err != nil {
		span.RecordError(err)
		log.Println(err)
		return nil, PublishModuleVersionError
	}

	// PUBLISH RELEASE
	parsedVersion, err := versions.ParseVersion(strings.ReplaceAll(request.Module.GetVersion(), "v", ""))
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	if parsedVersion.GreaterThan(DevelopmentVersion) && s.ReleaseService != nil {
		moduleAddress := strings.Split(request.Module.GetName(), "/")
		orgName := moduleAddress[0]

		if _, err := s.ReleaseService.Publish(ctx, &releasePkg.PublishRequest{
			Name:         request.Module.GetName(),
			Version:      request.Module.GetVersion(),
			Type:         "module",
			Organization: orgName,
		}); err != nil {
			span.RecordError(err)
		}
	}

	log.Println("Module version published.")
	return VersionPublished, nil
}

// ListModuleVersions Retrieve all versions of a given module and return an array of versions.
// Only versions that have been published should be reported
func (s *VersionManagerService) ListModuleVersions(ctx context.Context, request *services.ListModuleVersionsRequest) (*services.ListModuleVersionsResponse, error) {
	projection := expression.NamesList(expression.Name("version"))
	filter := expression.And(
		expression.Name("name").Equal(expression.Value(request.Module)),
		expression.Name("published_on").AttributeExists())
	expr, err := expression.NewBuilder().WithProjection(projection).WithFilter(filter).Build()
	if err != nil {
		log.Printf("Expression Builder failed creation: %v", err)
		return nil, err
	}

	scanQueryInputs := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(VersionsTableName),
	}

	response, err := s.Db.Scan(ctx, scanQueryInputs)
	if err != nil {
		log.Printf("ScanInput failed: %v", err)
		return nil, err
	}

	grpcResponse := services.ListModuleVersionsResponse{}
	if response.Items != nil {
		for _, item := range response.Items {
			moduleVersion := ModuleVersion{}
			if err3 := attributevalue.UnmarshalMap(item, &moduleVersion); err3 != nil {
				log.Printf("UnmarshalMap failed: %v", err3)
				return nil, err3
			}
			grpcResponse.Versions = append(grpcResponse.Versions, moduleVersion.Version)
		}
	}
	var semverList versions.List
	for _, moduleVersion := range grpcResponse.Versions {
		parsedVersion, err := versions.ParseVersion(moduleVersion)

		if err != nil {
			log.Printf("Skipping invalid semantic version: %v", moduleVersion)
		} else {
			semverList = append(semverList, parsedVersion)
		}

	}
	semverList.Sort()

	var sortedVersions []string
	for _, moduleVersion := range semverList {
		sortedVersions = append(sortedVersions, moduleVersion.String())
	}
	grpcResponse.Versions = sortedVersions

	return &grpcResponse, nil
}

// GetModuleVersionsSchema returns CreateTableInput that can be used to create table if it does not exist
func GetModuleVersionsSchema(table string) *dynamodb.CreateTableInput {
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
