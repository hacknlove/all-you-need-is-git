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

Use `--dwp-remote <name>` to scan a remote instead of local branches.

If `--dwp-remote` is omitted, AYNIG checks the latest commit trailer `dwp-source: git:<name>` and uses that remote when present.

Example:

```bash
aynig run --dwp-remote origin
```

This is useful when runners operate as distributed workers and the remote branch is the source of truth.

## Role-specific commands

Use `--role <name>` (or set `AYNIG_ROLE`) to run commands from a role-specific directory.

Resolution order:

1. `.dwp/roles/<role>/command/<state>`
2. `.dwp/command/<state>`
