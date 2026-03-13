package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var version = "dev"

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
	Version:      version,
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

// resolveVaultDir returns the nd vault directory.
//
// In Paivot-managed repos, the live vault is shared across branches and
// worktrees under the repository's git common dir. Other repos fall back to the
// nearest local .vault.
func resolveVaultDir() string {
	if vaultDir != "" {
		return vaultDir
	}
	if override := strings.TrimSpace(os.Getenv("ND_VAULT_DIR")); override != "" {
		return filepath.Clean(override)
	}

	dir, _ := os.Getwd()
	if path, err := resolvePaivotVaultDir(dir); err == nil {
		return path
	}

	if path, err := resolveLocalVaultDir(dir); err == nil {
		return path
	}

	return ".vault"
}

func parentDir(s string) string {
	parent := filepath.Dir(s)
	if parent == "." {
		return s
	}
	return parent
}

func errorf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "nd: "+format+"\n", args...)
}

func resolvePaivotVaultDir(start string) (string, error) {
	if !isPaivotManaged(start) {
		return "", fmt.Errorf("not a paivot-managed repo")
	}

	commonDir, err := gitCommonDir(start)
	if err != nil {
		return "", err
	}

	return filepath.Join(commonDir, "paivot", "nd-vault"), nil
}

func isPaivotManaged(start string) bool {
	dir := filepath.Clean(start)
	for {
		candidates := []string{
			filepath.Join(dir, ".vault", "knowledge", ".settings.yaml"),
			filepath.Join(dir, ".vault", "knowledge"),
			filepath.Join(dir, ".vault", ".dispatcher-state.json"),
			filepath.Join(dir, ".vault", ".piv-loop-state.json"),
		}
		for _, candidate := range candidates {
			if _, err := os.Stat(candidate); err == nil {
				return true
			}
		}

		parent := parentDir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return false
}

func resolveLocalVaultDir(start string) (string, error) {
	dir := filepath.Clean(start)
	for {
		candidate := filepath.Join(dir, ".vault")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, nil
		}
		parent := parentDir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("no .vault found")
}

func gitCommonDir(start string) (string, error) {
	repoRoot, err := findRepoRoot(start)
	if err != nil {
		return "", err
	}

	gitPath := filepath.Join(repoRoot, ".git")
	info, err := os.Stat(gitPath)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return filepath.Clean(gitPath), nil
	}

	data, err := os.ReadFile(gitPath)
	if err != nil {
		return "", err
	}

	line := strings.TrimSpace(string(data))
	const prefix = "gitdir:"
	if !strings.HasPrefix(line, prefix) {
		return "", fmt.Errorf("%s does not contain a gitdir pointer", gitPath)
	}

	gitDir := strings.TrimSpace(strings.TrimPrefix(line, prefix))
	if !filepath.IsAbs(gitDir) {
		gitDir = filepath.Join(repoRoot, gitDir)
	}
	gitDir = filepath.Clean(gitDir)

	commonDirPath := filepath.Join(gitDir, "commondir")
	if data, err := os.ReadFile(commonDirPath); err == nil {
		commonDir := strings.TrimSpace(string(data))
		if commonDir != "" {
			if !filepath.IsAbs(commonDir) {
				commonDir = filepath.Join(gitDir, commonDir)
			}
			return filepath.Clean(commonDir), nil
		}
	}

	return gitDir, nil
}

func findRepoRoot(start string) (string, error) {
	dir := filepath.Clean(start)
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parent := parentDir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("no git repo found")
}
