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

Use `--aynig-remote <name>` to scan a remote instead of local branches.

Example:

```bash
aynig run --aynig-remote origin
```

This is useful when runners operate as distributed workers and the remote branch is the source of truth.
