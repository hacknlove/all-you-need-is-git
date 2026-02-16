---
title: Concurrency Patterns
description: Fan-out, fan-in, and multi-repo coordination without a control plane.
---

These patterns are conventions implemented in your workflow commands, not kernel features.

## Fan-out and fan-in

Fan-out runs parallel work across branches or repos. Fan-in aggregates results.

Suggested trailers:

```
aynig-correlation-id: 2026-02-16-build-42
aynig-parent-run: run-123
```

Fan-out steps:

1. Create N branches or repo runs with the same correlation id.
2. Each branch progresses independently.

Fan-in steps:

1. A coordinator command reads refs and commit trailers to detect completion.
2. It commits a summary state when all children are complete.

## Multi-branch coordination

Use a stable correlation id and a shared convention for completion:

```
aynig-state: complete
aynig-correlation-id: 2026-02-16-build-42
```

The coordinator can poll for branches reaching `complete` and proceed.

## Multi-repo coordination

Without a control plane, rely on Git-native signals:

- A coordinating repo records the list of participant runs and correlation id.
- Each repo writes a final commit with the same correlation id.
- A coordinator command reads those repos and commits the aggregated result.

This keeps orchestration in your scripts and preserves the kernel boundary.
