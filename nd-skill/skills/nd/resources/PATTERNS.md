# Usage Patterns for AI Agents

Common patterns for AI agents using nd in coding workflows.

## Pattern: Claim and Complete

The basic work cycle:

```bash
nd ready                                    # Find work
nd show PROJ-a3f                            # Read context
nd update PROJ-a3f --status=in_progress     # Claim
# ... implement ...
nd close PROJ-a3f --reason="Implemented and tested"
nd ready                                    # Check what unblocked
```

## Pattern: Discovery During Work

While implementing one issue, you find another problem:

```bash
# Working on PROJ-a3f, discover a bug
nd create "Input not sanitized in auth handler" --type=bug --priority=1
# Returns: PROJ-c4d

# Link discovery to current work
nd dep add PROJ-c4d PROJ-a3f    # New bug depends on current work? Or:
# If the bug BLOCKS current work:
nd dep add PROJ-a3f PROJ-c4d    # Current work depends on the bug fix

# Assess: is it a blocker?
# YES: pause PROJ-a3f, work on PROJ-c4d first
# NO: continue PROJ-a3f, PROJ-c4d persists for later
```

## Pattern: Systematic Exploration

Research or investigation work:

```bash
nd create "Investigate caching options" --type=task --priority=2 -d "Evaluate Redis vs in-memory vs file-based caching"
nd update PROJ-b7c --status=in_progress
# Research...
nd update PROJ-b7c --append-notes="Redis: best for distributed. In-memory: simplest. Chose Redis for scalability."
nd create "Implement Redis cache layer" --type=feature --priority=1
nd close PROJ-b7c --reason="Research complete. Created follow-up: PROJ-d9e"
```

## Pattern: Bug Investigation

```bash
nd create "Login fails with special characters" --type=bug --priority=1
nd update PROJ-a3f --status=in_progress
# Investigate...
nd update PROJ-a3f --append-notes="Root cause: SQL injection in WHERE clause. Fix: use parameterized queries."
# Fix...
nd close PROJ-a3f --reason="Fixed with parameterized queries. Added integration test."
```

## Pattern: Refactoring with Dependencies

```bash
# Plan the refactoring steps
nd create "Extract auth middleware" --type=task --priority=2        # PROJ-a3f
nd create "Move user model to shared package" --type=task --priority=2  # PROJ-b7c
nd create "Update all route handlers" --type=task --priority=2     # PROJ-c4d
nd create "Remove deprecated auth code" --type=chore --priority=3  # PROJ-d9e

# Add execution ordering
nd dep add PROJ-b7c PROJ-a3f    # Move model needs middleware first
nd dep add PROJ-c4d PROJ-b7c    # Update routes needs model moved
nd dep add PROJ-d9e PROJ-c4d    # Cleanup needs routes updated

# Work through in order
nd ready    # Shows PROJ-a3f (foundation)
# ... work through each, nd ready advances the front
```

## Pattern: Multi-Issue Close

After a productive session:

```bash
nd close PROJ-a3f PROJ-b7c PROJ-c4d --reason="Batch: auth refactoring complete"
```

## Pattern: Notes as Breadcrumbs

Leave breadcrumbs during work for compaction survival:

```bash
# After each significant step
nd update PROJ-a3f --append-notes="Step 1 done: extracted middleware to pkg/auth/middleware.go"
nd update PROJ-a3f --append-notes="Step 2 done: all 12 handlers updated to use new middleware"
nd update PROJ-a3f --append-notes="BLOCKED: found circular import between pkg/auth and pkg/user"
```

After compaction, `nd show PROJ-a3f` gives the full trail.

## Pattern: Using nd with Git Workflow

Issues are files -- they live in your git repo:

```bash
# Create branch for a feature
git checkout -b feature/auth

# Work with nd as usual
nd create "Auth feature" --type=feature --priority=1
nd update PROJ-a3f --status=in_progress
# ... implement ...
nd close PROJ-a3f --reason="Done"

# Commit everything (code + issues)
git add .vault/issues/
git add src/
git commit -m "feat: implement auth system"
git push
```

Issue changes show up in PRs alongside code changes. Reviewers can see what work was tracked.

## Anti-Patterns

**Don't**: Create issues for every tiny task within a session.
**Do**: Use TodoWrite for within-session checklists, nd for cross-session work.

**Don't**: Leave issues in_progress across sessions without notes.
**Do**: Always update notes before ending a session.

**Don't**: Create dependencies for loosely related work.
**Do**: Only add deps when there's a true execution ordering requirement.

**Don't**: Skip `nd ready` at session start.
**Do**: Always check ready queue first -- it's the source of truth for what to work on.
