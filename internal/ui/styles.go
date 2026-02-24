package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func init() {
	if !ShouldUseColor() {
		lipgloss.SetColorProfile(termenv.Ascii)
	} else {
		lipgloss.SetColorProfile(termenv.TrueColor)
	}
}

// Ayu theme color palette
var (
	ColorMuted = lipgloss.AdaptiveColor{
		Light: "#828c99",
		Dark:  "#6c7680",
	}
	ColorAccent = lipgloss.AdaptiveColor{
		Light: "#399ee6",
		Dark:  "#59c2ff",
	}

	// Status colors
	ColorStatusInProgress = lipgloss.AdaptiveColor{
		Light: "#f2ae49",
		Dark:  "#ffb454",
	}
	ColorStatusClosed = lipgloss.AdaptiveColor{
		Light: "#9099a1",
		Dark:  "#8090a0",
	}
	ColorStatusBlocked = lipgloss.AdaptiveColor{
		Light: "#f07171",
		Dark:  "#f26d78",
	}

	// Priority colors
	ColorPriorityP0 = lipgloss.AdaptiveColor{
		Light: "#f07171",
		Dark:  "#f07178",
	}
	ColorPriorityP1 = lipgloss.AdaptiveColor{
		Light: "#ff8f40",
		Dark:  "#ff8f40",
	}
	ColorPriorityP2 = lipgloss.AdaptiveColor{
		Light: "#e6b450",
		Dark:  "#e6b450",
	}

	// Type colors
	ColorTypeBug = lipgloss.AdaptiveColor{
		Light: "#f07171",
		Dark:  "#f26d78",
	}
	ColorTypeEpic = lipgloss.AdaptiveColor{
		Light: "#d2a6ff",
		Dark:  "#d2a6ff",
	}
)

// Styles
var (
	MutedStyle  = lipgloss.NewStyle().Foreground(ColorMuted)
	AccentStyle = lipgloss.NewStyle().Foreground(ColorAccent)
	BoldStyle   = lipgloss.NewStyle().Bold(true)

	StatusInProgressStyle = lipgloss.NewStyle().Foreground(ColorStatusInProgress)
	StatusClosedStyle     = lipgloss.NewStyle().Foreground(ColorStatusClosed)
	StatusBlockedStyle    = lipgloss.NewStyle().Foreground(ColorStatusBlocked)

	PriorityP0Style = lipgloss.NewStyle().Foreground(ColorPriorityP0).Bold(true)
	PriorityP1Style = lipgloss.NewStyle().Foreground(ColorPriorityP1)
	PriorityP2Style = lipgloss.NewStyle().Foreground(ColorPriorityP2)

	TypeBugStyle  = lipgloss.NewStyle().Foreground(ColorTypeBug)
	TypeEpicStyle = lipgloss.NewStyle().Foreground(ColorTypeEpic)
)

// Status icons
const (
	StatusIconOpen       = "○"
	StatusIconInProgress = "◐"
	StatusIconBlocked    = "●"
	StatusIconClosed     = "✓"
	StatusIconDeferred   = "❄"
	PriorityIcon         = "●"
)

// RenderStatusIcon returns the icon for a status with coloring.
func RenderStatusIcon(status string) string {
	switch status {
	case "open":
		return StatusIconOpen
	case "in_progress":
		return StatusInProgressStyle.Render(StatusIconInProgress)
	case "blocked":
		return StatusBlockedStyle.Render(StatusIconBlocked)
	case "closed":
		return StatusClosedStyle.Render(StatusIconClosed)
	case "deferred":
		return MutedStyle.Render(StatusIconDeferred)
	default:
		return "?"
	}
}

// RenderStatus renders a status string with coloring.
func RenderStatus(status string) string {
	switch status {
	case "in_progress":
		return StatusInProgressStyle.Render(status)
	case "blocked":
		return StatusBlockedStyle.Render(status)
	case "closed":
		return StatusClosedStyle.Render(status)
	case "deferred":
		return MutedStyle.Render(status)
	default:
		return status
	}
}

// RenderPriority renders priority with icon and color.
func RenderPriority(priority int) string {
	label := fmt.Sprintf("%s P%d", PriorityIcon, priority)
	switch priority {
	case 0:
		return PriorityP0Style.Render(label)
	case 1:
		return PriorityP1Style.Render(label)
	case 2:
		return PriorityP2Style.Render(label)
	default:
		return label
	}
}

// RenderType renders an issue type with coloring.
func RenderType(issueType string) string {
	switch issueType {
	case "bug":
		return TypeBugStyle.Render(issueType)
	case "epic":
		return TypeEpicStyle.Render(issueType)
	default:
		return issueType
	}
}

// RenderID renders an issue ID (standard text).
func RenderID(id string) string {
	return id
}

// RenderMuted renders text in muted gray.
func RenderMuted(s string) string {
	return MutedStyle.Render(s)
}

// RenderBold renders text in bold.
func RenderBold(s string) string {
	return BoldStyle.Render(s)
}

// RenderAccent renders text with accent color.
func RenderAccent(s string) string {
	return AccentStyle.Render(s)
}

// RenderClosedLine renders an entire line in the closed/dimmed style.
func RenderClosedLine(line string) string {
	return StatusClosedStyle.Render(line)
}
