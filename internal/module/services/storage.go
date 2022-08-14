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
	DefaultBucketName                    = "terrarium-dev" //TODO: rename to terrarium or terrarium-modules?
	DefaultStorageServiceDefaultEndpoint = "storage:3001"
	DefaultChunkSize                     = 64 * 1024 // 64 KB
)

var BucketName string = DefaultBucketName
var StorageServiceEndpoint string = DefaultStorageServiceDefaultEndpoint
var ChunkSize = DefaultChunkSize

type StorageService struct {
	UnimplementedStorageServer
	S3 s3iface.S3API
}

// Upload Source Zip to storage
func (s *StorageService) UploadSourceZip(server Storage_UploadSourceZipServer) error {
	//TODO: send file hash as metadata and verify
	zip := []byte{}
	if md, ok := metadata.FromIncomingContext(server.Context()); ok {
		apiKey := md["api_key"]
		log.Println(apiKey)
	}

	for {
		req, err := server.Recv()

		if err == io.EOF {
			log.Printf("Received file with lenght: %v", len(zip))
			in := &s3.PutObjectInput{
				Bucket: aws.String(BucketName),
				Key:    aws.String(fmt.Sprintf("%s.zip", "123")),
				Body:   bytes.NewReader(zip),
			}

			if _, err := s.S3.PutObject(in); err != nil {
				log.Println(err)
				return err
			}

			return server.SendAndClose(ZipUploaded)
		}

		if err != nil {
			return err
		}

		log.Printf("Recieved %v bytes", len(req.ZipDataChunk))
		zip = append(zip, req.ZipDataChunk...)
	}
}

// func (s *StorageService) UploadSourceZip2(server Storage_UploadSourceZipServer) error {
// 	f, err := os.CreateTemp("/tmp", "upload*.zip")

// 	if err != nil {
// 		return err
// 	}

// 	for {
// 		chunk, err := server.Recv()
// 		var sessionKey  string

// 		if chunk != nil {
// 			sessionKey = chunk.GetSessionKey()
// 			if _, err := f.Write(chunk.GetZipDataChunk()); err != nil {
// 				return err
// 			}
// 		}

// 		if err == io.EOF {
// 			f.Seek(0, 0)
// 			in := &s3.PutObjectInput{
// 				Bucket: aws.String(BucketName),
// 				Key:    aws.String(fmt.Sprintf("%s.zip", sessionKey)),
// 				Body:   f,
// 			}
// 			if _, err := s.S3.PutObject(in); err != nil {
// 				return err
// 			}

// 			if err = server.SendAndClose(ZipUploaded); err != nil {
// 				return err
// 			}

// 			return nil
// 		}

// 		if err != nil {
// 			server.SendAndClose(ZipUploadFailed)
// 			return err
// 		}
// 	}
// }

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

	if n == int(*out.ContentLength) {
		res := &terrarium.SourceZipResponse{}
		for i := 0; i < int(*out.ContentLength); i += ChunkSize {
			if i+ChunkSize > int(*out.ContentLength) {
				res.ZipDataChunk = buf[i:int(*out.ContentLength)]
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
