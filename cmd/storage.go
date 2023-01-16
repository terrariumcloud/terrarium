package cmd

import (
	storage2 "github.com/terrariumcloud/terrarium/internal/module/services/storage"
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
	storageServiceCmd.Flags().StringVarP(&storage2.BucketName, "bucket", "b", storage2.DefaultBucketName, "Module bucket name")
}

func runStorageService(cmd *cobra.Command, args []string) {

	storageServiceServer := &storage2.StorageService{
		Client:     storage.NewS3Client(awsAccessKey, awsSecretKey, awsRegion),
		BucketName: storage2.BucketName,
		Region:     awsRegion,
	}

	startGRPCService("storage-s3", storageServiceServer)
}
