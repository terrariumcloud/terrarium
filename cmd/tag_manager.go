package cmd

import (
	"github.com/terrariumcloud/terrarium/internal/module/services"
	"github.com/terrariumcloud/terrarium/internal/storage"

	"github.com/spf13/cobra"
)

var tagManagerCmd = &cobra.Command{
	Use:   "tag-manager",
	Short: "Starts the Terrarium GRPC Tag Manager service",
	Long:  "Runs the Terrarium GRPC Tag Manager server.",
	Run:   runTagManager,
}

func init() {
	rootCmd.AddCommand(tagManagerCmd)
	tagManagerCmd.Flags().StringVarP(&services.TagTableName, "table", "t", services.DefaultTagTableName, "Module tags table name")
}

func runTagManager(cmd *cobra.Command, args []string) {

	tagManagerServer := &services.TagManagerService{
		Db:     storage.NewDynamoDbClient(awsAccessKey, awsSecretKey, awsRegion),
		Table:  services.TagTableName,
		Schema: services.GetTagsSchema(services.TagTableName),
	}

	startGRPCService("Terrarium GRPCTtag Manager service", tagManagerServer)
}
