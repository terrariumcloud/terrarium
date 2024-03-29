package cmd

import (
	"github.com/spf13/cobra"
	"github.com/terrariumcloud/terrarium/internal/module/services/version_manager"
	"github.com/terrariumcloud/terrarium/internal/release/services/release"
	"github.com/terrariumcloud/terrarium/internal/storage"
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
	versionManagerCmd.Flags().StringVarP(&release.ReleaseServiceEndpoint, "release", "", release.DefaultReleaseServiceEndpoint, "GRPC Endpoint for Release Service")
}

func runVersionManager(cmd *cobra.Command, args []string) {

	versionManagerServer := &version_manager.VersionManagerService{
		Db:             storage.NewDynamoDbClient(awsSessionConfig),
		Table:          version_manager.VersionsTableName,
		Schema:         version_manager.GetModuleVersionsSchema(version_manager.VersionsTableName),
		ReleaseService: release.NewPublisherGrpcClient(release.ReleaseServiceEndpoint),
	}

	startGRPCService("version-manager", versionManagerServer)
}
