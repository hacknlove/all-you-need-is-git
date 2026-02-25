---
title: Worktrees and Isolation
description: Why AYNIG runs commands inside ephemeral worktrees.
---

AYNIG runs each execution inside a dedicated **Git worktree**.

For every actionable commit, AYNIG:

- creates a worktree checked out on the target branch,
- runs the command inside that worktree,
- removes it once execution finishes.

Commands run with the working directory set to the worktree. Worktrees are created under `.worktrees/` and reused across runs when possible.

This guarantees:

- clean, reproducible snapshots
- the main working directory is not modified
- concurrent executions do not interfere

The default worktree location is:

```text
.worktrees/
```

It is added to `.gitignore` by `aynig init`.
