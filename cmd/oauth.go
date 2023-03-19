/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/terrariumcloud/terrarium/internal/oauth/services/authorization"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/jwt"
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
	err := jwt.CreatePKI()
	if err != nil {
		log.Fatalf("failed creating PKI keys: %s", err)
	}
	token := jwt.NewJWT([]string{})
	issued, err := token.Sign()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(issued)
	err = token.Signature().Verify()
	if err != nil {
		log.Fatal(err)
	}
	startGRPCService("oauth", &authServer)
}
