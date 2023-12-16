package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"io"
)

// moduleDepsCmd represents the moduleDeps command
var moduleDepsCmd = &cobra.Command{
	Use:   "module-deps",
	Short: "List the module dependencies of the specified module.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		conn, client, err := getModuleConsumerClient()
		if err != nil {
			printErrorAndExit("Failed to connect to terrarium", err, 1)
		}
		defer func() { _ = conn.Close() }()
		req := module.RetrieveModuleDependenciesRequest{
			Module: &module.Module{
				Name:    args[0],
				Version: args[1],
			},
		}
		responseClient, err := client.RetrieveModuleDependencies(context.Background(), &req)
		if err != nil {
			printErrorAndExit("Failed to retrieve module dependencies", err, 1)
		}
		for {
			response, err := responseClient.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				printErrorAndExit("Retrieving dependencies failed", err, 1)
			}
			fmt.Printf("%s:%s:\n", response.Module.Name, response.Module.Version)
			for _, module := range response.Dependencies {
				fmt.Printf("    - %s:%s\n", module.Name, module.Version)
			}
		}
	},
}

func init() {
	moduleCmd.AddCommand(moduleDepsCmd)
}
