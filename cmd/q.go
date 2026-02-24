package cmd

import (
	"fmt"
	"strings"

	"github.com/RamXX/nd/internal/store"
	"github.com/spf13/cobra"
)

var qCmd = &cobra.Command{
	Use:   "q [title]",
	Short: "Quick capture -- create issue, print only ID",
	Long:  "Create an issue and print only its ID to stdout, enabling ISSUE=$(nd q \"title\").",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := strings.Join(args, " ")
		issueType, _ := cmd.Flags().GetString("type")
		priority, _ := cmd.Flags().GetInt("priority")
		assignee, _ := cmd.Flags().GetString("assignee")
		labelsStr, _ := cmd.Flags().GetString("labels")
		description, _ := cmd.Flags().GetString("description")
		parent, _ := cmd.Flags().GetString("parent")
		bodyFile, _ := cmd.Flags().GetString("body-file")

		if bodyFile != "" {
			body, err := readBodyFile(bodyFile)
			if err != nil {
				return err
			}
			description = body
		}

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
		} else {
			fmt.Println(issue.ID)
		}
		return nil
	},
}

func init() {
	qCmd.Flags().StringP("type", "t", "task", "issue type")
	qCmd.Flags().IntP("priority", "p", 2, "priority 0-4")
	qCmd.Flags().String("assignee", "", "assignee")
	qCmd.Flags().String("labels", "", "comma-separated labels")
	qCmd.Flags().StringP("description", "d", "", "issue description")
	qCmd.Flags().String("parent", "", "parent issue ID")
	qCmd.Flags().String("body-file", "", "read description from file (- for stdin)")
	rootCmd.AddCommand(qCmd)
}
