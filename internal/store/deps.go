package store

import (
	"fmt"
	"strings"
)

// AddDependency adds a dependency: issue depends on depID (depID blocks issue).
// Updates both sides: adds depID to issue's blocked_by, and issue to depID's blocks.
func (s *Store) AddDependency(issueID, depID string) error {
	if issueID == depID {
		return fmt.Errorf("an issue cannot depend on itself")
	}

	// Read both to validate they exist.
	issue, err := s.ReadIssue(issueID)
	if err != nil {
		return fmt.Errorf("issue %s: %w", issueID, err)
	}
	dep, err := s.ReadIssue(depID)
	if err != nil {
		return fmt.Errorf("dependency %s: %w", depID, err)
	}

	// Update issue's blocked_by if not already present.
	changed := false
	if !contains(issue.BlockedBy, depID) {
		newList := append(issue.BlockedBy, depID)
		if err := s.setListProperty(issueID, "blocked_by", newList); err != nil {
			return err
		}
		changed = true
	}

	// Update dep's blocks if not already present.
	if !contains(dep.Blocks, issueID) {
		newList := append(dep.Blocks, issueID)
		if err := s.setListProperty(depID, "blocks", newList); err != nil {
			return err
		}
		changed = true
	}

	if changed {
		_ = s.UpdateLinksSection(issueID)
		_ = s.UpdateLinksSection(depID)
		_ = s.appendHistory(issueID, fmt.Sprintf("dep_added: blocked_by %s", depID))
		_ = s.appendHistory(depID, fmt.Sprintf("dep_added: blocks %s", issueID))
	}

	return nil
}

// RemoveDependency removes a dependency between two issues.
// The relationship is preserved in was_blocked_by for historical record.
func (s *Store) RemoveDependency(issueID, depID string) error {
	issue, err := s.ReadIssue(issueID)
	if err != nil {
		return fmt.Errorf("issue %s: %w", issueID, err)
	}
	dep, err := s.ReadIssue(depID)
	if err != nil {
		return fmt.Errorf("dependency %s: %w", depID, err)
	}

	// Archive: add depID to issue's was_blocked_by if it was actually blocking.
	if contains(issue.BlockedBy, depID) && !contains(issue.WasBlockedBy, depID) {
		newWas := append(issue.WasBlockedBy, depID)
		if err := s.setListProperty(issueID, "was_blocked_by", newWas); err != nil {
			return err
		}
	}

	// Remove depID from issue's blocked_by.
	newBlockedBy := remove(issue.BlockedBy, depID)
	if err := s.setListProperty(issueID, "blocked_by", newBlockedBy); err != nil {
		return err
	}

	// Remove issueID from dep's blocks.
	newBlocks := remove(dep.Blocks, issueID)
	if err := s.setListProperty(depID, "blocks", newBlocks); err != nil {
		return err
	}

	// Update Links sections for both sides.
	_ = s.UpdateLinksSection(issueID)
	_ = s.UpdateLinksSection(depID)

	_ = s.appendHistory(issueID, fmt.Sprintf("dep_removed: was_blocked_by %s", depID))
	_ = s.appendHistory(depID, fmt.Sprintf("dep_removed: no_longer_blocks %s", issueID))

	return nil
}

func (s *Store) setListProperty(id, key string, vals []string) error {
	if len(vals) == 0 {
		return s.vault.PropertyRemove(id, key)
	}
	value := fmt.Sprintf("[%s]", strings.Join(vals, ", "))
	return s.vault.PropertySet(id, key, value)
}

func contains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

// AddRelated adds a bidirectional related link between two issues.
func (s *Store) AddRelated(issueID, relatedID string) error {
	if issueID == relatedID {
		return fmt.Errorf("an issue cannot relate to itself")
	}

	issue, err := s.ReadIssue(issueID)
	if err != nil {
		return fmt.Errorf("issue %s: %w", issueID, err)
	}
	rel, err := s.ReadIssue(relatedID)
	if err != nil {
		return fmt.Errorf("related %s: %w", relatedID, err)
	}

	changed := false
	if !contains(issue.Related, relatedID) {
		newList := append(issue.Related, relatedID)
		if err := s.setListProperty(issueID, "related", newList); err != nil {
			return err
		}
		changed = true
	}

	if !contains(rel.Related, issueID) {
		newList := append(rel.Related, issueID)
		if err := s.setListProperty(relatedID, "related", newList); err != nil {
			return err
		}
		changed = true
	}

	if changed {
		_ = s.UpdateLinksSection(issueID)
		_ = s.UpdateLinksSection(relatedID)
	}

	return nil
}

// RemoveRelated removes a bidirectional related link between two issues.
func (s *Store) RemoveRelated(issueID, relatedID string) error {
	issue, err := s.ReadIssue(issueID)
	if err != nil {
		return fmt.Errorf("issue %s: %w", issueID, err)
	}
	rel, err := s.ReadIssue(relatedID)
	if err != nil {
		return fmt.Errorf("related %s: %w", relatedID, err)
	}

	newRelA := remove(issue.Related, relatedID)
	if err := s.setListProperty(issueID, "related", newRelA); err != nil {
		return err
	}

	newRelB := remove(rel.Related, issueID)
	if err := s.setListProperty(relatedID, "related", newRelB); err != nil {
		return err
	}

	_ = s.UpdateLinksSection(issueID)
	_ = s.UpdateLinksSection(relatedID)

	return nil
}

// ResolveDependentsOf removes the closed issue from the blocked_by lists of all
// issues it blocks. Returns the list of issue IDs that were unblocked.
// Individual removal errors are logged but do not fail the cascade.
func (s *Store) ResolveDependentsOf(id string) ([]string, error) {
	issue, err := s.ReadIssue(id)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", id, err)
	}

	var unblocked []string
	for _, blockedID := range issue.Blocks {
		if err := s.RemoveDependency(blockedID, id); err != nil {
			// Log but don't fail -- best-effort cascade.
			continue
		}
		unblocked = append(unblocked, blockedID)
	}
	return unblocked, nil
}

func remove(ss []string, s string) []string {
	var result []string
	for _, v := range ss {
		if v != s {
			result = append(result, v)
		}
	}
	return result
}
