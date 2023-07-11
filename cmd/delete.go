package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// DeleteResponse represents a response to the remove webhook request in Webflow.
// See: https://developers.webflow.com/reference/removewebhook.
type DeleteResponse struct {
	Deleted int
}

// deleteCmd represents the command to remove a webhook for a site in Webflow.
var deleteCmd = &cobra.Command{
	Use:   "delete [site_id] [webhook_id]",
	Short: "delete a webhook for a site",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		siteId := args[0]
		webhookId := args[1]
		client := &http.Client{}

		url := fmt.Sprintf("https://api.webflow.com/sites/%s/webhooks/%s", siteId, webhookId)
		request, err := http.NewRequest("DELETE", url, nil)
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
	},
}

func init() {
	webhooksCmd.AddCommand(deleteCmd)
}
