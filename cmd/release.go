package cmd

import (
	"github.com/terrariumcloud/terrarium/internal/release/services/release"

	"github.com/spf13/cobra"
)

var releaseServiceCmd = &cobra.Command{
	Use:   "publish",
	Short: "Starts the Terrarium GRPC Release service",
	Long:  "Runs the Terrarium GRPC Release server.",
	Run:   runReleaseService,
}

func init() {
	rootCmd.AddCommand(releaseServiceCmd)
	// TODO: database
	//releaseServiceCmd.Flags().StringVarP(&release.ReleaseTableName, "table", "t", release.DefaultReleaseTableName, "Releases table name")
}

func runReleaseService(cmd *cobra.Command, args []string) {
	releaseServiceServer := &release.ReleaseService{
		// Db:     storage.NewDynamoDbClient(awsAccessKey, awsSecretKey, awsRegion),
		// Table:  release.ReleaseTableName,
		// Schema: release.GetModulesSchema(release.ReleaseTableName),
	}

	startGRPCService("release", releaseServiceServer)
}
