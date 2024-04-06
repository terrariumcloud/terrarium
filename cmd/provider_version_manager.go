package cmd

import (
	"github.com/terrariumcloud/terrarium/internal/provider/services/version_manager"
	"github.com/terrariumcloud/terrarium/internal/storage"

	"github.com/spf13/cobra"
)

var providerVersionManagerServiceCmd = &cobra.Command{
	Use:   "provider-version-manager",
	Short: "Starts the Terrarium GRPC Provider Version Manager service",
	Long:  "Runs the Terrarium GRPC Provider Version Manager server.",
	Run:   runProviderVersionManagerService,
}

func init() {
	rootCmd.AddCommand(providerVersionManagerServiceCmd)
	providerVersionManagerServiceCmd.Flags().StringVarP(&version_manager.VersionsTableName, "table", "t", version_manager.DefaultProviderVersionsTableName, "Provider Version Manager table name")
}

func runProviderVersionManagerService(cmd *cobra.Command, args []string) {

	versionManagerServiceServer := &version_manager.VersionManagerService{
		Db:     storage.NewDynamoDbClient(awsSessionConfig),
		Table:  version_manager.VersionsTableName,
		Schema: version_manager.GetProviderVersionsSchema(version_manager.VersionsTableName),
	}

	startGRPCService("provider-version-manager", versionManagerServiceServer)
}
