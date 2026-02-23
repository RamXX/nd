package cmd

import (
	"os"
	"strings"

	"github.com/RamXX/nd/internal/format"
	"github.com/RamXX/nd/internal/graph"
	"github.com/RamXX/nd/internal/model"
	"github.com/RamXX/nd/internal/store"
	"github.com/spf13/cobra"
)

var readyCmd = &cobra.Command{
	Use:   "ready",
	Short: "Show actionable issues (no blockers)",
	RunE: func(cmd *cobra.Command, args []string) error {
		assignee, _ := cmd.Flags().GetString("assignee")
		sortBy, _ := cmd.Flags().GetString("sort")
		limit, _ := cmd.Flags().GetInt("limit")

		s, err := store.Open(resolveVaultDir())
		if err != nil {
			return err
		}

		all, err := s.ListIssues(store.FilterOptions{})
		if err != nil {
			return err
		}

		g := graph.Build(all)
		ready := g.Ready()

		// Filter by assignee if specified.
		if assignee != "" {
			var tmp []*model.Issue
			for _, r := range ready {
				if strings.EqualFold(r.Assignee, assignee) {
					tmp = append(tmp, r)
				}
			}
			ready = tmp
		}

		sortReady(ready, sortBy)

		if limit > 0 && len(ready) > limit {
			ready = ready[:limit]
		}

		if jsonOut {
			return format.JSON(os.Stdout, ready)
		}
		format.Table(os.Stdout, ready)
		return nil
	},
}

func sortReady(issues []*model.Issue, sortBy string) {
	less := func(a, b *model.Issue) bool {
		switch sortBy {
		case "created":
			return a.CreatedAt.Before(b.CreatedAt)
		case "updated":
			return a.UpdatedAt.After(b.UpdatedAt)
		case "id":
			return a.ID < b.ID
		default: // priority
			return a.Priority < b.Priority
		}
	}
	for i := 1; i < len(issues); i++ {
		for j := i; j > 0 && less(issues[j], issues[j-1]); j-- {
			issues[j], issues[j-1] = issues[j-1], issues[j]
		}
	}
}

func init() {
	readyCmd.Flags().String("assignee", "", "filter by assignee")
	readyCmd.Flags().String("sort", "priority", "sort by: priority, created, updated, id")
	readyCmd.Flags().IntP("limit", "n", 0, "max results")
	rootCmd.AddCommand(readyCmd)
}
