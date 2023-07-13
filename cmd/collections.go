package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
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
	Id           string `json:"_id"`
	LastUpdated  string
	CreatedOn    string
	Name         string
	Slug         string
	SingularName string
}

// ListCollectionsResponse represents a response to the list collections request in Webflow.
// See: https://developers.webflow.com/reference/list-collections
type ListCollectionsResponse []Collection

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
		for _, collection := range response {
			v := reflect.ValueOf(collection)
			for i := 0; i < v.NumField(); i++ {
				name := v.Type().Field(i).Name
				value := v.Field(i).Interface()
				fmt.Fprintf(w, "%s:\t%s\n", name, value)
			}
			fmt.Fprint(w, "\n")
		}
		w.Flush()
	},
}

// GetCollectionResponse represents a response to the get collection request in Webflow.
// See: https://developers.webflow.com/reference/get-collection
type GetCollectionResponse struct {
	*Collection
	Fields []struct {
		Id       string
		Slug     string
		Name     string
		Archived bool `json:"_archived"`
		Draft    bool `json:"_draft"`
		Editable bool
		Required bool
	}
}

// getCollectionCmd represents the command to get detailed info on a collection in Webflow
var getCollectionCmd = &cobra.Command{
	Use:   "get [collection_id]",
	Short: "get details information for a collection",
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
		for _, val := range response.Fields {
			fmt.Fprintf(w, "id:\t%s\n", val.Id)
			fmt.Fprintf(w, "name:\t%s\n", val.Name)
			fmt.Fprintf(w, "slug:\t%s\n", val.Slug)
			fmt.Fprintf(w, "draft:\t%v\n", val.Draft)
			fmt.Fprintf(w, "archived:\t%v\n", val.Archived)
			fmt.Fprintf(w, "editable:\t%v\n", val.Editable)
			fmt.Fprintf(w, "required:\t%v\n", val.Required)
			fmt.Fprint(w, "\n")
		}
		fmt.Fprint(w, "\n")
		w.Flush()
	},
}
