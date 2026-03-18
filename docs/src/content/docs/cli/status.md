---
title: CLI — status
description: Show the current AYNIG state for the repo.
---

```bash
aynig status [options]
```

`aynig status` prints the current branch state as derived from `HEAD` and the
available command files in the repository.

It is a read-only command intended for quick inspection and debugging.

## Behavior

- Prints the current branch name and `HEAD` commit hash.
- Reads the latest commit trailers and reports `dwp-state` and `dwp-run-id`.
- When the current state is `working`, also reports `dwp-origin-state` when present.
- Computes lease status as `active`, `expired`, `unknown`, or `n/a`.
- Resolves the command path for the current state, or for `dwp-origin-state`
  when the current state is `working`.
- Uses `command: lease` when the current commit is `working` but there is no
  origin state to resolve.
- Prefers `.dwp/roles/<role>/command/<state>` over `.dwp/command/<state>` when
  a role is provided.

## Options

### `--role <name>`

Checks the role-specific command directory first when resolving the command for
the current state. If omitted, AYNIG falls back to `AYNIG_ROLE` when it is set.

Example:

```bash
aynig status --role reviewer
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
command-path: <path>
```

Some lines appear only when the related data exists.

## Examples

Inspect the current branch:

```bash
aynig status
```

Resolve commands using a role-specific command directory:

```bash
aynig status --role reviewer
```

Use the environment variable form for role selection:

```bash
AYNIG_ROLE=reviewer aynig status
```
