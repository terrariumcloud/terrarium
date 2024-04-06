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
	DefaultProviderVersionsTableName      = "terrarium-providers"
	DefaultProviderVersionManagerEndpoint = "provider_version_manager:3001"
)

var (
	VersionsTableName      = DefaultProviderVersionsTableName
	VersionManagerEndpoint = DefaultProviderVersionManagerEndpoint

	ProviderRegistered = &terrarium.Response{Message: "Provider registered successfully."}
	VersionPublished   = &terrarium.Response{Message: "Version published."}
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
	Name          string                    `json:"name" bson:"name" dynamodbav:"name"`
	Version       string                    `json:"version" bson:"version" dynamodbav:"version"`
	Protocols     []string                  `json:"protocols" bson:"protocols" dynamodbav:"protocols"`
	Platforms     []*terrarium.PlatformItem `json:"platforms" bson:"platforms" dynamodbav:"platforms"`
	Description   string                    `json:"description" bson:"description" dynamodbav:"description"`
	SourceRepoUrl string                    `json:"source_repo_url" bson:"source_repo_url" dynamodbav:"source_repo_url"`
	Maturity      string                    `json:"maturity" bson:"maturity" dynamodbav:"maturity"`
	CreatedOn     string                    `json:"created_on" bson:"created_on" dynamodbav:"created_on"`
	ModifiedOn    string                    `json:"modified_on" bson:"modified_on" dynamodbav:"modified_on"`
	PublishedOn   string                    `json:"published_on" bson:"published_on" dynamodbav:"published_on"`
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

	return map[string]types.AttributeValue{
		"name":    providerName,
		"version": providerVersion,
	}, nil
}

// AbortProviderVersion removes a Version of a Provider.
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
func (s *VersionManagerService) PublishVersion(ctx context.Context, request *services.TerminateVersionRequest) (*terrarium.Response, error) {
	log.Println("Publishing provider version.")

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("provider.name", request.Provider.GetName()),
		attribute.String("provider.version", request.Provider.GetVersion()),
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

	res, err := s.Db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(VersionsTableName),
		Key: map[string]types.AttributeValue{
			"name":    providerID,
			"version": providerVersion,
		},
	})
	if err != nil {
		log.Println(err)
		span.RecordError(err)
		return nil, ProviderGetError
	}

	if res.Item == nil {
		provider := Provider{
			Name:          request.GetName(),
			Version:       request.GetVersion(),
			Protocols:     request.GetProtocols(),
			Platforms:     request.GetPlatforms(),
			Description:   request.GetDescription(),
			SourceRepoUrl: request.GetSourceRepoUrl(),
			Maturity:      request.GetMaturity().String(),
			CreatedOn:     time.Now().UTC().String(),
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
		update.Set(expression.Name("platforms"), expression.Value(request.GetPlatforms()))
		update.Set(expression.Name("protocols"), expression.Value(request.GetProtocols()))
		update.Set(expression.Name("modified_on"), expression.Value(time.Now().UTC().String()))

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
	projection := expression.NamesList(expression.Name("version"), expression.Name("protocols"), expression.Name("platforms"))

	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(projection).Build()
	if err != nil {
		log.Printf("Expression Builder failed creation: %v", err)
		return nil, err
	}

	scanInputs := &dynamodb.ScanInput{
		TableName:                 aws.String(VersionsTableName),
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
			versionItem, err := unmarshalProviderVersionItem(item)
			if err != nil {
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

	response, err := s.Db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(VersionsTableName),
		Key: map[string]types.AttributeValue{
			"name":    providerID,
			"version": providerVersion,
		},
	})
	if err != nil {
		log.Println(err)
		return nil, ProviderGetError
	} else {
		providerMetadata, err := unmarshalProviderMetadata(response.Item, request.Os, request.Arch)
		if err != nil {
			log.Println(err)
			return nil, MarshalProviderError
		}
		return providerMetadata, nil
	}
}

func (s *VersionManagerService) ListProviders(ctx context.Context, request *services.ListProvidersRequest) (*services.ListProvidersResponse, error) {

	// Initialize a map to store providers uniquely
	uniqueProviders := make(map[string]*services.ListProviderItem)

	projection := expression.NamesList(expression.Name("name"), expression.Name("description"), expression.Name("maturity"), expression.Name("source_repo_url"))

	expr, err := expression.NewBuilder().WithProjection(projection).Build()
	if err != nil {
		log.Printf("Expression Builder failed creation: %v", err)
		return nil, err
	}

	scanInputs := &dynamodb.ScanInput{
		TableName:                aws.String(VersionsTableName),
		ProjectionExpression:     expr.Projection(),
		ExpressionAttributeNames: expr.Names(),
	}

	response, err := s.Db.Scan(ctx, scanInputs)
	if err != nil {
		log.Printf("ScanInput failed: %v", err)
		return nil, err
	}

	if response.Items != nil {
		for _, item := range response.Items {
			if providerMetadata, err := unmarshalProvider(item); err != nil {
				return nil, err
			} else {
				key := providerMetadata.Name
				// Check if the provider already exists in the map
				if _, ok := uniqueProviders[key]; !ok {
					// Add the provider to the map if it doesn't exist
					uniqueProviders[key] = providerMetadata
				}
			}
		}
	}

	providersList := make([]*services.ListProviderItem, 0, len(uniqueProviders))
	for _, provider := range uniqueProviders {
		providersList = append(providersList, provider)
	}

	grpcResponse := services.ListProvidersResponse{
		Providers: providersList,
	}

	return &grpcResponse, nil
}

func (s *VersionManagerService) GetProvider(ctx context.Context, request *services.ProviderName) (*services.GetProviderResponse, error) {

	filter := expression.Name("name").Equal(expression.Value(request.GetProvider()))
	projection := expression.NamesList(expression.Name("name"), expression.Name("description"), expression.Name("maturity"), expression.Name("source_repo_url"))
	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(projection).Build()
	if err != nil {
		log.Printf("Expression Builder failed creation: %v", err)
		return nil, err
	}

	scanInputs := &dynamodb.ScanInput{
		ProjectionExpression:      expr.Projection(),
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
		Organization:  providerAddress[0],
		Name:          providerAddress[1],
		Description:   provider.Description,
		SourceRepoUrl: provider.SourceRepoUrl,
		Maturity:      terrarium.Maturity(terrarium.Maturity_value[provider.Maturity]),
	}

	return &result, nil
}

func unmarshalProviderMetadata(item map[string]types.AttributeValue, os, arch string) (*services.PlatformMetadataResponse, error) {
	provider := Provider{}
	if err := attributevalue.UnmarshalMap(item, &provider); err != nil {
		log.Printf("UnmarshalMap failed: %v", err)
		return nil, err
	}

	for _, platform := range provider.Platforms {
		if platform.Os == os && platform.Arch == arch {

			var gpgPublicKeys []*services.GPGPublicKey
			for _, key := range platform.SigningKeys.GpgPublicKeys {
				gpgPublicKeys = append(gpgPublicKeys, &services.GPGPublicKey{
					KeyId:          key.KeyId,
					AsciiArmor:     key.AsciiArmor,
					TrustSignature: key.TrustSignature,
					Source:         key.Source,
					SourceUrl:      key.SourceUrl,
				})
			}

			return &services.PlatformMetadataResponse{
				Protocols:           provider.Protocols,
				Os:                  platform.Os,
				Arch:                platform.Arch,
				Filename:            platform.Filename,
				DownloadUrl:         platform.DownloadUrl,
				ShasumsUrl:          platform.ShasumsUrl,
				ShasumsSignatureUrl: platform.ShasumsSignatureUrl,
				Shasum:              platform.Shasum,
				SigningKeys: &services.SigningKeys{
					GpgPublicKeys: gpgPublicKeys,
				},
			}, nil
		}
	}

	err := fmt.Errorf("requested os '%s' and arch '%s' doesn't exist", os, arch)
	return nil, err
}

func unmarshalProviderVersionItem(item map[string]types.AttributeValue) (*services.VersionItem, error) {
	provider := Provider{}
	if err := attributevalue.UnmarshalMap(item, &provider); err != nil {
		log.Printf("UnmarshalMap failed: %v", err)
		return nil, err
	}

	var platforms []*services.Platform
	for _, platformItem := range provider.Platforms {
		platform := &services.Platform{
			Os:   platformItem.Os,
			Arch: platformItem.Arch,
		}
		platforms = append(platforms, platform)
	}

	result := services.VersionItem{
		Version:   provider.Version,
		Protocols: provider.Protocols,
		Platforms: platforms,
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
