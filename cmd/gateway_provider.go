package cmd

import (
	"github.com/terrariumcloud/terrarium/internal/provider/services/gateway"
	"github.com/terrariumcloud/terrarium/internal/provider/services/version_manager"

	"github.com/spf13/cobra"
)

var gatewayProviderCmd = &cobra.Command{
	Use:   "gateway-provider",
	Short: "Starts the Terrarium GRPC Gateway service",
	Long:  "Runs the Terrarium GRPC Gateway server.",
	Run:   runProviderGateway,
}

func init() {
	rootCmd.AddCommand(gatewayProviderCmd)
	gatewayCmd.Flags().StringVarP(&version_manager.VersionManagerEndpoint, "provider-version-manager", "", version_manager.DefaultProviderVersionManagerEndpoint, "GRPC Endpoint for Module Version Manager Service")
}

func runProviderGateway(cmd *cobra.Command, args []string) {

	gatewayServer := gateway.New(
		version_manager.NewVersionManagerGrpcClient(version_manager.VersionManagerEndpoint),
	)

	startProviderGRPCService("api-gateway-provider", gatewayServer)
}