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
- `--log-level <level>`: `debug`, `info`, `warn`, or `error` (default)

In `--use-remote` mode, `--current-branch` resolves against the upstream branch of your local current branch (for example `origin/main`). If no upstream exists, `only` runs zero branches.

Log level precedence: `--log-level` > `aynig-log-level` trailer > `AYNIG_LOG_LEVEL` env.

The runner executes commands stored in `.aynig/command/<state>` based on the `aynig-state` trailer in the latest commit.
