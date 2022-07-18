package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/services"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/services/creation"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var creationServiceCmd = &cobra.Command{
	Use:   "creation-service",
	Short: "Starts the Terrarium GRPC Creation Service",
	Long:  `Runs the GRPC Creation Service server.`,
	Run:   runCreationService,
}

func init() {
	rootCmd.AddCommand(creationServiceCmd)
	creationServiceCmd.Flags().StringVarP(&creation.TableName, "table", "", creation.DefaultTableName, "Table name")
}

func runCreationService(cmd *cobra.Command, args []string) {
	log.Println("Starting Terrarium GRPC Creation Service")

	a := fmt.Sprintf("%s:%s", address, port)
	listener, err := net.Listen("tcp", a)
	if err != nil {
		log.Fatalf("Failed to start Terrarium GRPC Creation Service: %s", err.Error())
	}

	var opts []grpc.ServerOption
	dynamodb := createDynamoDbClient()
	creationServiceServer := &creation.CreationService{
		Db: dynamodb,
	}

	grpcServer := grpc.NewServer(opts...)
	services.RegisterCreatorServer(grpcServer, creationServiceServer)

	log.Printf("Listening at %s", a)
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
