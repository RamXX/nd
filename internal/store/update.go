package store

import (
	"fmt"
	"strings"
	"time"

	"github.com/RamXX/nd/internal/enforce"
	"github.com/RamXX/nd/internal/model"
	"github.com/RamXX/vlt"
)

// UpdateField updates a single frontmatter field on an issue.
func (s *Store) UpdateField(id, field, value string) error {
	if err := s.vault.PropertySet(id, field, value); err != nil {
		return fmt.Errorf("set %s on %s: %w", field, id, err)
	}
	return s.touchUpdatedAt(id)
}

// UpdateStatus changes the status of an issue with validation.
func (s *Store) UpdateStatus(id string, newStatus model.Status) error {
	issue, err := s.ReadIssue(id)
	if err != nil {
		return err
	}

	// Validate transition.
	if issue.Status == model.StatusClosed && newStatus != model.StatusOpen {
		return fmt.Errorf("closed issues can only be reopened (set to open)")
	}

	if s.config.StatusFSM {
		if err := s.validateFSMTransition(issue.Status, newStatus); err != nil {
			return err
		}
	}

	if err := s.vault.PropertySet(id, "status", string(newStatus)); err != nil {
		return err
	}
	return s.touchUpdatedAt(id)
}

// CloseIssue closes an issue with an optional reason.
func (s *Store) CloseIssue(id, reason string) error {
	issue, err := s.ReadIssue(id)
	if err != nil {
		return err
	}
	if issue.Status == model.StatusClosed {
		return fmt.Errorf("issue %s is already closed", id)
	}

	if s.config.StatusFSM {
		if err := s.validateFSMTransition(issue.Status, model.StatusClosed); err != nil {
			return err
		}
	}

	now := time.Now().UTC().Format(time.RFC3339)
	if err := s.vault.PropertySet(id, "status", "closed"); err != nil {
		return err
	}
	if err := s.vault.PropertySet(id, "closed_at", now); err != nil {
		return err
	}
	if reason != "" {
		if err := s.vault.PropertySet(id, "close_reason", fmt.Sprintf("%q", reason)); err != nil {
			return err
		}
	}
	return s.touchUpdatedAt(id)
}

// ReopenIssue changes a closed issue back to open.
func (s *Store) ReopenIssue(id string) error {
	issue, err := s.ReadIssue(id)
	if err != nil {
		return err
	}
	if issue.Status != model.StatusClosed {
		return fmt.Errorf("issue %s is not closed (status: %s)", id, issue.Status)
	}

	if err := s.vault.PropertySet(id, "status", "open"); err != nil {
		return err
	}
	// Clear closed_at and close_reason.
	_ = s.vault.PropertyRemove(id, "closed_at")
	_ = s.vault.PropertyRemove(id, "close_reason")
	return s.touchUpdatedAt(id)
}

// AppendNotes appends text to the Notes section.
func (s *Store) AppendNotes(id, content string) error {
	return s.vault.Patch(id, vlt.PatchOptions{
		Heading:    "## Notes",
		Content:    content + "\n",
		Timestamps: false,
	})
}

// UpdateBody replaces the body and recalculates the content hash.
func (s *Store) UpdateBody(id, body string) error {
	if err := s.vault.Write(id, body, false); err != nil {
		return err
	}
	hash := enforce.ComputeContentHash(body)
	if err := s.vault.PropertySet(id, "content_hash", fmt.Sprintf("%q", hash)); err != nil {
		return err
	}
	return s.touchUpdatedAt(id)
}

// UpdateLinksSection rebuilds the ## Links section from frontmatter relationships.
func (s *Store) UpdateLinksSection(id string) error {
	issue, err := s.ReadIssue(id)
	if err != nil {
		return err
	}

	content := buildLinksSection(issue)

	// Check if body already has a ## Links section.
	if !strings.Contains(issue.Body, "\n## Links\n") {
		// Insert ## Links before ## Comments.
		if idx := strings.Index(issue.Body, "\n## Comments\n"); idx >= 0 {
			newBody := issue.Body[:idx] + "\n## Links\n\n" + issue.Body[idx:]
			if err := s.vault.Write(id, newBody, false); err != nil {
				return err
			}
		}
	}

	// Patch the Links section content.
	if err := s.vault.Patch(id, vlt.PatchOptions{
		Heading:    "## Links",
		Content:    content,
		Timestamps: false,
	}); err != nil {
		return err
	}

	// Recompute content hash after Links section update.
	updated, err := s.ReadIssue(id)
	if err != nil {
		return err
	}
	hash := enforce.ComputeContentHash(updated.Body)
	return s.vault.PropertySet(id, "content_hash", fmt.Sprintf("%q", hash))
}

// SetParent sets the parent of an issue and updates the Links section.
func (s *Store) SetParent(id, parentID string) error {
	if parentID == "" {
		if err := s.vault.PropertyRemove(id, "parent"); err != nil {
			return err
		}
	} else {
		if err := s.vault.PropertySet(id, "parent", parentID); err != nil {
			return err
		}
	}
	if err := s.touchUpdatedAt(id); err != nil {
		return err
	}
	return s.UpdateLinksSection(id)
}

// RefreshAfterEdit recomputes the content hash and updates the Links section
// after a manual edit. Call this after an external editor modifies the file.
func (s *Store) RefreshAfterEdit(id string) error {
	if err := s.UpdateLinksSection(id); err != nil {
		return err
	}
	return s.touchUpdatedAt(id)
}

// DeferIssue sets the issue status to deferred with an optional until date.
func (s *Store) DeferIssue(id, until string) error {
	issue, err := s.ReadIssue(id)
	if err != nil {
		return err
	}
	if issue.Status == model.StatusClosed {
		return fmt.Errorf("cannot defer closed issue %s", id)
	}

	if err := s.vault.PropertySet(id, "status", "deferred"); err != nil {
		return err
	}
	if until != "" {
		if err := s.vault.PropertySet(id, "defer_until", until); err != nil {
			return err
		}
	}
	return s.touchUpdatedAt(id)
}

// UnDeferIssue restores a deferred issue to open.
func (s *Store) UnDeferIssue(id string) error {
	issue, err := s.ReadIssue(id)
	if err != nil {
		return err
	}
	if issue.Status != model.StatusDeferred {
		return fmt.Errorf("issue %s is not deferred (status: %s)", id, issue.Status)
	}

	if err := s.vault.PropertySet(id, "status", "open"); err != nil {
		return err
	}
	_ = s.vault.PropertyRemove(id, "defer_until")
	return s.touchUpdatedAt(id)
}

func (s *Store) touchUpdatedAt(id string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	return s.vault.PropertySet(id, "updated_at", now)
}

// validateFSMTransition enforces the FSM transition rules.
// The engine is generic -- all behavior is driven by configuration:
//   - status.sequence: forward +1 only, backward any
//   - status.exit_rules: restrict exits from specific statuses to listed targets
//   - Off-sequence statuses are unrestricted (escape hatch)
func (s *Store) validateFSMTransition(from, to model.Status) error {
	// Check exit rules first -- these override sequence logic.
	exitRules := s.ExitRules()
	if allowed, ok := exitRules[from]; ok {
		for _, a := range allowed {
			if a == to {
				return nil
			}
		}
		targets := make([]string, len(allowed))
		for i, a := range allowed {
			targets[i] = string(a)
		}
		return fmt.Errorf("FSM: cannot transition from %s to %s; allowed targets: %s",
			from, to, strings.Join(targets, ", "))
	}

	seq := s.StatusSequence()
	if len(seq) == 0 {
		return nil
	}

	fromIdx := indexInSequence(seq, from)
	toIdx := indexInSequence(seq, to)

	// Both in sequence: forward must be +1, backward is always allowed.
	if fromIdx >= 0 && toIdx >= 0 {
		if toIdx > fromIdx {
			if toIdx != fromIdx+1 {
				return fmt.Errorf("FSM: cannot skip from %s to %s; next step is %s", from, to, seq[fromIdx+1])
			}
		}
		return nil
	}

	// One or both off-sequence: allow (escape hatch for custom statuses like rejected).
	return nil
}

func indexInSequence(seq []model.Status, st model.Status) int {
	for i, s := range seq {
		if s == st {
			return i
		}
	}
	return -1
}
