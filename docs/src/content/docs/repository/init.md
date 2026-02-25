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
- `.aynig/CONTRACT.md`
- `.worktrees/` for ephemeral worktrees
- `.worktrees/` and `.aynig/logs/` entries in `.gitignore`
