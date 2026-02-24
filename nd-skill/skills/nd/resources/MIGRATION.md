# Migration from Beads (bd)

How to migrate an existing beads database to nd.

## Prerequisites

- nd installed and in PATH
- Access to the beads JSONL export file (`.beads/issues.jsonl`)

## Export from Beads

If you don't already have a JSONL export:

```bash
# In your beads project directory
bd list --json > /tmp/beads-export.jsonl
```

Or use the existing JSONL file that beads auto-maintains:

```bash
ls .beads/issues.jsonl
```

## Import into nd

```bash
# Initialize nd vault
nd init --prefix=PROJ --vault=.vault

# Import beads issues
nd import --from-beads .beads/issues.jsonl

# Verify
nd stats
nd doctor
```

## What Gets Imported

| Field | Imported | Notes |
|-------|----------|-------|
| Title | Yes | Required; issues without titles are skipped |
| Description | Yes | Placed in ## Description section |
| Type | Yes | Mapped to nd types (default: task) |
| Priority | Yes | 0-4 or P0-P4 format accepted |
| Assignee | Yes | |
| Status | Yes | open, in_progress, blocked, deferred, closed, custom |
| Labels | Yes | Array of strings |
| Notes | Yes | Appended to ## Notes section |
| Design | Yes | Patched into ## Design section |
| Timestamps | Yes | created_at, updated_at, closed_at preserved |
| Close reason | Yes | |

## What Does NOT Get Imported

| Field | Reason |
|-------|--------|
| External refs | Not applicable to nd |
| Molecules/chemistry | Not applicable to nd |
| Dolt-specific metadata | Not applicable |

## Multi-Pass Import

The migration runs three passes:

### Pass 1: Create Issues

All issues from the JSONL export are created as vault markdown files. Fields, timestamps, labels, notes, design, and acceptance criteria are preserved verbatim. Original IDs are kept.

### Pass 2: Wire Dependencies

After all issues exist, dependencies are wired:

- **Parent-child** from explicit beads deps, dotted IDs (e.g., `EPIC-abc.3`), and cross-references in descriptions
- **Blocks/blocked_by** from beads blocking deps
- **Related** from beads `discovered-from`, `related`, and `relates-to` links
- Parents with children are auto-promoted to `epic` type

### Pass 3: Infer Execution Trajectories

Since beads never had `follows`/`led_to` chains, and nd's auto-follows only fires on live `in_progress` transitions (which never happen during migration since issues arrive pre-closed), Pass 3 reconstructs temporal execution chains from `closed_at` timestamps:

1. **Sibling chains under shared parents**: For each parent with closed children, sort children by `closed_at` ascending and chain consecutive pairs via `follows`/`led_to`. This captures the execution order within an epic.

2. **Related orphan chains**: For closed issues that share `related` edges and have no parent, chain them temporally via `follows` based on which was closed first.

3. **Epic-to-epic chains**: For closed epics, sort by the `closed_at` of their last closed child and chain them in execution order. This captures the project-level trajectory across milestones.

After migration, `nd path` shows the full execution history as connected chains instead of disconnected islands.

## Post-Import Steps

1. **Verify count**: `nd stats` should show issue counts matching `bd stats`
2. **Run doctor**: `nd doctor --fix` to fix any content hash mismatches
3. **Verify trajectories**: `nd path` should show execution chains
4. **Test workflow**: `nd ready`, `nd blocked`, `nd show <id>` to verify everything works

## Coexistence

nd and beads can coexist in the same project. nd uses `.vault/` while beads uses `.beads/`. You can run both simultaneously during a transition period.

## Differences to Be Aware Of

| Aspect | beads (bd) | nd |
|--------|-----------|-----|
| IDs | `bd-HASH` (5+ chars) | `PREFIX-HASH` (3 chars) |
| Storage | Dolt SQL database | Markdown files |
| Sync | `bd dolt push/pull` | `git push/pull` |
| Compact | `bd admin compact` | Not yet available |
| Quick capture | `bd q "title"` | `nd q "title"` |

After migration, you may want to update any scripts or hooks that reference `bd` commands to use `nd` equivalents.
