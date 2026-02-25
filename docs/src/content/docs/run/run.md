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
- `--log-level <level>`: `debug`, `info`, `warn`, or `error` (default)

In `--use-remote` mode, `--current-branch` resolves against the upstream branch of your local current branch (for example `origin/main`). If no upstream exists, `only` runs zero branches.

Log level precedence: `--log-level` > `aynig-log-level` trailer > `AYNIG_LOG_LEVEL` env.

Example:

```bash
aynig run --use-remote origin --current-branch include
```

## How dispatch works

Commands live under `.aynig/command/<state>`. The `aynig-state` trailer selects which command is executed.

Branch logs use the resolved log level after trailers are parsed. Early branch logs are buffered and flushed once the level is known.

Command stdout/stderr is written to `.aynig/logs/<commit-hash>.log`, where `<commit-hash>` is the commit that triggered the command.
