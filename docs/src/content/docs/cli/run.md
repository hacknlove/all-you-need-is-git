---
title: run
description: Run AYNIG for the current repository.
---

```bash
aynig run [options]
```

Options:

- `-w, --worktree <path>`: worktree directory (default: `.worktrees`)
- `--use-remote <name>`: use remote branches instead of local
- `--current-branch <mode>`: `skip` (default), `include`, or `only`

The runner executes commands stored in `.aynig/command/<state>` based on the `aynig-state` trailer in the latest commit.
