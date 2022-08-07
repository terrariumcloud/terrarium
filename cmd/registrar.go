package cmd

import (
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/module/services"
	"github.com/terrariumcloud/terrarium-grpc-gateway/internal/storage"

	"github.com/spf13/cobra"
)

var registrarServiceCmd = &cobra.Command{
	Use:   "registrar",
	Short: "Starts the Terrarium GRPC Registrar service",
	Long:  "Runs the Terrarium GRPC Registrar server.",
	Run:   runRegistrarService,
}

func init() {
	rootCmd.AddCommand(registrarServiceCmd)
	registrarServiceCmd.Flags().StringVarP(&services.RegistrarTableName, "table", "t", services.DefaultRegistrarTableName, "Module Registrar table name")
}

func runRegistrarService(cmd *cobra.Command, args []string) {

	registrarServiceServer := &services.RegistrarService{
		Db: storage.NewDynamoDbClient(awsAccessKey, awsSecretKey, awsRegion),
	}

	startService("Terrarium GRPC Registrar service", registrarServiceServer)
}
