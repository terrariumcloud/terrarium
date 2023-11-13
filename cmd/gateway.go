package cmd

import (
	"github.com/terrariumcloud/terrarium/internal/module/services/dependency_manager"
	"github.com/terrariumcloud/terrarium/internal/module/services/gateway"
	"github.com/terrariumcloud/terrarium/internal/module/services/registrar"
	"github.com/terrariumcloud/terrarium/internal/module/services/storage"
	"github.com/terrariumcloud/terrarium/internal/module/services/tag_manager"
	"github.com/terrariumcloud/terrarium/internal/module/services/version_manager"
	"github.com/terrariumcloud/terrarium/internal/release/services/release"

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
	gatewayCmd.Flags().StringVarP(&registrar.RegistrarServiceEndpoint, "registrar", "", registrar.DefaultRegistrarServiceEndpoint, "GRPC Endpoint for Registrar Service")
	gatewayCmd.Flags().StringVarP(&dependency_manager.DependencyManagerEndpoint, "dependency-manager", "", dependency_manager.DefaultDependencyManagerEndpoint, "GRPC Endpoint for Dependency Manager Service")
	gatewayCmd.Flags().StringVarP(&version_manager.VersionManagerEndpoint, "version-manager", "", version_manager.DefaultVersionManagerEndpoint, "GRPC Endpoint for Version Manager Service")
	gatewayCmd.Flags().StringVarP(&storage.StorageServiceEndpoint, "storage", "", storage.DefaultStorageServiceDefaultEndpoint, "GRPC Endpoint for Storage Service")
	gatewayCmd.Flags().StringVarP(&tag_manager.TagManagerEndpoint, "tag-manager", "", tag_manager.DefaultTagManagerEndpoint, "GRPC Endpoint for Tag Service")
	gatewayCmd.Flags().StringVarP(&release.ReleaseServiceEndpoint, "tag", "", release.DefaultReleaseServiceEndpoint, "GRPC Endpoint for Release Service")
}

func runGateway(cmd *cobra.Command, args []string) {

	gatewayServer := &gateway.TerrariumGrpcGateway{}

	startGRPCService("api-gateway", gatewayServer)
}
