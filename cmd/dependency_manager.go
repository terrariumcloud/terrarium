package cmd

import (
	"github.com/terrariumcloud/terrarium/internal/module/services/dependency_manager"
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
	dependencyManagerCmd.Flags().StringVarP(&dependency_manager.ModuleDependenciesTableName, "module-table", "m", dependency_manager.DefaultModuleDependenciesTableName, "Module dependencies table name")
	dependencyManagerCmd.Flags().StringVarP(&dependency_manager.ContainerDependenciesTableName, "container-table", "c", dependency_manager.DefaultContainerDependenciesTableName, "Module dependencies table name")
}

func runDependencyManager(cmd *cobra.Command, args []string) {

	dependencyServiceServer := &dependency_manager.DependencyManagerService{
		Db:              storage.NewDynamoDbClient(awsSessionConfig),
		ModuleTable:     dependency_manager.ModuleDependenciesTableName,
		ModuleSchema:    dependency_manager.GetDependenciesSchema(dependency_manager.ModuleDependenciesTableName),
		ContainerTable:  dependency_manager.ContainerDependenciesTableName,
		ContainerSchema: dependency_manager.GetDependenciesSchema(dependency_manager.ContainerDependenciesTableName),
	}

	startGRPCService("dependency-manager", dependencyServiceServer)
}
