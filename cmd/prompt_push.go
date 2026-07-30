package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zokhcat/neatcli/internal/api"
	"github.com/zokhcat/neatcli/internal/config"
	"github.com/zokhcat/neatcli/internal/workspace"
)

var promptPushCmd = &cobra.Command{
	Use:   "push <name>",
	Short: "Push a local prompt file to Neatlogs as a new version",
	Long: `Push a prompt from .neatlogs/prompts/<name>.yaml to Neatlogs.

This creates a new version on the server. Use --label to set labels,
or the labels from the YAML file are used.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if cfg.APIKey == "" {
			return fmt.Errorf("no API key configured; run `neatcli init` or set NEATLOGS_API_KEY")
		}

		ws, err := workspace.Open()
		if err != nil {
			return err
		}

		name := args[0]
		pf, err := ws.ReadPromptFile(name)
		if err != nil {
			return fmt.Errorf("read local prompt: %w", err)
		}

		labels, _ := cmd.Flags().GetStringSlice("labels")
		if len(labels) == 0 {
			labels = pf.Labels
		}
		msg, _ := cmd.Flags().GetString("message")

		var content string
		var messages []map[string]string

		if pf.Type == "chat" && len(pf.Messages) > 0 {
			for _, m := range pf.Messages {
				messages = append(messages, map[string]string{"role": m.Role, "content": m.Content})
			}
		} else {
			content = pf.Content
		}

		req := &api.SaveVersionRequest{
			PromptName:    pf.Name,
			Content:       content,
			Messages:      messages,
			Labels:        labels,
			CommitMessage: msg,
		}

		client := api.New(cfg.BaseURL, cfg.APIKey)
		p, err := client.SaveAsVersion(req)
		if err != nil {
			return fmt.Errorf("push failed: %w", err)
		}

		fmt.Printf("Pushed %q as v%d\n", p.Name, p.Version)
		if len(labels) > 0 {
			fmt.Printf("Labels: %s\n", joinLabels(labels))
		}
		return nil
	},
}

func init() {
	promptPushCmd.Flags().StringSlice("labels", nil, "Override labels (default: from YAML)")
	promptPushCmd.Flags().String("message", "", "Commit message")
}
