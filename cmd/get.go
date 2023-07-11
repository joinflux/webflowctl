package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// GetResponse represents a response to the list request in Webflow.
// See: https://developers.webflow.com/reference/list-webhooks.
type GetResponse Webhook

// getCmd represents the command to list webhooks for a site in Webflow.
var getCmd = &cobra.Command{
	Use:   "get [site_id] [webhook_id]",
	Short: "get a webhook for a site",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		siteId := args[0]
		webhookId := args[1]
		client := &http.Client{}

		url := fmt.Sprintf("https://api.webflow.com/sites/%s/webhooks/%s", siteId, webhookId)
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

		var response GetResponse
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

func init() {
	webhooksCmd.AddCommand(getCmd)
}
