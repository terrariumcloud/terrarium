package cmd

import (
	"github.com/terrariumcloud/terrarium/internal/module/services/registrar"
	"github.com/terrariumcloud/terrarium/internal/storage"

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
	registrarServiceCmd.Flags().StringVarP(&registrar.RegistrarTableName, "table", "t", registrar.DefaultRegistrarTableName, "Module Registrar table name")
}

func runRegistrarService(cmd *cobra.Command, args []string) {

	registrarServiceServer := &registrar.RegistrarService{
		Db:     storage.NewDynamoDbClient(awsSessionConfig),
		Table:  registrar.RegistrarTableName,
		Schema: registrar.GetModulesSchema(registrar.RegistrarTableName),
	}

	startGRPCService("registrar", registrarServiceServer)
}
