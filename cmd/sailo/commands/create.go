package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create <task-description>",
	Short: "Create an isolated workspace for an AI agent",
	Long: `Creates a new isolated Docker workspace with:
  - A fresh container (using existing Dockerfile or auto-detected base image)
  - A shallow git clone on a new branch
  - Non-conflicting port mappings
  - SSH agent forwarding for git operations

Example:
  sailo create "add dark mode to settings page" --from=main
  sailo create "fix pagination bug" --from=feature/users`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		task := args[0]
		from, _ := cmd.Flags().GetString("from")

		fmt.Fprintf(cmd.OutOrStdout(), "sailo create: not yet implemented\n")
		fmt.Fprintf(cmd.OutOrStdout(), "\n")
		fmt.Fprintf(cmd.OutOrStdout(), "  Task:   %s\n", task)
		fmt.Fprintf(cmd.OutOrStdout(), "  From:   %s\n", from)
		fmt.Fprintf(cmd.OutOrStdout(), "\n")
		fmt.Fprintf(cmd.OutOrStdout(), "Will create isolated workspace with container + git clone + port mapping\n")
		return nil
	},
}

func init() {
	createCmd.Flags().String("from", "main", "base branch to create workspace from")
	createCmd.Flags().String("image", "", "override base Docker image")
	rootCmd.AddCommand(createCmd)
}
