package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start <workspace-id>",
	Short: "Resume a stopped workspace",
	Long: `Restarts a previously stopped workspace container.

Example:
  sailo start ws-7f3a`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		wsID := args[0]
		if err := deps.manager.Start(cmd.Context(), wsID); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Workspace %s started\n", wsID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
