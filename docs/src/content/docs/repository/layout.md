---
title: Repository Layout
description: Understand the folders AYNIG manages.
---

## .aynig/

AYNIG stores repo-level configuration and helper scripts here. By default it includes a starter `COMMANDS.md`, a `clean` script, and `CONTRACT.md`.

## .aynig/command/

Executable commands keyed by state name:

```
.aynig/command/build
.aynig/command/review
.aynig/command/test
```

The `aynig-state` trailer selects which command runs.

## .worktrees/

AYNIG runs each command inside a dedicated Git worktree. Worktrees are created under this directory and reused across runs when possible. This directory is ignored via `.gitignore`.

## .aynig/logs/

Command stdout/stderr logs are written here as `<commit-hash>.log`. This directory is ignored via `.gitignore`.
