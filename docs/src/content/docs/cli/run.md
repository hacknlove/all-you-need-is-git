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
- Resolves the remote from `--remote` when provided.
- If `--remote` is omitted, AYNIG checks the latest commit trailer
  `dwp-source: git:<name>` and uses that remote when present.
- `--current-branch` controls whether the current branch is skipped,
  included, or used as the only branch.
- In remote mode, current-branch resolution is based on the upstream branch of
  the current local branch, such as `origin/main`.
- `--role` and `ROLE` make AYNIG check
  `.aynig/roles/<role>/command/<state>` before `.aynig/command/<state>`.
- Commands are resolved from a trusted commands ref (the default branch unless
  `--commands-ref` overrides it), never from the event branch, so a branch
  cannot redefine the commands a runner executes.
- Log level precedence is `--log-level` > `dwp-log-level` trailer >
  `LOG_LEVEL` > default `error`.
- Command stdout and stderr are written to separate files:
  `.aynig/logs/<commit-hash>.stdout.log` and `.aynig/logs/<commit-hash>.stderr.log`.
- The command receives those paths in `STDOUT_LOG_PATH` and `STDERR_LOG_PATH`.
- AYNIG watches stdout for lines that begin with `SET_STATE ` and applies the
  last valid one after the command exits successfully.
- If the command exits non-zero, AYNIG marks the branch as `stalled` with a
  diagnostic body that includes the exit code and recent log lines.
- If the command exits zero without a valid `SET_STATE`, AYNIG writes a fresh
  `working` commit and keeps the lease alive.

## Options

### `-w, --worktree <path>`

Sets the directory where AYNIG creates and manages worktrees. The default is
`.worktrees`.

Example:

```bash
aynig run --worktree .aynig-worktrees
```

### `--remote <name>`

Runs against branches discovered from the named remote instead of using only
local branches.

Example:

```bash
aynig run --remote origin
```

### `--role <name>`

Prefers role-specific commands from `.aynig/roles/<name>/command` when present.

Example:

```bash
aynig run --role reviewer
```

### `--commands-ref <ref>`

Selects the branch, tag, or commit whose `.aynig` commands are executed. The
default is the repository's default branch (`<remote>/HEAD` in remote mode,
`master` or `main` locally). Use `same` to resolve commands from each event
branch, which restores the pre-trusted-ref behavior and lets any branch
redefine its own commands.

Examples:

```bash
# pin commands to an immutable commit (strongest guarantee)
aynig run --commands-ref 3f2a9c1

# develop a workflow on a feature branch
aynig run --commands-ref my-feature --current-branch only

# resolve commands from each event branch (trusted environments only)
aynig run --commands-ref same
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
aynig run --remote origin
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
ROLE=reviewer aynig run
```
