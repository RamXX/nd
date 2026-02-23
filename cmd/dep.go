package cmd

import (
	"fmt"
	"os"

	"github.com/RamXX/nd/internal/format"
	"github.com/RamXX/nd/internal/store"
	"github.com/spf13/cobra"
)

var depCmd = &cobra.Command{
	Use:   "dep",
	Short: "Manage dependencies between issues",
}

var depAddCmd = &cobra.Command{
	Use:   "add <issue> <depends-on>",
	Short: "Add dependency (issue depends on depends-on)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		issueID, depID := args[0], args[1]
		s, err := store.Open(resolveVaultDir())
		if err != nil {
			return err
		}
		if err := s.AddDependency(issueID, depID); err != nil {
			return err
		}
		if !quiet {
			fmt.Printf("%s now depends on %s\n", issueID, depID)
		}
		return nil
	},
}

var depRmCmd = &cobra.Command{
	Use:   "rm <issue> <depends-on>",
	Short: "Remove dependency",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		issueID, depID := args[0], args[1]
		s, err := store.Open(resolveVaultDir())
		if err != nil {
			return err
		}
		if err := s.RemoveDependency(issueID, depID); err != nil {
			return err
		}
		if !quiet {
			fmt.Printf("Removed dependency: %s no longer depends on %s\n", issueID, depID)
		}
		return nil
	},
}

var depListCmd = &cobra.Command{
	Use:   "list <id>",
	Short: "List dependencies of an issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		s, err := store.Open(resolveVaultDir())
		if err != nil {
			return err
		}
		issue, err := s.ReadIssue(id)
		if err != nil {
			return err
		}

		if len(issue.BlockedBy) > 0 {
			fmt.Printf("%s depends on:\n", id)
			for _, depID := range issue.BlockedBy {
				dep, err := s.ReadIssue(depID)
				if err != nil {
					fmt.Printf("  %s (unreadable)\n", depID)
					continue
				}
				format.Short(os.Stdout, dep)
			}
		}
		if len(issue.Blocks) > 0 {
			fmt.Printf("%s blocks:\n", id)
			for _, bID := range issue.Blocks {
				b, err := s.ReadIssue(bID)
				if err != nil {
					fmt.Printf("  %s (unreadable)\n", bID)
					continue
				}
				format.Short(os.Stdout, b)
			}
		}
		if len(issue.BlockedBy) == 0 && len(issue.Blocks) == 0 {
			fmt.Printf("%s has no dependencies\n", id)
		}
		return nil
	},
}

func init() {
	depCmd.AddCommand(depAddCmd)
	depCmd.AddCommand(depRmCmd)
	depCmd.AddCommand(depListCmd)
	rootCmd.AddCommand(depCmd)
}
