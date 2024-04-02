package version_manager

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/terrariumcloud/terrarium/internal/provider/services"
	"github.com/terrariumcloud/terrarium/internal/storage"
	terrarium "github.com/terrariumcloud/terrarium/pkg/terrarium/provider"

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
	DefaultVersionsTableName      = "terrarium-providers"
	DefaultVersionManagerEndpoint = "version_manager:3001"
)

var (
	VersionsTableName      = DefaultVersionsTableName
	VersionManagerEndpoint = DefaultVersionManagerEndpoint

	ProviderRegistered = &terrarium.Response{Message: "Provider registered successfully."}
	VersionPublished   = &terrarium.Response{Message: "Version published."}
	ProviderAborted    = &terrarium.Response{Message: "Provider aborted."}
	VersionAborted     = &terrarium.Response{Message: "Version aborted."}

	ProviderVersionsTableInitializationError = status.Error(codes.Unknown, "Failed to initialize table for provider versions.")
	AbortProviderVersionError                = status.Error(codes.Unknown, "Failed to abort provider version.")
	PublishProviderVersionError              = status.Error(codes.Unknown, "Failed to publish provider version.")
	ProviderGetError                         = status.Error(codes.Unknown, "Failed to check if provider already exists.")
	ProviderUpdateError                      = status.Error(codes.Unknown, "Failed to update provider.")
	ProviderRegisterError                    = status.Error(codes.Unknown, "Failed to register provider.")
	ExpressionBuildError                     = status.Error(codes.Unknown, "Failed to build update expression.")
	MarshalProviderError                     = status.Error(codes.Unknown, "Failed to marshal provider.")
)

type VersionManagerService struct {
	services.UnimplementedVersionManagerServer
	Db     storage.DynamoDBTableCreator
	Table  string
	Schema *dynamodb.CreateTableInput
}

type Provider struct {
	Name                string   `json:"name" bson:"name" dynamodbav:"name"`
	Version             string   `json:"version" bson:"version" dynamodbav:"version"`
	Protocols           []string `json:"protocols" bson:"protocols" dynamodbav:"protocols"`
	OS                  string   `json:"os" bson:"os" dynamodbav:"os"`
	Arch                string   `json:"arch" bson:"arch" dynamodbav:"arch"`
	Filename            string   `json:"filename" bson:"filename" dynamodbav:"filename"`
	DownloadURL         string   `json:"download_url" bson:"download_url" dynamodbav:"download_url"`
	ShasumsURL          string   `json:"shasums_url" bson:"shasums_url" dynamodbav:"shasums_url"`
	ShasumsSignatureURL string   `json:"shasums_signature_url" bson:"shasums_signature_url" dynamodbav:"shasums_signature_url"`
	Shasum              string   `json:"shasum" bson:"shasum" dynamodbav:"shasum"`
	KeyID               string   `json:"key_id" bson:"key_id" dynamodbav:"key_id"`
	ASCIIArmor          string   `json:"ascii_armor" bson:"ascii_armor" dynamodbav:"ascii_armor"`
	TrustSignature      string   `json:"trust_signature" bson:"trust_signature" dynamodbav:"trust_signature"`
	Source              string   `json:"source" bson:"source" dynamodbav:"source"`
	SourceURL           string   `json:"source_url" bson:"source_url" dynamodbav:"source_url"`
	Description         string   `json:"description" bson:"description" dynamodbav:"description"`
	SourceRepoUrl       string   `json:"source_repo_url" bson:"source_repo_url" dynamodbav:"source_repo_url"`
	Maturity            string   `json:"maturity" bson:"maturity" dynamodbav:"maturity"`
	CreatedOn           string   `json:"created_on" bson:"created_on" dynamodbav:"created_on"`
	ModifiedOn          string   `json:"modified_on" bson:"modified_on" dynamodbav:"modified_on"`
	PublishedOn         string   `json:"published_on" bson:"published_on" dynamodbav:"published_on"`
}

// RegisterWithServer Registers VersionManagerService with grpc server
func (s *VersionManagerService) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	if err := storage.InitializeDynamoDb(s.Table, s.Schema, s.Db); err != nil {
		log.Println(err)
		return ProviderVersionsTableInitializationError
	}
	services.RegisterVersionManagerServer(grpcServer, s)

	return nil
}

func (s *VersionManagerService) GetProviderKey(provider *terrarium.Provider) (map[string]types.AttributeValue, error) {

	providerName, err := attributevalue.Marshal(provider.GetName())
	if err != nil {
		return map[string]types.AttributeValue{}, err
	}

	providerVersion, err := attributevalue.Marshal(provider.GetVersion())
	if err != nil {
		return map[string]types.AttributeValue{}, err
	}

	os, err := attributevalue.Marshal(provider.GetOs())
	if err != nil {
		return map[string]types.AttributeValue{}, err
	}

	arch, err := attributevalue.Marshal(provider.GetArch())
	if err != nil {
		return map[string]types.AttributeValue{}, err
	}

	return map[string]types.AttributeValue{
		"name":    providerName,
		"version": providerVersion,
		"os":      os,
		"arch":    arch,
	}, nil
}

// AbortProvider Removes a Provider with Version Manager service
func (s *VersionManagerService) AbortProvider(ctx context.Context, request *services.TerminateProviderRequest) (*terrarium.Response, error) {
	log.Println("Aborting provider.")

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("provider.name", request.Provider.GetName()),
	)

	providerID, err := attributevalue.Marshal(request.Provider.GetName())
	if err != nil {
		span.RecordError(err)
		log.Println(err)
		return nil, ProviderGetError
	}

	in := &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"name": providerID,
		},
		TableName: aws.String(VersionsTableName),
	}

	if _, err := s.Db.DeleteItem(ctx, in); err != nil {
		span.RecordError(err)
		log.Println(err)
		return nil, AbortProviderVersionError
	}

	log.Println("Provider aborted.")
	return ProviderAborted, nil
}

func (s *VersionManagerService) AbortProviderVersion(ctx context.Context, request *services.TerminateVersionRequest) (*terrarium.Response, error) {
	log.Println("Aborting provider version.")

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("provider.name", request.Provider.GetName()),
		attribute.String("provider.version", request.Provider.GetVersion()),
	)

	providerID, err := attributevalue.Marshal(request.Provider.GetName())
	if err != nil {
		span.RecordError(err)
		log.Println(err)
		return nil, ProviderGetError
	}

	providerVersion, err := attributevalue.Marshal(request.Provider.GetVersion())
	if err != nil {
		span.RecordError(err)
		log.Println(err)
		return nil, ProviderGetError
	}

	in := &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"name":    providerID,
			"version": providerVersion,
		},
		TableName: aws.String(VersionsTableName),
	}

	if _, err := s.Db.DeleteItem(ctx, in); err != nil {
		span.RecordError(err)
		log.Println(err)
		return nil, AbortProviderVersionError
	}

	log.Println("Provider version aborted.")
	return VersionAborted, nil
}

// PublishVersion Updates Provider Version to published with Version Manager service
func (s *VersionManagerService) PublishVersion(ctx context.Context, request *services.PublishVersionRequest) (*terrarium.Response, error) {
	log.Println("Publishing provider version.")

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("provider.name", request.Provider.GetName()),
		attribute.String("provider.version", request.Provider.GetVersion()),
		attribute.String("provider.os",request.Provider.GetOs()),
		attribute.String("provider.arch",request.Provider.GetArch()),
	)

	providerKey, err := s.GetProviderKey(request.GetProvider())
	if err != nil {
		span.RecordError(err)
		log.Println(err)
		return nil, PublishProviderVersionError
	}

	publishOn, err := attributevalue.Marshal(time.Now().UTC().String())
	if err != nil {
		span.RecordError(err)
		log.Println(err)
		return nil, PublishProviderVersionError
	}

	in := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":published_on": publishOn,
		},
		Key:              providerKey,
		TableName:        aws.String(VersionsTableName),
		UpdateExpression: aws.String("set published_on = :published_on"),
	}

	if _, err := s.Db.UpdateItem(ctx, in); err != nil {
		span.RecordError(err)
		log.Println(err)
		return nil, PublishProviderVersionError
	}

	log.Println("Provider version published.")
	return VersionPublished, nil
}

// Register new Provider in Terrarium
func (s *VersionManagerService) Register(ctx context.Context, request *terrarium.RegisterProviderRequest) (*terrarium.Response, error) {
	log.Println("Registering new provider.")

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("provider.name", request.GetName()),
		attribute.String("provider.version", request.GetVersion()),
		attribute.String("provider.os",request.GetOs()),
		attribute.String("provider.arch",request.GetArch()),
	)

	providerID, err := attributevalue.Marshal(request.GetName())
	if err != nil {
		log.Println(err)
		span.RecordError(err)
		return nil, ProviderGetError
	}

	providerVersion, err := attributevalue.Marshal(request.GetVersion())
	if err != nil {
		log.Println(err)
		span.RecordError(err)
		return nil, ProviderGetError
	}

	providerOs, err := attributevalue.Marshal(request.GetOs())
	if err != nil {
		log.Println(err)
		span.RecordError(err)
		return nil, ProviderGetError
	}

	providerArch, err := attributevalue.Marshal(request.GetArch())
	if err != nil {
		log.Println(err)
		span.RecordError(err)
		return nil, ProviderGetError
	}

	res, err := s.Db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(VersionsTableName),
		Key: map[string]types.AttributeValue{
			"name":    providerID,
			"version": providerVersion,
			"os":      providerOs,
			"arch":    providerArch,
		},
	})
	if err != nil {
		log.Println(err)
		span.RecordError(err)
		return nil, ProviderGetError
	}

	if res.Item == nil {
		provider := Provider{
			Name:                request.GetName(),
			Version:             request.GetVersion(),
			Protocols:           request.GetProtocols(),
			OS:                  request.GetOs(),
			Arch:                request.GetArch(),
			Filename:            request.GetFilename(),
			DownloadURL:         request.GetDownloadUrl(),
			ShasumsURL:          request.GetShasum(),
			ShasumsSignatureURL: request.GetShasumsSignatureUrl(),
			Shasum:              request.GetShasum(),
			KeyID:               request.GetKeyId(),
			ASCIIArmor:          request.GetAsciiArmor(),
			TrustSignature:      request.GetTrustSignature(),
			Source:              request.GetSource(),
			SourceURL:           request.GetSourceUrl(),
			Description:         request.GetDescription(),
			SourceRepoUrl:       request.GetSourceRepoUrl(),
			Maturity:            request.GetMaturity().String(),
			CreatedOn:           time.Now().UTC().String(),
		}

		providerItem, err := attributevalue.MarshalMap(provider)
		if err != nil {
			log.Println(err)
			span.RecordError(err)
			return nil, MarshalProviderError
		}

		in := &dynamodb.PutItemInput{
			Item:      providerItem,
			TableName: aws.String(VersionsTableName),
		}

		if _, err = s.Db.PutItem(ctx, in); err != nil {
			log.Println(err)
			span.RecordError(err)
			return nil, ProviderRegisterError
		}
	} else {
		update := expression.Set(expression.Name("description"), expression.Value(request.GetDescription()))
		update.Set(expression.Name("source_repo_url"), expression.Value(request.GetSourceRepoUrl()))
		update.Set(expression.Name("maturity"), expression.Value(request.GetMaturity().String()))
		update.Set(expression.Name("modified_on"), expression.Value(time.Now().UTC().String()))
		update.Set(expression.Name("filename"), expression.Value(request.GetFilename()))
		update.Set(expression.Name("download_url"), expression.Value(request.GetDownloadUrl()))
		update.Set(expression.Name("shasums_url"), expression.Value(request.GetShasumsUrl()))
		update.Set(expression.Name("shasums_signature_url"), expression.Value(request.GetShasumsSignatureUrl()))
		update.Set(expression.Name("shasum"), expression.Value(request.GetShasum()))
		update.Set(expression.Name("key_id"), expression.Value(request.GetKeyId()))
		update.Set(expression.Name("ascii_armor"), expression.Value(request.GetAsciiArmor()))
		update.Set(expression.Name("trust_signature"), expression.Value(request.GetTrustSignature()))
		update.Set(expression.Name("source"), expression.Value(request.GetSource()))
		update.Set(expression.Name("source_url"), expression.Value(request.GetSourceUrl()))

		expr, err := expression.NewBuilder().WithUpdate(update).Build()
		if err != nil {
			log.Println(err)
			span.RecordError(err)
			return nil, ExpressionBuildError
		}

		in := &dynamodb.UpdateItemInput{
			TableName: aws.String(VersionsTableName),
			Key: map[string]types.AttributeValue{
				"name":    providerID,
				"version": providerVersion,
				"os":      providerOs,
				"arch":    providerArch,
			},
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
		}

		_, err = s.Db.UpdateItem(ctx, in)

		if err != nil {
			log.Println(err)
			span.RecordError(err)
			return nil, ProviderUpdateError
		}
	}

	log.Println("New provider registered.")
	return ProviderRegistered, nil
}

// ListProviderVersions Retrieve all versions of a given provider and return an array of versionItems.
// Only versions that have been published should be reported.
func (s *VersionManagerService) ListProviderVersions(ctx context.Context, request *services.ProviderName) (*services.ProviderVersionsResponse, error) {

	filter := expression.And(
		expression.Name("name").Equal(expression.Value(request.GetProvider())),
		expression.Name("published_on").AttributeExists())
	projection := expression.NamesList(expression.Name("version"), expression.Name("protocols"), expression.Name("os"), expression.Name("arch"))

	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(projection).Build()
	if err != nil {
		log.Printf("Expression Builder failed creation: %v", err)
		return nil, err
	}

	scanInputs := &dynamodb.ScanInput{
		TableName:                 aws.String(VersionsTableName), // CHECK LATERRRR
		ProjectionExpression:      expr.Projection(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}

	response, err := s.Db.Scan(ctx, scanInputs)
	if err != nil {
		log.Printf("ScanInput failed: %v", err)
		return nil, err
	}

	grpcResponse := services.ProviderVersionsResponse{}

	if response.Items != nil {
		for _, item := range response.Items {
			versionItem := &services.VersionItem{}
			if err := attributevalue.UnmarshalMap(item, &versionItem); err != nil {
				log.Printf("UnmarshalMap failed: %v", err)
				return nil, err
			}
			grpcResponse.Versions = append(grpcResponse.Versions, versionItem)
		}
	}

	// Validate and sort semantic versions
	var semverList versions.List
	for _, providerVersion := range grpcResponse.Versions {
		parsedVersion, err := versions.ParseVersion(providerVersion.Version)
		if err != nil {
			log.Printf("Skipping invalid semantic version: %v", providerVersion.Version)
		} else {
			semverList = append(semverList, parsedVersion)
		}
	}
	semverList.Sort()

	// Update grpcResponse with sorted versions
	var sortedVersions []*services.VersionItem

	for _, parsedVersion := range semverList {
		for _, providerVersion := range grpcResponse.Versions {
			if providerVersion.Version == parsedVersion.String() {
				sortedVersions = append(sortedVersions, providerVersion)
				break
			}
		}
	}
	grpcResponse.Versions = sortedVersions

	return &grpcResponse, nil
}

func (s *VersionManagerService) GetVersionData(ctx context.Context, request *services.VersionDataRequest) (*services.PlatformMetadataResponse, error) {

	providerID, err := attributevalue.Marshal(request.GetName())
	if err != nil {
		log.Println(err)
		return nil, ProviderGetError
	}

	providerVersion, err := attributevalue.Marshal(request.GetVersion())
	if err != nil {
		log.Println(err)
		return nil, ProviderGetError
	}

	providerOS, err := attributevalue.Marshal(request.GetOs())
	if err != nil {
		log.Println(err)
		return nil, ProviderGetError
	}

	providerArch, err := attributevalue.Marshal(request.GetArch())
	if err != nil {
		log.Println(err)
		return nil, ProviderGetError
	}

	response, err := s.Db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(VersionsTableName),
		Key: map[string]types.AttributeValue{
			"name":    providerID,
			"version": providerVersion,
			"os":      providerOS,
			"arch":    providerArch,
		},
	})
	if err != nil {
		log.Println(err)
		return nil, ProviderGetError
	} else {
		providerMetadata := &services.PlatformMetadataResponse{}
		err = attributevalue.UnmarshalMap(response.Item, &providerMetadata)
		if err != nil {
			log.Println(err)
			return nil, MarshalProviderError
		}
		return providerMetadata, nil
	}
}

func (s *VersionManagerService) ListProviders(ctx context.Context, request *services.ListProvidersRequest) (*services.ListProvidersResponse, error) {

	scanInputs := &dynamodb.ScanInput{
		TableName: aws.String(VersionsTableName),
	}

	response, err := s.Db.Scan(ctx, scanInputs)
	if err != nil {
		log.Printf("ScanInput failed: %v", err)
		return nil, err
	}

	grpcResponse := services.ListProvidersResponse{}

	if response.Items != nil {
		for _, item := range response.Items {
			if providerMetadata, err := unmarshalProvider(item); err != nil {
				return nil, err
			} else {
				grpcResponse.Providers = append(grpcResponse.Providers, providerMetadata)
			}
		}
	}

	return &grpcResponse, nil
}

func (s *VersionManagerService) GetProvider(ctx context.Context, request *services.ProviderName) (*services.GetProviderResponse, error) {

	filter := expression.Name("name").Equal(expression.Value(request.GetProvider()))
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		log.Printf("Expression Builder failed creation: %v", err)
		return nil, err
	}

	scanInputs := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(VersionsTableName),
	}

	response, err := s.Db.Scan(ctx, scanInputs)
	if err != nil {
		log.Printf("ScanInput failed: %v", err)
		return nil, err
	}

	if response.Items == nil || len(response.Items) < 1 {
		return nil, fmt.Errorf("provider not found '%v'", request.GetProvider())
	}

	grpcResponse := services.GetProviderResponse{}

	if providerMetadata, err := unmarshalProvider(response.Items[0]); err != nil {
		return nil, err
	} else {
		grpcResponse.Provider = providerMetadata
		return &grpcResponse, nil
	}
}

func unmarshalProvider(item map[string]types.AttributeValue) (*services.ListProviderItem, error) {
	provider := Provider{}
	if err := attributevalue.UnmarshalMap(item, &provider); err != nil {
		log.Printf("UnmarshalMap failed: %v", err)
		return nil, err
	}
	providerAddress := strings.Split(provider.Name, "/")

	result := services.ListProviderItem{
		Organization: providerAddress[0],
		Name:         providerAddress[1],
		Description:  provider.Description,
		SourceUrl:    provider.Source,
		Maturity:     terrarium.Maturity(terrarium.Maturity_value[provider.Maturity]),
	}

	return &result, nil
}

// GetProviderVersionsSchema returns CreateTableInput that can be used to create table if it does not exist
func GetProviderVersionsSchema(table string) *dynamodb.CreateTableInput {
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
			{
				AttributeName: aws.String("os"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("arch"),
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
			{
				AttributeName: aws.String("os"),
				KeyType:       types.KeyTypeRange,
			},
			{
				AttributeName: aws.String("arch"),
				KeyType:       types.KeyTypeRange,
			},
		},
		TableName:   aws.String(table),
		BillingMode: types.BillingModePayPerRequest,
	}
}
