package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var traceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent traces",
	Long: `List recent agent traces. Requires the Neatlogs management API.
Note: trace querying via the management API may not be available yet.
Check the Neatlogs dashboard at https://app.neatlogs.com for full trace views.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Trace list: coming soon — available in the Neatlogs dashboard at https://app.neatlogs.com")
		return nil
	},
}
