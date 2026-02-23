package store

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/RamXX/nd/internal/model"
)

// FilterOptions controls which issues ListIssues returns.
type FilterOptions struct {
	Status   string
	Type     string
	Assignee string
	Label    string
	Sort     string // "priority", "created", "updated", "id" (default)
	Limit    int
}

// ListIssues reads all issues from the vault and applies filters.
func (s *Store) ListIssues(opts FilterOptions) ([]*model.Issue, error) {
	files, err := s.vault.Files("issues", "md")
	if err != nil {
		return nil, err
	}

	type result struct {
		issue *model.Issue
		err   error
	}

	results := make([]result, len(files))
	var wg sync.WaitGroup

	for i, f := range files {
		wg.Add(1)
		go func(idx int, file string) {
			defer wg.Done()
			// Extract ID from filename (strip issues/ prefix and .md suffix).
			base := filepath.Base(file)
			id := strings.TrimSuffix(base, ".md")
			issue, err := s.ReadIssue(id)
			results[idx] = result{issue: issue, err: err}
		}(i, f)
	}
	wg.Wait()

	var issues []*model.Issue
	for _, r := range results {
		if r.err != nil {
			continue // skip unreadable issues
		}
		if matchesFilter(r.issue, opts) {
			issues = issues[:len(issues):len(issues)]
			issues = append(issues, r.issue)
		}
	}

	sortIssues(issues, opts.Sort)

	if opts.Limit > 0 && len(issues) > opts.Limit {
		issues = issues[:opts.Limit]
	}
	return issues, nil
}

func matchesFilter(issue *model.Issue, opts FilterOptions) bool {
	if opts.Status != "" {
		st, err := model.ParseStatus(opts.Status)
		if err != nil {
			return false
		}
		if issue.Status != st {
			return false
		}
	}
	if opts.Type != "" {
		t, err := model.ParseIssueType(opts.Type)
		if err != nil {
			return false
		}
		if issue.Type != t {
			return false
		}
	}
	if opts.Assignee != "" && !strings.EqualFold(issue.Assignee, opts.Assignee) {
		return false
	}
	if opts.Label != "" {
		found := false
		for _, l := range issue.Labels {
			if strings.EqualFold(l, opts.Label) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func sortIssues(issues []*model.Issue, sortBy string) {
	switch sortBy {
	case "priority":
		sortByFunc(issues, func(a, b *model.Issue) bool {
			return a.Priority < b.Priority
		})
	case "created":
		sortByFunc(issues, func(a, b *model.Issue) bool {
			return a.CreatedAt.Before(b.CreatedAt)
		})
	case "updated":
		sortByFunc(issues, func(a, b *model.Issue) bool {
			return a.UpdatedAt.After(b.UpdatedAt)
		})
	default: // "id"
		sortByFunc(issues, func(a, b *model.Issue) bool {
			return a.ID < b.ID
		})
	}
}

func sortByFunc(issues []*model.Issue, less func(a, b *model.Issue) bool) {
	// Simple insertion sort -- issue counts are small.
	for i := 1; i < len(issues); i++ {
		for j := i; j > 0 && less(issues[j], issues[j-1]); j-- {
			issues[j], issues[j-1] = issues[j-1], issues[j]
		}
	}
}
