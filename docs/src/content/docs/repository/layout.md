---
title: Repository Layout
description: Recommended layout for AYNIG-enabled repos.
---

A minimal AYNIG repository typically includes:

```text
.aynig/
  command/
    <state>
  COMMANDS.md   # optional, documents available states
.worktrees/     # ephemeral; created/cleaned by AYNIG
```

The `aynig-state` trailer selects which command runs.

## .worktrees/

AYNIG runs each command inside a dedicated Git worktree. Worktrees are created under this directory and reused across runs when possible. This directory is ignored via `.gitignore`.

## .aynig/logs/

Command stdout/stderr logs are written here as `<commit-hash>.log`. This directory is ignored via `.gitignore`.
Keep `.worktrees/` ignored. Commands should not create/manage worktrees manually.
