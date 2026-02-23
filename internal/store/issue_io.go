package store

import (
	"fmt"
	"strings"
	"time"

	"github.com/RamXX/nd/internal/enforce"
	"github.com/RamXX/nd/internal/idgen"
	"github.com/RamXX/nd/internal/model"
	"github.com/RamXX/vlt"
	"gopkg.in/yaml.v3"
)

// CreateIssue generates an ID, serializes the issue to markdown, and writes it to the vault.
func (s *Store) CreateIssue(title, description, issueType string, priority int, assignee string, labels []string, parent string) (*model.Issue, error) {
	itype, err := model.ParseIssueType(issueType)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	id, err := idgen.GenerateID(s.config.Prefix, title, s.IssueExists)
	if err != nil {
		return nil, fmt.Errorf("generate ID: %w", err)
	}

	body := buildBody(description)
	issue := &model.Issue{
		ID:          id,
		Title:       title,
		Status:      model.StatusOpen,
		Priority:    model.Priority(priority),
		Type:        itype,
		Assignee:    assignee,
		Labels:      labels,
		Parent:      parent,
		CreatedAt:   now,
		CreatedBy:   s.config.CreatedBy,
		UpdatedAt:   now,
		ContentHash: enforce.ComputeContentHash(body),
		Body:        body,
	}

	if err := issue.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	content := serializeIssue(issue)
	path := fmt.Sprintf("issues/%s.md", id)

	if err := s.vault.Create(id, path, content, true, false); err != nil {
		return nil, fmt.Errorf("create file: %w", err)
	}

	issue.FilePath = path
	return issue, nil
}

// ReadIssue reads and deserializes an issue by ID.
func (s *Store) ReadIssue(id string) (*model.Issue, error) {
	content, err := s.vault.Read(id, "")
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", id, err)
	}
	issue, err := deserializeIssue(content)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", id, err)
	}
	issue.FilePath = fmt.Sprintf("issues/%s.md", id)
	return issue, nil
}

func buildBody(description string) string {
	var sb strings.Builder
	sb.WriteString("\n## Description\n")
	if description != "" {
		sb.WriteString(description)
		sb.WriteString("\n")
	}
	sb.WriteString("\n## Acceptance Criteria\n\n")
	sb.WriteString("\n## Design\n\n")
	sb.WriteString("\n## Notes\n\n")
	sb.WriteString("\n## Comments\n")
	return sb.String()
}

// serializeIssue converts an Issue to frontmatter + body markdown.
func serializeIssue(issue *model.Issue) string {
	fm := marshalFrontmatter(issue)
	return fmt.Sprintf("---\n%s---\n%s", fm, issue.Body)
}

func marshalFrontmatter(issue *model.Issue) string {
	// Use a map to control field ordering via manual construction.
	// yaml.Marshal on the struct would work but doesn't guarantee order.
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("id: %s\n", issue.ID))
	sb.WriteString(fmt.Sprintf("title: %q\n", issue.Title))
	sb.WriteString(fmt.Sprintf("status: %s\n", issue.Status))
	sb.WriteString(fmt.Sprintf("priority: %d\n", issue.Priority))
	sb.WriteString(fmt.Sprintf("type: %s\n", issue.Type))
	if issue.Assignee != "" {
		sb.WriteString(fmt.Sprintf("assignee: %s\n", issue.Assignee))
	}
	if len(issue.Labels) > 0 {
		sb.WriteString(fmt.Sprintf("labels: [%s]\n", strings.Join(issue.Labels, ", ")))
	}
	if issue.Parent != "" {
		sb.WriteString(fmt.Sprintf("parent: %s\n", issue.Parent))
	}
	writeStringList(&sb, "blocks", issue.Blocks)
	writeStringList(&sb, "blocked_by", issue.BlockedBy)
	writeStringList(&sb, "related", issue.Related)
	sb.WriteString(fmt.Sprintf("created_at: %s\n", issue.CreatedAt.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("created_by: %s\n", issue.CreatedBy))
	sb.WriteString(fmt.Sprintf("updated_at: %s\n", issue.UpdatedAt.Format(time.RFC3339)))
	if issue.ClosedAt != "" {
		sb.WriteString(fmt.Sprintf("closed_at: %s\n", issue.ClosedAt))
	}
	if issue.CloseReason != "" {
		sb.WriteString(fmt.Sprintf("close_reason: %q\n", issue.CloseReason))
	}
	sb.WriteString(fmt.Sprintf("content_hash: %q\n", issue.ContentHash))
	return sb.String()
}

func writeStringList(sb *strings.Builder, key string, vals []string) {
	if len(vals) == 0 {
		return
	}
	sb.WriteString(fmt.Sprintf("%s: [%s]\n", key, strings.Join(vals, ", ")))
}

// deserializeIssue parses frontmatter + body markdown into an Issue.
func deserializeIssue(content string) (*model.Issue, error) {
	yamlStr, bodyStart, found := vlt.ExtractFrontmatter(content)
	if !found {
		return nil, fmt.Errorf("no frontmatter found")
	}

	var issue model.Issue
	if err := yaml.Unmarshal([]byte(yamlStr), &issue); err != nil {
		return nil, fmt.Errorf("unmarshal frontmatter: %w", err)
	}

	// Extract body: everything after the closing ---.
	lines := strings.SplitAfter(content, "\n")
	if bodyStart < len(lines) {
		issue.Body = strings.Join(lines[bodyStart:], "")
	}
	return &issue, nil
}
