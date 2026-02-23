package format

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/RamXX/nd/internal/model"
)

// Table renders a list of issues as a text table.
func Table(w io.Writer, issues []*model.Issue) {
	if len(issues) == 0 {
		fmt.Fprintln(w, "No issues found.")
		return
	}

	// Column widths.
	idW, titleW, statusW, prioW, typeW, assigneeW := 12, 40, 12, 4, 10, 12

	// Header.
	fmt.Fprintf(w, "%-*s %-*s %-*s %-*s %-*s %-*s\n",
		idW, "ID", titleW, "TITLE", statusW, "STATUS", prioW, "PRI", typeW, "TYPE", assigneeW, "ASSIGNEE")
	fmt.Fprintln(w, strings.Repeat("-", idW+titleW+statusW+prioW+typeW+assigneeW+5))

	for _, issue := range issues {
		title := issue.Title
		if len(title) > titleW {
			title = title[:titleW-3] + "..."
		}
		fmt.Fprintf(w, "%-*s %-*s %-*s %-*s %-*s %-*s\n",
			idW, issue.ID,
			titleW, title,
			statusW, issue.Status,
			prioW, issue.Priority.Short(),
			typeW, issue.Type,
			assigneeW, issue.Assignee)
	}
	fmt.Fprintf(w, "\n%d issue(s)\n", len(issues))
}

// Detail renders a single issue in full detail.
func Detail(w io.Writer, issue *model.Issue) {
	fmt.Fprintf(w, "%s: %s\n", issue.ID, issue.Title)
	fmt.Fprintf(w, "Status:   %s\n", issue.Status)
	fmt.Fprintf(w, "Priority: %s\n", issue.Priority)
	fmt.Fprintf(w, "Type:     %s\n", issue.Type)
	if issue.Assignee != "" {
		fmt.Fprintf(w, "Assignee: %s\n", issue.Assignee)
	}
	if len(issue.Labels) > 0 {
		fmt.Fprintf(w, "Labels:   %s\n", strings.Join(issue.Labels, ", "))
	}
	if issue.Parent != "" {
		fmt.Fprintf(w, "Parent:   %s\n", issue.Parent)
	}
	if len(issue.Blocks) > 0 {
		fmt.Fprintf(w, "Blocks:   %s\n", strings.Join(issue.Blocks, ", "))
	}
	if len(issue.BlockedBy) > 0 {
		fmt.Fprintf(w, "Blocked:  %s\n", strings.Join(issue.BlockedBy, ", "))
	}
	if len(issue.Related) > 0 {
		fmt.Fprintf(w, "Related:  %s\n", strings.Join(issue.Related, ", "))
	}
	fmt.Fprintf(w, "Created:  %s by %s\n", issue.CreatedAt.Format("2006-01-02 15:04"), issue.CreatedBy)
	fmt.Fprintf(w, "Updated:  %s\n", issue.UpdatedAt.Format("2006-01-02 15:04"))
	if issue.ClosedAt != "" {
		fmt.Fprintf(w, "Closed:   %s\n", issue.ClosedAt)
	}
	if issue.CloseReason != "" {
		fmt.Fprintf(w, "Reason:   %s\n", issue.CloseReason)
	}
	fmt.Fprintf(w, "Hash:     %s\n", issue.ContentHash)
	if issue.Body != "" {
		fmt.Fprintf(w, "\n%s", issue.Body)
	}
}

// JSON outputs issues as JSON.
func JSON(w io.Writer, issues []*model.Issue) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(issues)
}

// JSONSingle outputs a single issue as JSON.
func JSONSingle(w io.Writer, issue *model.Issue) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(issue)
}

// Short renders a one-line summary of an issue.
func Short(w io.Writer, issue *model.Issue) {
	fmt.Fprintf(w, "%s [%s] %s (%s)\n", issue.ID, issue.Status, issue.Title, issue.Priority.Short())
}
