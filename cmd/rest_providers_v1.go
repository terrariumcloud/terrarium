package cmd

import (
	providersv1 "github.com/terrariumcloud/terrarium/internal/restapi/providers/v1"

	"github.com/terrariumcloud/terrarium/internal/provider/services"

	"github.com/spf13/cobra"
)

var mountPathProviders string

var providersV1Cmd = &cobra.Command{
	Use:   "providers.v1",
	Short: "Starts the Terrarium REST API service implementing a read only version of the provider.v1 registry protocol",
	Long:  "Runs the Terrarium REST server for the implementation of the provider.v1 protocol",
	Run:   runRESTProvidersV1Server,
}

func init() {
	providersV1Cmd.Flags().StringVarP(
		&mountPathProviders,
		"mount-path",
		"m",
		"providers",
		"Mount path for the rest API server used to process request relative to a particular URL in a reverse proxy type setup",
	)
	rootCmd.AddCommand(providersV1Cmd)
}

func runRESTProvidersV1Server(cmd *cobra.Command, args []string) {
	vm, err := services.NewJSONFileProviderVersionManager()
	if err != nil {
		panic(err)
	}
	restAPIServer := providersv1.New(vm)
	startRESTAPIService("rest-providers-v1", mountPathProviders, restAPIServer)
}
