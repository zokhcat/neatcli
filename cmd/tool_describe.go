package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zokhcat/neatcli/internal/api"
	"github.com/zokhcat/neatcli/internal/config"
	"github.com/zokhcat/neatcli/internal/workspace"
)

var toolDescribeCmd = &cobra.Command{
	Use:   "describe <name>",
	Short: "Show a tool's description and schema",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		ws, wsErr := workspace.Open()
		if wsErr == nil {
			if tf, err := ws.ReadToolFile(name); err == nil {
				fmt.Printf("Name:        %s\n", tf.Name)
				fmt.Printf("Description: %s\n", tf.Description)
				if tf.Schema != "" {
					fmt.Println("--- Schema ---")
					fmt.Println(tf.Schema)
				}
				return nil
			}
		}

		cfg, err := config.Load()
		if err == nil && cfg.APIKey != "" {
			client := api.New(cfg.BaseURL, cfg.APIKey)
			t, err := client.GetTool(name)
			if err == nil {
				fmt.Printf("Name:        %s\n", t.Name)
				fmt.Printf("Description: %s\n", t.Description)
				if t.Schema != "" {
					fmt.Println("--- Schema ---")
					fmt.Println(t.Schema)
				}
				return nil
			}
		}

		return fmt.Errorf("tool %q not found", name)
	},
}
