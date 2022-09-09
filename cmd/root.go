package cmd

import (
	"fmt"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/restapi"
	"log"
	"net"
	"net/http"

	services "github.com/terrariumcloud/terrarium-grpc-gateway/internal/module/services"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	defaultEndpoint = "0.0.0.0:3001"
)

var endpoint string = defaultEndpoint
var awsAccessKey string
var awsSecretKey string
var awsRegion string

var rootCmd = &cobra.Command{
	Use:   "terrarium",
	Short: "Terrarium Services",
	Long:  "Runs backend that exposes Terrarium Services",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", defaultEndpoint, "Endpoint")
	rootCmd.PersistentFlags().StringVarP(&awsAccessKey, "aws-access-key-id", "k", "", "AWS Access Key (required)")
	rootCmd.MarkPersistentFlagRequired("aws-access-key-id")
	rootCmd.PersistentFlags().StringVarP(&awsSecretKey, "aws-secret-access-key", "s", "", "AWS Secret Key (required)")
	rootCmd.MarkPersistentFlagRequired("aws-secret-access-key")
	rootCmd.PersistentFlags().StringVarP(&awsRegion, "aws-region", "r", "", "AWS Region (required)")
	rootCmd.MarkPersistentFlagRequired("aws-region")
}

func startGRPCService(name string, service services.Service) {
	log.Printf("Starting %s", name)

	listener, err := net.Listen("tcp4", endpoint)
	if err != nil {
		log.Fatalf("Failed to start: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	if err := service.RegisterWithServer(grpcServer); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}

	log.Printf("Listening at %s", endpoint)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed: %v", err)
	}
}

func startRESTAPIService(name, mountPath string, rootHandler restapi.RESTAPIHandler) {
	log.Printf("Starting %s", name)
	log.Println(fmt.Sprintf("Listening on %s", endpoint))
	log.Fatal(http.ListenAndServe(endpoint, rootHandler.GetHttpHandler(mountPath)))
}

// Execute root command
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
