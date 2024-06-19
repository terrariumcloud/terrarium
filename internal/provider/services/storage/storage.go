package storage

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/terrariumcloud/terrarium/internal/provider/services"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/terrariumcloud/terrarium/internal/storage"
	terrarium "github.com/terrariumcloud/terrarium/pkg/terrarium/provider"

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

	BinaryZipUploaded = &terrarium.Response{Message: "Binary zip uploaded successfully."}
	ShasumUploaded = &terrarium.Response{Message: "Shasum file uploaded successfully."}
	ShasumSigUploaded = &terrarium.Response{Message: "Shasum signature uploaded successfully."}


	BucketInitializationError = status.Error(codes.Unknown, "Failed to initialize bucket for storage.")
	DownloadSourceZipError    = status.Error(codes.Unknown, "Failed to download source zip.")
	SendSourceZipError        = status.Error(codes.Unknown, "Failed to send source zip.")
	SendShasumError           = status.Error(codes.Unknown, "Failed to send shasum file.")
	DownloadShasumError       = status.Error(codes.Unknown, "Failed to download shasum.")
	UploadBinaryZipError      = status.Error(codes.Unknown, "Failed to upload binary zip.")
	ReceiveBinaryZipError     = status.Error(codes.Unknown, "Failed to receive binary zip.")
	UploadShasumError      	  = status.Error(codes.Unknown, "Failed to upload shasum file.")
	ReceiveShasumError     	  = status.Error(codes.Unknown, "Failed to receive shasum file.")
	UploadShasumSigError      	  = status.Error(codes.Unknown, "Failed to upload shasum signature file.")
	ReceiveShasumSigError     	  = status.Error(codes.Unknown, "Failed to receive shasum signature file.")
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
	fileLocation := ResolveS3Locations(request.Provider.GetName(), request.Provider.GetVersion(), filename)

	in := &s3.GetObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(fileLocation),
	}

	out, err := s.Client.GetObject(ctx, in)
	if err != nil {
		span.RecordError(err)
		log.Println("Error downloading source zip for provider binary", err)
		return DownloadSourceZipError
	}

	buf := make([]byte, ChunkSize)
	res := &services.SourceZipResponse{}

	for {
		n, err := out.Body.Read(buf)
		if err != nil && err != io.EOF {
			span.RecordError(err)
			log.Println("Failed to download source zip", err)
			return DownloadSourceZipError
		}
		if n == 0 {
			break
		}
		res.ZipDataChunk = buf[:n]
		if err := server.Send(res); err != nil {
			span.RecordError(err)
			log.Println("Failed to send source zip", err)
			return SendSourceZipError
		}
	}

	log.Println("Source zip downloaded.")
	return nil

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
	suffix := fmt.Sprintf("terraform-provider-%s_%s_SHA256SUMS", providerAddress[1], request.GetProvider().GetVersion())
	fileLocation := ResolveS3Locations(request.Provider.GetName(), request.Provider.GetVersion(), suffix)

	in := &s3.GetObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(fileLocation),
	}

	out, err := s.Client.GetObject(ctx, in)
	if err != nil {
		span.RecordError(err)
		log.Println("Error downloading shasum file", err)
		return DownloadShasumError
	}

	buf := make([]byte, ChunkSize)
	res := &services.DownloadShasumResponse{}

	for {
		n, err := out.Body.Read(buf)
		if err != nil && err != io.EOF {
			span.RecordError(err)
			log.Println("Failed to download shasum file", err)
			return DownloadShasumError
		}
		if n == 0 {
			break
		}

		res.ShasumDataChunk = buf[:n]
		if err := server.Send(res); err != nil {
			span.RecordError(err)
			log.Println("Failed to send shasum file", err)
			return SendShasumError
		}
	}

	log.Println("Shasum file downloaded.")
	return nil
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
	suffix := fmt.Sprintf("terraform-provider-%s_%s_SHA256SUMS.sig", providerAddress[1], request.GetProvider().GetVersion())
	fileLocation := ResolveS3Locations(request.Provider.GetName(), request.Provider.GetVersion(), suffix)

	in := &s3.GetObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(fileLocation),
	}

	out, err := s.Client.GetObject(ctx, in)
	if err != nil {
		span.RecordError(err)
		log.Println("Error downloading shasum signature file", err)
		return DownloadShasumError
	}

	buf := make([]byte, ChunkSize)
	res := &services.DownloadShasumResponse{}

	for {
		n, err := out.Body.Read(buf)
		if err != nil && err != io.EOF {
			span.RecordError(err)
			log.Println("Failed to download shasum signature file", err)
			return DownloadShasumError
		}
		if n == 0 {
			break
		}

		res.ShasumDataChunk = buf[:n]
		if err := server.Send(res); err != nil {
			span.RecordError(err)
			log.Println("Failed to send shasum signature file", err)
			return SendShasumError
		}
	}

	log.Println("Shasum signature file downloaded.")
	return nil
}

func ResolveS3Locations(providerID, providerVersion, value string) string {
	fileLocation := fmt.Sprintf("%s/%s/%s", providerID, providerVersion, value)
	return fileLocation
}

// Upload Provider Binary Zip to storage
func (s *StorageService) UploadProviderBinaryZip(server services.Storage_UploadProviderBinaryZipServer) error {
	log.Println("Uploading provider binary zip.")

	binary_zip := []byte{}
	var fileLocation string

	ctx := server.Context()
	span := trace.SpanFromContext(ctx)

	for {
		req, err := server.Recv()

		span.SetAttributes(
			attribute.String("provider.name", req.GetProvider().GetName()),
			attribute.String("provider.version", req.GetProvider().GetVersion()),
			attribute.String("provider.os", req.GetOs()),
			attribute.String("provider.arch", req.GetArch()),
		)

		providerAddress := strings.Split(req.GetProvider().GetName(), "/")
		if fileLocation == "" && req != nil {
			filename := fmt.Sprintf("terraform-provider-%s_%s_%s_%s.zip", providerAddress[1], req.GetProvider().GetVersion(), req.GetOs(), req.GetArch())
			fileLocation = ResolveS3Locations(req.Provider.GetName(), req.GetProvider().GetVersion(), filename)
		}

		if err == io.EOF {
			log.Printf("Received file with total length: %v", len(binary_zip))

			in := &s3.PutObjectInput{
				Bucket: aws.String(BucketName),
				Key:    aws.String(fileLocation),
				Body:   bytes.NewReader(binary_zip),
			}

			if _, err := s.Client.PutObject(ctx, in); err != nil {
				span.RecordError(err)
				log.Println(err)
				return UploadBinaryZipError
			}

			log.Println("Binary zip uploaded successfully.")
			return server.SendAndClose(BinaryZipUploaded)
		}

		if err != nil {
			log.Println(err)
			return ReceiveBinaryZipError
		}

		log.Printf("Received %v bytes", len(req.ZipDataChunk))
		binary_zip = append(binary_zip, req.ZipDataChunk...)

	}
}

// Upload Shasum to storage
func (s *StorageService) UploadShasum(server services.Storage_UploadShasumServer) error {
	log.Println("Uploading shasum file.")

	shasum := []byte{}
	var fileLocation string

	ctx := server.Context()
	span := trace.SpanFromContext(ctx)

	for {
		req, err := server.Recv()

		span.SetAttributes(
			attribute.String("provider.name", req.GetProvider().GetName()),
			attribute.String("provider.version", req.GetProvider().GetVersion()),
		)

		providerAddress := strings.Split(req.GetProvider().GetName(), "/")
		if fileLocation == "" && req != nil {
			filename := fmt.Sprintf("terraform-provider-%s_%s_SHA256SUMS", providerAddress[1], req.GetProvider().GetVersion())
			fileLocation = ResolveS3Locations(req.GetProvider().GetName(), req.GetProvider().GetVersion(), filename)
		}

		if err == io.EOF {
			log.Printf("Received file with total length: %v", len(shasum))

			in := &s3.PutObjectInput{
				Bucket: aws.String(BucketName),
				Key:    aws.String(fileLocation),
				Body:   bytes.NewReader(shasum),
			}

			if _, err := s.Client.PutObject(ctx, in); err != nil {
				span.RecordError(err)
				log.Println(err)
				return UploadShasumError
			}

			log.Println("Shasum file uploaded successfully.")
			return server.SendAndClose(ShasumUploaded)
		}

		if err != nil {
			log.Println(err)
			return ReceiveShasumError
		}

		log.Printf("Received %v bytes", len(req.ShasumDataChunk))
		shasum = append(shasum, req.ShasumDataChunk...)
	}
}

// Upload Shasum Signature to storage
func (s *StorageService) UploadShasumSignature(server services.Storage_UploadShasumSignatureServer) error {
	log.Println("Uploading shasum signature.")

	shasum_sig := []byte{}
	var fileLocation string

	ctx := server.Context()
	span := trace.SpanFromContext(ctx)

	for {
		req, err := server.Recv()

		span.SetAttributes(
			attribute.String("provider.name", req.GetProvider().GetName()),
			attribute.String("provider.version", req.GetProvider().GetVersion()),
		)

		providerAddress := strings.Split(req.GetProvider().GetName(), "/")
		if fileLocation == "" && req != nil {
			filename := fmt.Sprintf("terraform-provider-%s_%s_SHA256SUMS.sig", providerAddress[1], req.GetProvider().GetVersion())
			fileLocation = ResolveS3Locations(req.GetProvider().GetName(), req.GetProvider().GetVersion(), filename)
		}

		if err == io.EOF {
			log.Printf("Received file with total length: %v", len(shasum_sig))

			in := &s3.PutObjectInput{
				Bucket: aws.String(BucketName),
				Key:    aws.String(fileLocation),
				Body:   bytes.NewReader(shasum_sig),
			}

			if _, err := s.Client.PutObject(ctx, in); err != nil {
				span.RecordError(err)
				log.Println(err)
				return UploadShasumSigError
			}

			log.Println("Shasum signature uploaded successfully.")
			return server.SendAndClose(ShasumSigUploaded)
		}

		if err != nil {
			log.Println(err)
			return ReceiveShasumSigError
		}

		log.Printf("Received %v bytes", len(req.ShasumDataChunk))
		shasum_sig = append(shasum_sig, req.ShasumDataChunk...)
	}
}
