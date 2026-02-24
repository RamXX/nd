package store

import (
	"os"
	"strings"
	"testing"

	"github.com/RamXX/nd/internal/model"
)

func TestInitAndOpen(t *testing.T) {
	dir := t.TempDir()

	// Init.
	s, err := Init(dir, "TST", "tester")
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	if s.Prefix() != "TST" {
		t.Errorf("prefix = %q, want TST", s.Prefix())
	}

	// Verify .nd.yaml exists.
	if _, err := os.Stat(dir + "/.nd.yaml"); err != nil {
		t.Fatalf(".nd.yaml missing: %v", err)
	}
	// Verify issues/ dir exists.
	if _, err := os.Stat(dir + "/issues"); err != nil {
		t.Fatalf("issues/ missing: %v", err)
	}

	// Reopen.
	s2, err := Open(dir)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if s2.Prefix() != "TST" {
		t.Errorf("reopened prefix = %q, want TST", s2.Prefix())
	}
}

func TestCreateAndReadIssue(t *testing.T) {
	dir := t.TempDir()
	s, err := Init(dir, "TST", "tester")
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	issue, err := s.CreateIssue("Fix the login bug", "Users can't login when password has special chars", "bug", 1, "alice", []string{"auth"}, "")
	if err != nil {
		t.Fatalf("CreateIssue: %v", err)
	}

	if issue.ID == "" {
		t.Fatal("issue ID is empty")
	}
	if issue.Title != "Fix the login bug" {
		t.Errorf("title = %q", issue.Title)
	}
	if issue.Status != "open" {
		t.Errorf("status = %q", issue.Status)
	}
	if issue.Priority != 1 {
		t.Errorf("priority = %d", issue.Priority)
	}

	// Read back.
	got, err := s.ReadIssue(issue.ID)
	if err != nil {
		t.Fatalf("ReadIssue: %v", err)
	}
	if got.ID != issue.ID {
		t.Errorf("read ID = %q, want %q", got.ID, issue.ID)
	}
	if got.Title != issue.Title {
		t.Errorf("read title = %q, want %q", got.Title, issue.Title)
	}
	if got.Assignee != "alice" {
		t.Errorf("assignee = %q, want alice", got.Assignee)
	}
	if len(got.Labels) != 1 || got.Labels[0] != "auth" {
		t.Errorf("labels = %v, want [auth]", got.Labels)
	}
	if got.ContentHash == "" {
		t.Error("content hash is empty")
	}

	// Verify file on disk.
	if _, err := os.Stat(dir + "/issues/" + issue.ID + ".md"); err != nil {
		t.Errorf("issue file missing: %v", err)
	}
}

func TestListIssues(t *testing.T) {
	dir := t.TempDir()
	s, err := Init(dir, "TST", "tester")
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	_, err = s.CreateIssue("Issue A", "", "task", 2, "", nil, "")
	if err != nil {
		t.Fatalf("create A: %v", err)
	}
	_, err = s.CreateIssue("Issue B", "", "bug", 0, "bob", nil, "")
	if err != nil {
		t.Fatalf("create B: %v", err)
	}

	// List all.
	all, err := s.ListIssues(FilterOptions{})
	if err != nil {
		t.Fatalf("ListIssues: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 issues, got %d", len(all))
	}

	// Filter by type.
	bugs, err := s.ListIssues(FilterOptions{Type: "bug"})
	if err != nil {
		t.Fatalf("ListIssues type=bug: %v", err)
	}
	if len(bugs) != 1 {
		t.Errorf("expected 1 bug, got %d", len(bugs))
	}

	// Filter by assignee.
	bobs, err := s.ListIssues(FilterOptions{Assignee: "bob"})
	if err != nil {
		t.Fatalf("ListIssues assignee=bob: %v", err)
	}
	if len(bobs) != 1 {
		t.Errorf("expected 1 assigned to bob, got %d", len(bobs))
	}
}

func TestIssueExists(t *testing.T) {
	dir := t.TempDir()
	s, err := Init(dir, "TST", "tester")
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	if s.IssueExists("TST-xxx") {
		t.Error("nonexistent issue should not exist")
	}

	issue, err := s.CreateIssue("Test", "", "task", 2, "", nil, "")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if !s.IssueExists(issue.ID) {
		t.Error("created issue should exist")
	}
}

func TestBuildLinksSection(t *testing.T) {
	tests := []struct {
		name   string
		issue  *model.Issue
		want   []string
		empty  bool
	}{
		{
			name:  "no relationships",
			issue: &model.Issue{ID: "TST-0001"},
			empty: true,
		},
		{
			name:  "parent only",
			issue: &model.Issue{ID: "TST-0001", Parent: "TST-epic"},
			want:  []string{"- Parent: [[TST-epic]]"},
		},
		{
			name:  "blocks only",
			issue: &model.Issue{ID: "TST-0001", Blocks: []string{"TST-b1", "TST-b2"}},
			want:  []string{"- Blocks: [[TST-b1]], [[TST-b2]]"},
		},
		{
			name:  "blocked_by only",
			issue: &model.Issue{ID: "TST-0001", BlockedBy: []string{"TST-c1"}},
			want:  []string{"- Blocked by: [[TST-c1]]"},
		},
		{
			name:  "related only",
			issue: &model.Issue{ID: "TST-0001", Related: []string{"TST-r1", "TST-r2"}},
			want:  []string{"- Related: [[TST-r1]], [[TST-r2]]"},
		},
		{
			name: "all relationships",
			issue: &model.Issue{
				ID:        "TST-0001",
				Parent:    "TST-epic",
				Blocks:    []string{"TST-b1"},
				BlockedBy: []string{"TST-c1"},
				Related:   []string{"TST-r1"},
			},
			want: []string{
				"- Parent: [[TST-epic]]",
				"- Blocks: [[TST-b1]]",
				"- Blocked by: [[TST-c1]]",
				"- Related: [[TST-r1]]",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildLinksSection(tt.issue)
			if tt.empty {
				if got != "" {
					t.Errorf("expected empty, got %q", got)
				}
				return
			}
			for _, w := range tt.want {
				if !strings.Contains(got, w) {
					t.Errorf("missing %q in:\n%s", w, got)
				}
			}
		})
	}
}

func TestUpdateLinksSection(t *testing.T) {
	dir := t.TempDir()
	s, err := Init(dir, "TST", "tester")
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	// Create two issues.
	a, err := s.CreateIssue("Issue A", "", "task", 2, "", nil, "")
	if err != nil {
		t.Fatalf("create A: %v", err)
	}
	b, err := s.CreateIssue("Issue B", "", "task", 2, "", nil, "")
	if err != nil {
		t.Fatalf("create B: %v", err)
	}

	// Add dependency: B depends on A.
	if err := s.AddDependency(b.ID, a.ID); err != nil {
		t.Fatalf("add dep: %v", err)
	}

	// Read B and check for wikilinks.
	bRead, err := s.ReadIssue(b.ID)
	if err != nil {
		t.Fatalf("read B: %v", err)
	}
	if !strings.Contains(bRead.Body, "[["+a.ID+"]]") {
		t.Errorf("B body should contain wikilink to A:\n%s", bRead.Body)
	}

	// Read A and check for wikilinks.
	aRead, err := s.ReadIssue(a.ID)
	if err != nil {
		t.Fatalf("read A: %v", err)
	}
	if !strings.Contains(aRead.Body, "[["+b.ID+"]]") {
		t.Errorf("A body should contain wikilink to B:\n%s", aRead.Body)
	}

	// Remove dependency and verify wikilinks disappear.
	if err := s.RemoveDependency(b.ID, a.ID); err != nil {
		t.Fatalf("remove dep: %v", err)
	}
	bAfter, _ := s.ReadIssue(b.ID)
	if strings.Contains(bAfter.Body, "[["+a.ID+"]]") {
		t.Errorf("B body should not contain wikilink to A after removal:\n%s", bAfter.Body)
	}
}

func TestListFilterParent(t *testing.T) {
	dir := t.TempDir()
	s, err := Init(dir, "TST", "tester")
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	epic, err := s.CreateIssue("Epic", "", "epic", 2, "", nil, "")
	if err != nil {
		t.Fatalf("create epic: %v", err)
	}
	_, err = s.CreateIssue("Child A", "", "task", 2, "", nil, epic.ID)
	if err != nil {
		t.Fatalf("create child A: %v", err)
	}
	_, err = s.CreateIssue("Child B", "", "task", 2, "", nil, epic.ID)
	if err != nil {
		t.Fatalf("create child B: %v", err)
	}
	_, err = s.CreateIssue("Orphan", "", "task", 2, "", nil, "")
	if err != nil {
		t.Fatalf("create orphan: %v", err)
	}

	// Filter by parent.
	children, err := s.ListIssues(FilterOptions{Parent: epic.ID})
	if err != nil {
		t.Fatalf("list parent=%s: %v", epic.ID, err)
	}
	if len(children) != 2 {
		t.Errorf("expected 2 children, got %d", len(children))
	}

	// Filter by no-parent.
	orphans, err := s.ListIssues(FilterOptions{NoParent: true})
	if err != nil {
		t.Fatalf("list no-parent: %v", err)
	}
	if len(orphans) != 2 { // epic + orphan
		t.Errorf("expected 2 no-parent issues, got %d", len(orphans))
	}
}

func TestDeleteIssue(t *testing.T) {
	dir := t.TempDir()
	s, err := Init(dir, "TST", "tester")
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	a, err := s.CreateIssue("Issue A", "", "task", 2, "", nil, "")
	if err != nil {
		t.Fatalf("create A: %v", err)
	}
	b, err := s.CreateIssue("Issue B", "", "task", 2, "", nil, "")
	if err != nil {
		t.Fatalf("create B: %v", err)
	}

	// Add dep: B depends on A.
	if err := s.AddDependency(b.ID, a.ID); err != nil {
		t.Fatalf("add dep: %v", err)
	}

	// Delete A (soft) -- should clean up B's blocked_by.
	modified, err := s.DeleteIssue(a.ID, false)
	if err != nil {
		t.Fatalf("delete A: %v", err)
	}
	if len(modified) == 0 {
		t.Error("expected modified issues from dep cleanup")
	}

	// A should no longer exist.
	if s.IssueExists(a.ID) {
		t.Error("deleted issue should not exist")
	}

	// B should have no blockers.
	bRead, err := s.ReadIssue(b.ID)
	if err != nil {
		t.Fatalf("read B: %v", err)
	}
	if len(bRead.BlockedBy) != 0 {
		t.Errorf("B should have no blockers after A deleted: %v", bRead.BlockedBy)
	}
}

func TestLinksMigration(t *testing.T) {
	dir := t.TempDir()
	s, err := Init(dir, "TST", "tester")
	if err != nil {
		t.Fatalf("Init: %v", err)
	}

	// Create issue and verify it has ## Links section by default.
	issue, err := s.CreateIssue("Test issue", "some desc", "task", 2, "", nil, "")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	read, err := s.ReadIssue(issue.ID)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.Contains(read.Body, "\n## Links\n") {
		t.Errorf("new issue should have ## Links section:\n%s", read.Body)
	}

	// Call UpdateLinksSection on issue with no relationships (should be idempotent).
	if err := s.UpdateLinksSection(issue.ID); err != nil {
		t.Fatalf("UpdateLinksSection: %v", err)
	}

	// Verify body still has ## Links.
	read2, _ := s.ReadIssue(issue.ID)
	if !strings.Contains(read2.Body, "\n## Links\n") {
		t.Errorf("after UpdateLinksSection, ## Links should still exist:\n%s", read2.Body)
	}
}
