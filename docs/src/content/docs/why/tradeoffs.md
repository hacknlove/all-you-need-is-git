---
title: Tradeoffs
description: What you give up to keep the kernel small and Git-native.
---

- The kernel does not retry or interpret semantics
- Long-running workflows can create many commits (manage history)
- Concurrency is per-branch (patterns are higher-level)
- You must design workflow policies (approvals, gates, notifications)
