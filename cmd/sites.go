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
	sitesCmd.AddCommand(getSitesCmd)
	sitesCmd.AddCommand(listDomainsCmd)
	sitesCmd.AddCommand(listSitesCmd)
	sitesCmd.AddCommand(publishSitesCmd)
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
	Id            string
	LastPublished string
	Name          string `json:"displayName"`
	PreviewUrl    string
	ShortName     string
	Timezone      string `json:"timeZone"`
}

// ListSitesResponse represents a response to the list sites request in Webflow.
// See: https://developers.webflow.com/reference/list-sites
type ListSitesResponse struct {
	Sites []Site
}

// listSitesCmd represents the command to list sites in Webflow.
var listSitesCmd = &cobra.Command{
	Use:   "list",
	Short: "list sites",
	Run: func(cmd *cobra.Command, args []string) {
		c := internal.NewClient(ApiToken)
		body, err := c.Get([]string{"sites"})
		if err != nil {
			log.Fatalf("Request failed: %v", err)
		}

		var response ListSitesResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatalf("Failed to unmarshal response body: %v", err)
		}
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", "Name", "ID", "Last Published", "Preview URL")
		for _, site := range response.Sites {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", site.Name, site.Id, site.LastPublished, site.PreviewUrl)
		}
		w.Flush()
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
		if err != nil {
			log.Fatalf("Request failed: %v", err)
		}

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

// PublishSiteResponse represents a response to the publish site request in Webflow.
// See: https://developers.webflow.com/reference/publish-site
type PublishSiteResponse struct {
	Queued bool
}

// publishSitesCmd represents the command to publish a site in Webflow.
var publishSitesCmd = &cobra.Command{
	Use:   "publish [site_id]",
	Short: "publish a site",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		siteId := args[0]
		c := internal.NewClient(ApiToken)
		body, err := c.Post([]string{"sites", siteId, "publish"}, nil)
		if err != nil {
			log.Fatalf("Request failed: %v", err)
		}

		var response PublishSiteResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatalf("Failed to unmarshal response body: %v", err)
		}
	},
}

type Domain struct {
	Id            string
	Url           string
	LastPublished string
}

// ListDomainsResponse represents a response to the list domains request in Webflow.
// See: https://developers.webflow.com/reference/list-domains
type ListDomainsResponse struct {
	CustomDomains []Domain `json:"custom_domains"`
}

// listDomainsCmd represents the command to list the domains for a site.
var listDomainsCmd = &cobra.Command{
	Use:   "domains [site_id]",
	Short: "list a site's domains",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		siteId := args[0]
		c := internal.NewClient(ApiToken)
		body, err := c.Get([]string{"sites", siteId, "custom_domains"})
		if err != nil {
			log.Fatalf("Request failed: %v", err)
		}

		var response ListDomainsResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatalf("Failed to unmarshal response body: %v", err)
		}

		for _, domain := range response.CustomDomains {
			fmt.Println(domain.Url)
		}
	},
}
