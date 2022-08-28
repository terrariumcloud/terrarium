package services_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/mocks"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/module/services"
	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"
	"google.golang.org/grpc"
)

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() error {
	return nil
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

	s3 := &mocks.MockS3{
		HeadBucketError:   errors.New("some error"),
		CreateBucketError: errors.New("some error"),
	}

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

	mus := &mocks.MockUploadSourceZipServer{RecvRequest: req, RecvMaxInvocations: 2}

	err := svc.UploadSourceZip(mus)

	if err != nil {
		t.Errorf("Expected no error, got %v.", err)
	}

	if mus.RecvInvocations != 1 {
		t.Errorf("Expected 1 call to Recv, got %v", mus.RecvInvocations)
	}

	if s3.PutObjectInvocations != 1 {
		t.Errorf("Expected 1 call to PutObject, got %v", s3.PutObjectInvocations)
	}

	if mus.SendAndCloseInvocations != 1 {
		t.Errorf("Expected 1 call to SendAndClose, got %v", mus.SendAndCloseInvocations)
	}

	if mus.SendAndCloseResponse != services.SourceZipUploaded {
		t.Errorf("Expected %v, got %v.", services.SourceZipUploaded, mus.SendAndCloseResponse)
	}
}

// This test checks if error is returned when PutObject fails
func TestUploadSourceZipWhenPutObjectErrors(t *testing.T) {
	t.Parallel()

	s3 := &mocks.MockS3{PutObjectError: errors.New("some error")}

	svc := &services.StorageService{S3: s3}

	req := &terrarium.UploadSourceZipRequest{
		Module:       &terrarium.Module{Name: "test", Version: "v1"},
		ZipDataChunk: make([]byte, 1000),
	}

	mus := &mocks.MockUploadSourceZipServer{RecvRequest: req, RecvMaxInvocations: 1}

	err := svc.UploadSourceZip(mus)

	if mus.RecvInvocations != 1 {
		t.Errorf("Expected 1 call to Recv, got %v", mus.RecvInvocations)
	}

	if s3.PutObjectInvocations != 1 {
		t.Errorf("Expected 1 call to PutObject, got %v", s3.PutObjectInvocations)
	}

	if mus.SendAndCloseInvocations != 0 {
		t.Errorf("Expected 0 calls to SendAndClose, got %v", mus.SendAndCloseInvocations)
	}

	if err != services.UploadSourceZipError {
		t.Errorf("Expected %v, got %v.", services.UploadSourceZipError, err)
	}
}

// This test checks if error is returned when Recv fails
func TestUploadSourceZipWhenRecvErrors(t *testing.T) {
	t.Parallel()

	s3 := &mocks.MockS3{}

	svc := &services.StorageService{S3: s3}

	mus := &mocks.MockUploadSourceZipServer{
		RecvError:          errors.New("some error"),
		RecvMaxInvocations: 1,
	}

	err := svc.UploadSourceZip(mus)

	if mus.RecvInvocations != 1 {
		t.Errorf("Expected 1 call to Recv, got %v", mus.RecvInvocations)
	}

	if s3.PutObjectInvocations != 0 {
		t.Errorf("Expected 0 calls to PutObject, got %v", s3.PutObjectInvocations)
	}

	if mus.SendAndCloseInvocations != 0 {
		t.Errorf("Expected 0 calls to SendAndClose, got %v", mus.SendAndCloseInvocations)
	}

	if err != services.RecieveSourceZipError {
		t.Errorf("Expected %v, got %v.", services.RecieveSourceZipError, err)
	}
}

// This test checks if correct response is returned when source zip is downloaded
func TestDownloadSourceZip(t *testing.T) {
	t.Parallel()

	buf := &ClosingBuffer{bytes.NewBuffer(make([]byte, 1000))}

	len := int64(1000)

	s3 := &mocks.MockS3{GetObjectOut: &s3.GetObjectOutput{Body: buf, ContentLength: &len}}

	svc := &services.StorageService{S3: s3}

	res := &terrarium.SourceZipResponse{ZipDataChunk: make([]byte, 1000)}

	mds := &mocks.MockDownloadSourceZipServer{SendResponse: res}

	req := &terrarium.DownloadSourceZipRequest{
		Module: &terrarium.Module{Name: "Test", Version: "v1"},
	}

	err := svc.DownloadSourceZip(req, mds)

	if err != nil {
		t.Errorf("Expected no error, got %v.", err)
	}

	if s3.GetObjectInvocations != 1 {
		t.Errorf("Expected 1 call to GetObject, got %v", s3.GetObjectInvocations)
	}

	if mds.SendInvocations != 1 {
		t.Errorf("Expected 1 call to Send, got %v", mds.SendInvocations)
	}

	if !bytes.Equal(mds.SendResponse.ZipDataChunk, res.ZipDataChunk) {
		t.Errorf("Expected same data to be returned.")
	}
}

// This test checks if error is returned when GetObject fails
func TestDownloadSourceZipWhenGetObjectErrors(t *testing.T) {
	t.Parallel()

	s3 := &mocks.MockS3{GetObjectError: errors.New("some error")}

	svc := &services.StorageService{S3: s3}

	mds := &mocks.MockDownloadSourceZipServer{}

	req := &terrarium.DownloadSourceZipRequest{
		Module: &terrarium.Module{Name: "Test", Version: "v1"},
	}

	err := svc.DownloadSourceZip(req, mds)

	if s3.GetObjectInvocations != 1 {
		t.Errorf("Expected 1 call to GetObject, got %v", s3.GetObjectInvocations)
	}

	if mds.SendInvocations != 0 {
		t.Errorf("Expected 0 call to Sends, got %v", mds.SendInvocations)
	}

	if err != services.DownloadSourceZipError {
		t.Errorf("Expected %v, got %v.", services.DownloadSourceZipError, err)
	}
}

// This test checks error is returned when Send fails
func TestDownloadSourceZipWhenSendErrors(t *testing.T) {
	t.Parallel()

	buf := &ClosingBuffer{bytes.NewBuffer(make([]byte, 1000))}

	len := int64(1000)

	s3 := &mocks.MockS3{GetObjectOut: &s3.GetObjectOutput{Body: buf, ContentLength: &len}}

	svc := &services.StorageService{S3: s3}

	mds := &mocks.MockDownloadSourceZipServer{SendError: errors.New("some error")}

	req := &terrarium.DownloadSourceZipRequest{
		Module: &terrarium.Module{Name: "Test", Version: "v1"},
	}

	err := svc.DownloadSourceZip(req, mds)

	if s3.GetObjectInvocations != 1 {
		t.Errorf("Expected 1 call to GetObject, got %v", s3.GetObjectInvocations)
	}

	if mds.SendInvocations != 1 {
		t.Errorf("Expected 1 call to Send, got %v", mds.SendInvocations)
	}

	if err != services.SendSourceZipError {
		t.Errorf("Expected %v, got %v.", services.SendSourceZipError, err)
	}
}

// This test checks if error is returned when wrong content lenght is read
func TestDownloadSourceZipWhenWrongContentLenght(t *testing.T) {
	t.Parallel()

	buf := &ClosingBuffer{bytes.NewBuffer(make([]byte, 1000))}

	len := int64(10000)

	s3 := &mocks.MockS3{GetObjectOut: &s3.GetObjectOutput{Body: buf, ContentLength: &len}}

	svc := &services.StorageService{S3: s3}

	res := &terrarium.SourceZipResponse{ZipDataChunk: make([]byte, 1000)}

	mds := &mocks.MockDownloadSourceZipServer{SendResponse: res}

	req := &terrarium.DownloadSourceZipRequest{
		Module: &terrarium.Module{Name: "Test", Version: "v1"},
	}

	err := svc.DownloadSourceZip(req, mds)

	if s3.GetObjectInvocations != 1 {
		t.Errorf("Expected 1 call to GetObject, got %v", s3.GetObjectInvocations)
	}

	if mds.SendInvocations != 0 {
		t.Errorf("Expected 0 calls to Send, got %v", mds.SendInvocations)
	}

	if err != services.ContentLenghtError {
		t.Errorf("Expected %v, got %v.", services.ContentLenghtError, err)
	}
}
