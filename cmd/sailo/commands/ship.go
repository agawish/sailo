package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var shipCmd = &cobra.Command{
	Use:   "ship <workspace-id>",
	Short: "Extract workspace changes into a pull request",
	Long: `Commits all changes in the workspace, pushes the branch, and creates
a pull request on GitHub.

Steps:
  1. Run tests (unless --skip-tests)
  2. Stage and commit all changes
  3. Push branch to remote
  4. Create PR via gh CLI
  5. Archive workspace

Example:
  sailo ship ws-7f3a
  sailo ship ws-7f3a --skip-tests`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		wsID := args[0]
		skipTests, _ := cmd.Flags().GetBool("skip-tests")

		fmt.Fprintf(cmd.OutOrStdout(), "sailo ship: not yet implemented\n")
		fmt.Fprintf(cmd.OutOrStdout(), "\n")
		fmt.Fprintf(cmd.OutOrStdout(), "  Workspace:   %s\n", wsID)
		fmt.Fprintf(cmd.OutOrStdout(), "  Skip tests:  %v\n", skipTests)
		return nil
	},
}

func init() {
	shipCmd.Flags().Bool("skip-tests", false, "skip running tests before shipping")
	rootCmd.AddCommand(shipCmd)
}
