package cmd

import (
	"fmt"
	"log"
	"net"
	"terrarium-grpc-gateway/internal/services"

	"terrarium-grpc-gateway/internal/services/session"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var sessionServiceCmd = &cobra.Command{
	Use:   "session-service",
	Short: "Starts the Terrarium GRPC Session Service",
	Long:  `Runs the GRPC Session Service server.`,
	Run:   runSessionService,
}

func init() {
	rootCmd.AddCommand(sessionServiceCmd)
	sessionServiceCmd.PersistentFlags().StringVarP(&session.SessionTableName, "session-table", "", session.DefaultSessionTableName, "Session table name")
	sessionServiceCmd.PersistentFlags().StringVarP(&session.PublishedModulesTableName, "published-table", "", session.DefaultPublishedModulesTableName, "Published Module table name")
	sessionServiceCmd.PersistentFlags().StringVarP(&awsAccessKey, "aws-access-key", "", "", "AWS Access Key")
	sessionServiceCmd.PersistentFlags().StringVarP(&awsSecretKey, "aws-secret-key", "", "", "AWS Secret Key")
	sessionServiceCmd.PersistentFlags().StringVarP(&awsRegion, "aws-region", "", "", "AWS Region")
}

func runSessionService(cmd *cobra.Command, args []string) {
	log.Println("Starting Terrarium GRPC Session Service")

	a := fmt.Sprintf("%s:%s", address, port)
	listener, err := net.Listen("tcp", a)
	if err != nil {
		log.Fatalf("Failed to start Terrarium GRPC Session Service: %s", err.Error())
	}

	var opts []grpc.ServerOption
	dynamodb := createDynamoDbClient()
	sessionServiceServer := &session.SessionService{
		Db: dynamodb,
	}

	grpcServer := grpc.NewServer(opts...)
	services.RegisterSessionManagerServer(grpcServer, sessionServiceServer)

	log.Printf("Listening at %s", a)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed: %s", err.Error())
	}
}
