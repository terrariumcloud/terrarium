package cmd

import (
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/module/services"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/storage"

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
	dependencyManagerCmd.Flags().StringVarP(&services.ModuleDependenciesTableName, "table", "t", services.DefaultModuleDependenciesTableName, "Module dependencies table name")
}

func runDependencyManager(cmd *cobra.Command, args []string) {

	dependencyServiceServer := &services.DependencyManagerService{
		Db:     storage.NewDynamoDbClient(awsAccessKey, awsSecretKey, awsRegion),
		Table:  services.ModuleDependenciesTableName,
		Schema: services.GetModuleDependenciesSchema(services.ModuleDependenciesTableName),
	}

	startGRPCService("Terrarium GRPC Dependency Resolver service", dependencyServiceServer)
}
