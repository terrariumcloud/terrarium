package cmd

import (
	"fmt"
	"log"
	"net"

	services "github.com/terrariumcloud/terrarium-grpc-gateway/internal/module/services"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var registrarServiceCmd = &cobra.Command{
	Use:   "registrar-service",
	Short: "Starts the Terrarium GRPC Registrar Service",
	Long:  `Runs the GRPC Registrar Service server.`,
	Run:   runRegistrarService,
}

func init() {
	rootCmd.AddCommand(registrarServiceCmd)
	registrarServiceCmd.Flags().StringVarP(&services.RegistrarTableName, "table", "", services.DefaultRegistrarTableName, "Table name")
}

func runRegistrarService(cmd *cobra.Command, args []string) {
	log.Println("Starting Terrarium GRPC Registrar Service")

	endpoint := fmt.Sprintf("%s:%s", address, port)
	listener, err := net.Listen("tcp", endpoint)
	if err != nil {
		log.Fatalf("Failed to start Terrarium GRPC Registrar Service: %s", err.Error())
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	dynamodb := createDynamoDbClient()
	registrarServiceServer := &services.RegistrarService{
		Db: dynamodb,
	}

	services.RegisterRegistrarServer(grpcServer, registrarServiceServer)

	log.Printf("Listening at %s", endpoint)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed: %s", err.Error())
	}
}

func createDynamoDbClient() dynamodbiface.DynamoDBAPI {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
	})
	if err != nil {
		log.Fatalf("Unable to create AWS Session: %s", err.Error())
	}

	// sess = session.Must(session.NewSessionWithOptions(session.Options{
	// 	SharedConfigState: session.SharedConfigEnable,
	// }))

	return dynamodb.New(sess)
}
