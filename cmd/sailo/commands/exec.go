package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec <workspace-id> -- <command>",
	Short: "Run a command inside a workspace",
	Long: `Executes a command inside a running workspace container.
Use this to attach an AI agent or run any command.

Example:
  sailo exec ws-7f3a -- claude-code
  sailo exec ws-7f3a -- cursor --folder /workspace
  sailo exec ws-7f3a -- bash
  sailo exec ws-7f3a -- npm test`,
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		wsID := args[0]

		var execArgs []string
		for i, arg := range args {
			if arg == "--" {
				execArgs = args[i+1:]
				break
			}
		}

		fmt.Fprintf(cmd.OutOrStdout(), "sailo exec: not yet implemented\n")
		fmt.Fprintf(cmd.OutOrStdout(), "\n")
		fmt.Fprintf(cmd.OutOrStdout(), "  Workspace: %s\n", wsID)
		if len(execArgs) > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "  Command:   %v\n", execArgs)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
