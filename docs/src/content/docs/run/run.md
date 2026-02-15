---
title: Run AYNIG
description: Execute AYNIG against a repository.
---

```bash
aynig run
```

## Options

- `-w, --worktree <path>`: worktree directory (default: `.worktrees`)
- `--use-remote <name>`: use remote branches instead of local
- `--current-branch <mode>`: `skip` (default), `include`, or `only`

Example:

```bash
aynig run --use-remote origin --current-branch include
```

## How dispatch works

Commands live under `.aynig/command/<state>`. The `aynig-state` trailer selects which command is executed.
