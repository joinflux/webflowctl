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
	collectionsCmd.AddCommand(listCollectionsCmd)
	collectionsCmd.AddCommand(getCollectionCmd)
	rootCmd.AddCommand(collectionsCmd)
}

// collectionsCmd represents the webhooks command
var collectionsCmd = &cobra.Command{
	Use:   "collections",
	Short: "Manage collections",
	Long:  "List and manage collections",
}

type Collection struct {
	Id           string
	LastUpdated  string
	CreatedOn    string
	Name         string `json:"displayName"`
	Slug         string
	SingularName string
}

// ListCollectionsResponse represents a response to the list collections request in Webflow.
// See: https://developers.webflow.com/reference/list-collections
type ListCollectionsResponse struct {
	Collections []Collection
}

// listCollectionsCmd represents the command to list all collections for a site in Webflow.
var listCollectionsCmd = &cobra.Command{
	Use:   "list [site_id]",
	Short: "list collections for a site",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		siteId := args[0]

		c := internal.NewClient(ApiToken)

		body, err := c.Get([]string{"sites", siteId, "collections"})
		if err != nil {
			log.Fatalf("Request failed: %v", err)
		}

		var response ListCollectionsResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatalf("Failed to unmarshal response body: %v", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", "ID", "Name", "Slug", "Last Updated")
		for _, collection := range response.Collections {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", collection.Id, collection.Name, collection.Slug, collection.LastUpdated)
		}
		w.Flush()
	},
}

// GetCollectionResponse represents a response to the get collection request in Webflow.
// See: https://developers.webflow.com/reference/get-collection
type GetCollectionResponse struct {
	*Collection
	Fields []struct {
		DisplayName string
		HelpText    string
		Id          string
		IsEditable  bool
		IsRequired  bool
		Slug        string
		Type        string
	}
}

// getCollectionCmd represents the command to get detailed info on a collection in Webflow
var getCollectionCmd = &cobra.Command{
	Use:   "get [collection_id]",
	Short: "get detailed information for a collection",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		collectionId := args[0]

		c := internal.NewClient(ApiToken)

		body, err := c.Get([]string{"collections", collectionId})
		if err != nil {
			log.Fatalf("Request failed: %v", err)
		}

		var response GetCollectionResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatalf("Failed to unmarshal response body: %v", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintf(w, "id:\t%s\n", response.Id)
		fmt.Fprintf(w, "name:\t%s\n", response.Name)
		fmt.Fprintf(w, "slug:\t%s\n", response.Slug)
		fmt.Fprintf(w, "singular name:\t%v\n", response.SingularName)
		fmt.Fprintf(w, "crated on:\t%v\n", response.CreatedOn)
		fmt.Fprintf(w, "last updated:\t%v\n", response.LastUpdated)
		fmt.Fprint(w, "\n\nFields:\n")
		fmt.Fprint(w, "ID\tName\tSlug\tEditable\tRequired\tType\tHelp Text\n")
		for _, val := range response.Fields {
			fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\n", val.Id, val.DisplayName, val.Slug, val.IsEditable, val.IsRequired, val.Type, val.HelpText)
		}
		fmt.Fprint(w, "\n")
		w.Flush()
	},
}
