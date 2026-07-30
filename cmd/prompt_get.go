package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zokhcat/neatcli/internal/api"
	"github.com/zokhcat/neatcli/internal/config"
)

var promptGetCmd = &cobra.Command{
	Use:   "get <name>",
	Short: "Get a prompt's details and content",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if cfg.APIKey == "" {
			return fmt.Errorf("no API key configured; run `neatcli init` or set NEATLOGS_API_KEY")
		}

		name := args[0]
		label, _ := cmd.Flags().GetString("label")
		version, _ := cmd.Flags().GetInt("version")

		client := api.New(cfg.BaseURL, cfg.APIKey)
		p, err := client.GetPrompt(name, label, version)
		if err != nil {
			return err
		}

		fmt.Printf("Name:    %s\n", p.Name)
		fmt.Printf("Version: %d\n", p.Version)
		fmt.Printf("Type:    %s\n", p.Type)
		fmt.Printf("Labels:  %s\n", joinLabels(p.Labels))
		fmt.Printf("Updated: %s\n", p.UpdatedAt)
		fmt.Println("---")
		if p.Type == "chat" && len(p.Messages) > 0 {
			for _, m := range p.Messages {
				fmt.Printf("[%s]\n%s\n\n", m["role"], m["content"])
			}
		} else {
			fmt.Println(p.Content)
		}
		return nil
	},
}

func init() {
	promptGetCmd.Flags().String("label", "", "Get version by label")
	promptGetCmd.Flags().Int("version", 0, "Get specific version number")
}
