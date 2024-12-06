package tag_manager

import (
	"context"
	"log"
	"time"

	"github.com/terrariumcloud/terrarium/internal/module/services"
	"github.com/terrariumcloud/terrarium/internal/module/services/registrar"

	"github.com/terrariumcloud/terrarium/internal/storage"
	terrarium "github.com/terrariumcloud/terrarium/pkg/terrarium/module"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	DefaultTagTableName       = "terrarium-module-tags-ciedev-4757"
	DefaultTagManagerEndpoint = "tag_manager:3001"
)

var (
	TagTableName                      = DefaultTagTableName
	TagManagerEndpoint                = DefaultTagManagerEndpoint
	TagPublished                      = &terrarium.Response{Message: "Tag published."}
	ModuleTagTableInitializationError = status.Error(codes.Unknown, "Failed to initialize table for tags.")
	MarshalModuleTagError             = status.Error(codes.Unknown, "Failed to marshal module tags for Dynamodb.")
	PublishModuleTagError             = status.Error(codes.Unknown, "Failed to publish module tag.")
	UpdateModuleTagError              = status.Error(codes.Unknown, "Failed to update module tag.")
	ConnectToTagManagerError          = status.Error(codes.Unknown, "Failed to connect to TagManager service.")
)

type TagManagerService struct {
	services.UnimplementedTagManagerServer
	Db     storage.DynamoDBTableCreator
	Table  string
	Schema *dynamodb.CreateTableInput
}

type ModuleTag struct {
	Name       string   `json:"name" bson:"name" dynamodbav:"name"`
	Tags       []string `json:"tags" bson:"tags" dynamodbav:"tags"`
	CreatedOn  string   `json:"created_on" bson:"created_on" dynamodbav:"created_on"`
	ModifiedOn string   `json:"modified_on" bson:"modified_on" dynamodbav:"modified_on"`
}

// RegisterWithServer registers TagManagerService with grpc server
func (s *TagManagerService) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	if err := storage.InitializeDynamoDb(s.Table, s.Schema, s.Db); err != nil {
		log.Println(err)
		return ModuleTagTableInitializationError
	}

	services.RegisterTagManagerServer(grpcServer, s)

	return nil
}

func (s *TagManagerService) PublishTag(ctx context.Context, request *terrarium.PublishTagRequest) (*terrarium.Response, error) {
	log.Println("Publish module tag.")

	name, err := attributevalue.Marshal(request.GetName())
	if err != nil {
		log.Println(err)
		return nil, registrar.ModuleGetError
	}

	res, err := s.Db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(TagTableName),
		Key: map[string]types.AttributeValue{
			"name": name,
		},
	})

	if err != nil {
		log.Println(err)
		return nil, registrar.ModuleGetError
	}

	if res.Item == nil {
		ms := ModuleTag{
			Name:       request.GetName(),
			Tags:       request.GetTags(),
			CreatedOn:  time.Now().UTC().String(),
			ModifiedOn: time.Now().UTC().String(),
		}

		av, err := attributevalue.MarshalMap(ms)
		if err != nil {
			log.Println(err)
			return nil, MarshalModuleTagError
		}

		in := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(TagTableName),
		}

		if _, err = s.Db.PutItem(ctx, in); err != nil {
			log.Println(err)
			return nil, PublishModuleTagError
		}
	} else {
		update := expression.Set(expression.Name("tags"), expression.Value(request.GetTags()))
		update.Set(expression.Name("modified_on"), expression.Value(time.Now().UTC().String()))
		expr, err := expression.NewBuilder().WithUpdate(update).Build()

		if err != nil {
			log.Println(err)
			return nil, registrar.ExpressionBuildError
		}

		in := &dynamodb.UpdateItemInput{
			TableName: aws.String(TagTableName),
			Key: map[string]types.AttributeValue{
				"name": name,
			},
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
		}

		_, err = s.Db.UpdateItem(ctx, in)

		if err != nil {
			log.Println(err)
			return nil, UpdateModuleTagError
		}
	}

	log.Println("Module tags published.")
	return TagPublished, nil
}

// GetTagsSchema returns CreateTableInput that can be used to create table if it does not exist
func GetTagsSchema(table string) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("name"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName:   aws.String(table),
		BillingMode: types.BillingModePayPerRequest,
	}
}
