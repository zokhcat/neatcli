package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zokhcat/neatcli/internal/config"
)

var rootCmd = &cobra.Command{
	Use:   "neatcli",
	Short: "CLI for Neatlogs — agent observability & prompt management",
	Long: `neatcli is a command-line tool for managing Neatlogs resources.

It provides git-like workflows for prompt versioning, tool description
management, and trace inspection — all from your terminal.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "init" {
			return nil
		}
		return config.InitViper()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var cfgFile string

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default ~/.neatlogs/config.yaml)")
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(promptCmd)
	rootCmd.AddCommand(toolCmd)
	rootCmd.AddCommand(traceCmd)
}
