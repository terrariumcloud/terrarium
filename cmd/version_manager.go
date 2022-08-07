package cmd

import (
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/module/services"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/storage"

	"github.com/spf13/cobra"
)

var versionManagerCmd = &cobra.Command{
	Use:   "version-manager",
	Short: "Starts the Terrarium GRPC Version Manager service",
	Long:  "Runs the Terrarium GRPC Version Manager server.",
	Run:   runVersionManager,
}

func init() {
	rootCmd.AddCommand(versionManagerCmd)
	versionManagerCmd.Flags().StringVarP(&services.VersionsTableName, "table", "t", services.DefaultVersionsTableName, "Module versions table name")
}

func runVersionManager(cmd *cobra.Command, args []string) {

	versionManagerServer := &services.VersionManagerService{
		Db: storage.NewDynamoDbClient(awsAccessKey, awsSecretKey, awsRegion),
	}

	startService("Terrarium GRPC Version Manager service", versionManagerServer)
}
