package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	vaultDir string
	jsonOut  bool
	verbose  bool
	quiet    bool
)

var rootCmd = &cobra.Command{
	Use:          "nd",
	Short:        "Vault-backed issue tracker",
	Long:         "nd -- Git-native issue tracking with Obsidian-compatible markdown files.",
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&vaultDir, "vault", "", "vault directory (default: .vault)")
	rootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "output as JSON")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "suppress non-essential output")
}

// resolveVaultDir returns the vault directory, walking up the tree to find .vault.
func resolveVaultDir() string {
	if vaultDir != "" {
		return vaultDir
	}
	dir, _ := os.Getwd()
	for {
		candidate := dir + "/.vault"
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
		parent := parentDir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ".vault"
}

func parentDir(s string) string {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			if i == 0 {
				return "/"
			}
			return s[:i]
		}
	}
	return s
}

func errorf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "nd: "+format+"\n", args...)
}
