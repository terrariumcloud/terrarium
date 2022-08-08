package services_test

import (
	"io"
	"os"

	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/google/uuid"
	services "github.com/terrariumcloud/terrarium-grpc-gateway/internal/module/services"
	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"
)

type fakeUploadServer struct {
	services.Storage_UploadSourceZipServer
	response          *terrarium.TransactionStatusResponse
	err               error
	numberOfRecvCalls int
}

func (fus *fakeUploadServer) SendAndClose(response *terrarium.TransactionStatusResponse) error {
	fus.response = response
	return fus.err
}

func (fus *fakeUploadServer) Recv() (*terrarium.UploadSourceZipRequest, error) {
	fus.numberOfRecvCalls++
	f, err := os.ReadFile("storage.go")
	if err != nil {
		return nil, err
	}

	chunk := &terrarium.UploadSourceZipRequest{
		SessionKey:   uuid.NewString(),
		ZipDataChunk: f,
	}

	if fus.numberOfRecvCalls > 1 {
		return chunk, io.EOF
	} else {
		return chunk, nil
	}
}

type fakeS3Service struct {
	s3iface.S3API
	err                  error
	numberOfPutItemCalls int
}

func (f *fakeS3Service) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	output := new(s3.PutObjectOutput)
	f.numberOfPutItemCalls++
	return output, f.err
}

func TestUploadSourceZip(t *testing.T) {
	t.Parallel()

	storageService := &services.StorageService{
		S3: &fakeS3Service{},
	}
	fus := &fakeUploadServer{}

	err := storageService.UploadSourceZip(fus)
	if err != nil {
		t.Errorf("Unable to upload file, %v", err)
	} else {
		t.Log("Successfully uploaded file.")
	}
}

func IgnoreTestUploadSourceZipE2E(t *testing.T) {
	t.Parallel()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	storageService := &services.StorageService{
		S3: s3.New(sess),
	}
	fus := &fakeUploadServer{}

	err := storageService.UploadSourceZip(fus)
	if err != nil {
		t.Errorf("Unable to upload file, %v", err)
	} else {
		t.Log("Successfully uploaded file.")
	}
}
