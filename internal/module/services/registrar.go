package services

import (
	"context"
	"log"
	"time"

	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
)

const (
	DefaultRegistrarTableName              = "terrarium-module-stream" //TODO: rename to terrarium-modules
	DefaultRegistrarServiceDefaultEndpoint = "registrar:3001"
)

var RegistrarTableName string = DefaultRegistrarTableName
var RegistrarServiceEndpoint string = DefaultRegistrarServiceDefaultEndpoint

type RegistrarService struct {
	UnimplementedRegistrarServer
	Db dynamodbiface.DynamoDBAPI
}

type ModuleStream struct { //TODO: rename to Module
	ID          interface{} `json:"id" bson:"_id" dynamodbav:"_id"`
	Name        string      `json:"name" bson:"name" dynamodbav:"name"`
	Description string      `json:"description" bson:"description" dynamodbav:"description"`
	SourceUrl   string      `json:"source_url" bson:"source_url" dynamodbav:"source_url"`
	Maturity    string      `json:"maturity" bson:"maturity" dynamodbav:"maturity"`
	CreatedOn   string      `json:"created_on" bson:"created_on" dynamodbav:"created_on"`
}

// Register new Module in Terrarium
func (s *RegistrarService) Register(ctx context.Context, request *terrarium.RegisterModuleRequest) (*terrarium.TransactionStatusResponse, error) {
	log.Println("Registering new module.")

	ms := ModuleStream{
		ID:          uuid.NewString(),
		Name:        request.GetName(),
		Description: request.GetDescription(),
		SourceUrl:   request.GetSourceUrl(),
		Maturity:    request.GetMaturity().String(),
		CreatedOn:   time.Now().UTC().String(),
	}

	av, err := dynamodbattribute.MarshalMap(ms)

	if err != nil {
		log.Printf("Failed to marshal: %s", err.Error())
		return MarshalModuleError, err
	}

	in := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(RegistrarTableName),
	}

	if _, err = s.Db.PutItem(in); err != nil {
		log.Printf("Failed to put item: %s", err.Error())
		return ModuleNotRegistered, err
	}

	log.Println("New module registered.")
	return ModuleRegistered, nil
}
