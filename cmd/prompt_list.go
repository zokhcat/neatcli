package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/zokhcat/neatcli/internal/config"
	"github.com/zokhcat/neatcli/internal/api"
)

var promptListCmd = &cobra.Command{
	Use:   "list",
	Short: "List prompts and their versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if cfg.APIKey == "" {
			return fmt.Errorf("no API key configured; run `neatcli init` or set NEATLOGS_API_KEY")
		}

		name, _ := cmd.Flags().GetString("name")
		label, _ := cmd.Flags().GetString("label")
		limit, _ := cmd.Flags().GetInt("limit")

		client := api.New(cfg.BaseURL, cfg.APIKey)
		resp, err := client.ListPrompts(name, label, limit, 0)
		if err != nil {
			return err
		}

		if len(resp.Items) == 0 {
			fmt.Println("No prompts found.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tVERSION\tTYPE\tLABELS\tUPDATED")
		for _, p := range resp.Items {
			labels := joinLabels(p.Labels)
			updated := p.UpdatedAt
			if updated == "" {
				updated = p.CreatedAt
			}
			if len(updated) > 19 {
				updated = updated[:19]
			}
			fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\n", p.Name, p.Version, p.Type, labels, updated)
		}
		w.Flush()
		return nil
	},
}

func joinLabels(labels []string) string {
	if len(labels) == 0 {
		return "-"
	}
	s := ""
	for i, l := range labels {
		if i > 0 {
			s += ", "
		}
		s += l
	}
	return s
}

func init() {
	promptListCmd.Flags().String("name", "", "Filter by prompt name")
	promptListCmd.Flags().String("label", "", "Filter by label")
	promptListCmd.Flags().Int("limit", 50, "Maximum results")
}
