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
//
// The following are available trigger types:
//
//   - form_submission
//   - site_publish
//   - ecomm_new_order
//   - ecomm_order_changed
//   - ecomm_inventory_changed
//   - memberships_user_account_added
//   - memberships_user_account_updated
//   - memberships_user_account_deleted
//   - collection_item_created
//   - collection_item_changed
//   - collection_item_deleted
//   - collection_item_unpublished
//
// See: https://developers.webflow.com/reference/create-webhook.
type CreateResponse struct {
	CreatedOn   string
	Id          string `json:"_id"`
	Site        string
	TriggerId   string
	TriggerType string
}

// createCmd represents the command to create a webhook for a site in Webflow.
var createCmd = &cobra.Command{
	Use:   "create [site_id] [trigger_type] [url]",
	Short: "create webhooks for a site",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		siteId := args[0]
		c := internal.NewClient(ApiToken)

		payload := strings.NewReader(fmt.Sprintf(`{
        "triggerType": "%s",
        "url": "%s"
    }`, args[1], args[2]))

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
