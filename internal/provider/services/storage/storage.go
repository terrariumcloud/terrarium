package storage

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/terrariumcloud/terrarium/internal/provider/services"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/terrariumcloud/terrarium/internal/storage"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DefaultBucketName                          = "terrarium-providers"
	DefaultStorageServiceDefaultEndpoint       = "storage:3001"
	DefaultChunkSize                     int64 = 64 * 1024 // 64 KB
)

var (
	BucketName             = DefaultBucketName
	StorageServiceEndpoint = DefaultStorageServiceDefaultEndpoint
	ChunkSize              = DefaultChunkSize

	BucketInitializationError = status.Error(codes.Unknown, "Failed to initialize bucket for storage.")
	DownloadSourceZipError    = status.Error(codes.Unknown, "Failed to download source zip.")
	SendSourceZipError        = status.Error(codes.Unknown, "Failed to send source zip.")
	SendShasumError           = status.Error(codes.Unknown, "Failed to send shasum file.")
	ContentLengthError        = status.Error(codes.Unknown, "Failed to read correct content length.")
	DownloadShasumError       = status.Error(codes.Unknown, "Failed to download shasum.")
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
		log.Println("Error initializing S3 bucket for Provider storage", err)
		return BucketInitializationError
	}

	services.RegisterStorageServer(grpcServer, s)

	return nil
}

// Download Source Zip from storage
func (s *StorageService) DownloadProviderSourceZip(request *services.DownloadSourceZipRequest, server services.Storage_DownloadProviderSourceZipServer) error {

	log.Println("Downloading source zip.")

	ctx := server.Context()
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("provider.name", request.GetProvider().GetName()),
		attribute.String("provider.version", request.GetProvider().GetVersion()),
		attribute.String("provider.os", request.GetProvider().GetOs()),
		attribute.String("provider.arch", request.GetProvider().GetArch()),
	)

	providerAddress := strings.Split(request.Provider.GetName(), "/")
	filename := fmt.Sprintf("terraform-provider-%s_%s_%s_%s.zip", providerAddress[1], request.GetProvider().GetVersion(), request.GetProvider().GetOs(), request.GetProvider().GetArch())
	fileLocation := services.ResolveS3Locations(request.Provider.GetName(), request.Provider.GetVersion(), filename)

	in := &s3.GetObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(fileLocation),
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
		res := &services.SourceZipResponse{}
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
	} else {
		span.RecordError(err)
		return ContentLengthError
	}
}

// Download Shasum from storage
func (s *StorageService) DownloadShasum(request *services.DownloadShasumRequest, server services.Storage_DownloadShasumServer) error {

	log.Println("Downloading shasum file.")

	ctx := server.Context()
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("provider.name", request.GetProvider().GetName()),
		attribute.String("provider.version", request.GetProvider().GetVersion()),
	)

	providerAddress := strings.Split(request.GetProvider().GetName(), "/")
	prefix := fmt.Sprintf("terraform-provider-%s_%s_SHA256SUMS", providerAddress[1], request.GetProvider().GetVersion())
	fileLocation := services.ResolveS3Locations(request.Provider.GetName(), request.Provider.GetVersion(), prefix)

	in := &s3.GetObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(fileLocation),
	}

	out, err := s.Client.GetObject(ctx, in)
	if err != nil {
		span.RecordError(err)
		log.Println(err)
		return DownloadShasumError
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(out.Body)
	bb := buf.Bytes()

	outContentLength := out.ContentLength

	if int64(len(bb)) == outContentLength {
		res := &services.DownloadShasumResponse{}
		for i := int64(0); i < outContentLength; i += ChunkSize {
			if i+ChunkSize > outContentLength {
				res.ShasumDataChunk = bb[i:outContentLength]
			} else {
				res.ShasumDataChunk = bb[i : i+ChunkSize]
			}

			if err := server.Send(res); err != nil {
				span.RecordError(err)
				log.Println(err)
				return SendShasumError
			}
		}

		log.Println("Shasum file downloaded.")
		return nil
	} else {
		span.RecordError(err)
		return ContentLengthError
	}
}

// Download Shasum Signature from storage
func (s *StorageService) DownloadShasumSignature(request *services.DownloadShasumRequest, server services.Storage_DownloadShasumSignatureServer) error {

	log.Println("Downloading shasum signature file.")

	ctx := server.Context()
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("provider.name", request.GetProvider().GetName()),
		attribute.String("provider.version", request.GetProvider().GetVersion()),
	)

	providerAddress := strings.Split(request.GetProvider().GetName(), "/")
	prefix := fmt.Sprintf("terraform-provider-%s_%s_SHA256SUMS.sig", providerAddress[1], request.GetProvider().GetVersion())
	fileLocation := services.ResolveS3Locations(request.Provider.GetName(), request.Provider.GetVersion(), prefix)

	in := &s3.GetObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(fileLocation),
	}

	out, err := s.Client.GetObject(ctx, in)
	if err != nil {
		span.RecordError(err)
		log.Println(err)
		return DownloadShasumError
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(out.Body)
	bb := buf.Bytes()

	outContentLength := out.ContentLength

	if int64(len(bb)) == outContentLength {
		res := &services.DownloadShasumResponse{}
		for i := int64(0); i < outContentLength; i += ChunkSize {
			if i+ChunkSize > outContentLength {
				res.ShasumDataChunk = bb[i:outContentLength]
			} else {
				res.ShasumDataChunk = bb[i : i+ChunkSize]
			}

			if err := server.Send(res); err != nil {
				span.RecordError(err)
				log.Println(err)
				return SendShasumError
			}
		}

		log.Println("Shasum signature file downloaded.")
		return nil
	} else {
		span.RecordError(err)
		return ContentLengthError
	}
}
