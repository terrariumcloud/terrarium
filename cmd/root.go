package cmd

import (
	"fmt"
	"log"
	"net"

	services "github.com/terrariumcloud/terrarium-grpc-gateway/internal/module/services"
	terrarium "github.com/terrariumcloud/terrarium-grpc-gateway/pkg/terrarium/module"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	defaultAddress = "0.0.0.0"
	defaultPort    = "3001"
)

var address string = defaultAddress
var port string = defaultPort
var awsAccessKey string
var awsSecretKey string
var awsRegion string

var rootCmd = &cobra.Command{
	Use:   "terrarium",
	Short: "Terrarium Services",
	Long:  "Runs GRPC server that exposes Terrarium Services",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&address, "address", "a", defaultAddress, "IP Address")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", defaultPort, "Port number")
	rootCmd.PersistentFlags().StringVarP(&awsAccessKey, "aws-access-key", "k", "", "AWS Access Key (required)")
	rootCmd.MarkPersistentFlagRequired("aws-access-key")
	rootCmd.PersistentFlags().StringVarP(&awsSecretKey, "aws-secret-key", "s", "", "AWS Secret Key (required)")
	rootCmd.MarkPersistentFlagRequired("aws-secret-key")
	rootCmd.PersistentFlags().StringVarP(&awsRegion, "aws-region", "r", "", "AWS Region (required)")
	rootCmd.MarkPersistentFlagRequired("aws-region")
}

func startService(name string, service interface{}) {
	log.Printf("Starting %s", name)

	endpoint := fmt.Sprintf("%s:%s", address, port)
	listener, err := net.Listen("tcp4", endpoint)
	if err != nil {
		log.Fatalf("Failed to start: %s", err.Error())
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	register(grpcServer, service)

	log.Printf("Listening at %s", endpoint)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed: %s", err.Error())
	}
}

func register(grpcServer grpc.ServiceRegistrar, service interface{}) {
	switch t := service.(type) {
	case services.RegistrarServer:
		services.RegisterRegistrarServer(grpcServer, service.(*services.RegistrarService))
	case services.VersionManagerServer:
		services.RegisterVersionManagerServer(grpcServer, service.(*services.VersionManagerService))
	case services.DependencyResolverServer:
		services.RegisterDependencyResolverServer(grpcServer, service.(*services.DependencyResolverService))
	case services.StorageServer:
		services.RegisterStorageServer(grpcServer, service.(*services.StorageService))
	case terrarium.PublisherServer:
		terrarium.RegisterPublisherServer(grpcServer, service.(*services.TerrariumGrpcGateway))
		terrarium.RegisterConsumerServer(grpcServer, service.(*services.TerrariumGrpcGateway))
		// TODO: fallthrough doesn't seem to work with type switching
		// fallthrough
	// case terrarium.ConsumerServer:
	// 	terrarium.RegisterConsumerServer(grpcServer, service.(*services.TerrariumGrpcGateway))
	default:
		log.Fatalf("Failed to register unknown service type: %v", t)
	}
}

// Execute root command
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
