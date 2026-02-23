package store

import (
	"os"
	"testing"
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
