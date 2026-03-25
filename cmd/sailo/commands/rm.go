package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm <workspace-id>",
	Short: "Remove a workspace and clean up resources",
	Long: `Removes a workspace by stopping its container, freeing allocated ports,
and optionally deleting the git branch.

Example:
  sailo rm ws-7f3a
  sailo rm ws-7f3a --keep-branch`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		wsID := args[0]
		if err := deps.manager.Remove(cmd.Context(), wsID); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Workspace %s removed\n", wsID)
		return nil
	},
}

func init() {
	rmCmd.Flags().Bool("keep-branch", false, "keep the git branch after removing workspace")
	rmCmd.Flags().BoolP("force", "f", false, "force removal without confirmation")
	rootCmd.AddCommand(rmCmd)
}
