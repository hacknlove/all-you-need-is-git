---
title: Repository Layout
description: Recommended layout for AYNIG-enabled repos.
---

A minimal AYNIG repository typically includes:

```text
.aynig/
  command/
    <state>
  roles/
    <role>/
      command/
        <state>
  COMMANDS.md   # optional, documents available states
.worktrees/     # ephemeral; created/cleaned by AYNIG
```

The `aynig-state` trailer selects which command runs. When `AYNIG_ROLE` (or `--role`) is set, AYNIG checks `.aynig/roles/<role>/command/<state>` first.

## .worktrees/

AYNIG runs each command inside a dedicated Git worktree. Worktrees are created under this directory and reused across runs when possible. This directory is ignored via `.gitignore`.

## .aynig/logs/

Command stdout/stderr logs are written here as `<commit-hash>.log`. This directory is ignored via `.gitignore`.
Keep `.worktrees/` ignored. Commands should not create/manage worktrees manually.
