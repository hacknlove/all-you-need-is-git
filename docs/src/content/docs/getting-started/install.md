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

## Install the CLI

```bash
curl -fsSL https://aynig.org/install.sh | bash
```

To install a specific version:

```bash
VERSION=v0.1.0 curl -fsSL https://aynig.org/install.sh | bash
```

If your PATH does not include `~/.local/bin`, add it to your shell profile.

## Packages

We also publish `.deb` and `.rpm` packages under `https://aynig.org/releases/<version>/`.
