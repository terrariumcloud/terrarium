package cmd

import (
	"github.com/terrariumcloud/terrarium/internal/provider/services/version_manager"
	providersv1 "github.com/terrariumcloud/terrarium/internal/restapi/providers/v1"

	"github.com/spf13/cobra"
)

var mountPathProviders string

var providersV1Cmd = &cobra.Command{
	Use:   "providers.v1",
	Short: "Starts the Terrarium REST API service implementing a read only version of the provider.v1 registry protocol",
	Long:  "Runs the Terrarium REST server for the implementation of the provider.v1 protocol",
	Run:   runRESTProvidersV1Server,
}

func init() {
	providersV1Cmd.Flags().StringVarP(
		&mountPathProviders,
		"mount-path",
		"m",
		"providers",
		"Mount path for the rest API server used to process request relative to a particular URL in a reverse proxy type setup",
	)
	providersV1Cmd.Flags().StringVarP(&version_manager.VersionManagerEndpoint, "provider-version-manager", "", version_manager.DefaultProviderVersionManagerEndpoint, "GRPC Endpoint for Version Manager Service")
	rootCmd.AddCommand(providersV1Cmd)
}

func runRESTProvidersV1Server(cmd *cobra.Command, args []string) {
	restAPIServer := providersv1.New(version_manager.NewVersionManagerGrpcClient(version_manager.VersionManagerEndpoint))
	startRESTAPIService("rest-providers-v1", mountPathProviders, restAPIServer)
}
