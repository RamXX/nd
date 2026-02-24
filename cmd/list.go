package cmd

import (
	"fmt"
	"os"
	"time"

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
		parent, _ := cmd.Flags().GetString("parent")
		noParent, _ := cmd.Flags().GetBool("no-parent")
		createdAfterStr, _ := cmd.Flags().GetString("created-after")
		createdBeforeStr, _ := cmd.Flags().GetString("created-before")
		updatedAfterStr, _ := cmd.Flags().GetString("updated-after")
		updatedBeforeStr, _ := cmd.Flags().GetString("updated-before")
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

		// Parse date filters.
		var createdAfter, createdBefore, updatedAfter, updatedBefore time.Time
		dateFmt := "2006-01-02"
		if createdAfterStr != "" {
			if t, err := time.Parse(dateFmt, createdAfterStr); err != nil {
				return fmt.Errorf("invalid --created-after date: %w", err)
			} else {
				createdAfter = t
			}
		}
		if createdBeforeStr != "" {
			if t, err := time.Parse(dateFmt, createdBeforeStr); err != nil {
				return fmt.Errorf("invalid --created-before date: %w", err)
			} else {
				createdBefore = t.Add(24*time.Hour - time.Nanosecond) // end of day
			}
		}
		if updatedAfterStr != "" {
			if t, err := time.Parse(dateFmt, updatedAfterStr); err != nil {
				return fmt.Errorf("invalid --updated-after date: %w", err)
			} else {
				updatedAfter = t
			}
		}
		if updatedBeforeStr != "" {
			if t, err := time.Parse(dateFmt, updatedBeforeStr); err != nil {
				return fmt.Errorf("invalid --updated-before date: %w", err)
			} else {
				updatedBefore = t.Add(24*time.Hour - time.Nanosecond) // end of day
			}
		}

		s, err := store.Open(resolveVaultDir())
		if err != nil {
			return err
		}

		issues, err := s.ListIssues(store.FilterOptions{
			Status:        status,
			Type:          issueType,
			Assignee:      assignee,
			Label:         label,
			Priority:      priority,
			Parent:        parent,
			NoParent:      noParent,
			CreatedAfter:  createdAfter,
			CreatedBefore: createdBefore,
			UpdatedAfter:  updatedAfter,
			UpdatedBefore: updatedBefore,
			Sort:          sort,
			Reverse:       reverse,
			Limit:         limit,
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
	listCmd.Flags().String("parent", "", "filter by parent issue ID")
	listCmd.Flags().Bool("no-parent", false, "show only issues with no parent")
	listCmd.Flags().String("created-after", "", "filter by created date (YYYY-MM-DD)")
	listCmd.Flags().String("created-before", "", "filter by created date (YYYY-MM-DD)")
	listCmd.Flags().String("updated-after", "", "filter by updated date (YYYY-MM-DD)")
	listCmd.Flags().String("updated-before", "", "filter by updated date (YYYY-MM-DD)")
	listCmd.Flags().String("sort", "priority", "sort by: priority, created, updated, id")
	listCmd.Flags().BoolP("reverse", "r", false, "reverse sort order")
	listCmd.Flags().Bool("all", false, "show all issues including closed")
	listCmd.Flags().IntP("limit", "n", 50, "max results (0 for unlimited)")
	rootCmd.AddCommand(listCmd)
}
