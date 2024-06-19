package storage

import (
	"bytes"
	"errors"
	mocks2 "github.com/terrariumcloud/terrarium/internal/storage/mocks"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	terrarium "github.com/terrariumcloud/terrarium/internal/provider/services"
	"github.com/terrariumcloud/terrarium/internal/provider/services/mocks"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/provider"
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
			t.Errorf("Expected 1 calls to CreateTable, got %v.", s3Client.CreateBucketInvocations)
		}
	})
}

// Test_DownloadProviderSourceZip checks:
// - if correct response is returned when source zip is downloaded
// - if error is returned when GetObject fails
// - if error is returned when Send fails
func Test_DownloadProviderSourceZip(t *testing.T) {
	t.Parallel()

	t.Run("When source zip is downloaded", func(t *testing.T) {
		var length int64 = 70000
		buf := &ClosingBuffer{bytes.NewBuffer(make([]byte, length))}

		s3Client := &mocks2.S3{GetObjectOut: &s3.GetObjectOutput{Body: buf, ContentLength: length}}

		svc := &StorageService{Client: s3Client}

		res := &terrarium.SourceZipResponse{ZipDataChunk: make([]byte, length)}

		mds := &mocks.MockDownloadProviderSourceZipServer{SendResponse: res}

		req := &terrarium.DownloadSourceZipRequest{
			Provider: &terrarium.ProviderRequest{Name: "TestOrg/TestProvider", Version: "v1", Os: "linux", Arch: "amd64"},
		}

		err := svc.DownloadProviderSourceZip(req, mds)

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}

		if s3Client.GetObjectInvocations != 1 {
			t.Errorf("Expected 1 call to GetObject, got %v", s3Client.GetObjectInvocations)
		}

		if mds.SendInvocations != 2 {
			t.Errorf("Expected 2 call to Send, got %v", mds.SendInvocations)
		}

		if !bytes.Equal(mds.TotalReceived, res.ZipDataChunk) {
			t.Errorf("Expected same data to be returned.")
		}
	})

	t.Run("when GetObject fails", func(t *testing.T) {
		s3Client := &mocks2.S3{GetObjectError: errors.New("some error")}

		svc := &StorageService{Client: s3Client}

		mds := &mocks.MockDownloadProviderSourceZipServer{}

		req := &terrarium.DownloadSourceZipRequest{
			Provider: &terrarium.ProviderRequest{Name: "TestOrg/TestProvider", Version: "v1", Os: "linux", Arch: "amd64"},
		}

		err := svc.DownloadProviderSourceZip(req, mds)

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
		var length int64 = 70000
		buf := &ClosingBuffer{bytes.NewBuffer(make([]byte, length))}

		s3Client := &mocks2.S3{GetObjectOut: &s3.GetObjectOutput{Body: buf, ContentLength: length}}

		svc := &StorageService{Client: s3Client}

		mds := &mocks.MockDownloadProviderSourceZipServer{SendError: errors.New("some error")}

		req := &terrarium.DownloadSourceZipRequest{
			Provider: &terrarium.ProviderRequest{Name: "TestOrg/TestProvider", Version: "v1", Os: "linux", Arch: "amd64"},
		}

		err := svc.DownloadProviderSourceZip(req, mds)

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
}

// Test_DownloadShasum checks:
// - if correct response is returned when shasum file is downloaded
// - if error is returned when GetObject fails
// - if error is returned when Send fails
func Test_DownloadShasum(t *testing.T) {
	t.Parallel()

	t.Run("when shasum file is downloaded", func(t *testing.T) {
		var length int64 = 70000
		buf := &ClosingBuffer{bytes.NewBuffer(make([]byte, length))}

		s3Client := &mocks2.S3{GetObjectOut: &s3.GetObjectOutput{Body: buf, ContentLength: length}}

		svc := &StorageService{Client: s3Client}

		res := &terrarium.DownloadShasumResponse{ShasumDataChunk: make([]byte, length)}

		mds := &mocks.MockDownloadProviderShasumServer{SendResponse: res}

		req := &terrarium.DownloadShasumRequest{
			Provider: &provider.Provider{Name: "TestOrg/TestProvider", Version: "v1"},
		}

		err := svc.DownloadShasum(req, mds)

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}

		if s3Client.GetObjectInvocations != 1 {
			t.Errorf("Expected 1 call to GetObject, got %v", s3Client.GetObjectInvocations)
		}

		if mds.SendInvocations != 2 {
			t.Errorf("Expected 2 call to Send, got %v", mds.SendInvocations)
		}

		if !bytes.Equal(mds.TotalReceived, res.ShasumDataChunk) {
			t.Errorf("Expected same data to be returned.")
		}
	})

	t.Run("when GetObject fails", func(t *testing.T) {
		s3Client := &mocks2.S3{GetObjectError: errors.New("some error")}

		svc := &StorageService{Client: s3Client}

		mds := &mocks.MockDownloadProviderShasumServer{}

		req := &terrarium.DownloadShasumRequest{
			Provider: &provider.Provider{Name: "TestOrg/TestProvider", Version: "v1"},
		}

		err := svc.DownloadShasum(req, mds)

		if s3Client.GetObjectInvocations != 1 {
			t.Errorf("Expected 1 call to GetObject, got %v", s3Client.GetObjectInvocations)
		}

		if mds.SendInvocations != 0 {
			t.Errorf("Expected 0 call to Sends, got %v", mds.SendInvocations)
		}

		if err != DownloadShasumError {
			t.Errorf("Expected %v, got %v.", DownloadShasumError, err)
		}
	})

	t.Run("when Send fails", func(t *testing.T) {
		var length int64 = 70000
		buf := &ClosingBuffer{bytes.NewBuffer(make([]byte, length))}

		s3Client := &mocks2.S3{GetObjectOut: &s3.GetObjectOutput{Body: buf, ContentLength: length}}

		svc := &StorageService{Client: s3Client}

		mds := &mocks.MockDownloadProviderShasumServer{SendError: errors.New("some error")}

		req := &terrarium.DownloadShasumRequest{
			Provider: &provider.Provider{Name: "TestOrg/TestProvider", Version: "v1"},
		}

		err := svc.DownloadShasum(req, mds)

		if s3Client.GetObjectInvocations != 1 {
			t.Errorf("Expected 1 call to GetObject, got %v", s3Client.GetObjectInvocations)
		}

		if mds.SendInvocations != 1 {
			t.Errorf("Expected 1 call to Send, got %v", mds.SendInvocations)
		}

		if err != SendShasumError {
			t.Errorf("Expected %v, got %v.", SendShasumError, err)
		}
	})
}

// Test_DownloadShasumSignature checks:
// - if correct response is returned when shasum signature file is downloaded
// - if error is returned when GetObject fails
// - if error is returned when Send fails
func Test_DownloadShasumSignature(t *testing.T) {
	t.Parallel()

	t.Run("when shasum signature file is downloaded", func(t *testing.T) {
		var length int64 = 70000
		buf := &ClosingBuffer{bytes.NewBuffer(make([]byte, length))}

		s3Client := &mocks2.S3{GetObjectOut: &s3.GetObjectOutput{Body: buf, ContentLength: length}}

		svc := &StorageService{Client: s3Client}

		res := &terrarium.DownloadShasumResponse{ShasumDataChunk: make([]byte, length)}

		mds := &mocks.MockDownloadProviderShasumSignatureServer{SendResponse: res}

		req := &terrarium.DownloadShasumRequest{
			Provider: &provider.Provider{Name: "TestOrg/TestProvider", Version: "v1"},
		}

		err := svc.DownloadShasumSignature(req, mds)

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}

		if s3Client.GetObjectInvocations != 1 {
			t.Errorf("Expected 1 call to GetObject, got %v", s3Client.GetObjectInvocations)
		}

		if mds.SendInvocations != 2 {
			t.Errorf("Expected 2 call to Send, got %v", mds.SendInvocations)
		}

		if !bytes.Equal(mds.TotalReceived, res.ShasumDataChunk) {
			t.Errorf("Expected same data to be returned.")
		}
	})

	t.Run("when GetObject fails", func(t *testing.T) {
		s3Client := &mocks2.S3{GetObjectError: errors.New("some error")}

		svc := &StorageService{Client: s3Client}

		mds := &mocks.MockDownloadProviderShasumSignatureServer{}

		req := &terrarium.DownloadShasumRequest{
			Provider: &provider.Provider{Name: "TestOrg/TestProvider", Version: "v1"},
		}

		err := svc.DownloadShasumSignature(req, mds)

		if s3Client.GetObjectInvocations != 1 {
			t.Errorf("Expected 1 call to GetObject, got %v", s3Client.GetObjectInvocations)
		}

		if mds.SendInvocations != 0 {
			t.Errorf("Expected 0 call to Sends, got %v", mds.SendInvocations)
		}

		if err != DownloadShasumError {
			t.Errorf("Expected %v, got %v.", DownloadShasumError, err)
		}
	})

	t.Run("when Send fails", func(t *testing.T) {
		var length int64 = 70000
		buf := &ClosingBuffer{bytes.NewBuffer(make([]byte, length))}

		s3Client := &mocks2.S3{GetObjectOut: &s3.GetObjectOutput{Body: buf, ContentLength: length}}

		svc := &StorageService{Client: s3Client}

		mds := &mocks.MockDownloadProviderShasumSignatureServer{SendError: errors.New("some error")}

		req := &terrarium.DownloadShasumRequest{
			Provider: &provider.Provider{Name: "TestOrg/TestProvider", Version: "v1"},
		}

		err := svc.DownloadShasumSignature(req, mds)

		if s3Client.GetObjectInvocations != 1 {
			t.Errorf("Expected 1 call to GetObject, got %v", s3Client.GetObjectInvocations)
		}

		if mds.SendInvocations != 1 {
			t.Errorf("Expected 1 call to Send, got %v", mds.SendInvocations)
		}

		if err != SendShasumError {
			t.Errorf("Expected %v, got %v.", SendShasumError, err)
		}
	})
}

// Test_UploadProviderBinaryZip checks:
// - if correct response is returned when the binary zip is uploaded
// - if error is returned when PutObject fails
// - if error is returned when Recv fails
func Test_UploadProviderBinaryZip(t *testing.T) {
	t.Parallel()

	t.Run("when binary zip is uploaded", func(t *testing.T) {
		s3Client := &mocks2.S3{}

		svc := &StorageService{Client: s3Client}

		req := &provider.UploadProviderBinaryZipRequest{
			Provider: &provider.Provider{Name: "TestOrg/TestProvider", Version: "v1"},
			Os: "linux",
			Arch: "amd64",
			ZipDataChunk: make([]byte, 1000),
		}

		mus := &mocks.MockUploadProviderBinaryZipServer{
			RecvRequest: req,
			RecvMaxInvocations: 2,
		}

		err := svc.UploadProviderBinaryZip(mus)

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

		if mus.SendAndCloseResponse != BinaryZipUploaded {
			t.Errorf("Expected %v, got %v.", BinaryZipUploaded, mus.SendAndCloseResponse)
		}
	})

	t.Run("when PutObject fails", func(t *testing.T) {
		s3Client := &mocks2.S3{PutObjectError: errors.New("some error")}

		svc := &StorageService{Client: s3Client}

		req := &provider.UploadProviderBinaryZipRequest{
			Provider:       &provider.Provider{Name: "TestOrg/TestProvider", Version: "v1"},
			Os: "linux",
			Arch: "amd64",
			ZipDataChunk: make([]byte, 1000),
		}

		mus := &mocks.MockUploadProviderBinaryZipServer{
			RecvRequest: req,
			RecvMaxInvocations: 2,
		}

		err := svc.UploadProviderBinaryZip(mus)

		if mus.RecvInvocations != 2 {
			t.Errorf("Expected 1 call to Recv, got %v", mus.RecvInvocations)
		}

		if s3Client.PutObjectInvocations != 1 {
			t.Errorf("Expected 1 call to PutObject, got %v", s3Client.PutObjectInvocations)
		}

		if mus.SendAndCloseInvocations != 0 {
			t.Errorf("Expected 0 calls to SendAndClose, got %v", mus.SendAndCloseInvocations)
		}

		if err != UploadBinaryZipError {
			t.Errorf("Expected %v, got %v.", UploadBinaryZipError, err)
		}
	})

	t.Run("when Recv fails", func(t *testing.T) {
		s3Client := &mocks2.S3{}

		svc := &StorageService{Client: s3Client}

		mus := &mocks.MockUploadProviderBinaryZipServer{
			RecvError:          errors.New("some error"),
			RecvMaxInvocations: 1,
		}

		err := svc.UploadProviderBinaryZip(mus)

		if mus.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", mus.RecvInvocations)
		}

		if s3Client.PutObjectInvocations != 0 {
			t.Errorf("Expected 0 calls to PutObject, got %v", s3Client.PutObjectInvocations)
		}

		if mus.SendAndCloseInvocations != 0 {
			t.Errorf("Expected 0 calls to SendAndClose, got %v", mus.SendAndCloseInvocations)
		}

		if err != ReceiveBinaryZipError {
			t.Errorf("Expected %v, got %v.", ReceiveBinaryZipError, err)
		}
	})
}

// Test_UploadShasum checks:
// - if correct response is returned when the shasum file is uploaded
// - if error is returned when PutObject fails
// - if error is returned when Recv fails
func Test_UploadShasum(t *testing.T) {
	t.Parallel()

	t.Run("when shasum file is uploaded", func(t *testing.T) {
		s3Client := &mocks2.S3{}

		svc := &StorageService{Client: s3Client}

		req := &provider.UploadShasumRequest{
			Provider: &provider.Provider{Name: "TestOrg/TestProvider", Version: "v1"},
			ShasumDataChunk: make([]byte, 1000),
		}

		mus := &mocks.MockUploadShasumServer{
			RecvRequest: req,
			RecvMaxInvocations: 2,
		}

		err := svc.UploadShasum(mus)

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

		if mus.SendAndCloseResponse != ShasumUploaded {
			t.Errorf("Expected %v, got %v.", ShasumUploaded, mus.SendAndCloseResponse)
		}
	})

	t.Run("when PutObject fails", func(t *testing.T) {
		s3Client := &mocks2.S3{PutObjectError: errors.New("some error")}

		svc := &StorageService{Client: s3Client}

		req := &provider.UploadShasumRequest{
			Provider: &provider.Provider{Name: "TestOrg/TestProvider", Version: "v1"},
			ShasumDataChunk: make([]byte, 1000),
		}

		mus := &mocks.MockUploadShasumServer{
			RecvRequest: req,
			RecvMaxInvocations: 2,
		}

		err := svc.UploadShasum(mus)

		if mus.RecvInvocations != 2 {
			t.Errorf("Expected 1 call to Recv, got %v", mus.RecvInvocations)
		}

		if s3Client.PutObjectInvocations != 1 {
			t.Errorf("Expected 1 call to PutObject, got %v", s3Client.PutObjectInvocations)
		}

		if mus.SendAndCloseInvocations != 0 {
			t.Errorf("Expected 0 calls to SendAndClose, got %v", mus.SendAndCloseInvocations)
		}

		if err != UploadShasumError {
			t.Errorf("Expected %v, got %v.", UploadShasumError, err)
		}
	})

	t.Run("when Recv fails", func(t *testing.T) {
		s3Client := &mocks2.S3{}

		svc := &StorageService{Client: s3Client}

		mus := &mocks.MockUploadShasumServer{
			RecvError:          errors.New("some error"),
			RecvMaxInvocations: 1,
		}

		err := svc.UploadShasum(mus)

		if mus.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", mus.RecvInvocations)
		}

		if s3Client.PutObjectInvocations != 0 {
			t.Errorf("Expected 0 calls to PutObject, got %v", s3Client.PutObjectInvocations)
		}

		if mus.SendAndCloseInvocations != 0 {
			t.Errorf("Expected 0 calls to SendAndClose, got %v", mus.SendAndCloseInvocations)
		}

		if err != ReceiveShasumError {
			t.Errorf("Expected %v, got %v.", ReceiveShasumError, err)
		}
	})
}

// Test_UploadShasumSignature checks:
// - if correct response is returned when the shasum signature file is uploaded
// - if error is returned when PutObject fails
// - if error is returned when Recv fails
func Test_UploadShasumSignature(t *testing.T) {
	t.Parallel()

	t.Run("when shasum signature file is uploaded", func(t *testing.T) {
		s3Client := &mocks2.S3{}

		svc := &StorageService{Client: s3Client}

		req := &provider.UploadShasumRequest{
			Provider: &provider.Provider{Name: "TestOrg/TestProvider", Version: "v1"},
			ShasumDataChunk: make([]byte, 1000),
		}

		mus := &mocks.MockUploadShasumSignatureServer{
			RecvRequest: req,
			RecvMaxInvocations: 2,
		}

		err := svc.UploadShasumSignature(mus)

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

		if mus.SendAndCloseResponse != ShasumSigUploaded {
			t.Errorf("Expected %v, got %v.", ShasumSigUploaded, mus.SendAndCloseResponse)
		}
	})

	t.Run("when PutObject fails", func(t *testing.T) {
		s3Client := &mocks2.S3{PutObjectError: errors.New("some error")}

		svc := &StorageService{Client: s3Client}

		req := &provider.UploadShasumRequest{
			Provider: &provider.Provider{Name: "TestOrg/TestProvider", Version: "v1"},
			ShasumDataChunk: make([]byte, 1000),
		}

		mus := &mocks.MockUploadShasumSignatureServer{
			RecvRequest: req,
			RecvMaxInvocations: 2,
		}

		err := svc.UploadShasumSignature(mus)

		if mus.RecvInvocations != 2 {
			t.Errorf("Expected 1 call to Recv, got %v", mus.RecvInvocations)
		}

		if s3Client.PutObjectInvocations != 1 {
			t.Errorf("Expected 1 call to PutObject, got %v", s3Client.PutObjectInvocations)
		}

		if mus.SendAndCloseInvocations != 0 {
			t.Errorf("Expected 0 calls to SendAndClose, got %v", mus.SendAndCloseInvocations)
		}

		if err != UploadShasumSigError {
			t.Errorf("Expected %v, got %v.", UploadShasumSigError, err)
		}
	})

	t.Run("when Recv fails", func(t *testing.T) {
		s3Client := &mocks2.S3{}

		svc := &StorageService{Client: s3Client}

		mus := &mocks.MockUploadShasumSignatureServer{
			RecvError:          errors.New("some error"),
			RecvMaxInvocations: 1,
		}

		err := svc.UploadShasumSignature(mus)

		if mus.RecvInvocations != 1 {
			t.Errorf("Expected 1 call to Recv, got %v", mus.RecvInvocations)
		}

		if s3Client.PutObjectInvocations != 0 {
			t.Errorf("Expected 0 calls to PutObject, got %v", s3Client.PutObjectInvocations)
		}

		if mus.SendAndCloseInvocations != 0 {
			t.Errorf("Expected 0 calls to SendAndClose, got %v", mus.SendAndCloseInvocations)
		}

		if err != ReceiveShasumSigError {
			t.Errorf("Expected %v, got %v.", ReceiveShasumSigError, err)
		}
	})
}
