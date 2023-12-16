package cmd

import (
	"github.com/terrariumcloud/terrarium/pkg/terrarium/release"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/spf13/cobra"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Commands for managing releases",
}

func init() {
	rootCmd.AddCommand(releaseCmd)
}

func getReleasePublisherClient() (*grpc.ClientConn, release.ReleasePublisherClient, error) {
	conn, err := grpc.Dial(terrariumEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	return conn, release.NewReleasePublisherClient(conn), nil
}
