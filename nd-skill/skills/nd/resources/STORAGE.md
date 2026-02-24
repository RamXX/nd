# Storage Format

How nd stores issues on disk. This is the canonical reference for the file format.

## Vault Layout

```
.vault/                    # Default location (override with --vault)
  .nd.yaml                 # Config: version, prefix, created_by
  issues/                  # Flat directory, one .md file per issue
    PROJ-a3f.md
    PROJ-b7c.md
  .trash/                  # Soft-deleted issues
  .vlt.lock                # Advisory file lock (managed by vlt)
```

## Config File (.nd.yaml)

```yaml
version: "1"
prefix: PROJ
created_by: alice
status_custom: "review,qa,rejected"
status_sequence: "open,in_progress,review,qa,closed"
status_fsm: true
status_exit_rules: "blocked:open,in_progress;rejected:in_progress"
```

Manage via `nd config set/get/list` or edit directly. See [CLI_REFERENCE.md](CLI_REFERENCE.md#configuration) for config keys.

## Issue File Format

Each issue is a markdown file with YAML frontmatter:

```yaml
---
id: PROJ-a3f
title: "Implement user authentication"
status: open
priority: 1
type: feature
assignee: alice
labels: [security, milestone]
parent: ""
blocks: [PROJ-d9e]
blocked_by: [PROJ-b3c]
related: [PROJ-f2a]
created_at: 2026-02-23T20:15:00Z
created_by: alice
updated_at: 2026-02-24T10:30:00Z
closed_at: ""
close_reason: ""
content_hash: "sha256:a3f8c9d2e1b4..."
---

## Description
Implement OAuth 2.0 authentication with JWT tokens.

## Acceptance Criteria
- [ ] Login endpoint returns JWT
- [ ] Token refresh works

## Design
Using bcrypt with 12 rounds per OWASP recommendation.

## Notes
Spike complete. Chose Authorization Code flow.

## Comments

### 2026-02-23T20:15:00Z alice
Started implementation.
```

## Frontmatter Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier: PREFIX-HASH |
| `title` | string (quoted) | Yes | Issue title |
| `status` | enum | Yes | open, in_progress, blocked, deferred, closed (+ custom) |
| `priority` | int (0-4) | Yes | 0=critical, 4=backlog |
| `type` | enum | Yes | bug, feature, task, epic, chore, decision |
| `assignee` | string | No | Assigned person |
| `labels` | string[] | No | Labels (inline YAML array) |
| `parent` | string | No | Parent issue ID |
| `blocks` | string[] | No | IDs this issue blocks |
| `blocked_by` | string[] | No | IDs blocking this issue |
| `related` | string[] | No | Related issue IDs |
| `was_blocked_by` | string[] | No | Historical blockers (after removal) |
| `defer_until` | string | No | Target date for deferred issues (YYYY-MM-DD) |
| `created_at` | RFC3339 | Yes | Creation timestamp |
| `created_by` | string | Yes | Creator name |
| `updated_at` | RFC3339 | Yes | Last update timestamp |
| `closed_at` | RFC3339 | No | When closed |
| `close_reason` | string (quoted) | No | Why closed |
| `content_hash` | string (quoted) | Yes | SHA-256 of body content |

## Body Sections

The body (below frontmatter) contains these standard sections:

| Section | Purpose | Modified by |
|---------|---------|-------------|
| `## Description` | What and why | `nd create -d`, `nd update -d` |
| `## Acceptance Criteria` | Definition of done | Manual edit |
| `## Design` | Design decisions, architecture | Manual edit, import |
| `## Notes` | Working notes | `nd update --append-notes` |
| `## Links` | Wikilinks derived from relationships | Auto-maintained by nd |
| `## Comments` | Timestamped discussion | `nd comments add` |

## ID Generation

IDs use the format `PREFIX-HASH`:
- PREFIX: configured in `.nd.yaml`
- HASH: 3 hex characters from SHA-256 of title + timestamp + nonce
- Child IDs: `PREFIX-HASH.N` (sequential, e.g., PROJ-a3f.1, PROJ-a3f.2)

## Content Hash

The `content_hash` field stores a SHA-256 hash of the body content (everything below the frontmatter). This allows `nd doctor` to detect if the body was modified outside of nd (or if a hash was not updated after a body change).

## File Operations

nd uses [vlt](https://github.com/RamXX/vlt) for all file I/O:
- `v.Create()` for new issues
- `v.Read(title, heading)` for reading (supports heading-scoped reads)
- `v.Write()` for body replacement
- `v.Append()` for adding comments
- `v.Patch()` for surgical section edits
- `v.PropertySet()`/`PropertyRemove()` for frontmatter field updates
- `v.Files("issues", "md")` for listing all issues
- `v.Delete()` for soft delete to `.trash/`

For advanced vault operations beyond what nd exposes, consult the **vlt skill**.
