package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/zokhcat/neatcli/internal/api"
	"github.com/zokhcat/neatcli/internal/config"
)

var promptLogCmd = &cobra.Command{
	Use:   "log <name>",
	Short: "Show version history for a prompt",
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

		client := api.New(cfg.BaseURL, cfg.APIKey)
		resp, err := client.ListPrompts(name, "", 100, 0)
		if err != nil {
			return err
		}

		if len(resp.Items) == 0 {
			fmt.Printf("No versions found for %q\n", name)
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "VERSION\tLABELS\tCREATED")
		for _, p := range resp.Items {
			created := p.CreatedAt
			if len(created) > 19 {
				created = created[:19]
			}
			fmt.Fprintf(w, "v%d\t%s\t%s\n", p.Version, joinLabels(p.Labels), created)
		}
		w.Flush()
		return nil
	},
}
