package cmd

import (
	providerStorage "github.com/terrariumcloud/terrarium/internal/provider/services/storage"
	"github.com/terrariumcloud/terrarium/internal/storage"

	"github.com/spf13/cobra"
)

var providerStorageServiceCmd = &cobra.Command{
	Use:   "provider-storage",
	Short: "Starts the Terrarium GRPC Provider Storage service",
	Long:  "Runs the Terrarium GRPC Provider Storage server.",
	Run:   runProviderStorageService,
}

func init() {
	rootCmd.AddCommand(providerStorageServiceCmd)
	providerStorageServiceCmd.Flags().StringVarP(&providerStorage.BucketName, "bucket", "b", providerStorage.DefaultBucketName, "Provider bucket name")
}

func runProviderStorageService(cmd *cobra.Command, args []string) {

	storageServiceServer := &providerStorage.StorageService{
		Client:     storage.NewS3Client(awsSessionConfig),
		BucketName: providerStorage.BucketName,
		Region:     awsSessionConfig.Region,
	}

	startGRPCService("provider-storage-s3", storageServiceServer)
}
