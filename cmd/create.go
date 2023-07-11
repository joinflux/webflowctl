package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

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
		client := &http.Client{}

		payloadString := fmt.Sprintf(`{
        "triggerType": "%s",
        "url": "%s"
    }`, args[1], args[2])
		payload := strings.NewReader(payloadString)

		url := fmt.Sprintf("https://api.webflow.com/sites/%s/webhooks", siteId)
		request, err := http.NewRequest("POST", url, payload)
		if err != nil {
			log.Fatalf("Failed to create request: %v", err)
		}

		request.Header.Add("authorization", "Bearer "+ApiToken)
		request.Header.Add("accept", "application/json")
		request.Header.Add("content-type", "application/json")

		resp, err := client.Do(request)
		if err != nil {
			log.Fatalf("Error sending request: %v", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			log.Fatalf("Error reading response payload: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("Failed to create webhook: (%s)\n", resp.Status)
			log.Println(string(body))
			os.Exit(1)
		}

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
