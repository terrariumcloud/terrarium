package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var terrariumEndpoint = "localhost:3001"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&terrariumEndpoint, "endpoint", terrariumEndpoint, "GRPC Endpoint for Terrarium.")
}

func printErrorAndExit(msg string, err error, exitCode int) {
	fmt.Fprintf(os.Stderr, "ERROR: %s", msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, ": %s\n", err)
	} else {
		fmt.Fprintln(os.Stderr, "")
	}
	os.Exit(exitCode)
}
