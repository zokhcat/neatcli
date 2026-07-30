package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zokhcat/neatcli/internal/api"
	"github.com/zokhcat/neatcli/internal/config"
	"github.com/zokhcat/neatcli/internal/workspace"
)

var promptPullCmd = &cobra.Command{
	Use:   "pull <name>",
	Short: "Pull a prompt from Neatlogs to local workspace",
	Long: `Pull the latest version of a prompt from Neatlogs and save it
to .neatlogs/prompts/<name>.yaml as a YAML file you can edit.

Use --label or --version to pull a specific variant.`,
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
		label, _ := cmd.Flags().GetString("label")
		version, _ := cmd.Flags().GetInt("version")

		client := api.New(cfg.BaseURL, cfg.APIKey)
		p, err := client.GetPrompt(name, label, version)
		if err != nil {
			return err
		}

		pf := &workspace.PromptFile{
			Name:   p.Name,
			Type:   p.Type,
			Labels: p.Labels,
			Content: p.Content,
		}
		for _, m := range p.Messages {
			pf.Messages = append(pf.Messages, workspace.Message{
				Role:    m["role"],
				Content: m["content"],
			})
		}

		if err := ws.WritePromptFile(pf); err != nil {
			return fmt.Errorf("write local file: %w", err)
		}

		labelStr := "latest"
		if label != "" {
			labelStr = label
		}
		fmt.Printf("Pulled %q (v%d, %s) to .neatlogs/prompts/%s.yaml\n", p.Name, p.Version, labelStr, p.Name)
		return nil
	},
}

func init() {
	promptPullCmd.Flags().String("label", "", "Pull version by label")
	promptPullCmd.Flags().Int("version", 0, "Pull specific version number")
}
