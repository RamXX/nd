# nd Skill Maintenance Guide

This file documents how the nd skill is structured and when to update it.

## Structure

```
nd-skill/
  .claude-plugin/
    plugin.json              # Plugin metadata, hooks (SessionStart, PreCompact)
  skills/
    nd/
      SKILL.md               # Main skill entry point (lean, ~100 lines)
      CLAUDE.md              # This file (maintenance guide)
      resources/
        CLI_REFERENCE.md     # Complete command syntax and flags
        WORKFLOWS.md         # Session start, compaction recovery, handoff
        ISSUE_CREATION.md    # When and how to create issues
        DEPENDENCIES.md      # Dependency semantics and epic planning
        EPICS.md             # Epic hierarchies and tree views
        STORAGE.md           # File format, frontmatter schema, vault layout
        MIGRATION.md         # Migrating from beads (bd)
        TROUBLESHOOTING.md   # Common problems and fixes
        PATTERNS.md          # Usage patterns for AI agents
```

## DRY Principle

The skill follows the same DRY pattern as the beads skill:

- **SKILL.md** contains only decision frameworks, session protocol, and resource index
- **nd prime** is the single source of truth for AI context (auto-loaded by hooks)
- **nd \<command\> --help** provides specific command usage
- **resources/** provide depth beyond what --help covers (conceptual frameworks, decision trees, advanced patterns)
- CLI syntax is NOT duplicated in SKILL.md (it's in CLI_REFERENCE.md and nd --help)

## What Belongs Where

| Content type | Location |
|---|---|
| When to use nd vs TodoWrite | SKILL.md |
| Session protocol (numbered steps) | SKILL.md |
| Resource index | SKILL.md |
| Complete command syntax | resources/CLI_REFERENCE.md |
| Conceptual frameworks | resources/ (topic-specific) |
| Decision trees | resources/ (topic-specific) |
| Error handling | resources/TROUBLESHOOTING.md |

## When to Update

### New nd command added
1. Add to resources/CLI_REFERENCE.md
2. If it introduces a new concept, create a resource file
3. Update SKILL.md resource table if new resource added
4. Bump version in SKILL.md frontmatter and plugin.json

### nd command changed
1. Update resources/CLI_REFERENCE.md
2. Update affected resource files if behavior changed

### New concept or pattern
1. Create resource file in resources/
2. Add to SKILL.md resource table
3. Cross-link from related resources

### Bug fix in nd (no new features)
1. No skill changes needed unless it affects documented behavior

## vlt Reference

nd is built on [vlt](https://github.com/RamXX/vlt). The vlt skill covers advanced vault operations that nd doesn't expose directly. Reference the vlt skill when agents need to:

- Manipulate frontmatter beyond what nd PropertySet supports
- Use wikilinks, backlinks, or link analysis
- Apply templates
- Work with tags, bookmarks, or daily notes
- Perform vault-wide search with custom options

Point agents to the vlt skill with: "For advanced vault operations, consult the **vlt skill**."
