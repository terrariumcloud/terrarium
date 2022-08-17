package cmd

import (
	services "github.com/terrariumcloud/terrarium-grpc-gateway/internal/module/services"

	"github.com/spf13/cobra"
)

var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "Starts the Terrarium GRPC Gateway service",
	Long:  "Runs the Terrarium GRPC Gateway server.",
	Run:   runGateway,
}

func init() {
	rootCmd.AddCommand(gatewayCmd)
	gatewayCmd.Flags().StringVarP(&services.RegistrarServiceEndpoint, "registrar", "", services.DefaultRegistrarServiceDefaultEndpoint, "GRPC Endpoint for Registrar Service")
	gatewayCmd.Flags().StringVarP(&services.DependencyManagerEndpoint, "dependency-manager", "", services.DefaultDependencyManagerEndpoint, "GRPC Endpoint for Dependency Resolver Service")
	gatewayCmd.Flags().StringVarP(&services.VersionManagerEndpoint, "version-manager", "", services.DefaultVersionManagerEndpoint, "GRPC Endpoint for Version Manager Service")
	gatewayCmd.Flags().StringVarP(&services.StorageServiceEndpoint, "storage", "", services.DefaultStorageServiceDefaultEndpoint, "GRPC Endpoint for Storage Service")
}

func runGateway(cmd *cobra.Command, args []string) {

	gatewayServer := &services.TerrariumGrpcGateway{}

	startService("Terrarium GRPC Gateway service", gatewayServer)
}
