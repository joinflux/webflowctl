package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/joinflux/webflowctl/internal"
	"github.com/spf13/cobra"
)

func init() {
	sitesCmd.AddCommand(listSitesCmd)
	sitesCmd.AddCommand(getSitesCmd)
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
	Short: "list sites",
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

// GetSitesResponse represents a response to the get site request in Webflow.
// See: https://developers.webflow.com/reference/get-site
type GetSitesResponse Site

// getSitesCmd represents the command to get a site in Webflow.
var getSitesCmd = &cobra.Command{
	Use:   "get [site_id]",
	Short: "get a site",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		siteId := args[0]
		c := internal.NewClient(ApiToken)
		body, err := c.Get([]string{"sites", siteId})

		var response GetSitesResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatalf("Failed to unmarshal response body: %v", err)
		}
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintf(w, "name:\t%s\n", response.Name)
		fmt.Fprintf(w, "id:\t%s\n", response.Id)
		fmt.Fprintf(w, "created on:\t%s\n", response.CreatedOn)
		fmt.Fprintf(w, "last published:\t%s\n", response.LastPublished)
		fmt.Fprintf(w, "preview url:\t%s\n", response.PreviewUrl)
		fmt.Fprintf(w, "short name:\t%s\n", response.ShortName)
		fmt.Fprintf(w, "timezone:\t%s\n", response.Timezone)
		w.Flush()
	},
}
