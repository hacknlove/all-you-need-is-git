---
title: CLI — init
description: Initialize a repository for AYNIG.
---

```bash
aynig init
```

`aynig init` bootstraps the current git repository with the files and folders
that AYNIG expects.

It is safe to re-run. Existing files are left in place and missing pieces are
created.

## Behavior

- Requires the current directory to be a git repository.
- Creates `.dwp/` if it does not already exist.
- Creates `.dwp/command/` if it does not already exist.
- Creates `.dwp/COMMANDS.md` only when `.dwp/` is newly created.
- Creates `.dwp/CONTRACT.md` when it is missing.
- Creates `.dwp/command/clean` when it is missing.
- Creates `.worktrees/` when it is missing.
- Ensures `.worktrees/` and `.dwp/logs/` are present in `.gitignore`.
- Skips existing files instead of overwriting them.

## Options

This command has no flags.

## Files Created

Depending on what already exists, `aynig init` may create:

- `.dwp/`
- `.dwp/COMMANDS.md`
- `.dwp/CONTRACT.md`
- `.dwp/command/`
- `.dwp/command/clean`
- `.worktrees/`
- `.gitignore` entries for `.worktrees/` and `.dwp/logs/`

## Examples

Initialize the current repository:

```bash
aynig init
```

Initialize after creating a new git repo:

```bash
git init && aynig init
```

Re-run initialization to fill in missing files without overwriting existing ones:

```bash
aynig init
```
