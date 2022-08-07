package cmd

import (
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/module/services"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/storage"

	"github.com/spf13/cobra"
)

var dependencyResolverCmd = &cobra.Command{
	Use:   "dependency-resolver",
	Short: "Starts the Terrarium GRPC Dependency Resolver service",
	Long:  "Runs the Terrarium GRPC Dependency Resolver server.",
	Run:   runDependencyResolver,
}

func init() {
	rootCmd.AddCommand(dependencyResolverCmd)
	dependencyResolverCmd.Flags().StringVarP(&services.ModuleDependenciesTableName, "table", "t", services.DefaultModuleDependenciesTableName, "Module dependencies table name")
}

func runDependencyResolver(cmd *cobra.Command, args []string) {

	dependencyServiceServer := &services.DependencyResolverService{
		Db: storage.NewDynamoDbClient(awsAccessKey, awsSecretKey, awsRegion),
	}

	startService("Terrarium GRPC Dependency Resolver service", dependencyServiceServer)
}
