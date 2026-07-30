package cmd

import "github.com/spf13/cobra"

var toolCmd = &cobra.Command{
	Use:   "tool",
	Short: "Manage tool descriptions",
	Long:  `List, describe, and update tool descriptions used by your agents.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	toolCmd.AddCommand(toolListCmd)
	toolCmd.AddCommand(toolDescribeCmd)
	toolCmd.AddCommand(toolUpdateCmd)
}
