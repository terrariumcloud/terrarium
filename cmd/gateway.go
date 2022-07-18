package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/services/creation"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/services/dependency"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/services/gateway"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/services/session"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/services/storage"

	"github.com/spf13/cobra"
	"github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium"
	"google.golang.org/grpc"
)

var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "Starts the Terrarium GRPC Gateway",
	Long:  `Runs the GRPC gateway server for accessing Terrarium services.`,
	Run:   runGateway,
}

func init() {
	rootCmd.AddCommand(gatewayCmd)
	gatewayCmd.Flags().StringVarP(&creation.CreationServiceEndpoint, "creation-service", "", creation.DefaultCreationServiceDefaultEndpoint, "GRPC Endpoint for Creation Service")
	gatewayCmd.Flags().StringVarP(&dependency.DependencyServiceEndpoint, "dependency-service", "", dependency.DefaultDependencyServiceDefaultEndpoint, "GRPC Endpoint for Dependency Service")
	gatewayCmd.Flags().StringVarP(&session.SessionServiceEndpoint, "session-service", "", session.DefaultSessionServiceDefaultEndpoint, "GRPC Endpoint for Session Service")
	gatewayCmd.Flags().StringVarP(&storage.StorageServiceEndpoint, "storage-service", "", storage.DefaultStorageServiceDefaultEndpoint, "GRPC Endpoint for Storage Service")
}

func runGateway(cmd *cobra.Command, args []string) {
	log.Println("Welcome to Terrarium GRPC Gateway")

	a := fmt.Sprintf("%s:%s", address, port)
	listener, err := net.Listen("tcp4", a)
	if err != nil {
		log.Fatalf("Failed to start Terrarium GRPC Gateway: %s", err.Error())
	}

	var opts []grpc.ServerOption
	gatewayServer := &gateway.TerrariumGrpcGateway{}

	grpcServer := grpc.NewServer(opts...)
	terrarium.RegisterPublisherServer(grpcServer, gatewayServer)
	terrarium.RegisterConsumerServer(grpcServer, gatewayServer)

	log.Printf("Listening at %s", a)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed: %s", err.Error())
	}
}
