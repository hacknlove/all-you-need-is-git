---
title: Run AYNIG
description: Execute one tick of the workflow for actionable branches.
---

`aynig run` scans branches (local by default) and inspects the **latest commit** of each branch.

> Note: by default, the runner **skips the current branch**. If you're testing on the branch you're currently on, pass `--current-branch only`.

If `HEAD` contains an `aynig-state:` trailer, AYNIG will dispatch to:

```text
.aynig/command/<state>
```

## Commit format

Example commit message:

```text
chore: whatever title for humans

Here goes the prompt to be processed by the agent.

aynig-state: some-state
foo: bar
baz: qux
```

## Options

- `-w, --worktree <path>` — worktree directory (default: `.worktrees`)
- `--aynig-remote <name>` — use remote branches instead of local
- `--current-branch <mode>` — `skip` (default), `include`, or `only`
- `--log-level <level>` — `debug`, `info`, `warn`, or `error` (default)

In `--aynig-remote` mode, `--current-branch` resolves against the upstream branch of your local current branch (for example `origin/main`). If no upstream exists, `only` runs zero branches.

If `--aynig-remote` is omitted, AYNIG also checks the latest commit trailer `aynig-remote: <name>` and uses that remote when present.

Log level precedence: `--log-level` > `aynig-log-level` trailer > `AYNIG_LOG_LEVEL` env.

## Environment variables

Commands receive metadata via env vars such as:

- `AYNIG_BODY`
- `AYNIG_COMMIT_HASH`
- `AYNIG_TRAILER_FOO`
- `AYNIG_TRAILER_BAZ`

(See also: Commands → Environment Variables.)

Branch logs use the resolved log level after trailers are parsed. Early branch logs are buffered and flushed once the level is known.

Command stdout/stderr is written to `.aynig/logs/<commit-hash>.log`, where `<commit-hash>` is the commit that triggered the command.
