---
title: Install
description: Set up AYNIG in a repository.
---

## Initialize a repo

```bash
aynig init
```

This creates:

- `.aynig/` with command metadata
- `.worktrees/` for ephemeral worktrees
- `.worktrees/` entry in `.gitignore`

## Install workflows from another repo

```bash
aynig install <repo> [ref] [subfolder]
```

Examples:

```bash
aynig install https://github.com/org/workflows
aynig install https://github.com/org/workflows main
aynig install https://github.com/org/workflows main path/to/.aynig
```
