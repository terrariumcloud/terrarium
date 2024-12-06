package storage

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/terrariumcloud/terrarium/internal/module/services"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/terrariumcloud/terrarium/internal/storage"
	terrarium "github.com/terrariumcloud/terrarium/pkg/terrarium/module"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DefaultBucketName                          = "terrarium-modules-ciedev-4757"
	DefaultStorageServiceDefaultEndpoint       = "storage:3001"
	DefaultChunkSize                     int64 = 64 * 1024 // 64 KB
)

var (
	BucketName             = DefaultBucketName
	StorageServiceEndpoint = DefaultStorageServiceDefaultEndpoint
	ChunkSize              = DefaultChunkSize

	SourceZipUploaded = &terrarium.Response{Message: "Source zip uploaded successfully."}

	BucketInitializationError = status.Error(codes.Unknown, "Failed to initialize bucket for storage.")
	UploadSourceZipError      = status.Error(codes.Unknown, "Failed to upload source zip.")
	RecieveSourceZipError     = status.Error(codes.Unknown, "Failed to recieve source zip.")
	DownloadSourceZipError    = status.Error(codes.Unknown, "Failed to download source zip.")
	SendSourceZipError        = status.Error(codes.Unknown, "Failed to send source zip.")
	ContentLenghtError        = status.Error(codes.Unknown, "Failed to read correct content lenght.")
)

type StorageService struct {
	services.UnimplementedStorageServer
	Client     storage.AWSS3BucketClient
	BucketName string
	Region     string
}

// Registers StorageService with grpc server
func (s *StorageService) RegisterWithServer(grpcServer grpc.ServiceRegistrar) error {
	if err := storage.InitializeS3Bucket(s.BucketName, s.Region, s.Client); err != nil {
		log.Println(err)
		return BucketInitializationError
	}

	services.RegisterStorageServer(grpcServer, s)

	return nil
}

// Upload Source Zip to storage
func (s *StorageService) UploadSourceZip(server services.Storage_UploadSourceZipServer) error {
	log.Println("Uploading source zip.")
	zip := []byte{}
	var filename string
	ctx := server.Context()
	span := trace.SpanFromContext(ctx)
	for {
		req, err := server.Recv()
		span.SetAttributes(
			attribute.String("module.name", req.GetModule().GetName()),
			attribute.String("module.version", req.GetModule().GetVersion()),
		)

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

			if _, err := s.Client.PutObject(ctx, in); err != nil {
				span.RecordError(err)
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
func (s *StorageService) DownloadSourceZip(request *terrarium.DownloadSourceZipRequest, server services.Storage_DownloadSourceZipServer) error {
	log.Println("Downloading source zip.")
	ctx := server.Context()
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("module.name", request.GetModule().GetName()),
		attribute.String("module.version", request.GetModule().GetVersion()),
	)
	filename := fmt.Sprintf("%s/%s.zip", request.GetModule().Name, request.Module.GetVersion())

	in := &s3.GetObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(filename),
	}

	out, err := s.Client.GetObject(ctx, in)
	if err != nil {
		span.RecordError(err)
		log.Println(err)
		return DownloadSourceZipError
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(out.Body)
	bb := buf.Bytes()

	outContentLength := out.ContentLength
	if int64(len(bb)) == outContentLength {
		res := &terrarium.SourceZipResponse{}
		for i := int64(0); i < outContentLength; i += ChunkSize {
			if i+ChunkSize > outContentLength {
				res.ZipDataChunk = bb[i:outContentLength]
			} else {
				res.ZipDataChunk = bb[i : i+ChunkSize]
			}

			if err := server.Send(res); err != nil {
				span.RecordError(err)
				log.Println(err)
				return SendSourceZipError
			}
		}

		log.Println("Source zip downloaded.")
		return nil
	} else if err != nil {
		// TODO: check if this is unreachable/dead code
		log.Println(err)
		span.RecordError(err)
		return err
	} else {
		span.RecordError(err)
		return ContentLenghtError
	}
}
