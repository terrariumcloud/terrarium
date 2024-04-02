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
			Name:                "test",
			Version:             "1.0.0",
			Protocols:           []string{"5.1", "5.2"},
			Os:                  "linux",
			Arch:                "amd64",
			Filename:            "test.tar.gz",
			DownloadUrl:         "https://example.com/test.tar.gz",
			Shasum:              "1234567890abcdef",
			ShasumsUrl:          "https://example.com/test.tar.gz.sha256",
			ShasumsSignatureUrl: "https://example.com/test.tar.gz.sha256.asc",
			KeyId:               "test-key",
			AsciiArmor:          "test-armor",
			TrustSignature:      "test-trust",
			Source:              "test-source",
			SourceUrl:           "https://example.com/test-source",
			Description:         "test provider",
			SourceRepoUrl:       "https://example.com/test-repo",
			Maturity:            terrarium.Maturity_ALPHA,
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
						"created_on":            services.MustMarshallString(emptyString, t),
						"modified_on":           services.MustMarshallString(emptyString, t),
						"os":                    services.MustMarshallString(emptyString, t),
						"arch":                  services.MustMarshallString(emptyString, t),
						"filename":              services.MustMarshallString(emptyString, t),
						"download_url":          services.MustMarshallString(emptyString, t),
						"shasums_url":           services.MustMarshallString(emptyString, t),
						"shasums_signature_url": services.MustMarshallString(emptyString, t),
						"shasum":                services.MustMarshallString(emptyString, t),
						"key_id":                services.MustMarshallString(emptyString, t),
						"ascii_armor":           services.MustMarshallString(emptyString, t),
						"trust_signature":       services.MustMarshallString(emptyString, t),
						"source":                services.MustMarshallString(emptyString, t),
						"source_repo_url":       services.MustMarshallString(emptyString, t),
					},
				},
			},
			UpdateItemOut: &dynamodb.UpdateItemOutput{},
		}

		svc := &VersionManagerService{Db: db}

		req := &terrarium.RegisterProviderRequest{
			Name:                "test",
			Version:             "1.0.0",
			Protocols:           []string{"5.1", "5.2"},
			Os:                  "linux",
			Arch:                "amd64",
			Filename:            "test.tar.gz",
			DownloadUrl:         "https://example.com/test.tar.gz",
			Shasum:              "1234567890abcdef",
			ShasumsUrl:          "https://example.com/test.tar.gz.sha256",
			ShasumsSignatureUrl: "https://example.com/test.tar.gz.sha256.asc",
			KeyId:               "test-key",
			AsciiArmor:          "test-armor",
			TrustSignature:      "test-trust",
			Source:              "test-source",
			SourceUrl:           "https://example.com/test-source",
			Description:         "test provider",
			SourceRepoUrl:       "https://example.com/test-repo",
			Maturity:            terrarium.Maturity_ALPHA,
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
			Name:                "test",
			Version:             "1.0.0",
			Protocols:           []string{"5.1", "5.2"},
			Os:                  "linux",
			Arch:                "amd64",
			Filename:            "test.tar.gz",
			DownloadUrl:         "https://example.com/test.tar.gz",
			Shasum:              "1234567890abcdef",
			ShasumsUrl:          "https://example.com/test.tar.gz.sha256",
			ShasumsSignatureUrl: "https://example.com/test.tar.gz.sha256.asc",
			KeyId:               "test-key",
			AsciiArmor:          "test-armor",
			TrustSignature:      "test-trust",
			Source:              "test-source",
			SourceUrl:           "https://example.com/test-source",
			Description:         "test provider",
			SourceRepoUrl:       "https://example.com/test-repo",
			Maturity:            terrarium.Maturity_ALPHA,
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
			Name:                "test",
			Version:             "1.0.0",
			Protocols:           []string{"5.1", "5.2"},
			Os:                  "linux",
			Arch:                "amd64",
			Filename:            "test.tar.gz",
			DownloadUrl:         "https://example.com/test.tar.gz",
			Shasum:              "1234567890abcdef",
			ShasumsUrl:          "https://example.com/test.tar.gz.sha256",
			ShasumsSignatureUrl: "https://example.com/test.tar.gz.sha256.asc",
			KeyId:               "test-key",
			AsciiArmor:          "test-armor",
			TrustSignature:      "test-trust",
			Source:              "test-source",
			SourceUrl:           "https://example.com/test-source",
			Description:         "test provider",
			SourceRepoUrl:       "https://example.com/test-repo",
			Maturity:            terrarium.Maturity_ALPHA,
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

	t.Run("when provider is aborted", func(t *testing.T) {
		db := &mocks.DynamoDB{}

		svc := &VersionManagerService{Db: db}

		req := &services.TerminateProviderRequest{Provider: &terrarium.ProviderName{Name: "test"}}

		res, err := svc.AbortProvider(context.TODO(), req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if db.DeleteItemInvocations != 1 {
			t.Errorf("Expected 1 call to DeleteItem, got %v", db.DeleteItemInvocations)
		}

		if db.TableName != VersionsTableName {
			t.Errorf("Expected tableName to be %v, got %v.", VersionsTableName, db.TableName)
		}

		if res != ProviderAborted {
			t.Errorf("Expected %v, got %v.", ProviderAborted, res)
		}
	})

	t.Run("when provider version is aborted", func(t *testing.T) {
		db := &mocks.DynamoDB{}

		svc := &VersionManagerService{Db: db}

		req := &services.TerminateVersionRequest{Provider: &terrarium.ProviderVersion{Name: "test", Version: "1.0.0"}}

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

	t.Run("when DeleteItem fails for provider", func(t *testing.T) {
		db := &mocks.DynamoDB{DeleteItemError: errors.New("some error")}

		svc := &VersionManagerService{Db: db}

		req := services.TerminateProviderRequest{Provider: &terrarium.ProviderName{Name: "test"}}

		res, err := svc.AbortProvider(context.TODO(), &req)

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

	t.Run("when DeleteItem fails for provider version", func(t *testing.T) {
		db := &mocks.DynamoDB{DeleteItemError: errors.New("some error")}

		svc := &VersionManagerService{Db: db}

		req := services.TerminateVersionRequest{Provider: &terrarium.ProviderVersion{Name: "test", Version: "v1.0.0"}}

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

		req := &services.PublishVersionRequest{Provider: &terrarium.Provider{Name: "test", Version: "v1.0.0", Os: "linux", Arch: "amd64"}}

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

		req := &services.PublishVersionRequest{Provider: &terrarium.Provider{Name: "test", Version: "v1.0.0", Os: "linux", Arch: "amd64"}}

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

func Test_GetVersionData(t *testing.T) {
	t.Parallel()

	db := &mocks.DynamoDB{
		GetItemOuts: []*dynamodb.GetItemOutput{
			{
				Item: map[string]types.AttributeValue{
					"name":                  &types.AttributeValueMemberS{Value: "test-org/provider"},
					"version":               &types.AttributeValueMemberS{Value: "1.0.1"},
					"protocols":             &types.AttributeValueMemberSS{Value: []string{"5.1", "5.2"}},
					"os":                    &types.AttributeValueMemberS{Value: "linux"},
					"arch":                  &types.AttributeValueMemberS{Value: "amd64"},
					"filename":              &types.AttributeValueMemberS{Value: "test-provider"},
					"download_url":          &types.AttributeValueMemberS{Value: "http://test.com/download"},
					"shasums_url":           &types.AttributeValueMemberS{Value: "http://test.com/shasums"},
					"shasums_signature_url": &types.AttributeValueMemberS{Value: "http://test.com/shasums/signature"},
					"shasum":                &types.AttributeValueMemberS{Value: "1234567890"},
					"key_id":                &types.AttributeValueMemberS{Value: "ABCD1234"},
					"ascii_armor":           &types.AttributeValueMemberS{Value: "-----BEGIN PGP PUBLIC KEY BLOCK-----\n..."},
					"trust_signature":       &types.AttributeValueMemberS{Value: ""},
					"source":                &types.AttributeValueMemberS{Value: "example@example.com"},
					"source_url":            &types.AttributeValueMemberS{Value: "http://example.com/key.asc"},
				},
			},
		},
	}

	// Create an instance of VersionManagerService with the mocked DynamoDB client
	svc := &VersionManagerService{Db: db}

	// Create a request object
	req := services.VersionDataRequest{
		Name:    "test-org/provider",
		Version: "1.0.1",
		Os:      "linux",
		Arch:    "amd64",
	}

	// Call the GetVersionData method with the request object
	res, err := svc.GetVersionData(context.TODO(), &req)

	// Check if there's an error
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check if the returned response matches the expected response
	expectedResponse := &services.PlatformMetadataResponse{
		Protocols:           []string{"5.1", "5.2"},
		Os:                  "linux",
		Arch:                "amd64",
		Filename:            "test-provider",
		DownloadUrl:         "http://test.com/download",
		ShasumsUrl:          "http://test.com/shasums",
		ShasumsSignatureUrl: "http://test.com/shasums/signature",
		Shasum:              "1234567890",
		SigningKeys: &services.SigningKeys{
			GpgPublicKeys: []*services.GPGPublicKey{
				{
					KeyId:          "ABCD1234",
					AsciiArmor:     "-----BEGIN PGP PUBLIC KEY BLOCK-----\n...",
					TrustSignature: "",
					Source:         "example@example.com",
					SourceUrl:      "http://example.com/key.asc",
				},
			},
		},
	}

	if !reflect.DeepEqual(res, expectedResponse) {
		t.Errorf("Response does not match, got %v, want %v", res, expectedResponse)
	}
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

// Test_ListModuleVersions checks:
// - if correct response is returned when versions are fetched
func Test_ListProviderVersions(t *testing.T) {
	t.Parallel()

	t.Run("Listing versions", func(t *testing.T) {
		db := &mocks.DynamoDB{
			ScanOut: &dynamodb.ScanOutput{
				Items: []map[string]types.AttributeValue{
					{
						"name":         &types.AttributeValueMemberS{Value: "test-org/test-provider"},
						"version":      &types.AttributeValueMemberS{Value: "1.0.0"},
						"protocols":    &types.AttributeValueMemberSS{Value: []string{"5.1", "5.2"}},
						"os":           &types.AttributeValueMemberS{Value: "linux"},
						"arch":         &types.AttributeValueMemberS{Value: "amd64"},
						"published_on": &types.AttributeValueMemberS{Value: "exists"},
					},
					{
						"name":         &types.AttributeValueMemberS{Value: "test-org2/test-provider2"},
						"version":      &types.AttributeValueMemberS{Value: "2.0.0"},
						"protocols":    &types.AttributeValueMemberSS{Value: []string{"4.1"}},
						"os":           &types.AttributeValueMemberS{Value: "windows"},
						"arch":         &types.AttributeValueMemberS{Value: "amd64"},
						"published_on": &types.AttributeValueMemberS{Value: "exists"},
					},
				},
			},
		}

		svc := &VersionManagerService{Db: db}

		req := &services.ProviderName{Provider: "test-org/test-provider"}
		res, err := svc.ListProviderVersions(context.TODO(), req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		expectedVersions := []*services.VersionItem{
			{
				Version:   "1.0.0",
				Protocols: []string{"5.1", "5.2"},
				Platforms: []*services.Platform{
					{
						Os:   "linux",
						Arch: "amd64",
					},
				},
			},
			{
				Version:   "2.0.0",
				Protocols: []string{"4.1"},
				Platforms: []*services.Platform{
					{
						Os:   "windows",
						Arch: "amd64",
					},
				},
			},
		}
		if !reflect.DeepEqual(res.Versions, expectedVersions) {
			t.Errorf("Versions do not match, got %v, want %v", res.Versions, expectedVersions)
		}
	})
}
