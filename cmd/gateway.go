package cmd

import (
	"github.com/terrariumcloud/terrarium/internal/common/gateway"
	"github.com/terrariumcloud/terrarium/internal/module/services/dependency_manager"
	"github.com/terrariumcloud/terrarium/internal/module/services/registrar"
	"github.com/terrariumcloud/terrarium/internal/module/services/storage"
	"github.com/terrariumcloud/terrarium/internal/module/services/tag_manager"
	"github.com/terrariumcloud/terrarium/internal/module/services/version_manager"
	providerStorage "github.com/terrariumcloud/terrarium/internal/provider/services/storage"
	providerVersionManager "github.com/terrariumcloud/terrarium/internal/provider/services/version_manager"
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
	gatewayCmd.Flags().StringVarP(&version_manager.VersionManagerEndpoint, "version-manager", "", version_manager.DefaultVersionManagerEndpoint, "GRPC Endpoint for Module Version Manager Service")
	gatewayCmd.Flags().StringVarP(&storage.StorageServiceEndpoint, "storage", "", storage.DefaultStorageServiceDefaultEndpoint, "GRPC Endpoint for Module Storage Service")
	gatewayCmd.Flags().StringVarP(&tag_manager.TagManagerEndpoint, "tag-manager", "", tag_manager.DefaultTagManagerEndpoint, "GRPC Endpoint for Tag Service")
	gatewayCmd.Flags().StringVarP(&release.ReleaseServiceEndpoint, "release", "", release.DefaultReleaseServiceEndpoint, "GRPC Endpoint for Release Service")
	gatewayCmd.Flags().StringVarP(&providerVersionManager.VersionManagerEndpoint, "provider-version-manager", "", providerVersionManager.DefaultProviderVersionManagerEndpoint, "GRPC Endpoint for Provider Version Manager Service")
	gatewayCmd.Flags().StringVarP(&providerStorage.StorageServiceEndpoint, "provider-storage", "", providerStorage.DefaultStorageServiceDefaultEndpoint, "GRPC Endpoint for Provider Storage Service")
}

func runGateway(cmd *cobra.Command, args []string) {

	gatewayServer := gateway.New(registrar.NewRegistrarGrpcClient(registrar.RegistrarServiceEndpoint),
		tag_manager.NewTagManagerGrpcClient(tag_manager.TagManagerEndpoint),
		version_manager.NewVersionManagerGrpcClient(version_manager.VersionManagerEndpoint),
		storage.NewStorageGrpcClient(storage.StorageServiceEndpoint),
		dependency_manager.NewDependencyManagerGrpcClient(dependency_manager.DependencyManagerEndpoint),
		release.NewPublisherGrpcClient(release.ReleaseServiceEndpoint),
		providerVersionManager.NewVersionManagerGrpcClient(providerVersionManager.VersionManagerEndpoint),
		providerStorage.NewStorageGrpcClient(providerStorage.StorageServiceEndpoint),
	)

	startGRPCService("api-gateway", gatewayServer)
}
