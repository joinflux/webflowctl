package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/joinflux/webflowctl/internal"
	"github.com/spf13/cobra"
)

// CreateResponse represents a response to the create request in Webflow.
// See: https://developers.webflow.com/reference/create-webhook.
type CreateResponse struct {
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

// createCmd represents the command to create a webhook for a site in Webflow.
var createCmd = &cobra.Command{
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

		var response CreateResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatalf("Failed to unmarshal response body: %v", err)
		}
		fmt.Printf("%s\n", response.Id)
	},
}

func init() {
	webhooksCmd.AddCommand(createCmd)
}
