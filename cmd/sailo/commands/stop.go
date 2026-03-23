package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop <workspace-id>",
	Short: "Stop a running workspace (preserves state)",
	Long: `Stops the workspace container but preserves all state.
Resume with 'sailo start'.

Example:
  sailo stop ws-7f3a`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		wsID := args[0]
		fmt.Fprintf(cmd.OutOrStdout(), "sailo stop: not yet implemented\n")
		fmt.Fprintf(cmd.OutOrStdout(), "\n")
		fmt.Fprintf(cmd.OutOrStdout(), "  Workspace: %s\n", wsID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
