package cmd

import "github.com/spf13/cobra"

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Manage prompts",
	Long:  `Create, list, version, diff, and promote prompts.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	promptCmd.AddCommand(promptListCmd)
	promptCmd.AddCommand(promptGetCmd)
	promptCmd.AddCommand(promptCreateCmd)
	promptCmd.AddCommand(promptPushCmd)
	promptCmd.AddCommand(promptPullCmd)
	promptCmd.AddCommand(promptDiffCmd)
	promptCmd.AddCommand(promptPromoteCmd)
	promptCmd.AddCommand(promptRollbackCmd)
	promptCmd.AddCommand(promptLogCmd)
}
