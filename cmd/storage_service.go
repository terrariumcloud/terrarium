package cmd

import (
	"fmt"
	services "github.com/terrariumcloud/terrarium-grpc-gateway/internal/module/services"
	"log"
	"net"

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
	storageServiceCmd.Flags().StringVarP(&services.BucketName, "bucket", "", services.DefaultBucketName, "Bucket name")
}

func runStorageService(cmd *cobra.Command, args []string) {
	log.Println("Starting Terrarium GRPC Storage Service")

	endpoint := fmt.Sprintf("%s:%s", address, port)
	listener, err := net.Listen("tcp", endpoint)
	if err != nil {
		log.Fatalf("Failed to start Terrarium GRPC Storage Service: %s", err.Error())
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	s3 := createS3Client()
	storageServiceServer := &services.StorageService{
		S3: s3,
	}

	services.RegisterStorageServer(grpcServer, storageServiceServer)

	log.Printf("Listening at %s", endpoint)
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
