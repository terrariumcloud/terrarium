/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/terrariumcloud/terrarium/pkg/terrarium/release"
	"strings"
)

var (
	releaseToPublish = release.PublishRequest{
		Type:         "",
		Organization: "",
		Name:         "",
		Version:      "",
		Description:  "",
		Links:        nil,
	}
	releaseLinks []string
)

// releasePublishCmd represents the publish command
var releasePublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish details of a new release.",
	Run: func(cmd *cobra.Command, args []string) {
		conn, client, err := getReleasePublisherClient()
		if err != nil {
			printErrorAndExit("Failed to connect to terrarium", err, 1)
		}
		defer func() { _ = conn.Close() }()

		for _, link := range releaseLinks {
			if parts := strings.SplitN(link, "=", 2); len(parts) == 1 {
				releaseToPublish.Links = append(releaseToPublish.Links, &release.Link{
					Title: "",
					Url:   parts[0],
				})
			} else {
				releaseToPublish.Links = append(releaseToPublish.Links, &release.Link{
					Title: parts[0],
					Url:   parts[1],
				})
			}
		}

		if _, err := client.Publish(context.Background(), &releaseToPublish); err != nil {
			printErrorAndExit("Failed to publish release", err, 1)
		}
		fmt.Println("Release published.")
	},
}

func init() {
	releaseCmd.AddCommand(releasePublishCmd)
	releasePublishCmd.Flags().StringVarP(&releaseToPublish.Name, "name", "n", "", "Name of the release.")
	releasePublishCmd.MarkFlagRequired("name")
	releasePublishCmd.Flags().StringVarP(&releaseToPublish.Version, "version", "v", "", "Version of the release.")
	releasePublishCmd.MarkFlagRequired("version")
	releasePublishCmd.Flags().StringVarP(&releaseToPublish.Organization, "org", "o", "", "Organization to which the release belongs.")
	releasePublishCmd.MarkFlagRequired("org")
	releasePublishCmd.Flags().StringVarP(&releaseToPublish.Type, "type", "t", "", "Type of release.")
	releasePublishCmd.MarkFlagRequired("type")
	releasePublishCmd.Flags().StringVarP(&releaseToPublish.Description, "description", "d", "", "Description of the release being published.")
	releasePublishCmd.Flags().StringSliceVarP(&releaseLinks, "link", "l", []string{}, `Links and optional titles eg.: --link "title=url" or --link "url"`)
}
