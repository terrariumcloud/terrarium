package cmd

import (
	"fmt"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/services"
	"log"
	"net"

	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/services/storage"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var storageServiceCmd = &cobra.Command{
	Use:   "storage-service",
	Short: "Starts the Terrarium GRPC Storage Service",
	Long:  `Runs the GRPC Storage Service server.`,
	Run:   runStorageService,
}

func init() {
	rootCmd.AddCommand(storageServiceCmd)
	storageServiceCmd.PersistentFlags().StringVarP(&storage.BucketName, "bucket", "", storage.DefaultBucketName, "Bucket name")
	storageServiceCmd.PersistentFlags().StringVarP(&awsAccessKey, "aws-access-key", "", "", "AWS Access Key")
	storageServiceCmd.PersistentFlags().StringVarP(&awsSecretKey, "aws-secret-key", "", "", "AWS Secret Key")
	storageServiceCmd.PersistentFlags().StringVarP(&awsRegion, "aws-region", "", "", "AWS Region")
}

func runStorageService(cmd *cobra.Command, args []string) {
	log.Println("Starting Terrarium GRPC Storage Service")

	a := fmt.Sprintf("%s:%s", address, port)
	listener, err := net.Listen("tcp", a)
	if err != nil {
		log.Fatalf("Failed to start Terrarium GRPC Storage Service: %s", err.Error())
	}

	var opts []grpc.ServerOption
	s3 := createS3Client()
	storageServiceServer := &storage.StorageService{
		S3: s3,
	}

	grpcServer := grpc.NewServer(opts...)
	services.RegisterStorageServer(grpcServer, storageServiceServer)

	log.Printf("Listening at %s", a)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed: %s", err.Error())
	}
}

func createS3Client() s3iface.S3API {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
	})
	if err != nil {
		log.Fatalf("Unable to create AWS Session: %s", err.Error())
	}

	// sess = session.Must(session.NewSessionWithOptions(session.Options{
	// 	SharedConfigState: session.SharedConfigEnable,
	// }))

	return s3.New(sess)

}
