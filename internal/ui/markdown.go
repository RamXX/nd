package ui

import (
	"os"

	"github.com/charmbracelet/glamour"
	"golang.org/x/term"
)

// RenderMarkdown renders markdown text using glamour with auto-detected style.
// Returns the rendered markdown or the original text if rendering fails.
// Word wraps at terminal width (capped at 100 chars for readability).
func RenderMarkdown(markdown string) string {
	if !ShouldUseColor() {
		return markdown
	}

	const maxReadableWidth = 100
	wrapWidth := 80
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		wrapWidth = w
	}
	if wrapWidth > maxReadableWidth {
		wrapWidth = maxReadableWidth
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(wrapWidth),
	)
	if err != nil {
		return markdown
	}

	rendered, err := renderer.Render(markdown)
	if err != nil {
		return markdown
	}

	return rendered
}
