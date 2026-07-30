package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/zokhcat/neatcli/internal/api"
	"github.com/zokhcat/neatcli/internal/config"
	"github.com/zokhcat/neatcli/internal/workspace"
)

var toolListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tools from the workspace and remote",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		ws, wsErr := workspace.Open()

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tSOURCE\tDESCRIPTION")

		if wsErr == nil {
			files, err := ws.ListToolFiles()
			if err == nil {
				for _, f := range files {
					tf, err := ws.ReadToolFile(f[:len(f)-5])
					if err == nil {
						desc := tf.Description
						if len(desc) > 60 {
							desc = desc[:60] + "..."
						}
						fmt.Fprintf(w, "%s\tlocal\t%s\n", tf.Name, desc)
					}
				}
			}
		}

		if cfg.APIKey != "" {
			client := api.New(cfg.BaseURL, cfg.APIKey)
			tools, err := client.ListTools()
			if err == nil {
				for _, t := range tools {
					desc := t.Description
					if len(desc) > 60 {
						desc = desc[:60] + "..."
					}
					source := t.Source
					if source == "" {
						source = "remote"
					}
					fmt.Fprintf(w, "%s\t%s\t%s\n", t.Name, source, desc)
				}
			}
		}

		w.Flush()
		return nil
	},
}
