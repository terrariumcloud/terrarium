package storage

import (
	"fmt"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/services"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium"
)

const (
	DefaultBucketName                    = "terrarium-dev"
	DefaultStorageServiceDefaultEndpoint = "storage_service:3001"
)

var BucketName string = DefaultBucketName
var StorageServiceEndpoint string = DefaultStorageServiceDefaultEndpoint

type StorageService struct {
	services.UnimplementedStorageServer
	S3 s3iface.S3API
}

func (s *StorageService) UploadSourceZip(server services.Storage_UploadSourceZipServer) error {
	f, err := os.CreateTemp("/tmp", "upload*.zip")
	if err != nil {
		return err
	}

	for {
		chunk, err := server.Recv()

		if chunk != nil {
			if _, err := f.Write(chunk.GetZipDataChunk()); err != nil {
				return err
			}
		}

		if err == io.EOF {
			f.Seek(0, 0)
			input := &s3.PutObjectInput{
				Bucket: aws.String(BucketName),
				Key:    aws.String(fmt.Sprintf("%s.zip", chunk.GetSessionKey())),
				Body:   f,
			}
			_, err := s.S3.PutObject(input)
			if err != nil {
				return err
			}

			err = server.SendAndClose(Ok("Uploaded successfully."))
			if err != nil {
				return err
			}

			return nil
		}

		if err != nil {
			server.SendAndClose(Error("Something went wrong."))
			return err
		}
	}
}

func (s *StorageService) DownloadSourceZip(request *terrarium.DownloadSourceZipRequest, server services.Storage_DownloadSourceZipServer) error {
	return nil
}

func Ok(message string) *terrarium.TransactionStatusResponse {
	return &terrarium.TransactionStatusResponse{
		Status:        terrarium.Status_OK,
		StatusMessage: message,
	}
}

func Error(message string) *terrarium.TransactionStatusResponse {
	return &terrarium.TransactionStatusResponse{
		Status:        terrarium.Status_UNKNOWN_ERROR,
		StatusMessage: message,
	}
}
