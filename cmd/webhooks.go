package cmd

import (
	"github.com/spf13/cobra"
)

// webhooksCmd represents the webhooks command
var webhooksCmd = &cobra.Command{
	Use:   "webhooks",
	Short: "Manage webhooks",
	Long: `List, create, delete and manage webhooks.`,
}

func init() {
	rootCmd.AddCommand(webhooksCmd)
}
