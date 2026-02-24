package model

import (
	"fmt"
	"strings"
	"time"
)

// Status represents the lifecycle state of an issue.
type Status string

const (
	StatusOpen       Status = "open"
	StatusInProgress Status = "in_progress"
	StatusBlocked    Status = "blocked"
	StatusDeferred   Status = "deferred"
	StatusClosed     Status = "closed"
)

var validStatuses = map[Status]bool{
	StatusOpen:       true,
	StatusInProgress: true,
	StatusBlocked:    true,
	StatusDeferred:   true,
	StatusClosed:     true,
}

func ParseStatus(s string) (Status, error) {
	st := Status(strings.ToLower(strings.TrimSpace(s)))
	if !validStatuses[st] {
		return "", fmt.Errorf("invalid status %q: must be one of open, in_progress, blocked, deferred, closed", s)
	}
	return st, nil
}

func (s Status) String() string { return string(s) }

// Priority ranges from 0 (critical) to 4 (backlog).
type Priority int

const (
	PriorityCritical Priority = 0
	PriorityHigh     Priority = 1
	PriorityMedium   Priority = 2
	PriorityLow      Priority = 3
	PriorityBacklog  Priority = 4
)

func ParsePriority(s string) (Priority, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	switch s {
	case "0", "P0":
		return PriorityCritical, nil
	case "1", "P1":
		return PriorityHigh, nil
	case "2", "P2":
		return PriorityMedium, nil
	case "3", "P3":
		return PriorityLow, nil
	case "4", "P4":
		return PriorityBacklog, nil
	default:
		return -1, fmt.Errorf("invalid priority %q: must be 0-4 or P0-P4", s)
	}
}

func (p Priority) String() string {
	switch p {
	case PriorityCritical:
		return "P0 (critical)"
	case PriorityHigh:
		return "P1 (high)"
	case PriorityMedium:
		return "P2 (medium)"
	case PriorityLow:
		return "P3 (low)"
	case PriorityBacklog:
		return "P4 (backlog)"
	default:
		return fmt.Sprintf("P%d", p)
	}
}

func (p Priority) Short() string {
	return fmt.Sprintf("P%d", p)
}

// IssueType classifies the nature of work.
type IssueType string

const (
	TypeBug      IssueType = "bug"
	TypeFeature  IssueType = "feature"
	TypeTask     IssueType = "task"
	TypeEpic     IssueType = "epic"
	TypeChore    IssueType = "chore"
	TypeDecision IssueType = "decision"
)

var validTypes = map[IssueType]bool{
	TypeBug:      true,
	TypeFeature:  true,
	TypeTask:     true,
	TypeEpic:     true,
	TypeChore:    true,
	TypeDecision: true,
}

func ParseIssueType(s string) (IssueType, error) {
	t := IssueType(strings.ToLower(strings.TrimSpace(s)))
	if !validTypes[t] {
		return "", fmt.Errorf("invalid type %q: must be one of bug, feature, task, epic, chore, decision", s)
	}
	return t, nil
}

func (t IssueType) String() string { return string(t) }

// Issue is the core data model for nd.
type Issue struct {
	ID          string    `yaml:"id"`
	Title       string    `yaml:"title"`
	Status      Status    `yaml:"status"`
	Priority    Priority  `yaml:"priority"`
	Type        IssueType `yaml:"type"`
	Assignee    string    `yaml:"assignee,omitempty"`
	Labels      []string  `yaml:"labels,omitempty"`
	Parent      string    `yaml:"parent,omitempty"`
	Blocks      []string  `yaml:"blocks,omitempty"`
	BlockedBy   []string  `yaml:"blocked_by,omitempty"`
	Related     []string  `yaml:"related,omitempty"`
	CreatedAt   time.Time `yaml:"created_at"`
	CreatedBy   string    `yaml:"created_by"`
	UpdatedAt   time.Time `yaml:"updated_at"`
	DeferUntil  string    `yaml:"defer_until,omitempty"`
	ClosedAt    string    `yaml:"closed_at,omitempty"`
	CloseReason string    `yaml:"close_reason,omitempty"`
	ContentHash string    `yaml:"content_hash"`

	// Runtime fields -- not serialized to YAML frontmatter.
	Body     string `yaml:"-"`
	FilePath string `yaml:"-"`
}

// Validate checks that required fields are populated and values are in range.
func (i *Issue) Validate() error {
	if i.ID == "" {
		return fmt.Errorf("issue ID is required")
	}
	if i.Title == "" {
		return fmt.Errorf("issue title is required")
	}
	if !validStatuses[i.Status] {
		return fmt.Errorf("invalid status %q", i.Status)
	}
	if i.Priority < 0 || i.Priority > 4 {
		return fmt.Errorf("priority must be 0-4, got %d", i.Priority)
	}
	if !validTypes[i.Type] {
		return fmt.Errorf("invalid type %q", i.Type)
	}
	if i.CreatedAt.IsZero() {
		return fmt.Errorf("created_at is required")
	}
	if i.CreatedBy == "" {
		return fmt.Errorf("created_by is required")
	}
	return nil
}

// IsOpen returns true if the issue is not closed.
func (i *Issue) IsOpen() bool {
	return i.Status != StatusClosed
}

// IsActionable returns true if the issue can be worked on (open or in_progress, not blocked).
func (i *Issue) IsActionable() bool {
	return (i.Status == StatusOpen || i.Status == StatusInProgress) && len(i.BlockedBy) == 0
}
