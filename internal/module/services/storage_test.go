package services_test

import (
	"errors"
	"io"
	"testing"

	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/mocks"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/module/services"
	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"
	"google.golang.org/grpc"
)

type MockUploadServer struct {
	services.Storage_UploadSourceZipServer
	SendAndCloseInvocations int
	SendAndCloseResponse    *terrarium.Response
	SendAndCloseError       error
	RecvInvocations         int
	RecvMaxInvocations      int
	RecvRequest             *terrarium.UploadSourceZipRequest
	RecvError               error
}

func (mus *MockUploadServer) SendAndClose(response *terrarium.Response) error {
	mus.SendAndCloseInvocations++
	mus.SendAndCloseResponse = response
	return mus.SendAndCloseError
}

func (mus *MockUploadServer) Recv() (*terrarium.UploadSourceZipRequest, error) {
	mus.RecvInvocations++

	if mus.RecvMaxInvocations >= mus.RecvInvocations {
		return nil, io.EOF
	}

	return mus.RecvRequest, mus.RecvError
}

// This test checks if there was no error
func TestRegisterStorageWithServer(t *testing.T) {
	t.Parallel()

	s3 := &mocks.MockS3{}

	ss := &services.StorageService{S3: s3}

	s := grpc.NewServer(*new([]grpc.ServerOption)...)

	err := ss.RegisterWithServer(s)

	if err != nil {
		t.Errorf("Expected no error, got %v.", err)
	}

	if s3.HeadBucketInvocations != 1 {
		t.Errorf("Expected 1 call to HeadBucket, got %v.", s3.HeadBucketInvocations)
	}

	if s3.CreateBucketInvocations != 0 {
		t.Errorf("Expected no calls to CreateBucket, got %v.", s3.CreateBucketInvocations)
	}
}

// This test checks if error is returned when Bucket initialization fails
func TestRegisterWithServerWhenStorageBucketInitializationErrors(t *testing.T) {
	t.Parallel()

	s3 := &mocks.MockS3{HeadBucketError: errors.New("some error"), CreateBucketError: errors.New("some error")}

	vms := &services.StorageService{S3: s3}

	s := grpc.NewServer(*new([]grpc.ServerOption)...)

	err := vms.RegisterWithServer(s)

	if err != services.BucketInitializationError {
		t.Errorf("Expected %v, got %v.", services.BucketInitializationError, err)
	}

	if s3.HeadBucketInvocations != 1 {
		t.Errorf("Expected 1 call to DescribeTable, got %v.", s3.HeadBucketInvocations)
	}

	if s3.CreateBucketInvocations != 1 {
		t.Errorf("Expected 0 calls to CreateTable, got %v.", s3.CreateBucketInvocations)
	}
}

// This test checks if correct response is returned when source zip is uploaded
func TestUploadSourceZip(t *testing.T) {
	t.Parallel()

	s3 := &mocks.MockS3{}

	svc := &services.StorageService{S3: s3}

	req := &terrarium.UploadSourceZipRequest{
		Module:       &terrarium.Module{Name: "test", Version: "v1"},
		ZipDataChunk: make([]byte, 1000),
	}

	mus := &MockUploadServer{RecvRequest: req, RecvMaxInvocations: 1}

	err := svc.UploadSourceZip(mus)

	if err != nil {
		t.Errorf("Expected no error, got %v.", err)
	}

	if mus.RecvInvocations != 1 {
		t.Errorf("Expected 1 call to Recv, got %v", mus.RecvInvocations)
	}

	if s3.PutItemInvocations != 1 {
		t.Errorf("Expected 1 call to PutItem, got %v", s3.PutItemInvocations)
	}

	if mus.SendAndCloseInvocations != 1 {
		t.Errorf("Expected 1 call to SendAndClose, got %v", mus.SendAndCloseInvocations)
	}

	if mus.SendAndCloseResponse != services.SourceZipUploaded {
		t.Errorf("Expected %v, got %v.", services.SourceZipUploaded, mus.SendAndCloseResponse)
	}
}
