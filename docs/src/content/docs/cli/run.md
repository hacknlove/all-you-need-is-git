---
title: CLI — run
description: Run AYNIG for the current repository.
---

```bash
aynig run [options]
```

Options:

- `-w, --worktree <path>` — worktree directory (default: `.worktrees`)
- `--aynig-remote <name>` — use remote branches instead of local (e.g. `origin`)
- `--current-branch <mode>` — `skip` (default), `include`, or `only`
- `--log-level <level>` — `debug`, `info`, `warn`, or `error` (default)

In `--aynig-remote` mode, `--current-branch` resolves against the upstream branch of your local current branch (for example `origin/main`). If no upstream exists, `only` runs zero branches.

If `--aynig-remote` is omitted, AYNIG also checks the latest commit trailer `aynig-remote: <name>` and uses that remote when present.

Log level precedence: `--log-level` > `aynig-log-level` trailer > `AYNIG_LOG_LEVEL` env.

The runner reads `HEAD`, dispatches to `.aynig/command/<state>`, and validates the new `HEAD`.

Command stdout/stderr is written to `.aynig/logs/<commit-hash>.log`, where `<commit-hash>` is the commit that triggered the command.
