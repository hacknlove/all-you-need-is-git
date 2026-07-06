---
title: CLI — status
description: Show the current AYNIG state for the repo.
---

```bash
aynig status [branch|pattern] [options]
```

`aynig status` prints the current branch state as derived from `HEAD`, or
inspects specific local branches without checking them out.

It is a read-only command intended for quick inspection and debugging.

## Behavior

- Prints the inspected branch name and `HEAD` commit hash.
- Reads the latest commit trailers and reports `dwp-state` and `dwp-run-id`.
- When the current state is `working`, also reports `dwp-origin-state` when present.
- Computes lease status as `active`, `expired`, `unknown`, or `n/a`.
- Resolves the command path for the current state, or for `dwp-origin-state`
  when the current state is `working`.
- Uses `command: lease` when the current commit is `working` but there is no
  origin state to resolve.
- Resolves commands from the same trusted commands ref as `aynig run`: the
  default branch unless `--commands-ref` selects another ref, or
  `--commands-ref same` selects each inspected branch. Command lookups read
  the committed tree of that ref, so the report matches what `run` would
  execute, not what happens to be in the current checkout.
- Prefers `.aynig/roles/<role>/command/<state>` over `.aynig/command/<state>` when
  a role is provided and the role command is executable.
- With `--branch` or a positional branch name, inspects that branch directly.
- With `--branch-pattern` or a positional glob such as `"1-*"`, prints one
  status block per matching branch.
- When no branch selector is provided, preserves the current `HEAD`-based
  behavior.

## Options

### `--role <name>`

Checks the role-specific command directory first when resolving the command for
the current state. If omitted, AYNIG falls back to `ROLE` when it is set.

Example:

```bash
aynig status --role reviewer
```

### `--branch <name>`

Inspects a specific local branch without checking it out.

Example:

```bash
aynig status --branch 1-bootstrap
```

### `--branch-pattern <pattern>`

Inspects all local branches that match the given glob pattern.

Example:

```bash
aynig status --branch-pattern "1-*"
```

### `--commands-ref <ref>`

Branch, tag, or commit whose `.aynig` commands are inspected. Uses the same
default as `aynig run`: the repository's default branch (`master` or `main`).
Pass `same` to resolve commands from each inspected branch instead.

Examples:

```bash
aynig status --branch 1-bootstrap --commands-ref release
aynig status --branch 1-bootstrap --commands-ref same
```

## Output

The command prints lines in this shape:

```text
branch: <branch>
head: <commit>
dwp-state: <state>
dwp-origin-state: <state>
dwp-run-id: <run-id>
lease: <status>
command: <exists|missing|lease>
command-path: <ref>:<path>
```

Some lines appear only when the related data exists.

`command-path` uses Git revision notation (for example
`master:.aynig/command/build`), pointing at the commands ref the command is
resolved from. You can inspect it with `git show <command-path>`.

## Examples

Inspect the current branch:

```bash
aynig status
```

Inspect one branch directly:

```bash
aynig status 1-bootstrap
```

Inspect multiple branches with a pattern:

```bash
aynig status "1-*"
```

Resolve commands using a role-specific command directory:

```bash
aynig status --role reviewer
```

Use the environment variable form for role selection:

```bash
ROLE=reviewer aynig status
```
