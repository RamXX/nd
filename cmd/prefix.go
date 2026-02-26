package cmd

import (
	"bufio"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// prefixFromName derives a short uppercase prefix from a project name.
//   - Split on hyphens/underscores, take first letter of each word: "my-project" -> "MP"
//   - Single word: take first 2-3 chars: "tminus" -> "TM", "nd" -> "ND"
//   - Capped at 4 characters.
func prefixFromName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}

	// Split on hyphens, underscores, dots.
	parts := regexp.MustCompile(`[-_.]`).Split(name, -1)
	// Filter empty parts.
	var words []string
	for _, p := range parts {
		if p != "" {
			words = append(words, p)
		}
	}

	var prefix string
	if len(words) > 1 {
		// Multi-word: take first letter of each word.
		for _, w := range words {
			prefix += string(w[0])
		}
	} else {
		// Single word: take first 2-3 chars depending on length.
		w := words[0]
		switch {
		case len(w) <= 2:
			prefix = w
		case len(w) <= 4:
			prefix = w[:2]
		default:
			prefix = w[:3]
		}
	}

	// Cap at 4 characters and uppercase.
	if len(prefix) > 4 {
		prefix = prefix[:4]
	}
	return strings.ToUpper(prefix)
}

// inferPrefixFromGitRemote parses the git remote origin URL and derives a prefix.
// Returns (prefix, source) where source describes where the prefix came from.
func inferPrefixFromGitRemote() (string, string) {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return "", ""
	}
	url := strings.TrimSpace(string(out))
	if url == "" {
		return "", ""
	}

	// Extract repo name from URL patterns:
	//   git@github.com:User/repo-name.git
	//   https://github.com/User/repo-name.git
	name := url
	// Strip trailing .git
	name = strings.TrimSuffix(name, ".git")
	// Take last path component.
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		name = name[idx+1:]
	} else if idx := strings.LastIndex(name, ":"); idx >= 0 {
		name = name[idx+1:]
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return "", ""
	}

	prefix := prefixFromName(name)
	if prefix == "" {
		return "", ""
	}
	return prefix, "git remote \"" + name + "\""
}

// inferPrefixFromDir derives a prefix from the current directory basename.
func inferPrefixFromDir() (string, string) {
	dir, err := os.Getwd()
	if err != nil {
		return "", ""
	}
	name := filepath.Base(dir)
	if name == "" || name == "." || name == "/" {
		return "", ""
	}
	prefix := prefixFromName(name)
	if prefix == "" {
		return "", ""
	}
	return prefix, "directory \"" + name + "\""
}

// inferPrefix tries git remote first, then falls back to directory name.
func inferPrefix() (string, string) {
	if prefix, source := inferPrefixFromGitRemote(); prefix != "" {
		return prefix, source
	}
	return inferPrefixFromDir()
}

// inferPrefixFromBeadsID extracts the prefix from a beads issue ID like "TM-a3f8" or "PROJ-abc.3".
func inferPrefixFromBeadsID(id string) string {
	idx := strings.Index(id, "-")
	if idx <= 0 {
		return ""
	}
	return strings.ToUpper(id[:idx])
}

// peekBeadsPrefix reads the first valid JSONL line with an "id" field and extracts the prefix.
func peekBeadsPrefix(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var raw map[string]any
		if json.Unmarshal([]byte(line), &raw) != nil {
			continue
		}
		if id, ok := raw["id"].(string); ok && id != "" {
			if prefix := inferPrefixFromBeadsID(id); prefix != "" {
				return prefix
			}
		}
	}
	return ""
}
