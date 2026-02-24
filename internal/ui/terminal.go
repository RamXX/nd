// Package ui provides terminal styling and output helpers for nd CLI.
package ui

import (
	"os"

	"golang.org/x/term"
)

// IsTerminal returns true if stdout is connected to a terminal (TTY).
func IsTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// ShouldUseColor determines if ANSI color codes should be used.
// Respects standard conventions:
//   - NO_COLOR: https://no-color.org/ - disables color if set
//   - CLICOLOR=0: disables color
//   - CLICOLOR_FORCE: forces color even in non-TTY
//   - Falls back to TTY detection
func ShouldUseColor() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("CLICOLOR") == "0" {
		return false
	}
	if os.Getenv("CLICOLOR_FORCE") != "" {
		return true
	}
	return IsTerminal()
}
