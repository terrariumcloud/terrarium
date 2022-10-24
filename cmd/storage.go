package cmd

import (
	"github.com/terrariumcloud/terrarium/internal/module/services"
	"github.com/terrariumcloud/terrarium/internal/storage"

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

	storageServiceServer := &services.StorageService{
		S3:         storage.NewS3Client(awsAccessKey, awsSecretKey, awsRegion),
		BucketName: services.BucketName,
		Region:     awsRegion,
	}

	startGRPCService("Terrarium GRPC Storage service", storageServiceServer)
}
