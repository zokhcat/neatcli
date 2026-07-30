package cmd

import (
	"fmt"
	"strings"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/spf13/cobra"
	"github.com/zokhcat/neatcli/internal/api"
	"github.com/zokhcat/neatcli/internal/config"
)

var promptDiffCmd = &cobra.Command{
	Use:   "diff <name> <version1> <version2>",
	Short: "Show differences between two prompt versions",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if cfg.APIKey == "" {
			return fmt.Errorf("no API key configured; run `neatcli init` or set NEATLOGS_API_KEY")
		}

		name := args[0]
		v1 := args[1]
		v2 := args[2]

		// Parse versions — support "latest" keyword
		ver1 := 0
		ver2 := 0
		label1 := ""
		label2 := ""

		if v1 == "latest" {
			label1 = "production"
		} else if v1[0] >= '0' && v1[0] <= '9' {
			fmt.Sscanf(v1, "%d", &ver1)
		} else {
			label1 = v1
		}
		if v2 == "latest" {
			label2 = "production"
		} else if v2[0] >= '0' && v2[0] <= '9' {
			fmt.Sscanf(v2, "%d", &ver2)
		} else {
			label2 = v2
		}

		client := api.New(cfg.BaseURL, cfg.APIKey)
		p1, err := client.GetPrompt(name, label1, ver1)
		if err != nil {
			return fmt.Errorf("get version %s: %w", v1, err)
		}
		p2, err := client.GetPrompt(name, label2, ver2)
		if err != nil {
			return fmt.Errorf("get version %s: %w", v2, err)
		}

		text1 := promptText(p1)
		text2 := promptText(p2)

		if text1 == text2 {
			fmt.Println("No differences.")
			return nil
		}

		edits := myers.ComputeEdits(span.URIFromPath(name), text1, text2)
		diff := gotextdiff.ToUnified(
			fmt.Sprintf("v%d (%s)", p1.Version, labelOrLatest(p1.Labels)),
			fmt.Sprintf("v%d (%s)", p2.Version, labelOrLatest(p2.Labels)),
			text1,
			edits,
		)
		fmt.Printf("%s", diff)
		return nil
	},
}

func promptText(p *api.Prompt) string {
	if p.Type == "chat" && len(p.Messages) > 0 {
		var parts []string
		for _, m := range p.Messages {
			parts = append(parts, fmt.Sprintf("[%s]\n%s", m["role"], m["content"]))
		}
		return strings.Join(parts, "\n\n")
	}
	return p.Content
}

func labelOrLatest(labels []string) string {
	for _, l := range labels {
		if l == "production" || l == "staging" {
			return l
		}
	}
	if len(labels) > 0 {
		return labels[0]
	}
	return "latest"
}
