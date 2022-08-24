package services

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/storage"
	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	grpc "google.golang.org/grpc"
)

const (
	DefaultBucketName                    = "terrarium-modules"
	DefaultStorageServiceDefaultEndpoint = "storage:3001"
	DefaultChunkSize                     = 64 * 1024 // 64 KB
)

var BucketName string = DefaultBucketName
var StorageServiceEndpoint string = DefaultStorageServiceDefaultEndpoint
var ChunkSize = DefaultChunkSize

type StorageService struct {
	UnimplementedStorageServer
	S3         s3iface.S3API
	BucketName string
	Region     string
}

// Registers StorageService with grpc server
func (s *StorageService) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	RegisterStorageServer(grpcServer, s)

	if err := storage.InitializeS3Bucket(s.BucketName, s.Region, s.S3); err != nil {
		return err
	}

	return nil
}

// Upload Source Zip to storage
func (s *StorageService) UploadSourceZip(server Storage_UploadSourceZipServer) error {
	zip := []byte{}
	var filename string

	for {
		req, err := server.Recv()

		if filename == "" && req != nil  {
			filename = fmt.Sprintf("%s_%s.zip", req.Module.GetName(), req.Module.GetVersion())
		}
		
		if err == io.EOF {
			log.Printf("Received file with total lenght: %v", len(zip))

			in := &s3.PutObjectInput{
				Bucket: aws.String(BucketName),
				Key:    aws.String(filename),
				Body:   bytes.NewReader(zip),
			}

			if _, err := s.S3.PutObject(in); err != nil {
				log.Println(err)
				return err
			}

			log.Println("Source zip uploaded successfully.")
			return server.SendAndClose(ZipUploaded)
		}

		if err != nil {
			return err
		}

		log.Printf("Recieved %v bytes", len(req.ZipDataChunk))
		zip = append(zip, req.ZipDataChunk...)
	}
}

// Download Source Zip from storage
func (s *StorageService) DownloadSourceZip(request *terrarium.DownloadSourceZipRequest, server Storage_DownloadSourceZipServer) error {
	log.Println("Downloading source zip.")
	filename := fmt.Sprintf("%s_%s.zip", request.GetModule().Name, request.Module.GetVersion())

	in := &s3.GetObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(filename),
	}

	out, err := s.S3.GetObject(in)

	if err != nil {
		log.Printf("Failed to get object: %s", err.Error())
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(out.Body)
	bb := buf.Bytes()

	outContentLength := int(*out.ContentLength)
	if len(bb) == outContentLength {
		res := &terrarium.SourceZipResponse{}
		for i := 0; i < outContentLength; i += ChunkSize {
			if i+ChunkSize > outContentLength {
				res.ZipDataChunk = bb[i:outContentLength]
			} else {
				res.ZipDataChunk = bb[i : i+ChunkSize]
			}

			if err := server.Send(res); err != nil {
				log.Printf("Failed to send: %s", err.Error())
				return err
			}
		}

		return nil
	} else if err != nil {
		log.Printf("Failed to read content: %s", err.Error())
		return err
	} else {
		return errors.New("unexpected content lenght")
	}
}
