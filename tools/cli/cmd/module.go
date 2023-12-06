/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/module"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// moduleCmd represents the module command
var moduleCmd = &cobra.Command{
	Use:   "module",
	Short: "Commands for managing modules",
}

func init() {
	rootCmd.AddCommand(moduleCmd)
}

func getModulePublisherClient() (*grpc.ClientConn, module.PublisherClient, error) {
	conn, err := grpc.Dial(terrariumEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	return conn, module.NewPublisherClient(conn), nil
}

func getModuleConsumerClient() (*grpc.ClientConn, module.ConsumerClient, error) {
	conn, err := grpc.Dial(terrariumEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	return conn, module.NewConsumerClient(conn), nil
}
