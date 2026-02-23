package cmd

import (
	"fmt"

	"github.com/RamXX/nd/internal/store"
	"github.com/spf13/cobra"
)

var closeCmd = &cobra.Command{
	Use:   "close <id> [id...]",
	Short: "Close one or more issues",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		reason, _ := cmd.Flags().GetString("reason")

		s, err := store.Open(resolveVaultDir())
		if err != nil {
			return err
		}

		var errors []string
		for _, id := range args {
			if err := s.CloseIssue(id, reason); err != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", id, err))
				continue
			}
			if !quiet {
				fmt.Printf("Closed %s\n", id)
			}
		}

		if len(errors) > 0 {
			for _, e := range errors {
				errorf("%s", e)
			}
			return fmt.Errorf("%d issue(s) failed to close", len(errors))
		}
		return nil
	},
}

func init() {
	closeCmd.Flags().String("reason", "", "close reason")
	closeCmd.Flags().Bool("force", false, "close even if blocked")
	rootCmd.AddCommand(closeCmd)
}
