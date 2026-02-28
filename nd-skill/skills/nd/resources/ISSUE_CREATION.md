# Issue Creation Guide

When and how to create issues with nd.

## When to Create Issues

**Create directly** (no need to ask):
- Clear bug found during implementation
- Obvious follow-up work with defined scope
- Technical debt with concrete remediation
- Acceptance criteria items that are clearly separate work

**Ask the user first**:
- Knowledge work with fuzzy scope
- User intent is unclear
- Multiple valid approaches exist
- Significant effort required

## Discovery Workflow

When encountering new work during implementation:

```
Discovery:
- [ ] Notice bug, improvement, or follow-up work
- [ ] Assess: Can defer or is blocker?
- [ ] Create issue: nd create "Title" --type=bug --priority=1  (or: nd create --title="Title" ...)
- [ ] Add dependency: nd dep add <new-id> <current-id>
- [ ] If blocker: pause current work, switch to new issue
- [ ] If deferrable: continue current work, new issue persists for later
```

**Pattern**: File issues immediately when context is fresh. Don't rely on memory.

## Creating Good Issues

### Title

Imperative form, specific, scannable:
- Good: "Add rate limiting to auth endpoints"
- Bad: "Auth stuff" or "Rate limiting issue"

### Description

What needs to happen and why:
```bash
nd create "Add rate limiting to auth endpoints" \
  --type=task --priority=1 \
  -d "Login and password reset endpoints have no rate limiting.
An attacker can brute-force credentials with unlimited attempts.
Add sliding window rate limit: 5 attempts per minute per IP."
```

### Priority Selection

| Priority | When to use |
|----------|-------------|
| P0 | Active security vulnerability, data loss, production down |
| P1 | Important feature for current milestone, significant bug |
| P2 | Standard work, default for most issues |
| P3 | Nice to have, polish, minor optimization |
| P4 | Ideas, future work, "someday" items |

### Type Selection

| Type | When to use |
|------|-------------|
| `bug` | Something is broken or behaves incorrectly |
| `feature` | New functionality that doesn't exist yet |
| `task` | General work (refactoring, docs, testing) |
| `epic` | Large initiative with 3+ child issues |
| `chore` | Maintenance, dependencies, tooling |
| `decision` | Architectural or design decision to document |

## Batch Creation

For complex features, create all issues then add dependencies:

```bash
# Create epic
nd create "Auth System" --type=epic --priority=1

# Create children (use --parent for hierarchy)
nd create "OAuth client setup" --type=task --priority=1 --parent=PROJ-a3f
nd create "Auth code flow" --type=feature --priority=1 --parent=PROJ-a3f
nd create "Token storage" --type=task --priority=1 --parent=PROJ-a3f
nd create "Login endpoints" --type=feature --priority=1 --parent=PROJ-a3f

# Add dependencies (think: "X needs Y")
nd dep add <login> <token-storage>
nd dep add <token-storage> <auth-flow>
nd dep add <auth-flow> <oauth-setup>

# Verify
nd epic tree PROJ-a3f
nd blocked
```

## Labels Convention

Use labels for cross-cutting concerns:
- `security`, `performance`, `ux` -- concern area
- `critical`, `urgent` -- attention flags
- `milestone-v1`, `milestone-v2` -- release targeting
- `spike` -- investigation/research work
