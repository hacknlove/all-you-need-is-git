---
title: Run AYNIG
description: Execute one workflow step for actionable branches.
---

`aynig run` scans branches (local by default) and inspects the **latest commit** of each branch.

> Note: by default, the runner **skips the current branch**. If you're testing on the branch you're currently on, pass `--current-branch only`.

If `HEAD` contains a `dwp-state:` trailer, AYNIG will run:

```text
.aynig/command/<state>
```

## Commit format

Example commit message:

```text
chore: whatever title for humans

Here goes the prompt to be processed by the agent.

dwp-state: some-state
foo: bar
baz: qux
```

## Options

- `-w, --worktree <path>` — worktree directory (default: `.worktrees`)
- `--remote <name>` — use remote branches instead of local
- `--current-branch <mode>` — `skip` (default), `include`, or `only`
- `--log-level <level>` — `debug`, `info`, `warn`, or `error` (default)

In `--remote` mode, `--current-branch` resolves against the upstream branch of your local current branch (for example `origin/main`). If no upstream exists, `only` runs zero branches.

If `--remote` is omitted, AYNIG checks the latest commit trailer `dwp-source: git:<name>` and uses that remote when present.

Log level precedence: `--log-level` > `dwp-log-level` trailer > `LOG_LEVEL` env.

## Environment variables

Commands receive metadata via env vars such as:

- `BODY`
- `COMMIT_HASH`
- `STDOUT_LOG_PATH`
- `STDERR_LOG_PATH`
- `FOO`
- `BAZ`

(See also: Commands → Environment Variables.)

Branch logs use the resolved log level after trailers are parsed. Early branch logs are buffered and flushed once the level is known.

Command stdout/stderr is written to `.aynig/logs/<commit-hash>.stdout.log` and
`.aynig/logs/<commit-hash>.stderr.log`, where `<commit-hash>` is the commit
that triggered the command.

AYNIG watches stdout for lines that begin with `SET_STATE ` and applies the
last valid one after the command exits successfully. stderr is not parsed for state
transitions.

If the command exits non-zero, AYNIG marks the branch as `stalled` and records
the exit code plus recent stdout/stderr lines in the commit body.

If the command exits zero without a valid `SET_STATE`, AYNIG refreshes the
`working` commit and keeps waiting in case the command spawned a follow-up
process.
