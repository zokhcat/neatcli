package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zokhcat/neatcli/internal/api"
	"github.com/zokhcat/neatcli/internal/config"
)

var promptCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new prompt",
	Long: `Create a new prompt with content from --file or --content.

If --type is "chat", use --messages to pass a JSON array of {"role","content"} objects.

Examples:
  neatcli prompt create my-prompt --file prompt.txt --labels production
  neatcli prompt create my-prompt --content "You are a {{role}} assistant" --labels staging`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if cfg.APIKey == "" {
			return fmt.Errorf("no API key configured; run `neatcli init` or set NEATLOGS_API_KEY")
		}

		name := args[0]
		promptType, _ := cmd.Flags().GetString("type")
		labels, _ := cmd.Flags().GetStringSlice("labels")
		content, _ := cmd.Flags().GetString("content")
		file, _ := cmd.Flags().GetString("file")
		commitMsg, _ := cmd.Flags().GetString("message")

		if len(labels) == 0 {
			labels = []string{"staging"}
		}

		if file != "" {
			data, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("read file: %w", err)
			}
			content = string(data)
		}

		if content == "" && promptType != "chat" {
			return fmt.Errorf("provide content with --content or --file")
		}

		req := &api.CreatePromptRequest{
			Name:          name,
			Type:          promptType,
			Content:       content,
			Labels:        labels,
			CommitMessage: commitMsg,
		}

		client := api.New(cfg.BaseURL, cfg.APIKey)
		p, err := client.CreatePrompt(req)
		if err != nil {
			return err
		}

		fmt.Printf("Created prompt %q (v%d)\n", p.Name, p.Version)
		return nil
	},
}

func init() {
	promptCreateCmd.Flags().String("type", "text", "Prompt type: text or chat")
	promptCreateCmd.Flags().StringSlice("labels", nil, "Labels to apply (e.g. production,staging)")
	promptCreateCmd.Flags().String("content", "", "Prompt content (for text type)")
	promptCreateCmd.Flags().String("file", "", "Read prompt content from file")
	promptCreateCmd.Flags().String("message", "", "Commit message for this version")
}
