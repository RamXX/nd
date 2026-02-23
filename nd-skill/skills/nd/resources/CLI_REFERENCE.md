# CLI Command Reference

**For:** AI agents and developers using the nd command-line interface
**Version:** 0.1.0+

## Quick Navigation

- [Initialization](#initialization)
- [Issue Management](#issue-management)
- [Finding Work](#finding-work)
- [Dependencies](#dependencies)
- [Labels and Comments](#labels-and-comments)
- [Epics](#epics)
- [Search and Stats](#search-and-stats)
- [AI Context and Health](#ai-context-and-health)
- [Global Flags](#global-flags)

## Initialization

```bash
# Create a new nd vault
nd init --prefix=PROJ                    # Required: issue ID prefix
nd init --prefix=PROJ --vault=/path      # Custom vault location
nd init --prefix=PROJ --author=alice     # Custom default author
```

## Issue Management

### Create

```bash
# Basic creation
nd create "Issue title" --type=task --priority=2

# Full options
nd create "Title" \
  --type=bug|feature|task|epic|chore|decision \
  --priority=0-4 \
  --assignee=alice \
  --labels=auth,urgent \
  --description="Detailed description" \
  --parent=PROJ-a3f

# Short flags
nd create "Title" -t bug -p 1 -d "Description"
```

Output: `Created PROJ-a3f: Issue title`
With `--quiet`: just the ID
With `--json`: `{"id":"PROJ-a3f"}`

### Show

```bash
nd show PROJ-a3f              # Full detail view
nd show PROJ-a3f --short      # One-line summary
nd show PROJ-a3f --json       # JSON output
```

### Update

```bash
# Update fields (at least one required)
nd update PROJ-a3f --status=in_progress
nd update PROJ-a3f --priority=0 --assignee=bob
nd update PROJ-a3f --title="New title"
nd update PROJ-a3f --type=bug
nd update PROJ-a3f --append-notes="Found the root cause"
nd update PROJ-a3f -d "New description"
```

### Close and Reopen

```bash
# Close one or more issues
nd close PROJ-a3f                         # Close single
nd close PROJ-a3f PROJ-b7c               # Close multiple (batch)
nd close PROJ-a3f --reason="Implemented"  # With reason

# Reopen a closed issue
nd reopen PROJ-a3f
```

### List

```bash
# List with filters
nd list                                   # All issues (sorted by priority)
nd list --status=open                     # Filter by status
nd list --type=bug                        # Filter by type
nd list --assignee=alice                  # Filter by assignee
nd list --label=critical                  # Filter by label
nd list --sort=created                    # Sort: priority, created, updated, id
nd list -n 10                             # Limit results
nd list --json                            # JSON output
```

## Finding Work

```bash
# Ready work (no blockers)
nd ready                                  # All ready issues
nd ready --assignee=alice                 # Filter by assignee
nd ready --sort=priority                  # Sort: priority, created, updated, id
nd ready -n 5                             # Top 5

# Blocked work
nd blocked                                # Show blocked issues
nd blocked --verbose                      # Include blocker details
```

## Dependencies

```bash
# Add dependency (A depends on B: B blocks A)
nd dep add PROJ-a3f PROJ-b7c

# Remove dependency
nd dep rm PROJ-a3f PROJ-b7c

# List all dependencies of an issue
nd dep list PROJ-a3f
```

**Dependency semantics**: `nd dep add A B` means "A depends on B" (B must complete before A). This updates both files:
- Adds B to A's `blocked_by`
- Adds A to B's `blocks`

Removing a dependency cleans both sides.

## Labels and Comments

### Labels

```bash
nd labels add PROJ-a3f security          # Add label
nd labels rm PROJ-a3f security           # Remove label
nd labels list                           # All labels with counts
```

### Comments

```bash
nd comments add PROJ-a3f "Comment text"  # Add timestamped comment
nd comments list PROJ-a3f                # View comments
```

Comments are appended to the `## Comments` section in the issue file with RFC3339 timestamp and author.

## Epics

```bash
# Epic progress summary
nd epic status PROJ-a3f
# Output: Children count, open/in_progress/blocked/closed, progress %

# Epic tree view
nd epic tree PROJ-a3f
# Output: Hierarchical tree with status markers
#   [ ] open  [>] in_progress  [!] blocked  [x] closed
```

## Search and Stats

```bash
# Full-text search across issues
nd search "authentication"                # Returns matching lines with context

# Project statistics
nd stats                                  # Text summary
nd stats --json                           # JSON (Total, Open, InProgress, etc.)
```

## AI Context and Health

```bash
# AI context output
nd prime                                  # Structured summary for AI
nd prime --json                           # Full project state as JSON

# Vault health check
nd doctor                                 # Validate integrity
nd doctor --fix                           # Auto-fix problems
```

Doctor checks:
1. Content hash integrity (SHA-256 of body vs stored hash)
2. Bidirectional dependency consistency
3. Reference validity (no orphan dep references)
4. Field validation (required fields, valid enums)

### Import

```bash
# Import from beads JSONL
nd import --from-beads .beads/issues.jsonl
```

## Global Flags

All commands support these flags:

```bash
--vault PATH     # Override vault directory (default: .vault, auto-discovered)
--json           # Output as JSON
--verbose        # Verbose output
--quiet          # Suppress non-essential output
```

Vault auto-discovery walks up the directory tree looking for `.vault/`.

## Priority System

| Value | Label | Use for |
|-------|-------|---------|
| 0 / P0 | Critical | Security, data loss, broken builds |
| 1 / P1 | High | Major features, important bugs |
| 2 / P2 | Medium | Standard work (default) |
| 3 / P3 | Low | Polish, optimization |
| 4 / P4 | Backlog | Future ideas |

## Status Values

- `open` -- Available to work on
- `in_progress` -- Currently being worked
- `blocked` -- Blocked by dependencies
- `deferred` -- Intentionally deferred
- `closed` -- Completed

## Issue Types

`bug`, `feature`, `task`, `epic`, `chore`, `decision`
