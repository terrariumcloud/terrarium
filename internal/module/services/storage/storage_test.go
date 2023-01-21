package storage

import (
	"bytes"
	"errors"
	"github.com/terrariumcloud/terrarium/internal/module/services/mocks"
	mocks2 "github.com/terrariumcloud/terrarium/internal/storage/mocks"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	terrarium "github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"google.golang.org/grpc"
)

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() error {
	return nil
}

// Test_RegisterStorageWithServer checks:
// - if there was no error with bucket init
// - if error was returned when bucket init fails
func Test_RegisterStorageWithServer(t *testing.T) {
	t.Parallel()

	t.Run("when bucket init is successful", func(t *testing.T) {
		s3Client := &mocks2.S3{}

		ss := &StorageService{Client: s3Client}

		s := grpc.NewServer(*new([]grpc.ServerOption)...)

		err := ss.RegisterWithServer(s)

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}

		if s3Client.HeadBucketInvocations != 1 {
			t.Errorf("Expected 1 call to HeadBucket, got %v.", s3Client.HeadBucketInvocations)
		}

		if s3Client.CreateBucketInvocations != 0 {
			t.Errorf("Expected no calls to CreateBucket, got %v.", s3Client.CreateBucketInvocations)
		}
	})

	t.Run("when bucket init fails", func(t *testing.T) {
		s3Client := &mocks2.S3{
			HeadBucketError:   errors.New("some error"),
			CreateBucketError: errors.New("some error"),
		}

		vms := &StorageService{Client: s3Client}

		s := grpc.NewServer(*new([]grpc.ServerOption)...)

		err := vms.RegisterWithServer(s)

		if err != BucketInitializationError {
			t.Errorf("Expected %v, got %v.", BucketInitializationError, err)
		}

		if s3Client.HeadBucketInvocations != 1 {
			t.Errorf("Expected 1 call to DescribeTable, got %v.", s3Client.HeadBucketInvocations)
		}

		if s3Client.CreateBucketInvocations != 1 {
			t.Errorf("Expected 0 calls to CreateTable, got %v.", s3Client.CreateBucketInvocations)
		}
	})
}

// Test_UploadSourceZip checks:
// - if correct response is returned when source zip is uploaded
// - if error is returned when PutObject fails
// - if error is returned when Recv fails
func Test_UploadSourceZip(t *testing.T) {
	t.Parallel()

	t.Run("when source zip is uploaded", func(t *testing.T) {
		s3Client := &mocks2.S3{}

		svc := &StorageService{Client: s3Client}

		req := &terrarium.UploadSourceZipRequest{
			Module:       &terrarium.Module{Name: "test", Version: "v1"},
			ZipDataChunk: make([]byte, 1000),
		}

		mus := &mocks.MockUploadSourceZipServer{RecvRequest: req, RecvMaxInvocations: 2}

		err := svc.UploadSourceZip(mus)

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}

		if mus.RecvInvocations != 2 {
			t.Errorf("Expected 1 call to Recv, got %v", mus.RecvInvocations)
		}

		if s3Client.PutObjectInvocations != 1 {
			t.Errorf("Expected 1 call to PutObject, got %v", s3Client.PutObjectInvocations)
		}

		if mus.SendAndCloseInvocations != 1 {
			t.Errorf("Expected 1 call to SendAndClose, got %v", mus.SendAndCloseInvocations)
		}

		if mus.SendAndCloseResponse != SourceZipUploaded {
			t.Errorf("Expected %v, got %v.", SourceZipUploaded, mus.SendAndCloseResponse)
		}
	})

	t.Run("when PutObject fails", func(t *testing.T) {
		s3Client := &mocks2.S3{PutObjectError: errors.New("some error")}

		svc := &StorageService{Client: s3Client}

		req := &terrarium.UploadSourceZipRequest{
			Module:       &terrarium.Module{Name: "test", Version: "v1"},
			ZipDataChunk: make([]byte, 1000),
		}

		mus := &mocks.MockUploadSourceZipServer{RecvRequest: req, RecvMaxInvocations: 1}

		err := svc.UploadSourceZip(mus)

		if mus.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", mus.RecvInvocations)
		}

		if s3Client.PutObjectInvocations != 1 {
			t.Errorf("Expected 1 call to PutObject, got %v", s3Client.PutObjectInvocations)
		}

		if mus.SendAndCloseInvocations != 0 {
			t.Errorf("Expected 0 calls to SendAndClose, got %v", mus.SendAndCloseInvocations)
		}

		if err != UploadSourceZipError {
			t.Errorf("Expected %v, got %v.", UploadSourceZipError, err)
		}
	})

	t.Run("when Recv fails", func(t *testing.T) {
		s3Client := &mocks2.S3{}

		svc := &StorageService{Client: s3Client}

		mus := &mocks.MockUploadSourceZipServer{
			RecvError:          errors.New("some error"),
			RecvMaxInvocations: 1,
		}

		err := svc.UploadSourceZip(mus)

		if mus.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", mus.RecvInvocations)
		}

		if s3Client.PutObjectInvocations != 0 {
			t.Errorf("Expected 0 calls to PutObject, got %v", s3Client.PutObjectInvocations)
		}

		if mus.SendAndCloseInvocations != 0 {
			t.Errorf("Expected 0 calls to SendAndClose, got %v", mus.SendAndCloseInvocations)
		}

		if err != RecieveSourceZipError {
			t.Errorf("Expected %v, got %v.", RecieveSourceZipError, err)
		}
	})
}

// Test_DownloadSourceZip checks:
// - if correct response is returned when source zip is downloaded
// - if error is returned when GetObject fails
// - if error is returned when Send fails
// - if error is returned when wrong content lenght is read
func Test_DownloadSourceZip(t *testing.T) {
	t.Parallel()

	t.Run("when source zip is downloaded", func(t *testing.T) {
		var length int64 = 1000
		buf := &ClosingBuffer{bytes.NewBuffer(make([]byte, length))}

		s3Client := &mocks2.S3{GetObjectOut: &s3.GetObjectOutput{Body: buf, ContentLength: length}}

		svc := &StorageService{Client: s3Client}

		res := &terrarium.SourceZipResponse{ZipDataChunk: make([]byte, length)}

		mds := &mocks.MockDownloadSourceZipServer{SendResponse: res}

		req := &terrarium.DownloadSourceZipRequest{
			Module: &terrarium.Module{Name: "Test", Version: "v1"},
		}

		err := svc.DownloadSourceZip(req, mds)

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}

		if s3Client.GetObjectInvocations != 1 {
			t.Errorf("Expected 1 call to GetObject, got %v", s3Client.GetObjectInvocations)
		}

		if mds.SendInvocations != 1 {
			t.Errorf("Expected 1 call to Send, got %v", mds.SendInvocations)
		}

		if !bytes.Equal(mds.SendResponse.ZipDataChunk, res.ZipDataChunk) {
			t.Errorf("Expected same data to be returned.")
		}
	})

	t.Run("when GetObject fails", func(t *testing.T) {
		s3Client := &mocks2.S3{GetObjectError: errors.New("some error")}

		svc := &StorageService{Client: s3Client}

		mds := &mocks.MockDownloadSourceZipServer{}

		req := &terrarium.DownloadSourceZipRequest{
			Module: &terrarium.Module{Name: "Test", Version: "v1"},
		}

		err := svc.DownloadSourceZip(req, mds)

		if s3Client.GetObjectInvocations != 1 {
			t.Errorf("Expected 1 call to GetObject, got %v", s3Client.GetObjectInvocations)
		}

		if mds.SendInvocations != 0 {
			t.Errorf("Expected 0 call to Sends, got %v", mds.SendInvocations)
		}

		if err != DownloadSourceZipError {
			t.Errorf("Expected %v, got %v.", DownloadSourceZipError, err)
		}
	})

	t.Run("when Send fails", func(t *testing.T) {
		var length int64 = 1000
		buf := &ClosingBuffer{bytes.NewBuffer(make([]byte, length))}

		s3Client := &mocks2.S3{GetObjectOut: &s3.GetObjectOutput{Body: buf, ContentLength: length}}

		svc := &StorageService{Client: s3Client}

		mds := &mocks.MockDownloadSourceZipServer{SendError: errors.New("some error")}

		req := &terrarium.DownloadSourceZipRequest{
			Module: &terrarium.Module{Name: "Test", Version: "v1"},
		}

		err := svc.DownloadSourceZip(req, mds)

		if s3Client.GetObjectInvocations != 1 {
			t.Errorf("Expected 1 call to GetObject, got %v", s3Client.GetObjectInvocations)
		}

		if mds.SendInvocations != 1 {
			t.Errorf("Expected 1 call to Send, got %v", mds.SendInvocations)
		}

		if err != SendSourceZipError {
			t.Errorf("Expected %v, got %v.", SendSourceZipError, err)
		}
	})

	t.Run("when wrong content length is read", func(t *testing.T) {
		buf := &ClosingBuffer{bytes.NewBuffer(make([]byte, 1000))}

		s3Client := &mocks2.S3{GetObjectOut: &s3.GetObjectOutput{Body: buf, ContentLength: 10001}}

		svc := &StorageService{Client: s3Client}

		res := &terrarium.SourceZipResponse{ZipDataChunk: make([]byte, 1000)}

		mds := &mocks.MockDownloadSourceZipServer{SendResponse: res}

		req := &terrarium.DownloadSourceZipRequest{
			Module: &terrarium.Module{Name: "Test", Version: "v1"},
		}

		err := svc.DownloadSourceZip(req, mds)

		if s3Client.GetObjectInvocations != 1 {
			t.Errorf("Expected 1 call to GetObject, got %v", s3Client.GetObjectInvocations)
		}

		if mds.SendInvocations != 0 {
			t.Errorf("Expected 0 calls to Send, got %v", mds.SendInvocations)
		}

		if err != ContentLenghtError {
			t.Errorf("Expected %v, got %v.", ContentLenghtError, err)
		}
	})
}
