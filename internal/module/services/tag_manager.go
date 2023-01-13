package services

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/terrariumcloud/terrarium/internal/storage"
	terrarium "github.com/terrariumcloud/terrarium/pkg/terrarium/module"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk_tracy "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	//"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	//"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const OtelTracerName = "terrarium-tracer"

const (
	DefaultTagTableName       = "terrarium-module-tags"
	DefaultTagManagerEndpoint = "tag_manager:3001"
)

var (
	TagTableName                      string = DefaultTagTableName
	TagManagerEndpoint                string = DefaultTagManagerEndpoint
	TagPublished                             = &terrarium.Response{Message: "Tag published."}
	ModuleTagTableInitializationError        = status.Error(codes.Unknown, "Failed to initialize table for tags.")
	MarshalModuleTagError                    = status.Error(codes.Unknown, "Failed to marshal module tags for Dynamodb.")
	PublishModuleTagError                    = status.Error(codes.Unknown, "Failed to publish module tag.")
	UpdateModuleTagError                     = status.Error(codes.Unknown, "Failed to update module tag.")
	ConnectToTagManagerError                 = status.Error(codes.Unknown, "Failed to connect to TagManager service.")
)

type TagManagerService struct {
	UnimplementedTagManagerServer
	Db     dynamodbiface.DynamoDBAPI
	Table  string
	Schema *dynamodb.CreateTableInput
}

type ModuleTag struct {
	Name       string   `json:"name" bson:"name" dynamodbav:"name"`
	Tags       []string `json:"tags" bson:"tags" dynamodbav:"tags"`
	CreatedOn  string   `json:"created_on" bson:"created_on" dynamodbav:"created_on"`
	ModifiedOn string   `json:"modified_on" bson:"modified_on" dynamodbav:"modified_on"`
}

func newExporter(w io.Writer) (sdk_tracy.SpanExporter, error) {
	log.Printf("Returning newExporter")
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		// Do not print timestamps for the demo.
	)
}

func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(OtelTracerName),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "dev"),
		),
	)
	log.Printf("Returning newResource")
	return r
}

// RegisterWithServer registers TagManagerService with grpc server
func (s *TagManagerService) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	if err := storage.InitializeDynamoDb(s.Table, s.Schema, s.Db); err != nil {
		log.Println(err)
		return ModuleTagTableInitializationError
	}

	RegisterTagManagerServer(grpcServer, s)

	return nil
}

func (s *TagManagerService) PublishTag(ctx context.Context, request *terrarium.PublishTagRequest) (*terrarium.Response, error) {
	log.Println("Publish module tag.")

	f, err := os.Create("traces.txt")
	if err != nil {
		log.Fatal(err)
		log.Printf("Failed 117")
	}
	defer f.Close()

	exp, err := newExporter(f)
	if err != nil {
		log.Printf("Failed 122")
		log.Fatal(err)
	}

	tp := sdk_tracy.NewTracerProvider(
		sdk_tracy.WithBatcher(exp),
		sdk_tracy.WithResource(newResource()),
	)
	// defer func() {
	// 	log.Printf("Failed 132")
	// 	if err := tp.Shutdown(context.Background()); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	log.Printf("Will execute TP")
	otel.SetTracerProvider(tp)

	var span trace.Span

	ctx, span = otel.Tracer(OtelTracerName).Start(ctx, "MiTracer")
	//defer span.End()

	if span != nil {
		log.Printf("Started Tracer")
	}

	/* … */

	res, err := s.Db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TagTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"name": {S: aws.String(request.GetName())},
		},
	})

	if err != nil {
		log.Println(err)
		return nil, ModuleGetError
	}

	if res.Item == nil {
		ms := ModuleTag{
			Name:       request.GetName(),
			Tags:       request.GetTags(),
			CreatedOn:  time.Now().UTC().String(),
			ModifiedOn: time.Now().UTC().String(),
		}

		av, err := dynamodbattribute.MarshalMap(ms)

		if err != nil {
			log.Println(err)
			return nil, MarshalModuleTagError
		}

		in := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(TagTableName),
		}

		if _, err = s.Db.PutItem(in); err != nil {
			log.Println(err)
			return nil, PublishModuleTagError
		}
	} else {
		update := expression.Set(expression.Name("tags"), expression.Value(request.GetTags()))
		update.Set(expression.Name("modified_on"), expression.Value(time.Now().UTC().String()))
		expr, err := expression.NewBuilder().WithUpdate(update).Build()

		if err != nil {
			log.Println(err)
			return nil, ExpressionBuildError
		}

		in := &dynamodb.UpdateItemInput{
			TableName: aws.String(TagTableName),
			Key: map[string]*dynamodb.AttributeValue{
				"name": {S: aws.String(request.GetName())}},
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
		}

		_, err = s.Db.UpdateItem(in)

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
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("name"),
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       aws.String("HASH"),
			},
		},
		TableName:   aws.String(table),
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
	}
}
