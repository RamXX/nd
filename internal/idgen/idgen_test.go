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
	// Prefix + dash + 4 base36 chars = len("TEST") + 1 + 4 = 9
	if len(id) != 9 {
		t.Errorf("ID %q should be 9 chars, got %d", id, len(id))
	}
	// Verify hash part is valid base36 (0-9, a-z)
	hash := id[5:]
	for _, c := range hash {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'z')) {
			t.Errorf("ID hash %q contains invalid base36 char %c", hash, c)
		}
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
	child := GenerateChildID("TEST-a3f8", 3)
	if child != "TEST-a3f8.3" {
		t.Errorf("expected TEST-a3f8.3, got %s", child)
	}
}

func TestEncodeBase36(t *testing.T) {
	// Known test: 3 bytes = 24 bits, should produce a 4-char base36 string
	data := []byte{0xff, 0xff, 0xff}
	result := EncodeBase36(data, 4)
	if len(result) != 4 {
		t.Errorf("expected 4 chars, got %d: %q", len(result), result)
	}
	for _, c := range result {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'z')) {
			t.Errorf("invalid base36 char %c in %q", c, result)
		}
	}

	// Zero input should produce all zeros
	zeros := EncodeBase36([]byte{0, 0, 0}, 4)
	if zeros != "0000" {
		t.Errorf("expected 0000 for zero input, got %q", zeros)
	}
}
