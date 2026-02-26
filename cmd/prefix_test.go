package cmd

import "testing"

func TestPrefixFromName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"tminus", "TMI"},
		{"my-project", "MP"},
		{"nd", "ND"},
		{"a", "A"},
		{"my_cool_app", "MCA"},
		{"foo-bar-baz-qux-extra", "FBBQ"},
		{"HexGraph", "HEX"},
		{"react2rlm", "REA"},
		{"a-b", "AB"},
		{"test.proj", "TP"},
		{"", ""},
		{"   ", ""},
		{"abc", "AB"},
		{"abcd", "AB"},
		{"abcde", "ABC"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := prefixFromName(tt.name)
			if got != tt.want {
				t.Errorf("prefixFromName(%q) = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}

func TestInferPrefixFromBeadsID(t *testing.T) {
	tests := []struct {
		id   string
		want string
	}{
		{"TM-a3f8", "TM"},
		{"PROJ-abc.3", "PROJ"},
		{"ND-001", "ND"},
		{"a-b", "A"},
		{"noid", ""},
		{"", ""},
		{"-bad", ""},
	}
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			got := inferPrefixFromBeadsID(tt.id)
			if got != tt.want {
				t.Errorf("inferPrefixFromBeadsID(%q) = %q, want %q", tt.id, got, tt.want)
			}
		})
	}
}
