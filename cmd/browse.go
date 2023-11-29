package cmd

import (
	"github.com/spf13/cobra"
	"github.com/terrariumcloud/terrarium/internal/module/services/registrar"
	"github.com/terrariumcloud/terrarium/internal/module/services/version_manager"
	"github.com/terrariumcloud/terrarium/internal/release/services/release"
	"github.com/terrariumcloud/terrarium/internal/restapi/browse"
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Starts the Terrarium service that provides the web UI and its backing API",
	Long:  "Runs the Terrarium REST server for the implementation of webui server and backing API",
	Run:   runBrowseServer,
}

func init() {
	browseCmd.Flags().StringVarP(&registrar.RegistrarServiceEndpoint, "registrar", "", registrar.DefaultRegistrarServiceEndpoint, "GRPC Endpoint for Registrar Service")
	browseCmd.Flags().StringVarP(&version_manager.VersionManagerEndpoint, "version-manager", "", version_manager.DefaultVersionManagerEndpoint, "GRPC Endpoint for Version Manager Service")
	browseCmd.Flags().StringVarP(&release.ReleaseServiceEndpoint, "release", "", release.DefaultReleaseServiceEndpoint, "GRPC Endpoint for Release Service")
	rootCmd.AddCommand(browseCmd)
}

func runBrowseServer(cmd *cobra.Command, args []string) {

	restAPIServer := browse.New(registrar.NewRegistrarGrpcClient(registrar.RegistrarServiceEndpoint),
		version_manager.NewVersionManagerGrpcClient(version_manager.VersionManagerEndpoint),
		release.NewBrowseGrpcClient(release.ReleaseServiceEndpoint))
	startRESTAPIService("browse", "", restAPIServer)
}
