package services

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"

	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"google.golang.org/grpc/metadata"
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

// Upload Source Zip to storage
func (s *StorageService) UploadSourceZip(server Storage_UploadSourceZipServer) error {
	zip := []byte{}
	var sessionKey string

	if md, ok := metadata.FromIncomingContext(server.Context()); ok {
		sessionKey = md.Get("session_key")[0]
	}

	for {
		req, err := server.Recv()

		if err == io.EOF {
			log.Printf("Received file with total lenght: %v", len(zip))
			in := &s3.PutObjectInput{
				Bucket: aws.String(BucketName),
				Key:    aws.String(fmt.Sprintf("%s.zip", sessionKey)),
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
	//TODO: fetch session key based on request data
	sessionKey := "123"

	in := &s3.GetObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(sessionKey),
	}

	out, err := s.S3.GetObject(in)

	if err != nil {
		log.Printf("Failed to get object: %s", err.Error())
		return err
	}

	buf := make([]byte, *out.ContentLength)

	n, err := out.Body.Read(buf)
        outContentLength := int(*out.ContentLength)
	if n == outContentLength {
		res := &terrarium.SourceZipResponse{}
		for i := 0; i < outContentLength; i += ChunkSize {
			if i+ChunkSize > outContentLength {
				res.ZipDataChunk = buf[i:outContentLength]
			} else {
				res.ZipDataChunk = buf[i : i+ChunkSize]
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
