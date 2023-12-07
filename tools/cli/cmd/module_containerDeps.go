/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"io"
)

// containerDepsCmd represents the containerDeps command
var containerDepsCmd = &cobra.Command{
	Use:   "container-deps",
	Short: "List the container dependencies for a module",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		conn, client, err := getModuleConsumerClient()
		if err != nil {
			printErrorAndExit("Failed to connect to terrarium", err, 1)
		}
		defer func() { _ = conn.Close() }()
		req := module.RetrieveContainerDependenciesRequestV2{
			Module: &module.Module{
				Name:    args[0],
				Version: args[1],
			},
		}
		responseClient, err := client.RetrieveContainerDependenciesV2(context.Background(), &req)
		if err != nil {
			printErrorAndExit("Failed to retrieve container dependencies", err, 1)
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
			for name, containerDetails := range response.Dependencies {
				fmt.Printf("    %s/%s:%s:\n", containerDetails.Namespace, name, containerDetails.Tag)

				for _, details := range containerDetails.Images {
					fmt.Printf("        - arch: %s\n", details.Arch)
					fmt.Printf("          image: %s\n", details.Image)
				}
			}
		}

	},
}

func init() {
	moduleCmd.AddCommand(containerDepsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// containerDepsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// containerDepsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
