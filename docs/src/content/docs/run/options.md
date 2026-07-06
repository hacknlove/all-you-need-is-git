---
title: Run Options
description: Practical flags for scanning branches and using remotes.
---

## Current branch handling

By default, `aynig run` skips the current branch.

Use `--current-branch` to control this behavior:

- `skip` (default): ignore the branch you are currently on
- `include`: scan the current branch in addition to others
- `only`: scan only the current branch

Examples:

```bash
# test a workflow on your current branch
aynig run --current-branch only
```

## Use remote branches

Use `--remote <name>` to scan a remote instead of local branches.

If `--remote` is omitted, AYNIG checks the latest commit trailer `dwp-source: git:<name>` and uses that remote when present.

Example:

```bash
aynig run --remote origin
```

This is useful when runners operate as distributed workers and the remote branch is the source of truth.

## Trusted commands ref

Commands are resolved from a trusted ref, not from the branch that emitted the
event. By default that is the repository's default branch (`<remote>/HEAD` in
remote mode, `master` or `main` locally).

This means a feature branch can *invoke* states but cannot *redefine* what a
state executes — changing a command requires landing it on the default branch.

Use `--commands-ref` to override:

```bash
# pin to an immutable commit
aynig run --commands-ref 3f2a9c1

# test workflow changes from a feature branch
aynig run --commands-ref my-feature --current-branch only

# resolve commands from each event branch (trusted environments only)
aynig run --commands-ref same
```

Commands receive `COMMANDS_PATH` pointing to the `.aynig` directory they were
resolved from, while `WORKTREE_PATH` and the working directory remain the
event branch's worktree.

## Role-specific commands

Use `--role <name>` (or set `ROLE`) to run commands from a role-specific directory.

Resolution order:

1. `.aynig/roles/<role>/command/<state>`
2. `.aynig/command/<state>`
