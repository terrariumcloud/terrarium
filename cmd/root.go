package cmd

import (
	"fmt"
	"log"
	"net"

	services "github.com/terrariumcloud/terrarium-grpc-gateway/internal/module/services"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/storage"
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
		log.Fatalf("Failed to start: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	if err := register(grpcServer, service); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}

	// storage.InitialiseDynamoDb(tableName, schema, db)

	log.Printf("Listening at %s", endpoint)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed: %v", err)
	}
}

func register(grpcServer grpc.ServiceRegistrar, service interface{}) error {
	switch t := service.(type) {
	case services.RegistrarServer:
		r := service.(*services.RegistrarService)
		services.RegisterRegistrarServer(grpcServer, r)
		if err := storage.InitializeDynamoDb(r.Table, r.Schema, r.Db); err != nil {
			return err
		}
	case services.VersionManagerServer:
		vms := service.(*services.VersionManagerService)
		services.RegisterVersionManagerServer(grpcServer, vms)
		if err := storage.InitializeDynamoDb(vms.Table, vms.Schema, vms.Db); err != nil {
			return err
		}
	case services.DependencyResolverServer:
		services.RegisterDependencyResolverServer(grpcServer, service.(*services.DependencyResolverService))
	case services.StorageServer:
		s := service.(*services.StorageService)
		services.RegisterStorageServer(grpcServer, s)
		if err := storage.InitializeS3Bucket(s.BucketName, s.Region, s.S3); err != nil {
			return err
		}
	case terrarium.PublisherServer:
		terrarium.RegisterPublisherServer(grpcServer, service.(*services.TerrariumGrpcGateway))
		terrarium.RegisterConsumerServer(grpcServer, service.(*services.TerrariumGrpcGateway))
		// TODO: fallthrough doesn't seem to work with type switching
		// fallthrough
	// case terrarium.ConsumerServer:
	// 	terrarium.RegisterConsumerServer(grpcServer, service.(*services.TerrariumGrpcGateway))
	default:
		return fmt.Errorf("failed to register unknown service type: %v", t)
	}
	return nil
}

// Execute root command
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
