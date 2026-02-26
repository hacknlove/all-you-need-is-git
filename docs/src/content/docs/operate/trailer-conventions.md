---
title: Trailer Conventions
description: Optional metadata conventions used by workflows.
---

These trailers are optional conventions. The kernel does not interpret them, but your commands can.

## Attempt and retry

```
aynig-attempt: 1
aynig-max-attempts: 5
aynig-next-on-failure: ops-triage
```

## Correlation and grouping

```
aynig-parent-run: run-123
aynig-correlation-id: 2026-02-16-build-42
```

## Checkpoints and notes

```
aynig-checkpoint: package-published
aynig-note: manual rollback completed
```

Use consistent, lowercase keys. Keep values short and script-friendly.

## Working lease context

```
aynig-origin-state: build
aynig-remote: origin
```

When a runner writes `aynig-state: working`, it must also include the original state so retries and tools can recover intent.
`aynig-remote` pins the remote name for follow-up state commits when remote mode is used.
