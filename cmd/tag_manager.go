package cmd

import (
	"github.com/terrariumcloud/terrarium/internal/module/services/tag_manager"
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
	tagManagerCmd.Flags().StringVarP(&tag_manager.TagTableName, "table", "t", tag_manager.DefaultTagTableName, "Module tags table name")
}

func runTagManager(cmd *cobra.Command, args []string) {

	tagManagerServer := &tag_manager.TagManagerService{
		Db:     storage.NewDynamoDbClient(awsSessionConfig),
		Table:  tag_manager.TagTableName,
		Schema: tag_manager.GetTagsSchema(tag_manager.TagTableName),
	}

	startGRPCService("tag-manager", tagManagerServer)
}
