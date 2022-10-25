package cmd

import (
	"github.com/spf13/cobra"
	"github.com/terrariumcloud/terrarium/internal/module/services"
	"github.com/terrariumcloud/terrarium/internal/restapi/browse"
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Starts the Terrarium service that provides the web UI and its backing API",
	Long:  "Runs the Terrarium REST server for the implementation of webui server and backing API",
	Run:   runBrowseServer,
}

func init() {
	//browseCmd.Flags().StringVarP(
	//	&mountPath,
	//	"mount-path",
	//	"m",
	//	"modules",
	//	"Mount path for the rest API server used to process request relative to a particular URL in a reverse proxy type setup",
	//)
	browseCmd.Flags().StringVarP(&services.RegistrarServiceEndpoint, "registrar", "", services.DefaultRegistrarServiceEndpoint, "GRPC Endpoint for Registrar Service")
	browseCmd.Flags().StringVarP(&services.VersionManagerEndpoint, "version-manager", "", services.DefaultVersionManagerEndpoint, "GRPC Endpoint for Version Manager Service")
	rootCmd.AddCommand(browseCmd)
}

func runBrowseServer(cmd *cobra.Command, args []string) {

	restAPIServer := browse.New()
	startRESTAPIService("Terrarium REST API Server for modules.v1 protocol", "", restAPIServer)
}
