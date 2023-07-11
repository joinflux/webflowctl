package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/joinflux/webflowctl/internal"
	"github.com/spf13/cobra"
)

// ListResponse represents a response to the list request in Webflow.
// See: https://developers.webflow.com/reference/list-webhooks.
type ListResponse []Webhook

// listCmd represents the command to list webhooks for a site in Webflow.
var listCmd = &cobra.Command{
	Use:   "list [site_id]",
	Short: "list webhooks for a site",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		siteId := args[0]

		c := internal.NewClient(ApiToken)
		body, err := c.Get([]string{"sites", siteId, "webhooks"})

		var response ListResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatalf("Failed to unmarshal response body")
		}
		for _, webhook := range response {
			fmt.Printf("%s\t%s\t%s\n", webhook.Id, webhook.TriggerType, webhook.Url)
		}
	},
}

func init() {
	webhooksCmd.AddCommand(listCmd)
}
