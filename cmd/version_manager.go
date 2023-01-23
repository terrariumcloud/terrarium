package cmd

import (
	"github.com/terrariumcloud/terrarium/internal/module/services/version_manager"
	"github.com/terrariumcloud/terrarium/internal/storage"

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
	versionManagerCmd.Flags().StringVarP(&version_manager.VersionsTableName, "table", "t", version_manager.DefaultVersionsTableName, "Module versions table name")
}

func runVersionManager(cmd *cobra.Command, args []string) {

	versionManagerServer := &version_manager.VersionManagerService{
		Db:     storage.NewDynamoDbClient(awsAccessKey, awsSecretKey, awsRegion),
		Table:  version_manager.VersionsTableName,
		Schema: version_manager.GetModuleVersionsSchema(version_manager.VersionsTableName),
	}

	startGRPCService("version-manager", versionManagerServer)
}
