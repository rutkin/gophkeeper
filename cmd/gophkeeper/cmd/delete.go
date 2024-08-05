package cmd

import (
	"github.com/spf13/cobra"
)

var deleteDataID string

func init() {
	deleteCmd.PersistentFlags().StringVar(&deleteDataID, "id", "", "data identificator")
	deleteCmd.MarkPersistentFlagRequired("id")
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete item",
	RunE: func(cmd *cobra.Command, args []string) error {
		return makeRequest(upstreamURL + "/api/keeper/delete/" + deleteDataID)
	},
}
