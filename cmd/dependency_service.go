package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/services"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/services/dependency"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var dependencyServiceCmd = &cobra.Command{
	Use:   "dependency-service",
	Short: "Starts the Terrarium GRPC Dependency Service",
	Long:  `Runs the GRPC Dependency Service server.`,
	Run:   runDependencyService,
}

func init() {
	rootCmd.AddCommand(dependencyServiceCmd)
	dependencyServiceCmd.Flags().StringVarP(&dependency.ModuleDependenciesTableName, "module-table", "", dependency.DefaultModuleDependenciesTableName, "Module dependencies table name")
	dependencyServiceCmd.Flags().StringVarP(&dependency.ContainerDependenciesTableName, "container-table", "", dependency.DefaultContainerDependenciesTableName, "Container dependencies table name")
}

func runDependencyService(cmd *cobra.Command, args []string) {
	log.Println("Starting Terrarium GRPC Dependency Service")

	a := fmt.Sprintf("%s:%s", address, port)
	listener, err := net.Listen("tcp", a)
	if err != nil {
		log.Fatalf("Failed to start Terrarium GRPC Dependency Service: %s", err.Error())
	}

	var opts []grpc.ServerOption
	dynamodb := createDynamoDbClient()
	dependencyServiceServer := &dependency.DependencyService{
		Db: dynamodb,
	}

	grpcServer := grpc.NewServer(opts...)
	services.RegisterDependencyResolverServer(grpcServer, dependencyServiceServer)

	log.Printf("Listening at %s", a)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed: %s", err.Error())
	}
}
