package cmd

import (
	"os"

	"github.com/RamXX/nd/internal/format"
	"github.com/RamXX/nd/internal/store"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List issues",
	RunE: func(cmd *cobra.Command, args []string) error {
		status, _ := cmd.Flags().GetString("status")
		issueType, _ := cmd.Flags().GetString("type")
		assignee, _ := cmd.Flags().GetString("assignee")
		label, _ := cmd.Flags().GetString("label")
		priority, _ := cmd.Flags().GetString("priority")
		sort, _ := cmd.Flags().GetString("sort")
		limit, _ := cmd.Flags().GetInt("limit")
		reverse, _ := cmd.Flags().GetBool("reverse")
		showAll, _ := cmd.Flags().GetBool("all")

		// Default: show non-closed issues (matching bd behavior).
		// Use --status=all to see everything, --status=closed for closed only.
		if !cmd.Flags().Changed("status") {
			if showAll {
				status = "all"
			} else {
				status = "!closed"
			}
		}

		// --all without explicit --limit removes the default cap.
		if showAll && !cmd.Flags().Changed("limit") {
			limit = 0
		}

		s, err := store.Open(resolveVaultDir())
		if err != nil {
			return err
		}

		issues, err := s.ListIssues(store.FilterOptions{
			Status:   status,
			Type:     issueType,
			Assignee: assignee,
			Label:    label,
			Priority: priority,
			Sort:     sort,
			Reverse:  reverse,
			Limit:    limit,
		})
		if err != nil {
			return err
		}

		if jsonOut {
			return format.JSON(os.Stdout, issues)
		}
		format.Table(os.Stdout, issues)
		return nil
	},
}

func init() {
	listCmd.Flags().StringP("status", "s", "", "filter by status")
	listCmd.Flags().String("type", "", "filter by type")
	listCmd.Flags().StringP("assignee", "a", "", "filter by assignee")
	listCmd.Flags().StringP("label", "l", "", "filter by label")
	listCmd.Flags().StringP("priority", "p", "", "filter by priority (0-4 or P0-P4)")
	listCmd.Flags().String("sort", "priority", "sort by: priority, created, updated, id")
	listCmd.Flags().BoolP("reverse", "r", false, "reverse sort order")
	listCmd.Flags().Bool("all", false, "show all issues including closed")
	listCmd.Flags().IntP("limit", "n", 50, "max results (0 for unlimited)")
	rootCmd.AddCommand(listCmd)
}
