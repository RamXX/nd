package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/RamXX/nd/internal/store"
	"github.com/RamXX/vlt"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import issues from beads JSONL",
	RunE: func(cmd *cobra.Command, args []string) error {
		fromBeads, _ := cmd.Flags().GetString("from-beads")
		if fromBeads == "" {
			return fmt.Errorf("--from-beads is required")
		}

		s, err := store.Open(resolveVaultDir())
		if err != nil {
			return err
		}

		f, err := os.Open(fromBeads)
		if err != nil {
			return fmt.Errorf("open %s: %w", fromBeads, err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		// Increase buffer for large lines.
		buf := make([]byte, 0, 256*1024)
		scanner.Buffer(buf, 1024*1024)

		imported, skipped := 0, 0
		for scanner.Scan() {
			line := scanner.Text()
			if strings.TrimSpace(line) == "" {
				continue
			}

			var raw map[string]any
			if err := json.Unmarshal([]byte(line), &raw); err != nil {
				skipped++
				continue
			}

			title, _ := raw["title"].(string)
			if title == "" {
				skipped++
				continue
			}

			// Extract fields with sensible defaults.
			desc, _ := raw["description"].(string)
			issueType := extractString(raw, "type", "task")
			priority := extractInt(raw, "priority", 2)
			assignee, _ := raw["assignee"].(string)
			status, _ := raw["status"].(string)

			// Extract labels.
			var labels []string
			if labelsRaw, ok := raw["labels"]; ok {
				if arr, ok := labelsRaw.([]any); ok {
					for _, v := range arr {
						if s, ok := v.(string); ok {
							labels = append(labels, s)
						}
					}
				}
			}

			issue, err := s.CreateIssue(title, desc, issueType, priority, assignee, labels, "")
			if err != nil {
				if !quiet {
					errorf("skip %q: %v", title, err)
				}
				skipped++
				continue
			}

			// Update status if not open.
			if status != "" && status != "open" {
				switch status {
				case "closed":
					closedAt, _ := raw["closed_at"].(string)
					reason, _ := raw["close_reason"].(string)
					_ = s.CloseIssue(issue.ID, reason)
					if closedAt != "" {
						_ = s.UpdateField(issue.ID, "closed_at", closedAt)
					}
				case "in_progress":
					_ = s.UpdateField(issue.ID, "status", "in_progress")
				case "blocked":
					_ = s.UpdateField(issue.ID, "status", "blocked")
				case "deferred":
					_ = s.UpdateField(issue.ID, "status", "deferred")
				}
			}

			// Preserve original timestamps if available.
			if createdAt, ok := raw["created_at"].(string); ok && createdAt != "" {
				_ = s.UpdateField(issue.ID, "created_at", createdAt)
			}
			if updatedAt, ok := raw["updated_at"].(string); ok && updatedAt != "" {
				_ = s.UpdateField(issue.ID, "updated_at", updatedAt)
			}

			// Import notes.
			if notes, ok := raw["notes"].(string); ok && notes != "" {
				_ = s.AppendNotes(issue.ID, notes)
			}

			// Import design.
			if design, ok := raw["design"].(string); ok && design != "" {
				_ = s.Vault().Patch(issue.ID, vlt.PatchOptions{
					Heading: "## Design",
					Content: design + "\n",
				})
			}

			imported++
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("scan: %w", err)
		}

		fmt.Printf("Imported %d issues (%d skipped)\n", imported, skipped)
		return nil
	},
}

func extractString(m map[string]any, key, fallback string) string {
	if v, ok := m[key].(string); ok && v != "" {
		return v
	}
	return fallback
}

func extractInt(m map[string]any, key string, fallback int) int {
	switch v := m[key].(type) {
	case float64:
		return int(v)
	case string:
		// Try parsing "P2" -> 2
		v = strings.TrimPrefix(strings.ToUpper(v), "P")
		for i := 0; i <= 4; i++ {
			if v == fmt.Sprintf("%d", i) {
				return i
			}
		}
	}
	return fallback
}


func init() {
	importCmd.Flags().String("from-beads", "", "path to beads JSONL file")
	rootCmd.AddCommand(importCmd)
}
