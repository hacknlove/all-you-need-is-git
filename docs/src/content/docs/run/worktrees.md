---
title: Worktrees and Isolation
description: How AYNIG isolates executions.
---

AYNIG runs each command inside a dedicated Git worktree. This is a core design choice:

- Commands always execute on a clean snapshot
- The main working directory is never modified
- Concurrent executions do not interfere

Commands run with the working directory set to the worktree. Worktrees are created under `.worktrees/` and reused across runs when possible.
