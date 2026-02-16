---
title: Retries
description: Retry patterns built from trailers and state transitions.
---

Retries are workflow policy, not kernel behavior. Use trailers to track attempts and make decisions in your command scripts.

## Recommended pattern

1. Read attempt counters from trailers.
2. Decide whether to retry or fail.
3. Commit a new state with updated trailers.

Example trailers:

```
aynig-attempt: 2
aynig-max-attempts: 5
```

## State-driven retries

Use explicit states to control retry logic.

```
aynig-state: retry
```

Command logic can map `retry` back to the same command with updated metadata.

## Exponential backoff via command scripts

The kernel does not wait or sleep, but your command can.

Pseudo-logic (helper functions are workflow-specific):

```
attempt=$(read_trailer aynig-attempt)
sleep_seconds=$((2 ** attempt))
sleep "$sleep_seconds"
```

When the command finishes, commit with:

```
aynig-state: working
aynig-attempt: 3
```

## Escalation path

If you hit the maximum:

```
aynig-state: failed
aynig-next-on-failure: ops-triage
```

This keeps policy in your workflow definition and avoids adding semantics to the kernel.
