/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/terrariumcloud/terrarium/internal/oauth/services/authorization"
)

// oauthCmd represents the oauth command
var oauthCmd = &cobra.Command{
	Use:   "oauth",
	Short: "Starts the OAuth GRPC service",
	Long:  `The OAuth GRPC service implements the OAuth 2.0 specification to secure the Terrarium REST APIs`,
	Run:   runOAuth,
}

func init() {
	rootCmd.AddCommand(oauthCmd)
}

func runOAuth(cmd *cobra.Command, args []string) {
	authServer := authorization.AuthorizationServer{}
	err := authServer.CreatePKI()
	if err != nil {
		log.Fatalf("failed creating PKI keys: %s", err)
	}
	startGRPCService("oauth", &authServer)
}
