---
title: Runbooks
description: Operational playbooks for failures, stalls, and takeovers.
---

This section provides practical, non-kernel runbooks. These steps use the existing contract and keep all policy in your workflows.

## Command failure handling

Goal: make failures explicit, reproducible, and recoverable.

1. Ensure your command writes a clear failure message to stdout or logs.
2. Transition to a failure state using the dispatch trailer:

```
aynig-state: failed
```

3. Attach context via trailers (optional):

```
aynig-failure-reason: tests failed
aynig-failure-step: integration
```

4. Create a follow-up commit that either retries or escalates (see retries).

## Stalled lease handling

A lease is considered stale if `HEAD` has not advanced within your agreed cadence.

1. Confirm the lease is actually stale by checking timestamps on recent commits.
2. If stale, mark the branch as stalled:

```
aynig-state: stalled
```

3. Record why you stalled the run with an optional trailer:

```
aynig-stalled-reason: no heartbeat for 30m
```

4. Notify the owning team or operator.

## Takeover process

Use this when the original runner is gone or wedged.

1. Fetch and inspect recent commits for ownership or context trailers.
2. If no active lease exists, create a takeover commit.
   If you use a takeover marker, treat it as a workflow-specific convention (see trailer conventions).

```
aynig-state: working
aynig-takeover: true
```

3. If the previous command must be re-run, set the target state and retry metadata.
4. Proceed with the command execution under your own runner identity.

## Human intervention flow

1. Halt automation by transitioning to a safe state such as `stalled` or `failed`.
2. Add a human-readable note in trailers or commit message body.
3. Document the manual action taken in a follow-up commit.
4. Resume the workflow by moving to the next state in your command map.

For kernel-level behavior, see `docs/src/content/docs/contract/index.md`.
