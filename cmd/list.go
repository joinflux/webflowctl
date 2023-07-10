package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// Webhook represents a webhook in Webflow
type Webhook struct {
	CreatedOn   string
	Id          string `json:"_id"`
	LastUsed    string
	Site        string
	TriggerId   string
	TriggerType string
	Url         string
}

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
		client := &http.Client{}

		url := fmt.Sprintf("https://api.webflow.com/sites/%s/webhooks", siteId)
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatalf("Failed to create request: %v", err)
		}

		request.Header.Add("authorization", "Bearer "+ApiToken)
		request.Header.Add("accept", "application/json")

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
			log.Printf("Failed to list webhooks: %s\n", resp.Status)
			log.Println(string(body))
			os.Exit(1)
		}

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
