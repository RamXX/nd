# nd - Vault-backed Issue Tracker

**nd** (short for `node` as node in a graph) is a Git-native issue tracker that stores issues as Obsidian-compatible markdown files with YAML frontmatter. No database server. No field size limits. Plain files you can read, grep, and version with git.

## Why nd Exists

We built [beads](https://github.com/steveyegge/beads) (`bd`) into the backbone of our AI-assisted development workflow. It worked.Persistent memory across compaction, dependency graphs, epic hierarchies. We love beads. But extremely fast development cycles, breaking changes, and the latest adoption of a new storage backend (Dolt SQL) was too much for us. It became the weakest link:

- **65KB TEXT field limit.** We routinely store 80-160KB in issue descriptions, design notes, and acceptance criteria. Dolt silently truncates at 65KB. Data loss you don't notice until you need it. We believe in agents that are single-use, so beads carry all context. This limitation was too much.
- **Running server required.** `dolt sql-server` must be running or bd falls back to embedded mode with different behavior. Configuration confusion is constant.
- **Migration headaches.** Schema migrations, repo fingerprint mismatches, shared server database confusion, JSONL import gaps. Every other session starts with `bd doctor --fix`.
- **Not inspectable.** Issues live in a binary database. You can't `cat` an issue, `grep` across your backlog, or diff changes in a PR review.

`nd` solves all of this by storing issues as plain markdown files in a directory. The storage layer is [vlt](https://github.com/RamXX/vlt), an Obsidian-compliant vault management library, used as an importable Go library. vlt handles file I/O, frontmatter parsing, search, file locking, and content patching. `nd` adds issue-tracker semantics on top.

### What vlt Provides

[vlt](https://github.com/RamXX/vlt) (`github.com/RamXX/vlt`) is an Obsidian-compatible vault CLI and Go library. nd imports it directly for:

| Capability | vlt API | nd usage |
|---|---|---|
| Open vault | `vlt.Open(dir)` | Store initialization |
| Create file | `v.Create()` | Issue creation |
| Read file | `v.Read(title, heading)` | Issue reads, section-scoped reads |
| Write file | `v.Write()` | Body replacement |
| Append | `v.Append()` | Adding comments |
| Patch | `v.Patch(title, PatchOptions)` | Surgical heading-level edits |
| Frontmatter | `v.PropertySet()`, `v.PropertyRemove()` | Single-field updates |
| Search | `v.Search()`, `v.SearchWithContext()` | Full-text search |
| File listing | `v.Files(folder, ext)` | Issue enumeration |
| Locking | `vlt.LockVault(dir, exclusive)` | Concurrent access safety |
| Delete | `v.Delete()` | Soft delete to .trash/ |

nd adds: issue model with validation, collision-resistant ID generation, dependency graph computation (ready/blocked/cycles), content hashing, epic tree traversal, and output formatting.

## Installation

```bash
git clone https://github.com/RamXX/nd.git
cd nd
make build
make install    # Installs to ~/.local/bin/nd
```

Requires Go 1.22+.

## Quick Start

```bash
# Initialize a vault in your project
nd init --prefix=PROJ

# Create issues
nd create "Implement user auth" --type=feature --priority=1 --assignee=alice
nd create "Fix login crash" --type=bug --priority=0 -d "App crashes on special chars"

# Find work
nd ready                    # Show actionable issues (no blockers)
nd list --status=open       # All open issues

# Work on something
nd update PROJ-a3f --status=in_progress
nd comments add PROJ-a3f "Root cause found: missing input sanitization"

# Manage dependencies
nd dep add PROJ-b7c PROJ-a3f    # PROJ-b7c depends on PROJ-a3f
nd blocked                       # See what's stuck

# Complete work
nd close PROJ-a3f --reason="Fixed with input validation"
nd ready                         # PROJ-b7c is now unblocked
```

## Storage Format

Issues live as markdown files in `.vault/issues/`:

```yaml
---
id: PROJ-a3f
title: "Implement user authentication"
status: open
priority: 1
type: feature
assignee: alice
labels: [security, milestone]
blocks: [PROJ-d9e]
blocked_by: [PROJ-b3c]
created_at: 2026-02-23T20:15:00Z
created_by: alice
updated_at: 2026-02-24T10:30:00Z
content_hash: "sha256:a3f8c9d2..."
---

## Description
Implement OAuth 2.0 authentication with JWT tokens...

## Acceptance Criteria
- [ ] Login endpoint returns JWT
- [ ] Token refresh works
- [ ] Rate limiting on auth endpoints

## Design
Using bcrypt with 12 rounds per OWASP recommendation...

## Notes
Spike complete. Chose Authorization Code flow over Implicit.

## Comments

### 2026-02-23T20:15:00Z alice
Started implementation. Base models done.
```

Every issue is a file you can read with `cat`, search with `grep`, diff with `git diff`, and edit with any text editor. No database required.

### Vault Layout

```
.vault/
  .nd.yaml            # Config: version, prefix, created_by
  issues/             # One .md file per issue
    PROJ-a3f.md
    PROJ-b7c.md
    PROJ-d9e.md
  .trash/             # Soft-deleted issues
  .vlt.lock           # Advisory file lock
```

## Command Reference

### Initialization

```bash
nd init --prefix=PROJ [--vault=PATH] [--author=NAME]
```

Creates the vault directory structure and `.nd.yaml` config. Prefix is required -- it becomes part of every issue ID (e.g., `PROJ-a3f`).

### Issue Creation

```bash
nd create "Title" [flags]
  -t, --type         bug|feature|task|epic|chore|decision (default: task)
  -p, --priority     0-4, where 0=critical (default: 2)
  -d, --description  Issue description body
  --assignee         Assignee name
  --labels           Comma-separated labels
  --parent           Parent issue ID (for epic children)
```

IDs are generated as `PREFIX-HASH` (3 hex chars from SHA-256). Children use dot notation: `PROJ-a3f.1`, `PROJ-a3f.2`.

### Listing and Filtering

```bash
nd list [flags]
  --status       Filter: open, in_progress, blocked, deferred, closed
  --type         Filter by issue type
  --assignee     Filter by assignee
  --label        Filter by label
  --sort         Sort by: priority (default), created, updated, id
  -n, --limit    Max results
```

### Viewing Issues

```bash
nd show <id> [--short] [--json]
```

`--short` gives a one-line summary. `--json` outputs the full issue as JSON.

### Updating Issues

```bash
nd update <id> [flags]
  --status          New status
  --title           New title
  --priority        New priority (0-4 or P0-P4)
  --assignee        New assignee
  --type            New type
  -d, --description New description
  --append-notes    Append text to Notes section
```

### Closing and Reopening

```bash
nd close <id> [id...] [--reason="explanation"]
nd reopen <id>
```

Close accepts multiple IDs for batch operations. Closing sets `closed_at` and optionally `close_reason`. Reopening clears both.

### Dependencies

```bash
nd dep add <issue> <depends-on>    # issue depends on depends-on
nd dep rm <issue> <depends-on>     # Remove dependency
nd dep list <id>                   # Show all deps for an issue
```

Dependencies are bidirectional: `nd dep add A B` adds B to A's `blocked_by` AND A to B's `blocks`. Removing cleans both sides.

### Finding Work

```bash
nd ready [--assignee=NAME] [--sort=priority] [-n LIMIT]
nd blocked [--verbose]
```

`ready` shows issues with no open blockers. `blocked` shows issues waiting on dependencies. With `--verbose`, blocked also shows which issues are blocking each one.

### Search

```bash
nd search <query>
```

Full-text search across all issue files. Returns matching lines with 2 lines of context. Delegates to vlt's search engine.

### Labels

```bash
nd labels add <id> <label>
nd labels rm <id> <label>
nd labels list                # All labels with counts
```

### Comments

```bash
nd comments add <id> "Comment text"
nd comments list <id>
```

Comments are appended to the `## Comments` section with RFC3339 timestamps and author attribution.

### Epics

```bash
nd epic status <id>    # Progress summary (open/closed/blocked counts, %)
nd epic tree <id>      # Hierarchical tree view with status markers
```

Epic children are found by matching the `parent` field. Tree view uses status markers: `[ ]` open, `[>]` in progress, `[!]` blocked, `[x]` closed.

### Statistics

```bash
nd stats [--json]
```

Aggregate counts by status, type, and priority.

### AI Context

```bash
nd prime [--json]
```

Outputs a structured summary for AI context injection: total counts, ready work, blocked work, in-progress items. JSON mode includes all issues.

### Import from Beads

```bash
nd import --from-beads <path-to-jsonl>
```

Parses a beads JSONL export and creates vault issue files. Preserves timestamps, status, labels, notes, and design content.

### Vault Health

```bash
nd doctor [--fix]
```

Validates:
1. Content hash integrity (SHA-256 of body matches stored hash)
2. Bidirectional dependency consistency (A blocks B <-> B blocked_by A)
3. Reference validity (no deps pointing to nonexistent issues)
4. Field validation (required fields present, enums valid)

With `--fix`, automatically repairs hash mismatches and broken dependency references.

### Global Flags

All commands support:

```
--vault PATH    Override vault directory (default: .vault, auto-discovered)
--json          Output as JSON
--verbose       Verbose output
--quiet         Suppress non-essential output
```

## Priority System

| Priority | Label | Use for |
|----------|-------|---------|
| P0 | Critical | Security, data loss, broken builds |
| P1 | High | Major features, important bugs |
| P2 | Medium | Standard work (default) |
| P3 | Low | Polish, optimization |
| P4 | Backlog | Future ideas |

## Status Lifecycle

```
open --> in_progress --> closed
  |          |
  v          v
deferred   blocked --> in_progress --> closed
                              |
                              v
                           closed
```

Closed issues can only transition back to `open` via `nd reopen`.

## Issue Types

| Type | Use for |
|------|---------|
| `bug` | Defects and broken behavior |
| `feature` | New functionality |
| `task` | General work items |
| `epic` | Large initiatives with child issues |
| `chore` | Maintenance, tooling, housekeeping |
| `decision` | Architectural decision records |

## Comparison with beads (bd)

### What nd has that bd has

| Capability | bd | nd |
|---|---|---|
| Create/show/list/update/close/reopen | Yes | Yes |
| Dependencies (blocks/blocked_by) | Yes | Yes |
| Ready/blocked work discovery | Yes | Yes |
| Epic hierarchies | Yes | Yes |
| Labels | Yes | Yes |
| Comments | Yes | Yes |
| Search | Yes | Yes |
| Stats | Yes | Yes |
| Import from JSONL | Yes | Yes |
| Doctor/health check | Yes | Yes |
| Prime (AI context) | Yes | Yes |
| JSON output | Yes | Yes |
| Prefix rename | Yes | -- |

### What nd does differently

| Aspect | bd | nd |
|---|---|---|
| Storage | Dolt SQL database | Markdown files (YAML frontmatter) |
| Server | Requires `dolt sql-server` | No server. Files on disk. |
| Field size | 65KB TEXT limit | Unlimited (it's a file) |
| Inspectability | `bd show` or SQL queries | `cat`, `grep`, `git diff`, any editor |
| Sync | `bd dolt push/pull` | `git push/pull` (files are in your repo) |
| Dependencies | Go binary, Dolt | Go binary only |
| Vault compatibility | None | Obsidian-compatible (open in Obsidian) |
| Config | Multiple config sources | Single `.nd.yaml` |
| Migrations | SQL schema migrations | None needed (markdown is the schema) |
| Content integrity | Trust the database | SHA-256 content hashing |

### What bd has that nd does not (yet)

| Feature | Status |
|---|---|
| `bd q` (quick capture) | Not yet. Use `nd create` with `--quiet`. |
| `bd edit` (open in $EDITOR) | Not yet. Edit the .md file directly. |
| `bd graph` (ASCII visualization) | Not yet. Use `nd epic tree` for hierarchies. |
| `bd stale` (stale issue detection) | Not yet. Use `nd list --sort=updated`. |
| Molecules/chemistry system | Not planned. Use templates externally. |
| `bd sync` (Dolt sync) | Replaced by `git push/pull`. |
| `bd admin compact` (summarization) | Not yet. |
| `bd merge` (duplicate merging) | Not yet. |
| `bd rename-prefix` | Not yet. |
| External references | Not yet. |
| Date range filters | Not yet. |
| Hooks (SessionStart/PreCompact) | Via nd-skill plugin. |

## Architecture

```
nd (Go CLI, cobra)
  |
  cmd/               -- One file per command
  |
  internal/
    model/           -- Issue struct, Status/Priority/Type enums, validation
    idgen/           -- SHA-256 collision-resistant ID generation
    store/           -- Wraps vlt.Vault for issue CRUD, deps, filtering
    graph/           -- In-memory dependency graph: ready, blocked, cycles, epics
    enforce/         -- Content hashing, validation rules
    format/          -- Table, detail, JSON, prime context output
```

**Dependencies**: `cobra`, `vlt`. That's it.

## Testing

```bash
make test     # Unit + integration tests
make vet      # Go vet
make build    # Build binary
```

Unit tests cover model validation, ID generation, content hashing, and graph traversal. Integration tests create real temp vaults and run full workflows (init -> create -> dep -> ready -> close -> stats) with no mocks.

## License

Apache License 2.0. See [LICENSE](LICENSE).
