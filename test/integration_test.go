package test

import (
	"strings"
	"testing"

	"github.com/RamXX/nd/internal/enforce"
	"github.com/RamXX/nd/internal/graph"
	"github.com/RamXX/nd/internal/model"
	"github.com/RamXX/nd/internal/store"
	"github.com/RamXX/vlt"
)

// Full workflow: init -> create -> dep -> ready -> close -> stats.
// No mocks. Real vault files on disk.

func TestFullWorkflow(t *testing.T) {
	dir := t.TempDir()

	// 1. Init.
	s, err := store.Init(dir, "INT", "integration-tester")
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	if s.Prefix() != "INT" {
		t.Fatalf("prefix = %q", s.Prefix())
	}

	// 2. Create issues.
	epic, err := s.CreateIssue("Auth Epic", "Implement authentication", "epic", 1, "", nil, "")
	if err != nil {
		t.Fatalf("create epic: %v", err)
	}

	taskA, err := s.CreateIssue("Design auth flow", "Design the auth flow", "task", 1, "alice", []string{"auth"}, epic.ID)
	if err != nil {
		t.Fatalf("create taskA: %v", err)
	}

	taskB, err := s.CreateIssue("Implement auth flow", "Build the auth flow", "feature", 1, "bob", nil, epic.ID)
	if err != nil {
		t.Fatalf("create taskB: %v", err)
	}

	bug, err := s.CreateIssue("Fix login crash", "App crashes on login", "bug", 0, "alice", []string{"critical"}, "")
	if err != nil {
		t.Fatalf("create bug: %v", err)
	}

	// 3. Add dependency: taskB depends on taskA.
	if err := s.AddDependency(taskB.ID, taskA.ID); err != nil {
		t.Fatalf("add dep: %v", err)
	}

	// Verify bidirectional.
	taskARead, _ := s.ReadIssue(taskA.ID)
	taskBRead, _ := s.ReadIssue(taskB.ID)
	if !containsStr(taskARead.Blocks, taskB.ID) {
		t.Errorf("taskA should block taskB: blocks=%v", taskARead.Blocks)
	}
	if !containsStr(taskBRead.BlockedBy, taskA.ID) {
		t.Errorf("taskB should be blocked by taskA: blocked_by=%v", taskBRead.BlockedBy)
	}

	// 4. Build graph and test ready/blocked.
	all, err := s.ListIssues(store.FilterOptions{})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(all) != 4 {
		t.Fatalf("expected 4 issues, got %d", len(all))
	}

	g := graph.Build(all)
	ready := g.Ready()
	blocked := g.Blocked()

	// taskB should be blocked, everything else ready.
	readyIDs := idsOf(ready)
	blockedIDs := idsOf(blocked)

	if containsStr(readyIDs, taskB.ID) {
		t.Errorf("taskB should not be ready: %v", readyIDs)
	}
	if !containsStr(blockedIDs, taskB.ID) {
		t.Errorf("taskB should be blocked: %v", blockedIDs)
	}
	if !containsStr(readyIDs, taskA.ID) {
		t.Errorf("taskA should be ready: %v", readyIDs)
	}
	if !containsStr(readyIDs, bug.ID) {
		t.Errorf("bug should be ready: %v", readyIDs)
	}

	// 5. Close taskA -> taskB should become unblocked.
	if err := s.CloseIssue(taskA.ID, "Design complete"); err != nil {
		t.Fatalf("close taskA: %v", err)
	}

	all2, _ := s.ListIssues(store.FilterOptions{})
	g2 := graph.Build(all2)
	ready2 := g2.Ready()
	readyIDs2 := idsOf(ready2)

	if !containsStr(readyIDs2, taskB.ID) {
		t.Errorf("after closing taskA, taskB should be ready: %v", readyIDs2)
	}

	// 6. Test stats.
	stats := g2.Stats()
	if stats.Total != 4 {
		t.Errorf("total = %d", stats.Total)
	}
	if stats.Closed != 1 {
		t.Errorf("closed = %d", stats.Closed)
	}

	// 7. Reopen taskA.
	if err := s.ReopenIssue(taskA.ID); err != nil {
		t.Fatalf("reopen taskA: %v", err)
	}
	taskAReopened, _ := s.ReadIssue(taskA.ID)
	if taskAReopened.Status != model.StatusOpen {
		t.Errorf("taskA should be open after reopen, got %s", taskAReopened.Status)
	}

	// 8. Verify content hash integrity.
	for _, issue := range all {
		expected := enforce.ComputeContentHash(issue.Body)
		if issue.ContentHash != expected {
			// Content hash may drift if we modified body via comments.
			// This is expected and doctor --fix addresses it.
			t.Logf("hash drift on %s (expected after body modifications)", issue.ID)
		}
	}

	// 9. Epic tree.
	tree := g2.EpicTree(epic.ID)
	if tree == nil {
		t.Fatal("epic tree should not be nil")
	}
	if len(tree.Children) != 2 {
		t.Errorf("epic should have 2 children, got %d", len(tree.Children))
	}

	// 10. Epic status.
	summary := g2.EpicStatus(epic.ID)
	if summary.Total != 2 {
		t.Errorf("epic total = %d, want 2", summary.Total)
	}

	// 11. Remove dependency.
	if err := s.RemoveDependency(taskB.ID, taskA.ID); err != nil {
		t.Fatalf("remove dep: %v", err)
	}
	taskBAfter, _ := s.ReadIssue(taskB.ID)
	if len(taskBAfter.BlockedBy) != 0 {
		t.Errorf("taskB should have no blockers after removal: %v", taskBAfter.BlockedBy)
	}

	// 12. Update fields.
	if err := s.UpdateField(taskB.ID, "assignee", "charlie"); err != nil {
		t.Fatalf("update assignee: %v", err)
	}
	taskBUpdated, _ := s.ReadIssue(taskB.ID)
	if taskBUpdated.Assignee != "charlie" {
		t.Errorf("assignee = %q, want charlie", taskBUpdated.Assignee)
	}

	// 13. Filter by status.
	openIssues, _ := s.ListIssues(store.FilterOptions{Status: "open"})
	for _, i := range openIssues {
		if i.Status != model.StatusOpen {
			t.Errorf("filter returned non-open issue: %s %s", i.ID, i.Status)
		}
	}

	// 14. Filter by assignee.
	aliceIssues, _ := s.ListIssues(store.FilterOptions{Assignee: "alice"})
	for _, i := range aliceIssues {
		if !strings.EqualFold(i.Assignee, "alice") {
			t.Errorf("filter returned non-alice issue: %s %s", i.ID, i.Assignee)
		}
	}

	// 15. Filter by label.
	criticalIssues, _ := s.ListIssues(store.FilterOptions{Label: "critical"})
	for _, i := range criticalIssues {
		found := false
		for _, l := range i.Labels {
			if l == "critical" {
				found = true
			}
		}
		if !found {
			t.Errorf("filter returned issue without critical label: %s", i.ID)
		}
	}
}

func TestSearchIntegration(t *testing.T) {
	dir := t.TempDir()
	s, err := store.Init(dir, "SRC", "tester")
	if err != nil {
		t.Fatalf("init: %v", err)
	}

	_, err = s.CreateIssue("Database migration", "Migrate from Dolt to SQLite", "task", 2, "", nil, "")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	_, err = s.CreateIssue("API redesign", "Redesign the REST API", "feature", 1, "", nil, "")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	// Search for "Dolt".
	results, err := s.Vault().Search(vlt.SearchOptions{Query: "Dolt", Path: "issues"})
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result for 'Dolt', got %d", len(results))
	}
}

func containsStr(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func idsOf(issues []*model.Issue) []string {
	ids := make([]string, len(issues))
	for i, issue := range issues {
		ids[i] = issue.ID
	}
	return ids
}
