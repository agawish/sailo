package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var previewCmd = &cobra.Command{
	Use:   "preview <workspace-id>",
	Short: "Open a workspace's mapped port in the browser",
	Long: `Opens the primary mapped port of a workspace in the default browser.

Example:
  sailo preview ws-7f3a
  sailo preview ws-7f3a --port 3000`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		wsID := args[0]
		fmt.Fprintf(cmd.OutOrStdout(), "sailo preview: not yet implemented\n")
		fmt.Fprintf(cmd.OutOrStdout(), "\n")
		fmt.Fprintf(cmd.OutOrStdout(), "  Workspace: %s\n", wsID)
		return nil
	},
}

func init() {
	previewCmd.Flags().Int("port", 0, "specific container port to preview")
	rootCmd.AddCommand(previewCmd)
}
