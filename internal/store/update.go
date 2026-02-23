package store

import (
	"fmt"
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

func (s *Store) touchUpdatedAt(id string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	return s.vault.PropertySet(id, "updated_at", now)
}
