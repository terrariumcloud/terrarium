package cmd

import (
	"github.com/terrariumcloud/terrarium/internal/module/services/storage"
	"github.com/terrariumcloud/terrarium/internal/module/services/version_manager"
	modulesv1 "github.com/terrariumcloud/terrarium/internal/restapi/modules/v1"

	"github.com/spf13/cobra"
)

var mountPath string

var modulesV1Cmd = &cobra.Command{
	Use:   "modules.v1",
	Short: "Starts the Terrarium REST API service implementing a read only version of the module.v1 registry protocol",
	Long:  "Runs the Terrarium REST server for the implementation of the module.v1 protocol",
	Run:   runRESTModulesV1Server,
}

func init() {
	modulesV1Cmd.Flags().StringVarP(
		&mountPath,
		"mount-path",
		"m",
		"modules",
		"Mount path for the rest API server used to process request relative to a particular URL in a reverse proxy type setup",
	)
	modulesV1Cmd.Flags().StringVarP(&version_manager.VersionManagerEndpoint, "version-manager", "", version_manager.DefaultVersionManagerEndpoint, "GRPC Endpoint for Version Manager Service")
	modulesV1Cmd.Flags().StringVarP(&storage.StorageServiceEndpoint, "storage", "", storage.DefaultStorageServiceDefaultEndpoint, "GRPC Endpoint for Storage Service")
	rootCmd.AddCommand(modulesV1Cmd)
}

func runRESTModulesV1Server(cmd *cobra.Command, args []string) {

	restAPIServer := modulesv1.New()
	startRESTAPIService("rest-modules-v1", mountPath, restAPIServer)
}
