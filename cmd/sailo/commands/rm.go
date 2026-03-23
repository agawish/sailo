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
		keepBranch, _ := cmd.Flags().GetBool("keep-branch")

		fmt.Fprintf(cmd.OutOrStdout(), "sailo rm: not yet implemented\n")
		fmt.Fprintf(cmd.OutOrStdout(), "\n")
		fmt.Fprintf(cmd.OutOrStdout(), "  Workspace:    %s\n", wsID)
		fmt.Fprintf(cmd.OutOrStdout(), "  Keep branch:  %v\n", keepBranch)
		return nil
	},
}

func init() {
	rmCmd.Flags().Bool("keep-branch", false, "keep the git branch after removing workspace")
	rmCmd.Flags().BoolP("force", "f", false, "force removal without confirmation")
	rootCmd.AddCommand(rmCmd)
}
