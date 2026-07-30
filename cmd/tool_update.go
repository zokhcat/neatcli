package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zokhcat/neatcli/internal/workspace"
)

var toolUpdateCmd = &cobra.Command{
	Use:   "update <name>",
	Short: "Update a tool's description locally",
	Long: `Update a tool description in the local workspace.

The tool description is saved as .neatlogs/tools/<name>.yaml.
Use --desc to set the description and --schema for the JSON schema.

Example:
  neatcli tool update get_account --desc "Fetch account by customer ID"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ws, err := workspace.Open()
		if err != nil {
			return err
		}

		name := args[0]
		desc, _ := cmd.Flags().GetString("desc")
		schema, _ := cmd.Flags().GetString("schema")

		tf, err := ws.ReadToolFile(name)
		if err != nil {
			tf = &workspace.ToolFile{Name: name}
		}

		if desc != "" {
			tf.Description = desc
		}
		if schema != "" {
			tf.Schema = schema
		}

		if err := ws.WriteToolFile(tf); err != nil {
			return fmt.Errorf("write tool file: %w", err)
		}

		fmt.Printf("Updated tool %q in workspace\n", name)
		return nil
	},
}

func init() {
	toolUpdateCmd.Flags().String("desc", "", "New tool description")
	toolUpdateCmd.Flags().String("schema", "", "New JSON schema for the tool")
}
