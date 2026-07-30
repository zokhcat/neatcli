package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var traceGetCmd = &cobra.Command{
	Use:   "get <trace-id>",
	Short: "View a trace's full details",
	Long: `View the full span tree of a trace, including LLM calls, tool calls,
inputs, outputs, timing, and token usage.

Note: requires the Neatlogs trace querying API (coming soon).
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Trace view: coming soon — available in the Neatlogs dashboard at https://app.neatlogs.com")
		return nil
	},
}
