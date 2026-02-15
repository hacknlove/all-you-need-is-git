---
title: Repository Layout
description: Understand the folders AYNIG manages.
---

## .aynig/

AYNIG stores repo-level configuration and helper scripts here. By default it includes a starter `COMMANDS.md` and a `clean` script.

## .aynig/command/

Executable commands keyed by state name:

```
.aynig/command/build
.aynig/command/review
.aynig/command/test
```

The `aynig-state` trailer selects which command runs.

## .worktrees/

AYNIG runs each command inside a dedicated Git worktree. Worktrees are created and removed automatically, and this directory is ignored via `.gitignore`.
