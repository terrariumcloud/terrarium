package cmd

import (
	modulesv1 "github.com/terrariumcloud/terrarium-grpc-gateway/internal/restapi/modules/v1"

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
	rootCmd.AddCommand(modulesV1Cmd)
}

func runRESTModulesV1Server(cmd *cobra.Command, args []string) {

	restAPIServer := modulesv1.New()
	startRESTAPIService("Terrarium REST API Server for modules.v1 protocol", mountPath, restAPIServer)
}
