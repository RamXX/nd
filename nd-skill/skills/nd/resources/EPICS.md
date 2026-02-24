# Epic Hierarchies

How to structure and track large initiatives with nd epics.

## What is an Epic

An epic is an issue of type `epic` that serves as a parent for child issues. Children reference the parent via the `parent` frontmatter field. Epics provide:

- Hierarchical organization for related work
- Aggregate progress tracking (% complete)
- Tree visualization of work breakdown

## Creating an Epic with Children

```bash
# Create the epic
nd create "Authentication System" --type=epic --priority=1

# Create children referencing the epic as parent
nd create "OAuth setup" --type=task --priority=1 --parent=PROJ-a3f
nd create "Login flow" --type=feature --priority=1 --parent=PROJ-a3f
nd create "Token refresh" --type=task --priority=2 --parent=PROJ-a3f
nd create "Auth tests" --type=task --priority=2 --parent=PROJ-a3f
```

## Viewing Epic Status

```bash
nd epic status PROJ-a3f
```

Output:
```
Epic: PROJ-a3f - Authentication System
Children: 4 total
  Open:        2
  In Progress: 1
  Blocked:     0
  Closed:      1
  Progress:    25%
```

Progress is `closed / total * 100`.

## Viewing Epic Tree

```bash
nd epic tree PROJ-a3f
```

Output:
```
[ ] PROJ-a3f Authentication System (P1)
  [x] PROJ-b7c OAuth setup (P1)
  [>] PROJ-d9e Login flow (P1)
  [ ] PROJ-f2a Token refresh (P2)
  [!] PROJ-c4d Auth tests (P2)
```

Status markers:
- `[ ]` -- open
- `[>]` -- in_progress
- `[!]` -- blocked
- `[x]` -- closed

## Combining Epics with Dependencies

Epics and dependencies are orthogonal. An epic provides organizational hierarchy; dependencies provide execution ordering. Use both:

```bash
# Create epic + children
nd create "Auth System" --type=epic --priority=1       # PROJ-a3f
nd create "OAuth setup" --parent=PROJ-a3f              # PROJ-b7c
nd create "Login flow" --parent=PROJ-a3f               # PROJ-d9e
nd create "Auth tests" --parent=PROJ-a3f               # PROJ-f2a

# Add execution ordering
nd dep add PROJ-d9e PROJ-b7c    # Login needs OAuth setup
nd dep add PROJ-f2a PROJ-d9e    # Tests need Login flow

# View: epic tree shows hierarchy, nd blocked shows execution order
nd epic tree PROJ-a3f
nd ready                        # Only PROJ-b7c is ready (foundation)
```

## Finding Close-Eligible Epics

```bash
nd epic close-eligible
```

Lists epics where all children are closed. These are candidates for closing the parent epic.

## Listing Children

```bash
nd children PROJ-a3f
```

Lists all child issues of a parent (equivalent to `nd list --parent=PROJ-a3f`).

## When to Use Epics

Use epics when:
- A feature has 3+ related child issues
- You want aggregate progress tracking
- Work needs hierarchical organization for navigation

Don't use epics when:
- Only 1-2 related issues (just use dependencies)
- Issues are loosely related (use labels instead)
- You need execution ordering only (use dependencies alone)
