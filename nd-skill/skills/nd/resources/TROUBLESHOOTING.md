# Troubleshooting

Common problems and fixes when using nd.

## Content Hash Mismatch

**Symptom**: `nd doctor` reports `[HASH] ID: content hash mismatch`

**Cause**: The body was modified (e.g., comment added, section edited, history entry appended) but the stored `content_hash` was not recalculated. This happens when comments are added via `nd comments add`, history entries are appended during state transitions, or when files are edited directly.

**Fix**:
```bash
nd doctor --fix
```

This recalculates all content hashes to match current body content.

## Bidirectional Dependency Inconsistency

**Symptom**: `nd doctor` reports `[SYNC] A blocks B, but B doesn't list A in blocked_by`

**Cause**: A dependency was partially written (crash, manual edit, concurrent access). nd dependencies are bidirectional -- both files must be updated.

**Fix**:
```bash
nd doctor --fix    # Restores bidirectional consistency
```

## Orphan Dependency References

**Symptom**: `nd doctor` reports `[DEP] A blocked by B, but B does not exist`

**Cause**: An issue was deleted or its ID was changed, but references to it remain in other issues.

**Fix**:
```bash
nd doctor --fix    # Removes references to nonexistent issues
```

## "vault already initialized"

**Symptom**: `nd init` fails with "vault already initialized"

**Cause**: `.vault/.nd.yaml` already exists.

**Fix**: This is expected if the vault was already initialized. Use `nd list` to verify it works. If the vault is corrupted, delete `.vault/` and re-initialize.

## "no frontmatter found"

**Symptom**: `nd show <id>` fails with "no frontmatter found"

**Cause**: The issue file is missing its YAML frontmatter (the `---` delimited block at the top).

**Fix**: Edit the file directly and add frontmatter. See [STORAGE.md](STORAGE.md) for the correct format.

## "issue not found"

**Symptom**: `nd show <id>` or `nd update <id>` can't find the issue.

**Cause**: The issue file doesn't exist at `.vault/issues/<id>.md`, or the vault directory is wrong.

**Fix**:
```bash
# Check vault location
nd list --vault=/path/to/vault

# List all issue files
ls .vault/issues/
```

## Links Section Inconsistency

**Symptom**: `nd doctor` reports `[LINKS] ID: Links section missing or inconsistent`

**Cause**: The `## Links` section (containing wikilinks like `[[PROJ-b7c]]`) doesn't match the frontmatter relationships (parent, blocks, blocked_by, related, follows, led_to).

**Fix**:
```bash
nd doctor --fix    # Regenerates Links sections from frontmatter
```

This also happens after `nd edit` if the Links section was manually changed. The `nd edit` command auto-refreshes Links on save.

## Empty Ready Queue

**Symptom**: `nd ready` shows "No issues found."

**Possible causes**:
1. All issues are closed
2. All open issues are blocked
3. No issues exist

**Diagnosis**:
```bash
nd stats           # Check total counts
nd blocked         # See what's stuck
nd list            # See all issues regardless of status
```

## Performance with Many Issues

nd reads all issue files for `list`, `ready`, `blocked`, and `stats` commands. For vaults with 500+ issues, this uses parallel goroutines for reading. If performance is a concern:

- Use `--limit` / `-n` flags to cap results
- Use `--status` filters to reduce the scan scope
- The scan is bounded by I/O, not CPU

## vlt Dependency

nd requires [vlt](https://github.com/RamXX/vlt) as a Go library dependency. It is compiled into the nd binary -- no separate vlt installation is needed. However, if you want to perform advanced vault operations (wikilink management, template application, bookmark manipulation), install the vlt CLI separately and consult the **vlt skill**.
