package cmd

import (
	"log"

	"github.com/joinflux/webflowctl/internal"
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
		c := internal.NewClient(ApiToken)

		_, err := c.Delete([]string{"sites", siteId, "webhooks", webhookId})
		if err != nil {
			log.Fatalf("Unable to delete webhook: %v", err)
		}
	},
}

func init() {
	webhooksCmd.AddCommand(deleteCmd)
}
