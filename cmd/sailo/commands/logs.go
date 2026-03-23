package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs <workspace-id>",
	Short: "Stream logs from a workspace",
	Long: `Streams stdout/stderr from a running workspace container.

Example:
  sailo logs ws-7f3a
  sailo logs ws-7f3a --follow`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		wsID := args[0]
		fmt.Fprintf(cmd.OutOrStdout(), "sailo logs: not yet implemented\n")
		fmt.Fprintf(cmd.OutOrStdout(), "\n")
		fmt.Fprintf(cmd.OutOrStdout(), "  Workspace: %s\n", wsID)
		return nil
	},
}

func init() {
	logsCmd.Flags().BoolP("follow", "f", false, "follow log output")
	logsCmd.Flags().Int("tail", 100, "number of lines to show from the end")
	rootCmd.AddCommand(logsCmd)
}
