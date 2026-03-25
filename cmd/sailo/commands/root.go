package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "sailo",
	Short: "sAIlo — workspace isolation layer for AI agents",
	Long: `sAIlo creates isolated Docker workspaces for AI coding agents.

Each workspace gets its own container, git clone, port range, and SSH
forwarding. Agent-agnostic — attach Claude Code, Cursor, Codex, or
any tool.

  sailo create "add dark mode" --from=main
  sailo exec ws-7f3a -- claude-code
  sailo ship ws-7f3a`,
	Version: version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")

	rootCmd.SetVersionTemplate(fmt.Sprintf("sAIlo %s\n", version))

	// Initialize dependencies for commands that need the full dep chain.
	// Commands like init, config, help, and version skip this.
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		switch cmd.Name() {
		case "version", "help", "init", "config", "show", "set", "sailo":
			return nil
		}
		return initDeps()
	}

	rootCmd.PersistentPostRun = func(cmd *cobra.Command, args []string) {
		if deps.store != nil {
			deps.store.Close()
		}
	}
}
