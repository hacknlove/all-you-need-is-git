---
title: CLI — run
description: Run AYNIG for the current repository.
---

```bash
aynig run [options]
```

`aynig run` inspects the current repository, resolves the command for the
current DWP state, and executes the workflow for the branches it selects.

It is the main entry point for processing AYNIG state transitions.

## Behavior

- Uses `.worktrees` by default for worktree management.
- Resolves the remote from `--dwp-remote` when provided.
- If `--dwp-remote` is omitted, AYNIG checks the latest commit trailer
  `dwp-source: git:<name>` and uses that remote when present.
- `--current-branch` controls whether the current branch is skipped,
  included, or used as the only branch.
- In remote mode, current-branch resolution is based on the upstream branch of
  the current local branch, such as `origin/main`.
- `--role` and `AYNIG_ROLE` make AYNIG check
  `.dwp/roles/<role>/command/<state>` before `.dwp/command/<state>`.
- Log level precedence is `--log-level` > `dwp-log-level` trailer >
  `AYNIG_LOG_LEVEL` > default `error`.
- Command stdout and stderr are written to `.dwp/logs/<commit-hash>.log`.

## Options

### `-w, --worktree <path>`

Sets the directory where AYNIG creates and manages worktrees. The default is
`.worktrees`.

Example:

```bash
aynig run --worktree .aynig-worktrees
```

### `--dwp-remote <name>`

Runs against branches discovered from the named remote instead of using only
local branches.

Example:

```bash
aynig run --dwp-remote origin
```

### `--role <name>`

Prefers role-specific commands from `.dwp/roles/<name>/command` when present.

Example:

```bash
aynig run --role reviewer
```

### `--current-branch <mode>`

Controls how the current branch is handled. Supported values are `skip`
(default), `include`, and `only`.

Example:

```bash
aynig run --current-branch include
```

### `--log-level <level>`

Sets the runner log level. Supported values are `debug`, `info`, `warn`, and
`error`.

Example:

```bash
aynig run --log-level debug
```

## Examples

Run with the default local-branch behavior:

```bash
aynig run
```

Run with a specific role:

```bash
aynig run --role reviewer
```

Run against remote branches:

```bash
aynig run --dwp-remote origin
```

Run only the current branch:

```bash
aynig run --current-branch only
```

Run with a custom worktree directory and verbose logging:

```bash
aynig run --worktree .aynig-worktrees --log-level debug
```

Use the environment variable form for role selection:

```bash
AYNIG_ROLE=reviewer aynig run
```
