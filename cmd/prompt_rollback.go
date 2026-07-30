package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zokhcat/neatcli/internal/api"
	"github.com/zokhcat/neatcli/internal/config"
)

var promptRollbackCmd = &cobra.Command{
	Use:   "rollback <name>",
	Short: "Rollback to the previous production version",
	Long: `Find the version that was previously labeled "production"
and move the label back to it. Shows a diff before applying.

Example:
  neatcli prompt rollback support-prompt`,
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
		client := api.New(cfg.BaseURL, cfg.APIKey)

		resp, err := client.ListPrompts(name, "", 100, 0)
		if err != nil {
			return err
		}

		var currentProdVersion int
		var prevProdVersion int
		var prevProdID string

		for _, p := range resp.Items {
			for _, l := range p.Labels {
				if l == "production" {
					currentProdVersion = p.Version
					break
				}
			}
		}

		if currentProdVersion == 0 {
			return fmt.Errorf("no version with 'production' label found for %q", name)
		}

		for _, p := range resp.Items {
			if p.Version < currentProdVersion {
				isLabeledProd := false
				for _, l := range p.Labels {
					if l == "production" {
						isLabeledProd = true
						break
					}
				}
				if !isLabeledProd || p.Version > prevProdVersion {
					// Find the one that USED to be production —
					// it's the one before current with no production label
					if p.Version != currentProdVersion {
						prevProdVersion = p.Version
						prevProdID = p.ID
					}
				}
			}
		}

		if prevProdVersion == 0 {
			return fmt.Errorf("no previous version found to rollback to")
		}

		force, _ := cmd.Flags().GetBool("force")
		if !force {
			fmt.Printf("Rollback %q from v%d to v%d? (use --force to confirm)\n", name, currentProdVersion, prevProdVersion)
			return nil
		}

		if err := client.UpdatePromptLabels(name, prevProdVersion, []string{"production"}); err != nil {
			return fmt.Errorf("rollback failed: %w", err)
		}

		fmt.Printf("Rolled back %q to v%d\n", name, prevProdVersion)
		_ = prevProdID
		return nil
	},
}

func init() {
	promptRollbackCmd.Flags().Bool("force", false, "Execute the rollback")
}
