package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage sailo configuration",
	Long: `View and modify sailo configuration.

Configuration files:
  Project: .sailo.yaml (checked into repo)
  User:    ~/.sailo/config.yaml

Example:
  sailo config show
  sailo config set defaults.from develop
  sailo config set defaults.port_range 4001-4999`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprintln(cmd.OutOrStdout(), "sailo config show: not yet implemented")
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprintf(cmd.OutOrStdout(), "sailo config set: not yet implemented\n")
		fmt.Fprintf(cmd.OutOrStdout(), "\n")
		fmt.Fprintf(cmd.OutOrStdout(), "  Key:   %s\n", args[0])
		fmt.Fprintf(cmd.OutOrStdout(), "  Value: %s\n", args[1])
		return nil
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(configCmd)
}
