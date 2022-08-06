package services

import (
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"
)

const (
	DefaultBucketName                    = "terrarium-dev"
	DefaultStorageServiceDefaultEndpoint = "storage_service:3001"
)

var BucketName string = DefaultBucketName
var StorageServiceEndpoint string = DefaultStorageServiceDefaultEndpoint

type StorageService struct {
	UnimplementedStorageServer
	S3 s3iface.S3API
}

func (s *StorageService) UploadSourceZip(server Storage_UploadSourceZipServer) error {
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

			err = server.SendAndClose(ZipUploaded)
			if err != nil {
				return err
			}

			return nil
		}

		if err != nil {
			server.SendAndClose(ZipUploadFailed)
			return err
		}
	}
}

func (s *StorageService) DownloadSourceZip(request *terrarium.DownloadSourceZipRequest, server Storage_DownloadSourceZipServer) error {

	sessionKey := "123"

	input := &s3.GetObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(sessionKey),
	}

	output, err := s.S3.GetObject(input)
	if err != nil {
		return err
	}

	buf := make([]byte, *output.ContentLength)

	n, err := output.Body.Read(buf)
	if n > 0 && err == io.EOF {
		response := &terrarium.SourceZipResponse{
			ZipDataChunk: buf,
		}

		if err := server.Send(response); err != nil {
			return err
		}

		return nil
	}

	return nil
}
