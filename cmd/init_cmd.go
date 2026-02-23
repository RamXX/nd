package cmd

import (
	"fmt"
	"os"
	"os/user"

	"github.com/RamXX/nd/internal/store"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new nd vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		prefix, _ := cmd.Flags().GetString("prefix")
		author, _ := cmd.Flags().GetString("author")
		dir := resolveVaultDir()

		if prefix == "" {
			return fmt.Errorf("--prefix is required")
		}
		if author == "" {
			u, err := user.Current()
			if err == nil {
				author = u.Username
			} else {
				author = "unknown"
			}
		}

		// Check if already initialized.
		if _, err := os.Stat(dir + "/.nd.yaml"); err == nil {
			return fmt.Errorf("vault already initialized at %s", dir)
		}

		s, err := store.Init(dir, prefix, author)
		if err != nil {
			return err
		}
		if !quiet {
			fmt.Printf("Initialized nd vault at %s (prefix: %s)\n", s.Dir(), prefix)
		}
		return nil
	},
}

func init() {
	initCmd.Flags().String("prefix", "", "issue ID prefix (required)")
	initCmd.Flags().String("author", "", "default author (defaults to OS user)")
	rootCmd.AddCommand(initCmd)
}
