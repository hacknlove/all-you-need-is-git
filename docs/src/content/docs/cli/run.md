---
title: CLI — run
description: Run AYNIG for the current repository.
---

```bash
aynig run [options]
```

Options:

- `-w, --worktree <path>` — worktree directory (default: `.worktrees`)
- `--dwp-remote <name>` — use remote branches instead of local (e.g. `origin`)
- `--role <name>` — prefer role-specific commands from `.dwp/roles/<name>/command`
- `--current-branch <mode>` — `skip` (default), `include`, or `only`
- `--log-level <level>` — `debug`, `info`, `warn`, or `error` (default)

In `--dwp-remote` mode, `--current-branch` resolves against the upstream branch of your local current branch (for example `origin/main`). If no upstream exists, `only` runs zero branches.

If `--dwp-remote` is omitted, AYNIG checks the latest commit trailer `dwp-source: git:<name>` and uses that remote when present.

When `--role` is provided (or `AYNIG_ROLE` is set), AYNIG looks for `.dwp/roles/<role>/command/<state>` first and falls back to `.dwp/command/<state>`.

Example:

```bash
AYNIG_ROLE=some-role aynig run
```

```bash
aynig run --role some-role
```

Log level precedence: `--log-level` > `dwp-log-level` trailer > `AYNIG_LOG_LEVEL` env.

The runner reads `HEAD`, runs `.dwp/command/<state>`, and checks the new `HEAD`.

Command stdout/stderr is written to `.dwp/logs/<commit-hash>.log`, where `<commit-hash>` is the commit that triggered the command.
