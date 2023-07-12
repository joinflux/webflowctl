package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// ApiToken is the Webflow API Token
var ApiToken string

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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "webflowctl",
	Short: "A command line tool to interact with the Webflow API",
	Long:  `A tool to help manage webhooks in the Webflow API`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		flagValue, _ := cmd.PersistentFlags().GetString("api-token")
		if flagValue == "" {
			envValue := os.Getenv("WEBFLOW_API_TOKEN")
			if envValue == "" {
				return fmt.Errorf("Missing API Token: --api-token (-a) [WEBFLOW_API_TOKEN]")
			}
			ApiToken = envValue
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&ApiToken, "api-token", "a", "", "Webflow API Token")
}
