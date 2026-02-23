package idgen

import (
	"strings"
	"testing"
)

func TestGenerateID(t *testing.T) {
	id, err := GenerateID("TEST", "My issue title", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(id, "TEST-") {
		t.Errorf("ID %q should start with TEST-", id)
	}
	// Prefix + dash + 3 hex chars = len("TEST") + 1 + 3 = 8
	if len(id) != 8 {
		t.Errorf("ID %q should be 8 chars, got %d", id, len(id))
	}
}

func TestGenerateIDCollision(t *testing.T) {
	seen := map[string]bool{}
	id, err := GenerateID("X", "title", func(candidate string) bool {
		if seen[candidate] {
			return true // simulate collision
		}
		seen[candidate] = true
		return false
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == "" {
		t.Error("expected non-empty ID")
	}
}

func TestGenerateIDExhausted(t *testing.T) {
	_, err := GenerateID("X", "title", func(string) bool { return true })
	if err == nil {
		t.Error("expected error when all IDs collide")
	}
}

func TestGenerateChildID(t *testing.T) {
	child := GenerateChildID("TEST-abc", 3)
	if child != "TEST-abc.3" {
		t.Errorf("expected TEST-abc.3, got %s", child)
	}
}
