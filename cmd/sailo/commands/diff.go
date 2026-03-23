package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff <workspace-id>",
	Short: "Show changes made in a workspace",
	Long: `Shows the git diff of all changes made inside a workspace since creation.

Example:
  sailo diff ws-7f3a
  sailo diff ws-7f3a --stat`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		wsID := args[0]
		fmt.Fprintf(cmd.OutOrStdout(), "sailo diff: not yet implemented\n")
		fmt.Fprintf(cmd.OutOrStdout(), "\n")
		fmt.Fprintf(cmd.OutOrStdout(), "  Workspace: %s\n", wsID)
		return nil
	},
}

func init() {
	diffCmd.Flags().Bool("stat", false, "show diffstat instead of full diff")
	rootCmd.AddCommand(diffCmd)
}
