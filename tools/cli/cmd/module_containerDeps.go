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
			panic(err)
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
			panic(err)
		}
		for {
			response, err := responseClient.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			fmt.Printf("%v\n", response.Dependencies)
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
