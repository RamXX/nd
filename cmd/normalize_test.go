package cmd

import (
	"strings"
	"testing"
)

func TestEscapeFalseTags(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "hex color",
			input: "color=#22c55e in the badge",
			want:  "color=`#22c55e` in the badge",
		},
		{
			name:  "short hex",
			input: "uses #fff for white",
			want:  "uses `#fff` for white",
		},
		{
			name:  "preprocessor if",
			input: "Platform guards (#if os(iOS)) for APIs",
			want:  "Platform guards (`#if` os(iOS)) for APIs",
		},
		{
			name:  "preprocessor endif",
			input: "wrapped in #endif block",
			want:  "wrapped in `#endif` block",
		},
		{
			name:  "canImport",
			input: "#if canImport(UIKit) import UIKit #endif",
			want:  "`#if` canImport(UIKit) import UIKit `#endif`",
		},
		{
			name:  "css selector",
			input: "site/index.html section#proof displays",
			want:  "site/index.html section`#proof` displays",
		},
		{
			name:  "standalone hash-word",
			input: "the -> #proof section here",
			want:  "the -> `#proof` section here",
		},
		{
			name:  "already in backticks",
			input: "uses `#if os(iOS)` for guard",
			want:  "uses `#if os(iOS)` for guard",
		},
		{
			name:  "heading preserved",
			input: "## Description\nsome text",
			want:  "## Description\nsome text",
		},
		{
			name:  "no false tags",
			input: "This is plain text with no issues",
			want:  "This is plain text with no issues",
		},
		{
			name:  "multiple on one line",
			input: "green=#22c55e, yellow=#eab308, blue=#3b82f6",
			want:  "green=`#22c55e`, yellow=`#eab308`, blue=`#3b82f6`",
		},
		{
			name:  "mixed backtick and bare",
			input: "see `#if` and also #endif here",
			want:  "see `#if` and also `#endif` here",
		},
		{
			name:  "hash route",
			input: "Router at #/accounts route",
			want:  "Router at `#/accounts` route",
		},
		{
			name:  "hash route with query",
			input: "navigate to #/scheduling?duration=60&contact=Alice",
			want:  "navigate to `#/scheduling?duration=60&contact=Alice`",
		},
		{
			name:  "hash route bare slash ignored",
			input: "redirect to #/ default",
			want:  "redirect to #/ default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := escapeFalseTags(tt.input)
			if got != tt.want {
				t.Errorf("\n got: %s\nwant: %s", got, tt.want)
			}
		})
	}
}

func TestNormalizeMarkdownIntegration(t *testing.T) {
	input := "## Notes\nPlatform guards (#if os(iOS)) work\n| col |\n| --- |"
	got := normalizeMarkdown(input)

	// Heading should get blank line after.
	if !strings.Contains(got, "## Notes\n\n") {
		t.Errorf("heading should have blank line after:\n%s", got)
	}
	// #if should be escaped.
	if !strings.Contains(got, "`#if`") {
		t.Errorf("#if should be escaped:\n%s", got)
	}
	// Table should have blank line before.
	if !strings.Contains(got, "work\n\n| col |") {
		t.Errorf("table should have blank line before:\n%s", got)
	}
}
