package model

import (
	"testing"
	"time"
)

func TestParseStatus(t *testing.T) {
	tests := []struct {
		input string
		want  Status
		err   bool
	}{
		{"open", StatusOpen, false},
		{"in_progress", StatusInProgress, false},
		{"BLOCKED", StatusBlocked, false},
		{"  deferred  ", StatusDeferred, false},
		{"closed", StatusClosed, false},
		{"invalid", "", true},
		{"", "", true},
	}
	for _, tt := range tests {
		got, err := ParseStatus(tt.input)
		if tt.err && err == nil {
			t.Errorf("ParseStatus(%q) expected error", tt.input)
		}
		if !tt.err && err != nil {
			t.Errorf("ParseStatus(%q) unexpected error: %v", tt.input, err)
		}
		if got != tt.want {
			t.Errorf("ParseStatus(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestParsePriority(t *testing.T) {
	tests := []struct {
		input string
		want  Priority
		err   bool
	}{
		{"0", PriorityCritical, false},
		{"P1", PriorityHigh, false},
		{"p2", PriorityMedium, false},
		{"3", PriorityLow, false},
		{"P4", PriorityBacklog, false},
		{"5", -1, true},
		{"high", -1, true},
	}
	for _, tt := range tests {
		got, err := ParsePriority(tt.input)
		if tt.err && err == nil {
			t.Errorf("ParsePriority(%q) expected error", tt.input)
		}
		if !tt.err && err != nil {
			t.Errorf("ParsePriority(%q) unexpected error: %v", tt.input, err)
		}
		if got != tt.want {
			t.Errorf("ParsePriority(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestParseIssueType(t *testing.T) {
	tests := []struct {
		input string
		want  IssueType
		err   bool
	}{
		{"bug", TypeBug, false},
		{"FEATURE", TypeFeature, false},
		{"task", TypeTask, false},
		{"epic", TypeEpic, false},
		{"chore", TypeChore, false},
		{"decision", TypeDecision, false},
		{"story", "", true},
	}
	for _, tt := range tests {
		got, err := ParseIssueType(tt.input)
		if tt.err && err == nil {
			t.Errorf("ParseIssueType(%q) expected error", tt.input)
		}
		if !tt.err && err != nil {
			t.Errorf("ParseIssueType(%q) unexpected error: %v", tt.input, err)
		}
		if got != tt.want {
			t.Errorf("ParseIssueType(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestValidate(t *testing.T) {
	valid := &Issue{
		ID:        "TEST-abc",
		Title:     "Test issue",
		Status:    StatusOpen,
		Priority:  PriorityMedium,
		Type:      TypeTask,
		CreatedAt: time.Now(),
		CreatedBy: "tester",
	}
	if err := valid.Validate(); err != nil {
		t.Errorf("valid issue returned error: %v", err)
	}

	// Missing ID.
	noID := *valid
	noID.ID = ""
	if err := noID.Validate(); err == nil {
		t.Error("expected error for missing ID")
	}

	// Invalid priority.
	badPri := *valid
	badPri.Priority = 5
	if err := badPri.Validate(); err == nil {
		t.Error("expected error for priority 5")
	}

	// Invalid type.
	badType := *valid
	badType.Type = "story"
	if err := badType.Validate(); err == nil {
		t.Error("expected error for invalid type")
	}
}

func TestParseStatusWithCustom(t *testing.T) {
	custom := []Status{"delivered", "accepted", "rejected"}

	// Built-in still works.
	st, err := ParseStatusWithCustom("open", custom)
	if err != nil {
		t.Fatalf("unexpected error for built-in: %v", err)
	}
	if st != StatusOpen {
		t.Errorf("got %q, want open", st)
	}

	// Custom status works.
	st, err = ParseStatusWithCustom("delivered", custom)
	if err != nil {
		t.Fatalf("unexpected error for custom: %v", err)
	}
	if st != "delivered" {
		t.Errorf("got %q, want delivered", st)
	}

	// Unknown fails.
	_, err = ParseStatusWithCustom("unknown", custom)
	if err == nil {
		t.Error("expected error for unknown status")
	}

	// Without custom, custom status fails.
	_, err = ParseStatusWithCustom("delivered", nil)
	if err == nil {
		t.Error("expected error for custom status without custom list")
	}
}

func TestValidateWithCustomStatus(t *testing.T) {
	custom := []Status{"delivered"}
	issue := &Issue{
		ID:        "TST-001",
		Title:     "Test",
		Status:    "delivered",
		Priority:  PriorityMedium,
		Type:      TypeTask,
		CreatedAt: time.Now(),
		CreatedBy: "tester",
	}

	// Should fail with Validate (no custom).
	if err := issue.Validate(); err == nil {
		t.Error("Validate should reject custom status without custom list")
	}

	// Should pass with ValidateWithCustom.
	if err := issue.ValidateWithCustom(custom); err != nil {
		t.Errorf("ValidateWithCustom should accept custom status: %v", err)
	}
}

func TestIsBuiltinStatus(t *testing.T) {
	for _, name := range []string{"open", "in_progress", "blocked", "deferred", "closed"} {
		if !IsBuiltinStatus(name) {
			t.Errorf("IsBuiltinStatus(%q) = false, want true", name)
		}
	}
	for _, name := range []string{"delivered", "accepted", "rejected", "unknown"} {
		if IsBuiltinStatus(name) {
			t.Errorf("IsBuiltinStatus(%q) = true, want false", name)
		}
	}
}

func TestIsActionable(t *testing.T) {
	issue := &Issue{Status: StatusOpen}
	if !issue.IsActionable() {
		t.Error("open issue with no blockers should be actionable")
	}

	issue.BlockedBy = []string{"X-001"}
	if issue.IsActionable() {
		t.Error("blocked issue should not be actionable")
	}

	issue.BlockedBy = nil
	issue.Status = StatusClosed
	if issue.IsActionable() {
		t.Error("closed issue should not be actionable")
	}
}
