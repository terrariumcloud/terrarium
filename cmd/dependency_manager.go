package cmd

import (
	"github.com/terrariumcloud/terrarium/internal/module/services"
	"github.com/terrariumcloud/terrarium/internal/storage"

	"github.com/spf13/cobra"
)

var dependencyManagerCmd = &cobra.Command{
	Use:   "dependency-manager",
	Short: "Starts the Terrarium GRPC Dependency Manager service",
	Long:  "Runs the Terrarium GRPC Dependency Manager server.",
	Run:   runDependencyManager,
}

func init() {
	rootCmd.AddCommand(dependencyManagerCmd)
	dependencyManagerCmd.Flags().StringVarP(&services.ModuleDependenciesTableName, "module-table", "m", services.DefaultModuleDependenciesTableName, "Module dependencies table name")
	dependencyManagerCmd.Flags().StringVarP(&services.ContainerDependenciesTableName, "container-table", "c", services.DefaultContainerDependenciesTableName, "Module dependencies table name")
}

func runDependencyManager(cmd *cobra.Command, args []string) {

	dependencyServiceServer := &services.DependencyManagerService{
		Db:              storage.NewDynamoDbClient(awsAccessKey, awsSecretKey, awsRegion),
		ModuleTable:     services.ModuleDependenciesTableName,
		ModuleSchema:    services.GetDependenciesSchema(services.ModuleDependenciesTableName),
		ContainerTable:  services.ContainerDependenciesTableName,
		ContainerSchema: services.GetDependenciesSchema(services.ContainerDependenciesTableName),
	}

	startGRPCService("Terrarium GRPC Dependency Resolver service", dependencyServiceServer)
}
