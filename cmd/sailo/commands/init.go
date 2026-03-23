package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize sailo in the current project",
	Long: `Scans the current project for existing Docker configuration
(Dockerfile, docker-compose.yml, devcontainer.json) and creates
a .sailo.yaml config file with sensible defaults.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprintln(cmd.OutOrStdout(), "sailo init: not yet implemented")
		fmt.Fprintln(cmd.OutOrStdout(), "")
		fmt.Fprintln(cmd.OutOrStdout(), "Will detect project configuration and create .sailo.yaml")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
