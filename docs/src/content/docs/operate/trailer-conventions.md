---
title: Trailer Conventions
description: Recommended trailer keys your workflows may use.
---

These are **conventions only**. AYNIG does not interpret them.

Suggested keys:

```
dwp-attempt: 1
dwp-max-attempts: 5
dwp-next-on-failure: ops-triage
```

## Correlation and grouping

```
dwp-parent-run: run-123
dwp-correlation-id: 2026-02-16-build-42
```

## Checkpoints and notes

```
dwp-checkpoint: package-published
dwp-note: manual rollback completed
```

Use consistent, lowercase keys. Keep values short and script-friendly.

## Working lease context

```
dwp-origin-state: build
dwp-source: git:origin
```

When a runner writes `dwp-state: working`, it must also include the original state so retries and tools can recover intent.
`dwp-source` pins the remote name for follow-up state commits when remote mode is used.
