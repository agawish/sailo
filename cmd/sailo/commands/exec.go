package commands

import (
	"fmt"

	"github.com/agawish/sailo/pkg/workspace"
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

		// Parse command after --
		var execArgs []string
		for i, arg := range args {
			if arg == "--" {
				execArgs = args[i+1:]
				break
			}
		}
		if len(execArgs) == 0 {
			execArgs = []string{"bash"}
		}

		// Look up workspace
		ws, err := deps.manager.Get(cmd.Context(), wsID)
		if err != nil {
			return err
		}
		if ws.State != workspace.StateRunning {
			return fmt.Errorf("workspace %s is %s, not running; use 'sailo start %s' first", wsID, ws.State, wsID)
		}

		// Check if container is actually running
		state, err := deps.container.InspectContainer(cmd.Context(), ws.ContainerID)
		if err == nil && !state.Running {
			return fmt.Errorf("workspace %s container has stopped unexpectedly; use 'sailo start %s' to restart", wsID, wsID)
		}

		return deps.container.ExecInteractive(cmd.Context(), ws.ContainerID, execArgs)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
