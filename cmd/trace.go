package cmd

import "github.com/spf13/cobra"

var traceCmd = &cobra.Command{
	Use:   "trace",
	Short: "Inspect traces",
	Long:  `List and view agent traces from Neatlogs.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	traceCmd.AddCommand(traceListCmd)
	traceCmd.AddCommand(traceGetCmd)
}
