package commands

import (
	"context"
	"fmt"
	"io"

	"github.com/agawish/sailo/pkg/workspace"
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
		img, _ := cmd.Flags().GetString("image")

		return runCreate(cmd.OutOrStdout(), deps.manager, cmd.Context(), workspace.CreateOptions{
			Task:       task,
			FromBranch: from,
			Image:      img,
		})
	},
}

func init() {
	createCmd.Flags().String("from", "main", "base branch to create workspace from")
	createCmd.Flags().String("image", "", "override base Docker image")
	rootCmd.AddCommand(createCmd)
}

func runCreate(out io.Writer, mgr *workspace.Manager, ctx context.Context, opts workspace.CreateOptions) error {
	ws, err := mgr.Create(ctx, opts)
	if err != nil {
		return err
	}

	fmt.Fprintln(out, "Workspace created!")
	fmt.Fprintln(out, "")
	fmt.Fprintf(out, "  ID:       %s\n", ws.ID)
	fmt.Fprintf(out, "  Branch:   %s\n", ws.Branch)
	fmt.Fprintf(out, "  Status:   %s\n", ws.State)
	if len(ws.Ports) > 0 {
		fmt.Fprintln(out, "  Ports:")
		for cp, hp := range ws.Ports {
			fmt.Fprintf(out, "    %d → localhost:%d\n", cp, hp)
		}
	}
	fmt.Fprintln(out, "")
	fmt.Fprintf(out, "  Attach:   sailo exec %s -- <agent>\n", ws.ID)
	fmt.Fprintf(out, "  Stop:     sailo stop %s\n", ws.ID)
	fmt.Fprintf(out, "  Remove:   sailo rm %s\n", ws.ID)
	return nil
}
