package commands

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/agawish/sailo/pkg/workspace"
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
		all, _ := cmd.Flags().GetBool("all")
		return runPs(cmd.OutOrStdout(), deps.manager, cmd.Context(), all)
	},
}

func init() {
	psCmd.Flags().BoolP("all", "a", false, "show all workspaces including archived")
	rootCmd.AddCommand(psCmd)
}

func runPs(out io.Writer, mgr *workspace.Manager, ctx context.Context, includeAll bool) error {
	workspaces, err := mgr.List(ctx, includeAll)
	if err != nil {
		return err
	}

	if len(workspaces) == 0 {
		fmt.Fprintln(out, "No workspaces found. Create one with: sailo create <task>")
		return nil
	}

	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tTASK\tSTATE\tPORTS\tBRANCH")

	for _, ws := range workspaces {
		task := ws.Task
		if len(task) > 35 {
			task = task[:32] + "..."
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			ws.ID, task, ws.State, formatPorts(ws.Ports), ws.Branch)
	}
	w.Flush()
	return nil
}

func formatPorts(ports map[int]int) string {
	if len(ports) == 0 {
		return "-"
	}
	var parts []string
	for cp, hp := range ports {
		parts = append(parts, fmt.Sprintf("%d→%d", cp, hp))
	}
	sort.Strings(parts)
	return strings.Join(parts, ",")
}
