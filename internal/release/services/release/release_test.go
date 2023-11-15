package release

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/terrariumcloud/terrarium/internal/release/services"
	"github.com/terrariumcloud/terrarium/internal/storage/mocks"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/release"
	"google.golang.org/grpc"
)

// Test_RegisterReleaseWithServer checks:
// - if there was no error with table init
// - if error is returned when Table initialization fails
func Test_RegisterReleaseWithServer(t *testing.T) {
	t.Parallel()

	t.Run("when there is no error with table init", func(t *testing.T) {
		db := &mocks.DynamoDB{}

		rs := &ReleaseService{
			Db: db,
		}

		s := grpc.NewServer(*new([]grpc.ServerOption)...)

		err := rs.RegisterWithServer(s)

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}

		if db.DescribeTableInvocations != 1 {
			t.Errorf("Expected 1 call to DescribeTable, got %v.", db.DescribeTableInvocations)
		}

		if db.CreateTableInvocations != 0 {
			t.Errorf("Expected no calls to CreateTable, got %v.", db.CreateTableInvocations)
		}
	})

	t.Run("when Table initialization fails", func(t *testing.T) {

		db := &mocks.DynamoDB{
			DescribeTableErrors: []error{errors.New("some error")},
		}

		rs := &ReleaseService{
			Db: db,
		}

		s := grpc.NewServer(*new([]grpc.ServerOption)...)

		err := rs.RegisterWithServer(s)

		if err != ReleaseTableInitializationError {
			t.Errorf("Expected %v, got %v.", ReleaseTableInitializationError, err)
		}

		if db.DescribeTableInvocations != 1 {
			t.Errorf("Expected 1 call to DescribeTable, got %v.", db.DescribeTableInvocations)
		}

		if db.CreateTableInvocations != 0 {
			t.Errorf("Expected 0 calls to CreateTable, got %v.", db.CreateTableInvocations)
		}
	})
}

// Test_PublishRelease checks:
// - if correct response is returned when release is published
// - if error is returned when PutItem fails
func Test_PublishRelease(t *testing.T) {
	t.Parallel()

	t.Run("when new release is published", func(t *testing.T) {
		db := &mocks.DynamoDB{}

		svc := &ReleaseService{Db: db}

		req := &release.PublishRequest{
			Type:         "test type",
			Organization: "test org",
			Name:         "test",
			Version:      "v1.0.0",
			Description:  "test desc",
			Links: []*release.Link{
				{
					Title: "test title",
					Url:   "http://www.google.com",
				},
			},
		}

		res, err := svc.Publish(context.TODO(), req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if db.PutItemInvocations != 1 {
			t.Errorf("Expected 1 call to PutItem, got %d", db.PutItemInvocations)
		}

		if db.TableName != ReleaseTableName {
			t.Errorf("Expected tableName to be %s, got %s", ReleaseTableName, db.TableName)
		}

		if res != ReleasePublished {
			t.Errorf("Expected %v, got %v.", ReleasePublished, res)
		}
	})

	t.Run("when PutItem fails", func(t *testing.T) {
		db := &mocks.DynamoDB{PutItemError: errors.New("some error")}

		svc := &ReleaseService{Db: db}

		req := &release.PublishRequest{
			Type:         "test type",
			Organization: "test org",
			Name:         "test",
			Version:      "v1.0.0",
			Description:  "test desc",
			Links: []*release.Link{
				{
					Title: "test title",
					Url:   "http://www.google.com",
				},
			},
		}

		res, err := svc.Publish(context.TODO(), req)

		if res != nil {
			t.Errorf("Expected no response, got %v", res)
		}

		if db.PutItemInvocations != 1 {
			t.Errorf("Expected 1 call to PutItem, got %v", db.PutItemInvocations)
		}

		if db.TableName != ReleaseTableName {
			t.Errorf("Expected tableName to be %v, got %v.", ReleaseTableName, db.TableName)
		}

		if err != PublishReleaseError {
			t.Errorf("Expected %v, got %v.", PublishReleaseError, err)
		}
	})

}

// Test_ListReleases checks:
// - if correct response is returned when releases are fetched
func Test_ListReleases(t *testing.T) {
	t.Parallel()

	t.Run("Listing releases", func(t *testing.T) {
		db := &mocks.DynamoDB{
			ScanOut: &dynamodb.ScanOutput{
				Items: []map[string]types.AttributeValue{
					{
						"CreatedAt":    &types.AttributeValueMemberS{Value: "2022-11-17 15:11:35.198401764 +0000 UTC"},
						"Type":         &types.AttributeValueMemberS{Value: "module"},
						"Organization": &types.AttributeValueMemberS{Value: "cie"},
						"Name":         &types.AttributeValueMemberS{Value: "test name"},
						"Version":      &types.AttributeValueMemberS{Value: "1.0.0"},
						"Description":  &types.AttributeValueMemberS{Value: "test desc"},
						"Links": &types.AttributeValueMemberL{
							Value: []types.AttributeValue{
								&types.AttributeValueMemberM{
									Value: map[string]types.AttributeValue{
										"Title": &types.AttributeValueMemberS{Value: "test title"},
										"URL":   &types.AttributeValueMemberS{Value: "https://example.com"},
									},
								},
							},
						},
					},
				},
			},
		}

		svc := &ReleaseService{Db: db}

		req := services.ListReleasesRequest{
			Types:         []string{"module"},
			Organizations: []string{"cie"},
		}
		res, err := svc.ListReleases(context.TODO(), &req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		expectedReleases := []*services.Release{
			{
				CreatedAt:    "2022-11-17 15:11:35.198401764 +0000 UTC",
				Type:         "module",
				Organization: "cie",
				Name:         "test name",
				Version:      "1.0.0",
				Description:  "test desc",
				Links: []*release.Link{
					{
						Title: "test title",
						Url:   "https://example.com",
					},
				},
			},
		}
		if !reflect.DeepEqual(res.Releases, expectedReleases) {
			t.Errorf("Releases do not match, got %v, want %v", res.Releases, expectedReleases)
		}

	})

}
