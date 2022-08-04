package cmd

import (
	"fmt"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/services"
	"log"
	"net"

	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/services/session"

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
	sessionServiceCmd.Flags().StringVarP(&session.SessionTableName, "session-table", "", session.DefaultSessionTableName, "Session table name")
	sessionServiceCmd.Flags().StringVarP(&session.PublishedModulesTableName, "published-table", "", session.DefaultPublishedModulesTableName, "Published Module table name")
}

func runSessionService(cmd *cobra.Command, args []string) {
	log.Println("Starting Terrarium GRPC Session Service")

	endpoint := fmt.Sprintf("%s:%s", address, port)
	listener, err := net.Listen("tcp", endpoint)
	if err != nil {
		log.Fatalf("Failed to start Terrarium GRPC Session Service: %s", err.Error())
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	dynamodb := createDynamoDbClient()
	sessionServiceServer := &session.SessionService{
		Db: dynamodb,
	}

	services.RegisterSessionManagerServer(grpcServer, sessionServiceServer)

	log.Printf("Listening at %s", endpoint)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed: %s", err.Error())
	}
}
