package cmd

import (
	"fmt"
	"strings"

	"github.com/RamXX/nd/internal/store"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [title]",
	Short: "Create a new issue",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := strings.Join(args, " ")
		issueType, _ := cmd.Flags().GetString("type")
		priority, _ := cmd.Flags().GetInt("priority")
		assignee, _ := cmd.Flags().GetString("assignee")
		labelsStr, _ := cmd.Flags().GetString("labels")
		description, _ := cmd.Flags().GetString("description")
		parent, _ := cmd.Flags().GetString("parent")

		if issueType == "" {
			issueType = "task"
		}

		var labels []string
		if labelsStr != "" {
			for _, l := range strings.Split(labelsStr, ",") {
				l = strings.TrimSpace(l)
				if l != "" {
					labels = append(labels, l)
				}
			}
		}

		s, err := store.Open(resolveVaultDir())
		if err != nil {
			return err
		}

		issue, err := s.CreateIssue(title, description, issueType, priority, assignee, labels, parent)
		if err != nil {
			return err
		}

		if jsonOut {
			fmt.Printf(`{"id":"%s"}`, issue.ID)
			fmt.Println()
		} else if !quiet {
			fmt.Printf("Created %s: %s\n", issue.ID, issue.Title)
		} else {
			fmt.Println(issue.ID)
		}
		return nil
	},
}

func init() {
	createCmd.Flags().StringP("type", "t", "task", "issue type (bug, feature, task, epic, chore, decision)")
	createCmd.Flags().IntP("priority", "p", 2, "priority 0-4 (0=critical)")
	createCmd.Flags().String("assignee", "", "assignee")
	createCmd.Flags().String("labels", "", "comma-separated labels")
	createCmd.Flags().StringP("description", "d", "", "issue description")
	createCmd.Flags().String("parent", "", "parent issue ID")
	rootCmd.AddCommand(createCmd)
}
