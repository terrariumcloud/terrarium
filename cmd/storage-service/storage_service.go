package main

import (
	"io"
	"os"
	"terrarium-grpc-gateway/internal/services"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium"
)

type StorageService struct {
	services.UnimplementedStorageServer
	uploader s3manager.Uploader
}

func (s *StorageService) UploadSourceZip(server services.Storage_UploadSourceZipServer) error {
	f, err := os.CreateTemp("/tmp", "upload*.zip")
	if err != nil {
		return err
	}
	for {
		chunk, err := server.Recv()
		if err != nil {
			return err
		}
		if err == io.EOF {
			_, err := s.uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String("terrarium-dev"),
				Key:    aws.String("test.file"),
				Body:   f,
			})
			if err != nil {
				return err
			}
			server.SendAndClose(&terrarium.TransactionStatusResponse{
				Status:        terrarium.Status_OK,
				StatusMessage: "Uploaded successfully.",
			})
			return nil
		}
		_, err = f.Write(chunk.ZipDataChunk)

		if err != nil {
			server.SendAndClose(&terrarium.TransactionStatusResponse{
				Status:        terrarium.Status_UNKNOWN_ERROR,
				StatusMessage: "Something went wrong.",
			})
			return err
		}
	}
}

func (s *StorageService) DownloadSourceZip(request *terrarium.DownloadSourceZipRequest, server services.Storage_DownloadSourceZipServer) error {
	return nil
}
