package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zokhcat/neatcli/internal/api"
	"github.com/zokhcat/neatcli/internal/config"
)

var promptPromoteCmd = &cobra.Command{
	Use:   "promote <name> <version>",
	Short: "Promote a prompt version to a label (e.g. production)",
	Long: `Move a label to a specific version.
Common labels: production, staging, candidate

Example:
  neatcli prompt promote support-prompt 3 --label production`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if cfg.APIKey == "" {
			return fmt.Errorf("no API key configured; run `neatcli init` or set NEATLOGS_API_KEY")
		}

		name := args[0]
		version := 0
		fmt.Sscanf(args[1], "%d", &version)
		if version <= 0 {
			return fmt.Errorf("invalid version number: %s", args[1])
		}

		labels, _ := cmd.Flags().GetStringSlice("label")
		if len(labels) == 0 {
			labels = []string{"production"}
		}

		client := api.New(cfg.BaseURL, cfg.APIKey)
		if err := client.UpdatePromptLabels(name, version, labels); err != nil {
			return err
		}

		labelStr := joinLabels(labels)
		fmt.Printf("Promoted %q v%d → %s\n", name, version, labelStr)
		return nil
	},
}

func init() {
	promptPromoteCmd.Flags().StringSlice("label", []string{"production"}, "Target label(s)")
}
