package cmd

import (
	"fmt"

	"github.com/RamXX/nd/internal/model"
	"github.com/RamXX/nd/internal/store"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update issue fields",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		s, err := store.Open(resolveVaultDir())
		if err != nil {
			return err
		}

		// Verify issue exists.
		if _, err := s.ReadIssue(id); err != nil {
			return fmt.Errorf("issue %s not found: %w", id, err)
		}

		changed := false

		if cmd.Flags().Changed("status") {
			v, _ := cmd.Flags().GetString("status")
			st, err := model.ParseStatus(v)
			if err != nil {
				return err
			}
			if err := s.UpdateStatus(id, st); err != nil {
				return err
			}
			changed = true
		}

		if cmd.Flags().Changed("title") {
			v, _ := cmd.Flags().GetString("title")
			if err := s.UpdateField(id, "title", fmt.Sprintf("%q", v)); err != nil {
				return err
			}
			changed = true
		}

		if cmd.Flags().Changed("priority") {
			v, _ := cmd.Flags().GetString("priority")
			p, err := model.ParsePriority(v)
			if err != nil {
				return err
			}
			if err := s.UpdateField(id, "priority", fmt.Sprintf("%d", p)); err != nil {
				return err
			}
			changed = true
		}

		if cmd.Flags().Changed("assignee") {
			v, _ := cmd.Flags().GetString("assignee")
			if err := s.UpdateField(id, "assignee", v); err != nil {
				return err
			}
			changed = true
		}

		if cmd.Flags().Changed("type") {
			v, _ := cmd.Flags().GetString("type")
			if _, err := model.ParseIssueType(v); err != nil {
				return err
			}
			if err := s.UpdateField(id, "type", v); err != nil {
				return err
			}
			changed = true
		}

		if cmd.Flags().Changed("append-notes") {
			v, _ := cmd.Flags().GetString("append-notes")
			if err := s.AppendNotes(id, v); err != nil {
				return err
			}
			changed = true
		}

		if cmd.Flags().Changed("description") {
			v, _ := cmd.Flags().GetString("description")
			if err := s.UpdateField(id, "description", v); err != nil {
				return err
			}
			changed = true
		}

		if !changed {
			return fmt.Errorf("no fields specified to update")
		}

		if !quiet {
			fmt.Printf("Updated %s\n", id)
		}
		return nil
	},
}

func init() {
	updateCmd.Flags().String("status", "", "new status")
	updateCmd.Flags().String("title", "", "new title")
	updateCmd.Flags().String("priority", "", "new priority (0-4 or P0-P4)")
	updateCmd.Flags().String("assignee", "", "new assignee")
	updateCmd.Flags().String("type", "", "new type")
	updateCmd.Flags().String("append-notes", "", "append text to Notes section")
	updateCmd.Flags().StringP("description", "d", "", "new description")
	rootCmd.AddCommand(updateCmd)
}
