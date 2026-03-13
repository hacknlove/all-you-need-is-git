---
title: Runner Contract
description: The minimum guarantees of the AYNIG runner.
---

# AYNIG Runner Contract (v0)

This document defines the minimum responsibilities of the AYNIG runner.
AYNIG does not implement workflows or policies; it provides a **deterministic, Git-based execution mechanism** that other tools can build on.

The canonical contract lives at the repository root in `CONTRACT.md`.

The source of truth is always the Git branch (local or remote in `remote` mode).

## 1. Model

AYNIG interprets each commit as a **state event**.

A commit contains:

- **title** → human-only (ignored by the system)
- **body** → prompt delivered to the command
- **trailers** → structured metadata

The mandatory trailer is:

```text
dwp-state: <state>
```

The `<state>` value is the key used to select the command to execute.

AYNIG:

1. reads `HEAD`
2. extracts trailers
3. selects the command
4. executes it
5. checks the result by looking only at the new `HEAD`

AYNIG never interprets business semantics.

## 2. Command selection

`dwp-state: <state>` → executable command.

AYNIG does not define what a state means; it only uses it to select a command.
Meaning belongs to your workflow.

## 3. Execution

The command receives:

- body (prompt)
- trailers
- commit hash
- runner configuration

Metadata is delivered as environment variables.

AYNIG:

- does not modify the repository during execution
- does not create final commits
- does not decide the next state

**Only the command advances the state machine.**

## 4. Working lease (one runner at a time)

AYNIG prevents two runners from working on the same branch at the same time by using Git commits.

Before executing, the runner creates a commit:

```text
dwp-state: working
```

and pushes it to the branch.

If the push fails (the branch advanced), another runner won the execution → abort.

This behaves like a **remote compare-and-swap** without external coordination.

### Reserved `working` trailers

```text
dwp-state: working
dwp-origin-state: <state>
dwp-run-id: <uuid>
dwp-runner-id: <host-id>
dwp-lease-seconds: <ttl>
```

Reason: enable distributed runners without local locks.

## 5. Lease and liveness

While executing, the command must renew the lease:

- all intermediate commits → `dwp-state: working`
- same `dwp-run-id`
- implicit heartbeat update (committer date)

AYNIG uses the **committer timestamp of HEAD** as the liveness signal.

Takeover is allowed when:

```text
HEAD == working
and
now > committer_date + lease-seconds (+grace)
```

Reason:

- prevent permanent blocking
- avoid dependence on local processes
- tolerate machine crashes

History is never scanned.

## 6. Valid completion

A step is valid when, after execution:

- `HEAD` contains `dwp-state: <state>`
- `state != working`

That commit is the **step output**.

AYNIG does not search previous commits nor attempt to reconstruct history.
It only observes the latest state.

Reason: avoid duplication, loops, and temporal ambiguity.

## 7. Takeover

If a runner finds:

```text
dwp-state: working
lease expired
```

it may recover the branch by creating:

```text
dwp-state: stalled
dwp-stalled-run: <run-id>
```

and continue evaluation.

Reason: self-healing system without external coordination or mandatory human intervention.

## 8. What AYNIG does not do

AYNIG does not:

- retry commands
- define workflows
- interpret states
- scan history
- decide merges
- resolve semantic conflicts
- guarantee task success

Those belong to tools and workflows built on top.

Reason: keep AYNIG small, deterministic, and universal.

## 9. System guarantees

AYNIG guarantees:

1. A single active executor per branch
2. Auditable execution (everything is a commit)
3. Recovery after crashes
4. HEAD-based determinism
5. Distributed compatibility without external services

## Summary

AYNIG turns Git into:

- a place to store workflow events
- a distributed lock
- a state machine

The runner acts as a simple loop: execute → check → observe.

All intelligence lives above it.
