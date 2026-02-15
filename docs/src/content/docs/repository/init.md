---
title: Initialize a Repository
description: Set up AYNIG inside a Git repository.
---

Run this inside the repository where you want to use AYNIG:

```bash
aynig init
```

This creates:

- `.aynig/` with starter files
- `.worktrees/` for ephemeral worktrees
- a `.worktrees/` entry in `.gitignore`
