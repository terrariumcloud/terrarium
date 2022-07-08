package main

import (
	"io"
	"os"
	"terrarium-grpc-gateway/internal/services"

	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/google/uuid"
	"github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium"
)

type fakeUploadServer struct {
	services.Storage_UploadSourceZipServer
	response *terrarium.TransactionStatusResponse
	err error
}

func (fus *fakeUploadServer) SendAndClose(response *terrarium.TransactionStatusResponse) error {
	fus.response = response
	return fus.err
}

func (fus *fakeUploadServer) Recv() (*terrarium.UploadSourceZipChunkRequest, error) {
	f, err := os.ReadFile("main.go")
	if err != nil {
		return nil, err
	}

	chunk := &terrarium.UploadSourceZipChunkRequest{
		SessionKey:   uuid.NewString(),
		ZipDataChunk: f,
	}

	return chunk, io.EOF
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
	t.Run("It creates entry in DynamoDB", func(t *testing.T) {
		storageService := &StorageService{
			s3: &fakeS3Service{},
		}
		fus := &fakeUploadServer{}

		err := storageService.UploadSourceZip(fus)
		if err != nil {
			t.Errorf("Unable to upload file, %v", err)
		} else {
			t.Log("Successfully uploaded file.\n")
		}
	})
}

func TestUploadSourceZipE2E(t *testing.T) {
	t.Run("It creates entry in DynamoDB", func(t *testing.T) {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		storageService := &StorageService{
			s3: s3.New(sess),
		}
		fus := &fakeUploadServer{}

		err := storageService.UploadSourceZip(fus)
		if err != nil {
			t.Errorf("Unable to upload file, %v", err)
		} else {
			t.Log("Successfully uploaded file.\n")
		}
	})
}
