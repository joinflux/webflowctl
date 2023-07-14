package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"os"
	"text/tabwriter"

	"github.com/joinflux/webflowctl/internal"
	"github.com/spf13/cobra"
)

func init() {
	webhooksCmd.AddCommand(createWebhooksCmd)
	webhooksCmd.AddCommand(deleteWebhooksCmd)
	webhooksCmd.AddCommand(getWebhooksCmd)
	webhooksCmd.AddCommand(listWebhooksCmd)
	rootCmd.AddCommand(webhooksCmd)
}

// webhooksCmd represents the webhooks command
var webhooksCmd = &cobra.Command{
	Use:   "webhooks",
	Short: "Manage webhooks",
	Long:  `List, create, delete and manage webhooks.`,
}

// CreateWebhooksResponse represents a response to the create request in Webflow.
// See: https://developers.webflow.com/reference/create-webhook.
type CreateWebhooksResponse struct {
	CreatedOn   string
	Id          string `json:"_id"`
	Site        string
	TriggerId   string
	TriggerType string
}

// TriggerTypes is a list of all available trigger types that can be created in Webflow.
var TriggerTypes = []string{
	"form_submission",
	"site_publish",
	"ecomm_new_order",
	"ecomm_order_changed",
	"ecomm_inventory_changed",
	"memberships_user_account_added",
	"memberships_user_account_updated",
	"memberships_user_account_deleted",
	"collection_item_created",
	"collection_item_changed",
	"collection_item_deleted",
	"collection_item_unpublished",
}

// createWebhooksCmd represents the command to create a webhook for a site in Webflow.
var createWebhooksCmd = &cobra.Command{
	Use:   "create [site_id] [trigger_type] [url]",
	Short: "create webhooks for a site",
	Args:  cobra.ExactArgs(3),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 1 {
			var candidates = []string{}
			// for autocompletion, we will suggest anything that contains the string
			// we are typing regardless of where in the string we match.
			// For example, if someone types: "item"
			// We will suggest:
			// - "collection_item_created"
			// - "collection_item_changed"
			// - "collection_item_deleted"
			// - "collection_item_unpublished"
			for _, value := range TriggerTypes {
				if strings.Contains(value, toComplete) {
					candidates = append(candidates, value)
				}
			}
			// if there are no valid suggestions, suggest all available options
			if len(candidates) == 0 {
				candidates = TriggerTypes
			}
			return candidates, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveError
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		triggerType := args[1]
		for _, b := range TriggerTypes {
			if b == triggerType {
				return nil
			}
		}
		return fmt.Errorf("unknown Trigger Type: '%s'.\ntrigger_type must be one of: %+q\n", triggerType, TriggerTypes)
	},
	Run: func(cmd *cobra.Command, args []string) {
		siteId := args[0]
		triggerType := args[1]
		url := args[2]

		c := internal.NewClient(ApiToken)

		payload := strings.NewReader(fmt.Sprintf(`{
        "triggerType": "%s",
        "url": "%s"
    }`, triggerType, url))

		body, err := c.Post([]string{"sites", siteId, "webhooks"}, payload)
		if err != nil {
			log.Fatalf("Request failed: %v", err)
		}

		var response CreateWebhooksResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatalf("Failed to unmarshal response body: %v", err)
		}
		fmt.Printf("%s\n", response.Id)
	},
}

// ListWebhooksResponse represents a response to the list request in Webflow.
// See: https://developers.webflow.com/reference/list-webhooks.
type ListWebhooksResponse []Webhook

// listWebhooksCmd represents the command to list webhooks for a site in Webflow.
var listWebhooksCmd = &cobra.Command{
	Use:   "list [site_id]",
	Short: "list webhooks for a site",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		siteId := args[0]

		c := internal.NewClient(ApiToken)
		body, err := c.Get([]string{"sites", siteId, "webhooks"})
		if err != nil {
			log.Fatalf("Request failed: %v", err)
		}

		var response ListWebhooksResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatalf("Failed to unmarshal response body")
		}
		for _, webhook := range response {
			fmt.Printf("%s\t%s\t%s\n", webhook.Id, webhook.TriggerType, webhook.Url)
		}
	},
}

// DeleteWebhooksResponse represents a response to the remove webhook request in Webflow.
// See: https://developers.webflow.com/reference/removewebhook.
type DeleteWebhooksResponse struct {
	Deleted int
}

// deleteWebhooksCmd represents the command to remove a webhook for a site in Webflow.
var deleteWebhooksCmd = &cobra.Command{
	Use:   "delete [site_id] [webhook_id]",
	Short: "delete a webhook for a site",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		siteId := args[0]
		webhookId := args[1]
		c := internal.NewClient(ApiToken)

		_, err := c.Delete([]string{"sites", siteId, "webhooks", webhookId})
		if err != nil {
			log.Fatalf("Unable to delete webhook: %v", err)
		}
	},
}

// GetWebhooksResponse represents a response to the list request in Webflow.
// See: https://developers.webflow.com/reference/list-webhooks.
type GetWebhooksResponse Webhook

// getWebhooksCmd represents the command to list webhooks for a site in Webflow.
var getWebhooksCmd = &cobra.Command{
	Use:   "get [site_id] [webhook_id]",
	Short: "get a webhook for a site",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		siteId := args[0]
		webhookId := args[1]

		client := internal.NewClient(ApiToken)
		body, err := client.Get([]string{"sites", siteId, "webhooks", webhookId})
		if err != nil {
			log.Fatalf("Request failed: %v", err)
		}

		var response GetWebhooksResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatalf("Failed to unmarshal response body")
		}
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Fprintf(w, "id:\t%s\n", response.Id)
		fmt.Fprintf(w, "created on:\t%s\n", response.CreatedOn)
		fmt.Fprintf(w, "last used:\t%s\n", response.LastUsed)
		fmt.Fprintf(w, "type:\t%s\n", response.TriggerType)
		fmt.Fprintf(w, "url:\t%s\n", response.Url)
		w.Flush()
	},
}
