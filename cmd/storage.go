package cmd

import (
	"fmt"

	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/module/services"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/storage"

	"github.com/spf13/cobra"
)

var storageServiceCmd = &cobra.Command{
	Use:   "storage",
	Short: "Starts the Terrarium GRPC Storage service",
	Long:  "Runs the Terrarium GRPC Storage server.",
	Run:   runStorageService,
}

func init() {
	rootCmd.AddCommand(storageServiceCmd)
	storageServiceCmd.Flags().StringVarP(&services.BucketName, "bucket", "b", services.DefaultBucketName, "Module bucket name")
}

func runStorageService(cmd *cobra.Command, args []string) {
	endpoint := fmt.Sprintf("%s:%s", address, port)

	storageServiceServer := &services.StorageService{
		S3: storage.NewS3Client(awsAccessKey, awsSecretKey, awsRegion),
	}

	startService("Terrarium GRPC Storage service", endpoint, storageServiceServer)
}
