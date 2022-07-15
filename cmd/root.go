package cmd

import (
	"github.com/spf13/cobra"
)

const (
	defaultAddress = "0.0.0.0"
	defaultPort    = "3001"
)

var address string
var port string
var awsAccessKey string
var awsSecretKey string
var awsRegion string

var rootCmd = &cobra.Command{
	Use:   "terrarium",
	Short: "Terrarium Services",
	Long:  `Runs GRPC server that exposes Terrarium Services`,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&address, "address", "a", defaultAddress, "IP Address")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", defaultPort, "Port number")
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
