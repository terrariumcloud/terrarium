package cmd

import (
	"github.com/spf13/cobra"
)

const (
	defaultAddress = "0.0.0.0"
	defaultPort    = "3001"
)

var address string = defaultAddress
var port string = defaultPort
var awsAccessKey string
var awsSecretKey string
var awsRegion string

var rootCmd = &cobra.Command{
	Use:   "terrarium",
	Short: "Terrarium Services",
	Long:  "Runs GRPC server that exposes Terrarium Services",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&address, "address", "a", defaultAddress, "IP Address")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", defaultPort, "Port number")
	rootCmd.PersistentFlags().StringVarP(&awsAccessKey, "aws-access-key", "", "", "AWS Access Key (required)")
	rootCmd.MarkPersistentFlagRequired("aws-access-key")
	rootCmd.PersistentFlags().StringVarP(&awsSecretKey, "aws-secret-key", "", "", "AWS Secret Key (required)")
	rootCmd.MarkPersistentFlagRequired("aws-secret-key")
	rootCmd.PersistentFlags().StringVarP(&awsRegion, "aws-region", "", "", "AWS Region (required)")
	rootCmd.MarkPersistentFlagRequired("aws-region")
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
