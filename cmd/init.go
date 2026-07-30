package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zokhcat/neatcli/internal/config"
	"github.com/zokhcat/neatcli/internal/workspace"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize neatcli config and workspace",
	Long: `Initialize neatcli by setting up your API key and creating a workspace.

This will:
1. Prompt for your Neatlogs API key (or use NEATLOGS_API_KEY env var)
2. Save the config to ~/.neatlogs/config.yaml
3. Create a .neatlogs/ directory in the current folder

You can get an API key from https://app.neatlogs.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := os.Getenv("NEATLOGS_API_KEY")
		if apiKey == "" {
			fmt.Print("Enter your Neatlogs API key: ")
			fmt.Scanln(&apiKey)
		}
		if apiKey == "" {
			return fmt.Errorf("API key is required (set NEATLOGS_API_KEY or enter it)")
		}

		baseURL, _ := cmd.Flags().GetString("base-url")
		if baseURL == "" {
			baseURL = "https://ingest.neatlogs.com"
		}

		project, _ := cmd.Flags().GetString("project")

		cfg := config.Config{
			APIKey:  apiKey,
			BaseURL: baseURL,
			Project: project,
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("save config: %w", err)
		}
		fmt.Printf("Config saved to ~/%s/%s\n", config.ConfigDir, config.ConfigFile)

		if _, err := workspace.Init(); err != nil {
			return fmt.Errorf("init workspace: %w", err)
		}
		fmt.Printf("Workspace initialized at ./%s/\n", workspace.DirName)
		fmt.Printf("  prompts/  — place prompt YAML files here\n")
		fmt.Printf("  tools/    — place tool description YAML files here\n")

		return nil
	},
}

func init() {
	initCmd.Flags().String("base-url", "https://ingest.neatlogs.com", "Neatlogs API base URL")
	initCmd.Flags().String("project", "", "Default project name")
}
