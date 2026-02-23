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
	if !contains(issue.BlockedBy, depID) {
		newList := append(issue.BlockedBy, depID)
		if err := s.setListProperty(issueID, "blocked_by", newList); err != nil {
			return err
		}
	}

	// Update dep's blocks if not already present.
	if !contains(dep.Blocks, issueID) {
		newList := append(dep.Blocks, issueID)
		if err := s.setListProperty(depID, "blocks", newList); err != nil {
			return err
		}
	}

	return nil
}

// RemoveDependency removes a dependency between two issues.
func (s *Store) RemoveDependency(issueID, depID string) error {
	issue, err := s.ReadIssue(issueID)
	if err != nil {
		return fmt.Errorf("issue %s: %w", issueID, err)
	}
	dep, err := s.ReadIssue(depID)
	if err != nil {
		return fmt.Errorf("dependency %s: %w", depID, err)
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

func remove(ss []string, s string) []string {
	var result []string
	for _, v := range ss {
		if v != s {
			result = append(result, v)
		}
	}
	return result
}
