# Dependency System

How dependencies work in nd and how to use them correctly.

## Semantics

`nd dep add A B` means **"A depends on B"** (B must complete before A can start).

This creates a bidirectional relationship:
- B is added to A's `blocked_by` list
- A is added to B's `blocks` list

When B is closed, A becomes unblocked. If A has no other open blockers, it appears in `nd ready`.

## Dependency Direction

Think in terms of **requirements**, not **sequence**:

```
CORRECT thinking:   "Tests NEED the feature"       --> nd dep add tests feature
INCORRECT thinking: "Feature BEFORE tests"          --> nd dep add feature tests  (WRONG!)
```

The mnemonic: **"X needs Y"** maps to `nd dep add X Y`.

### Common Mistake: Inverted Dependencies

```
Situation: Feature must be done before tests can run.

Wrong:  nd dep add feature tests     # "feature depends on tests" -- backwards!
Right:  nd dep add tests feature     # "tests depend on feature" -- correct!
```

Verify with `nd blocked`:
- Tests blocked BY feature: correct
- Feature blocked BY tests: inverted, fix it

## Ready Fronts

A **Ready Front** is the set of issues with all dependencies satisfied -- what can be worked on right now. As issues close, the front advances.

```
Ready Front 1:  Setup (foundation, no deps)
Ready Front 2:  Core logic (needs setup)
Ready Front 3:  API + UI (parallel, both need core logic)
Ready Front 4:  Integration tests (needs API and UI)
```

`nd ready` always shows the current front. Closing work in the current front advances the front automatically.

## Epic Planning with Dependencies

Walk **backward** from the goal to get correct dependencies:

```
Start: "What's the final deliverable?"
       --> Integration tests passing

"What does that need?"
       --> Streaming support, Header display

"What do those need?"
       --> Message rendering

"What does that need?"
       --> Buffer layout (foundation, no deps)
```

This produces:
```bash
nd dep add integration streaming
nd dep add integration header
nd dep add streaming messages
nd dep add header messages
nd dep add messages buffer
# buffer has no deps -- it's Ready Front 1
```

## Removing Dependencies

```bash
nd dep rm A B    # A no longer depends on B
```

Cleans both sides: removes B from A's `blocked_by` AND A from B's `blocks`.

## Listing Dependencies

```bash
nd dep list PROJ-a3f
# Output:
# PROJ-a3f depends on:
#   PROJ-b7c [open] OAuth setup (P1)
# PROJ-a3f blocks:
#   PROJ-d9e [open] Integration tests (P2)
```

## Related Links

Related links are **soft, bidirectional connections** that don't affect the ready queue. Use them for context and discoverability.

```bash
nd dep relate PROJ-a3f PROJ-b7c     # Add bidirectional related link
nd dep unrelate PROJ-a3f PROJ-b7c   # Remove bidirectional related link
```

**When to use related vs blocks:**
- `blocks`: "A literally cannot proceed until B is done" (hard blocker)
- `related`: "A and B are connected but either can proceed independently" (soft link)

Examples of related: alternative approaches, documentation alongside code, refactoring efforts in similar areas.

## Cycle Detection

```bash
nd dep cycles     # Detect dependency cycles
```

A cycle means no issue in the cycle can ever become ready. The command returns the cycle path so you can identify which dependency to remove.

`nd doctor` also detects cycles as part of its validation checks.

## Dependency Tree

```bash
nd dep tree PROJ-a3f     # Show dependency tree from issue
```

Displays the forward dependency tree (what the issue blocks, recursively) as an ASCII tree in the terminal.
