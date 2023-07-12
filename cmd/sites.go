package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/joinflux/webflowctl/internal"
	"github.com/spf13/cobra"
)

func init() {
	sitesCmd.AddCommand(listSitesCmd)
	rootCmd.AddCommand(sitesCmd)
}

// sitesCmd represents the sites command
var sitesCmd = &cobra.Command{
	Use:   "sites",
	Short: "Manage sites",
	Long:  `List, create, delete and manage sites.`,
}

type Site struct {
	CreatedOn     string
	Id            string `json:"_id"`
	LastPublished string
	Name          string
	PreviewUrl    string
	ShortName     string
	Timezone      string
}

// ListSitesResponse represents a response to the list sites request in Webflow.
// See: https://developers.webflow.com/reference/list-sites
type ListSitesResponse []Site

// listSitesCmd represents the command to list sites in Webflow.
var listSitesCmd = &cobra.Command{
	Use:   "list",
	Short: "create webhooks for a site",
	Run: func(cmd *cobra.Command, args []string) {
		c := internal.NewClient(ApiToken)
		body, err := c.Get([]string{"sites"})

		var response ListSitesResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatalf("Failed to unmarshal response body: %v", err)
		}
		for _, site := range response {
			fmt.Printf("%s\t%s\t%s\t%s\n", site.Name, site.Id, site.LastPublished, site.PreviewUrl)
		}
	},
}
