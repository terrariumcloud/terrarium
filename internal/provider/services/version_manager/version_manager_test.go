package version_manager

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/terrariumcloud/terrarium/internal/provider/services"
	"github.com/terrariumcloud/terrarium/internal/storage/mocks"
	terrarium "github.com/terrariumcloud/terrarium/pkg/terrarium/provider"

	"google.golang.org/grpc"
)

// Test_RegisterProvider checks:
// - if correct response is returned when Provider is registered
// - if there was no error when version already exists
// - if error was returned when GetItem fails
// - if error is returned when marshal fails
// - if error is returned when PutItem fails
func Test_RegisterProvider(t *testing.T) {
	t.Parallel()

	t.Run("when new version is created", func(t *testing.T) {
		db := &mocks.DynamoDB{
			GetItemOuts: []*dynamodb.GetItemOutput{{}},
		}

		svc := &VersionManagerService{Db: db}

		req := &terrarium.RegisterProviderRequest{
			Name:          "test-org/test-provider2",
			Version:       "2.0.0",
			Protocols:     []string{"5.1", "5.2"},
			Platforms: []*terrarium.PlatformItem{
				{
					Os:                  "linux",
					Arch:                "amd64",
					Filename:            "example-provider-linux-amd64",
					DownloadUrl:         "https://example.com/download/linux-amd64",
					ShasumsUrl:          "https://example.com/shasums/linux-amd64",
					ShasumsSignatureUrl: "https://example.com/shasums-signature/linux-amd64",
					Shasum:              "1234567890",
					SigningKeys: &terrarium.SigningKeys{
						GpgPublicKeys: []*terrarium.GPGPublicKey{
							{
								KeyId:          "keyid1",
								AsciiArmor:     "asciiarmor1",
								TrustSignature: "trustsignature1",
								Source:         "source1",
								SourceUrl:      "https://example.com/sourceurl1",
							},
						},
					},
				},
				{
					Os:                  "windows",
					Arch:                "amd64",
					Filename:            "example-provider-windows-amd64",
					DownloadUrl:         "https://example.com/download/windows-amd64",
					ShasumsUrl:          "https://example.com/shasums/windows-amd64",
					ShasumsSignatureUrl: "https://example.com/shasums-signature/windows-amd64",
					Shasum:              "0987654321",
					SigningKeys: &terrarium.SigningKeys{
						GpgPublicKeys: []*terrarium.GPGPublicKey{
							{
								KeyId:          "keyid3_updated",
								AsciiArmor:     "asciiarmor3",
								TrustSignature: "trustsignature3",
								Source:         "source3",
								SourceUrl:      "https://example.com/sourceurl3",
							},
						},
					},
				},
			},
			Description:   "Example Provider Description",
			SourceRepoUrl: "https://example.com/source-repo",
			Maturity:      terrarium.Maturity_ALPHA,
		}

		res, err := svc.Register(context.TODO(), req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if db.GetItemInvocations != 1 {
			t.Errorf("Expected 1 call to GetItem, got %d", db.GetItemInvocations)
		}

		if db.PutItemInvocations != 1 {
			t.Errorf("Expected 1 call to PutItem, got %d", db.PutItemInvocations)
		}

		if db.TableName != VersionsTableName {
			t.Errorf("Expected tableName to be %s, got %s", VersionsTableName, db.TableName)
		}

		if res != ProviderRegistered {
			t.Errorf("Expected %v, got %v.", ProviderRegistered, res)
		}
	})

	t.Run("when version already exists", func(t *testing.T) {
		name := "test"
		emptyString := ""
		db := &mocks.DynamoDB{
			GetItemOuts: []*dynamodb.GetItemOutput{
				{
					Item: map[string]types.AttributeValue{
						"name":                  services.MustMarshallString(name, t),
						"description":           services.MustMarshallString(emptyString, t),
						"source_url":            services.MustMarshallString(emptyString, t),
						"maturity":              services.MustMarshallString(emptyString, t),
						"modified_on":           services.MustMarshallString(emptyString, t),
						"source_repo_url":       services.MustMarshallString(emptyString, t),
						"protocols":            services.MustMarshallString(emptyString, t),
						"platforms":            services.MustMarshallString(emptyString, t),
					},
				},
			},
			UpdateItemOut: &dynamodb.UpdateItemOutput{},
		}

		svc := &VersionManagerService{Db: db}

		req := &terrarium.RegisterProviderRequest{
			Name:          "test-org/test-provider2",
			Version:       "2.0.0",
			Protocols:     []string{"5.1", "5.2"},
			Platforms: []*terrarium.PlatformItem{
				{
					Os:                  "linux",
					Arch:                "amd64",
					Filename:            "example-provider-linux-amd64",
					DownloadUrl:         "https://example.com/download/linux-amd64",
					ShasumsUrl:          "https://example.com/shasums/linux-amd64",
					ShasumsSignatureUrl: "https://example.com/shasums-signature/linux-amd64",
					Shasum:              "1234567890",
					SigningKeys: &terrarium.SigningKeys{
						GpgPublicKeys: []*terrarium.GPGPublicKey{
							{
								KeyId:          "keyid1",
								AsciiArmor:     "asciiarmor1",
								TrustSignature: "trustsignature1",
								Source:         "source1",
								SourceUrl:      "https://example.com/sourceurl1",
							},
						},
					},
				},
				{
					Os:                  "windows",
					Arch:                "amd64",
					Filename:            "example-provider-windows-amd64",
					DownloadUrl:         "https://example.com/download/windows-amd64",
					ShasumsUrl:          "https://example.com/shasums/windows-amd64",
					ShasumsSignatureUrl: "https://example.com/shasums-signature/windows-amd64",
					Shasum:              "0987654321",
					SigningKeys: &terrarium.SigningKeys{
						GpgPublicKeys: []*terrarium.GPGPublicKey{
							{
								KeyId:          "keyid3_updated",
								AsciiArmor:     "asciiarmor3",
								TrustSignature: "trustsignature3",
								Source:         "source3",
								SourceUrl:      "https://example.com/sourceurl3",
							},
						},
					},
				},
			},
			Description:   "Example Provider Description",
			SourceRepoUrl: "https://example.com/source-repo",
			Maturity:      terrarium.Maturity_ALPHA,
		}

		res, err := svc.Register(context.TODO(), req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if db.GetItemInvocations != 1 {
			t.Errorf("Expected 1 call to GetItem, got %d", db.GetItemInvocations)
		}

		if db.PutItemInvocations != 0 {
			t.Errorf("Expected 0 calls to PutItem, got %d", db.PutItemInvocations)
		}

		if db.UpdateItemInvocations != 1 {
			t.Errorf("Expected 1 call to UpdateItem, got %d", db.UpdateItemInvocations)
		}

		if db.TableName != VersionsTableName {
			t.Errorf("Expected tableName to be %s, got %s", VersionsTableName, db.TableName)
		}

		if res != ProviderRegistered {
			t.Errorf("Expected %v, got %v.", ProviderRegistered, res)
		}
	})

	t.Run("when GetItem fails", func(t *testing.T) {
		db := &mocks.DynamoDB{
			GetItemErrors: []error{errors.New("some error")},
		}

		svc := &VersionManagerService{Db: db}

		req := &terrarium.RegisterProviderRequest{
			Name:          "test-org/test-provider2",
			Version:       "2.0.0",
			Protocols:     []string{"5.1", "5.2"},
			Platforms: []*terrarium.PlatformItem{
				{
					Os:                  "linux",
					Arch:                "amd64",
					Filename:            "example-provider-linux-amd64",
					DownloadUrl:         "https://example.com/download/linux-amd64",
					ShasumsUrl:          "https://example.com/shasums/linux-amd64",
					ShasumsSignatureUrl: "https://example.com/shasums-signature/linux-amd64",
					Shasum:              "1234567890",
					SigningKeys: &terrarium.SigningKeys{
						GpgPublicKeys: []*terrarium.GPGPublicKey{
							{
								KeyId:          "keyid1",
								AsciiArmor:     "asciiarmor1",
								TrustSignature: "trustsignature1",
								Source:         "source1",
								SourceUrl:      "https://example.com/sourceurl1",
							},
						},
					},
				},
				{
					Os:                  "windows",
					Arch:                "amd64",
					Filename:            "example-provider-windows-amd64",
					DownloadUrl:         "https://example.com/download/windows-amd64",
					ShasumsUrl:          "https://example.com/shasums/windows-amd64",
					ShasumsSignatureUrl: "https://example.com/shasums-signature/windows-amd64",
					Shasum:              "0987654321",
					SigningKeys: &terrarium.SigningKeys{
						GpgPublicKeys: []*terrarium.GPGPublicKey{
							{
								KeyId:          "keyid3_updated",
								AsciiArmor:     "asciiarmor3",
								TrustSignature: "trustsignature3",
								Source:         "source3",
								SourceUrl:      "https://example.com/sourceurl3",
							},
						},
					},
				},
			},
			Description:   "Example Provider Description",
			SourceRepoUrl: "https://example.com/source-repo",
			Maturity:      terrarium.Maturity_ALPHA,
		}

		res, err := svc.Register(context.TODO(), req)

		if err == nil {
			t.Errorf("Expected an error")
		}

		if db.GetItemInvocations != 1 {
			t.Errorf("Expected 1 call to GetItem, got %d", db.GetItemInvocations)
		}

		if db.PutItemInvocations != 0 {
			t.Errorf("Expected 0 calls to PutItem, got %d", db.PutItemInvocations)
		}

		if res != nil {
			t.Errorf("Expected no response, got %v.", res)
		}
	})

	t.Run("when PutItem fails", func(t *testing.T) {
		db := &mocks.DynamoDB{
			GetItemOuts:  []*dynamodb.GetItemOutput{{}},
			PutItemError: errors.New("some error"),
		}

		svc := &VersionManagerService{Db: db}

		req := &terrarium.RegisterProviderRequest{
			Name:          "test-org/test-provider2",
			Version:       "2.0.0",
			Protocols:     []string{"5.1", "5.2"},
			Platforms: []*terrarium.PlatformItem{
				{
					Os:                  "linux",
					Arch:                "amd64",
					Filename:            "example-provider-linux-amd64",
					DownloadUrl:         "https://example.com/download/linux-amd64",
					ShasumsUrl:          "https://example.com/shasums/linux-amd64",
					ShasumsSignatureUrl: "https://example.com/shasums-signature/linux-amd64",
					Shasum:              "1234567890",
					SigningKeys: &terrarium.SigningKeys{
						GpgPublicKeys: []*terrarium.GPGPublicKey{
							{
								KeyId:          "keyid1",
								AsciiArmor:     "asciiarmor1",
								TrustSignature: "trustsignature1",
								Source:         "source1",
								SourceUrl:      "https://example.com/sourceurl1",
							},
						},
					},
				},
				{
					Os:                  "windows",
					Arch:                "amd64",
					Filename:            "example-provider-windows-amd64",
					DownloadUrl:         "https://example.com/download/windows-amd64",
					ShasumsUrl:          "https://example.com/shasums/windows-amd64",
					ShasumsSignatureUrl: "https://example.com/shasums-signature/windows-amd64",
					Shasum:              "0987654321",
					SigningKeys: &terrarium.SigningKeys{
						GpgPublicKeys: []*terrarium.GPGPublicKey{
							{
								KeyId:          "keyid3_updated",
								AsciiArmor:     "asciiarmor3",
								TrustSignature: "trustsignature3",
								Source:         "source3",
								SourceUrl:      "https://example.com/sourceurl3",
							},
						},
					},
				},
			},
			Description:   "Example Provider Description",
			SourceRepoUrl: "https://example.com/source-repo",
			Maturity:      terrarium.Maturity_ALPHA,
		}

		res, err := svc.Register(context.TODO(), req)

		if res != nil {
			t.Errorf("Expected no response, got %v", err)
		}

		if db.PutItemInvocations != 1 {
			t.Errorf("Expected 1 call to PutItem, got %d", db.PutItemInvocations)
		}

		if db.TableName != VersionsTableName {
			t.Errorf("Expected tableName to be %s, got %s", VersionsTableName, db.TableName)
		}

		if err != ProviderRegisterError {
			t.Errorf("Expected %v, got %v.", ProviderRegisterError, err)
		}
	})
}

// Test_RegisterVersionManagerWithServer checks:
// - if there was no error with table init
// - if error is returned when Table initialization fails
func Test_RegisterVersionManagerWithServer(t *testing.T) {
	t.Parallel()

	t.Run("when table init is successful", func(t *testing.T) {
		db := &mocks.DynamoDB{}

		vms := &VersionManagerService{Db: db}

		s := grpc.NewServer(*new([]grpc.ServerOption)...)

		err := vms.RegisterWithServer(s)

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
			CreateTableError:    errors.New("some error"),
		}

		vms := &VersionManagerService{Db: db}

		s := grpc.NewServer(*new([]grpc.ServerOption)...)

		err := vms.RegisterWithServer(s)

		if err != ProviderVersionsTableInitializationError {
			t.Errorf("Expected %v, got %v.", ProviderVersionsTableInitializationError, err)
		}

		if db.DescribeTableInvocations != 1 {
			t.Errorf("Expected 1 call to DescribeTable, got %v.", db.DescribeTableInvocations)
		}

		if db.CreateTableInvocations != 1 {
			t.Errorf("Expected 1 calls to CreateTable, got %v.", db.CreateTableInvocations)
		}
	})
}

// Test_AbortProvider checks:
// - if correct response is returned when provider is aborted
// - if correct response is returned when version is aborted
// - if error is returned when DeleteItem fails for provider
// - if error is returned when DeleteItem fails for provider version
func Test_AbortProvider(t *testing.T) {
	t.Parallel()

	t.Run("when provider version is aborted", func(t *testing.T) {
		db := &mocks.DynamoDB{}

		svc := &VersionManagerService{Db: db}

		req := &services.TerminateVersionRequest{Provider: &terrarium.Provider{Name: "test-org/test-provider2", Version: "2.0.0"}}

		res, err := svc.AbortProviderVersion(context.TODO(), req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if db.DeleteItemInvocations != 1 {
			t.Errorf("Expected 1 call to DeleteItem, got %v", db.DeleteItemInvocations)
		}

		if db.TableName != VersionsTableName {
			t.Errorf("Expected tableName to be %v, got %v.", VersionsTableName, db.TableName)
		}

		if res != VersionAborted {
			t.Errorf("Expected %v, got %v.", VersionAborted, res)
		}
	})

	t.Run("when DeleteItem fails for provider version", func(t *testing.T) {
		db := &mocks.DynamoDB{DeleteItemError: errors.New("some error")}

		svc := &VersionManagerService{Db: db}

		req := services.TerminateVersionRequest{Provider: &terrarium.Provider{Name: "test-org/test-provider2", Version: "2.0.0"}}

		res, err := svc.AbortProviderVersion(context.TODO(), &req)

		if res != nil {
			t.Errorf("Expected no response, got %v", res)
		}

		if db.DeleteItemInvocations != 1 {
			t.Errorf("Expected 1 call to DeleteItem, got %v", db.DeleteItemInvocations)
		}

		if db.TableName != VersionsTableName {
			t.Errorf("Expected tableName to be %v, got %v.", VersionsTableName, db.TableName)
		}

		if err != AbortProviderVersionError {
			t.Errorf("Expected %v, got %v.", AbortProviderVersionError, err)
		}
	})
}

// Test_PublishVersion checks:
// - if correct response is returned when version is published
// - if error is returned when UpdateItem fails
func Test_PublishVersion(t *testing.T) {
	t.Parallel()

	t.Run("when version is published", func(t *testing.T) {
		db := &mocks.DynamoDB{}

		svc := &VersionManagerService{Db: db}

		req := &services.TerminateVersionRequest{Provider: &terrarium.Provider{Name: "test-org/test-provider2", Version: "2.0.0"}}

		res, err := svc.PublishVersion(context.TODO(), req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if db.UpdateItemInvocations != 1 {
			t.Errorf("Expected 1 call to UpdateItem, got %v", db.UpdateItemInvocations)
		}

		if db.TableName != VersionsTableName {
			t.Errorf("Expected tableName to be %v, got %v.", VersionsTableName, db.TableName)
		}

		if res != VersionPublished {
			t.Errorf("Expected %v, got %v.", VersionPublished, res)
		}
	})

	t.Run("when UpdateItem fails", func(t *testing.T) {
		db := &mocks.DynamoDB{UpdateItemError: errors.New("some error")}

		svc := &VersionManagerService{Db: db}

		req := &services.TerminateVersionRequest{Provider: &terrarium.Provider{Name: "test-org/test-provider2", Version: "2.0.0"}}

		res, err := svc.PublishVersion(context.TODO(), req)

		if res != nil {
			t.Errorf("Expected no response, got %v", res)
		}

		if db.UpdateItemInvocations != 1 {
			t.Errorf("Expected 1 call to UpdateItem, got %v", db.UpdateItemInvocations)
		}

		if db.TableName != VersionsTableName {
			t.Errorf("Expected tableName to be %v, got %v.", VersionsTableName, db.TableName)
		}

		if err != PublishProviderVersionError {
			t.Errorf("Expected %v, got %v.", PublishProviderVersionError, err)
		}
	})
}

func Test_GetProvider(t *testing.T) {
	t.Parallel()

	t.Run("Get existing provider", func(t *testing.T) {
		db := &mocks.DynamoDB{
			ScanOut: &dynamodb.ScanOutput{
				Items: []map[string]types.AttributeValue{
					{
						"name":            &types.AttributeValueMemberS{Value: "test-org/test-provider"},
						"description":     &types.AttributeValueMemberS{Value: "Test Description"},
						"source_repo_url": &types.AttributeValueMemberS{Value: "http://test.com/provider"},
						"maturity":        &types.AttributeValueMemberS{Value: "ALPHA"},
					},
				},
			},
		}

		svc := &VersionManagerService{Db: db}

		req := services.ProviderName{Provider: "test-provider"}
		res, err := svc.GetProvider(context.TODO(), &req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		expectedResponse := &services.GetProviderResponse{
			Provider: &services.ListProviderItem{
				Organization:  "test-org",
				Name:          "test-provider",
				Description:   "Test Description",
				SourceRepoUrl: "http://test.com/provider",
				Maturity:      terrarium.Maturity_ALPHA,
			},
		}
		if !reflect.DeepEqual(res, expectedResponse) {
			t.Errorf("Response does not match, got %v, want %v", res, expectedResponse)
		}
	})

	t.Run("Get non-existing provider", func(t *testing.T) {
		db := &mocks.DynamoDB{
			ScanOut: &dynamodb.ScanOutput{},
		}

		svc := &VersionManagerService{Db: db}

		req := services.ProviderName{Provider: "non-existing-provider"}
		res, err := svc.GetProvider(context.TODO(), &req)

		if res != nil {
			t.Errorf("Expected nil response for non-existing provider, got %v", res)
		}

		expectedErr := fmt.Errorf("provider not found 'non-existing-provider'")
		if err.Error() != expectedErr.Error() {
			t.Errorf("Expected error '%v', got '%v'", expectedErr, err)
		}
	})
}

func Test_ListProviders(t *testing.T) {
	t.Parallel()

	t.Run("List providers", func(t *testing.T) {
		db := &mocks.DynamoDB{
			ScanOut: &dynamodb.ScanOutput{
				Items: []map[string]types.AttributeValue{
					{
						"name":            &types.AttributeValueMemberS{Value: "test-org/test-provider"},
						"description":     &types.AttributeValueMemberS{Value: "Test Description"},
						"source_repo_url": &types.AttributeValueMemberS{Value: "http://test.com/provider"},
						"maturity":        &types.AttributeValueMemberS{Value: "ALPHA"},
					},
					{
						"name":            &types.AttributeValueMemberS{Value: "test-org2/test-provider2"},
						"description":     &types.AttributeValueMemberS{Value: "Test Description2"},
						"source_repo_url": &types.AttributeValueMemberS{Value: "http://test.com/provider2"},
						"maturity":        &types.AttributeValueMemberS{Value: "BETA"},
					},
				},
			},
		}

		svc := &VersionManagerService{Db: db}

		req := services.ListProvidersRequest{}
		res, err := svc.ListProviders(context.TODO(), &req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		expectedResponse := &services.ListProvidersResponse{
			Providers: []*services.ListProviderItem{
				{
					Organization:  "test-org",
					Name:          "test-provider",
					Description:   "Test Description",
					SourceRepoUrl: "http://test.com/provider",
					Maturity:      terrarium.Maturity_ALPHA,
				},
				{
					Organization:  "test-org2",
					Name:          "test-provider2",
					Description:   "Test Description2",
					SourceRepoUrl: "http://test.com/provider2",
					Maturity:      terrarium.Maturity_BETA,
				},
			},
		}
		if !reflect.DeepEqual(res, expectedResponse) {
			t.Errorf("Response does not match, got %v, want %v", res, expectedResponse)
		}
	})
}

