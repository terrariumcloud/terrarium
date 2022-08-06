package cmd

import (
	"fmt"
	"log"
	"net"

	services "github.com/terrariumcloud/terrarium-grpc-gateway/internal/module/services"

	"github.com/spf13/cobra"
	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"
	"google.golang.org/grpc"
)

var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "Starts the Terrarium GRPC Gateway",
	Long:  "Runs the GRPC gateway server for accessing Terrarium services.",
	Run:   runGateway,
}

func init() {
	rootCmd.AddCommand(gatewayCmd)
	gatewayCmd.Flags().StringVarP(&services.RegistrarServiceEndpoint, "registrar-service", "", services.DefaultRegistrarServiceDefaultEndpoint, "GRPC Endpoint for Creation Service")
	gatewayCmd.Flags().StringVarP(&services.DependencyServiceEndpoint, "dependency-service", "", services.DefaultDependencyServiceDefaultEndpoint, "GRPC Endpoint for Dependency Service")
	gatewayCmd.Flags().StringVarP(&services.SessionServiceEndpoint, "version-service", "", services.DefaultVersionServiceDefaultEndpoint, "GRPC Endpoint for Session Service")
	gatewayCmd.Flags().StringVarP(&services.StorageServiceEndpoint, "storage-service", "", services.DefaultStorageServiceDefaultEndpoint, "GRPC Endpoint for Storage Service")
}

func runGateway(cmd *cobra.Command, args []string) {
	log.Println("Welcome to Terrarium GRPC Gateway")

	endpoint := fmt.Sprintf("%s:%s", address, port)
	listener, err := net.Listen("tcp4", endpoint)
	if err != nil {
		log.Fatalf("Failed to start Terrarium GRPC Gateway: %s", err.Error())
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	gatewayServer := &services.TerrariumGrpcGateway{}

	terrarium.RegisterPublisherServer(grpcServer, gatewayServer)
	terrarium.RegisterConsumerServer(grpcServer, gatewayServer)

	log.Printf("Listening at %s", endpoint)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed: %s", err.Error())
	}
}
