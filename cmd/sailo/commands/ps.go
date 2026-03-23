package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List all workspaces",
	Long: `Shows all sailo workspaces with their status, ports, and branches.

  ID       TASK                          STATUS    PORT   BRANCH
  ws-7f3a  add dark mode to settings     running   3007   sailo/ws-7f3a/dark-mode
  ws-9b1c  fix pagination bug            stopped   3008   sailo/ws-9b1c/fix-pagination`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprintln(cmd.OutOrStdout(), "sailo ps: not yet implemented")
		fmt.Fprintln(cmd.OutOrStdout(), "")
		fmt.Fprintln(cmd.OutOrStdout(), "No workspaces found. Create one with: sailo create <task>")
		return nil
	},
}

func init() {
	psCmd.Flags().BoolP("all", "a", false, "show all workspaces including archived")
	rootCmd.AddCommand(psCmd)
}
