package services

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/terrariumcloud/terrarium/internal/storage"
	terrarium "github.com/terrariumcloud/terrarium/pkg/terrarium/module"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DefaultBucketName                    = "terrarium-modules"
	DefaultStorageServiceDefaultEndpoint = "storage:3001"
	DefaultChunkSize                     = 64 * 1024 // 64 KB
)

var (
	BucketName             string = DefaultBucketName
	StorageServiceEndpoint string = DefaultStorageServiceDefaultEndpoint
	ChunkSize                     = DefaultChunkSize

	SourceZipUploaded = &terrarium.Response{Message: "Source zip uploaded successfully."}

	BucketInitializationError = status.Error(codes.Unknown, "Failed to initialize bucket for storage.")
	UploadSourceZipError      = status.Error(codes.Unknown, "Failed to upload source zip.")
	RecieveSourceZipError     = status.Error(codes.Unknown, "Failed to recieve source zip.")
	DownloadSourceZipError    = status.Error(codes.Unknown, "Failed to download source zip.")
	SendSourceZipError        = status.Error(codes.Unknown, "Failed to send source zip.")
	ContentLenghtError        = status.Error(codes.Unknown, "Failed to read correct content lenght.")
)

type StorageService struct {
	UnimplementedStorageServer
	S3         s3iface.S3API
	BucketName string
	Region     string
}

// Registers StorageService with grpc server
func (s *StorageService) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	if err := storage.InitializeS3Bucket(s.BucketName, s.Region, s.S3); err != nil {
		log.Println(err)
		return BucketInitializationError
	}

	RegisterStorageServer(grpcServer, s)

	return nil
}

// Upload Source Zip to storage
func (s *StorageService) UploadSourceZip(server Storage_UploadSourceZipServer) error {
	log.Println("Uploading source zip.")
	zip := []byte{}
	var filename string

	for {
		req, err := server.Recv()

		if filename == "" && req != nil {
			filename = fmt.Sprintf("%s/%s.zip", req.Module.GetName(), req.Module.GetVersion())
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
				return UploadSourceZipError
			}

			log.Println("Source zip uploaded successfully.")
			return server.SendAndClose(SourceZipUploaded)
		}

		if err != nil {
			log.Println(err)
			return RecieveSourceZipError
		}

		log.Printf("Recieved %v bytes", len(req.ZipDataChunk))
		zip = append(zip, req.ZipDataChunk...)
	}
}

// Download Source Zip from storage
func (s *StorageService) DownloadSourceZip(request *terrarium.DownloadSourceZipRequest, server Storage_DownloadSourceZipServer) error {
	log.Println("Downloading source zip.")
	filename := fmt.Sprintf("%s/%s.zip", request.GetModule().Name, request.Module.GetVersion())

	in := &s3.GetObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(filename),
	}

	out, err := s.S3.GetObject(in)

	if err != nil {
		log.Println(err)
		return DownloadSourceZipError
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
				log.Println(err)
				return SendSourceZipError
			}
		}

		log.Println("Source zip downloaded.")
		return nil
	} else if err != nil { // TODO: check if this is unreachable/dead code
		log.Println(err)
		return err
	} else {
		return ContentLenghtError
	}
}
