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
		sort, _ := cmd.Flags().GetString("sort")
		limit, _ := cmd.Flags().GetInt("limit")

		// Default: show non-closed issues (matching bd behavior).
		// Use --status=all to see everything, --status=closed for closed only.
		if !cmd.Flags().Changed("status") {
			status = "!closed"
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
			Sort:     sort,
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
	listCmd.Flags().String("status", "", "filter by status")
	listCmd.Flags().String("type", "", "filter by type")
	listCmd.Flags().String("assignee", "", "filter by assignee")
	listCmd.Flags().String("label", "", "filter by label")
	listCmd.Flags().String("sort", "priority", "sort by: priority, created, updated, id")
	listCmd.Flags().IntP("limit", "n", 0, "max results")
	rootCmd.AddCommand(listCmd)
}
