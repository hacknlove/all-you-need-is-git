---
title: Kernel Contract
description: The minimum guarantees of the AYNIG runner.
---

# AYNIG Kernel Contract (v0)

This document defines the minimum responsibilities of the AYNIG runner.
AYNIG does not implement workflows or policies; it provides a **deterministic, Git-based distributed execution mechanism** on top of which other layers can build solutions.

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
aynig-state: <state>
```

The `<state>` value is the dispatch key of the command to execute.

AYNIG:

1. reads `HEAD`
2. extracts trailers
3. resolves the command
4. executes
5. validates the result by looking only at the new `HEAD`

AYNIG never interprets business semantics.

## 2. Dispatch

`aynig-state: <state>` → executable command.

AYNIG does not define what a state means; it only uses it as a selector.
Semantics belong to upper layers (frameworks, policies, profiles).

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

## 4. Distributed mutual exclusion (working lease)

AYNIG implements a distributed lock using Git.

Before executing, the runner creates a commit:

```text
aynig-state: working
```

and pushes it to the branch.

If the push fails (the branch advanced), another runner won the execution → abort.

This behaves as a **remote compare-and-swap** without external coordination.

### Reserved `working` trailers

```text
aynig-state: working
aynig-origin-state: <state>
aynig-run-id: <uuid>
aynig-runner-id: <host-id>
aynig-lease-seconds: <ttl>
```

Reason: enable distributed runners without local locks.

## 5. Lease and liveness

While executing, the command must renew the lease:

- all intermediate commits → `aynig-state: working`
- same `aynig-run-id`
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

A tick is valid when, after execution:

- `HEAD` contains `aynig-state: <state>`
- `state != working`

That commit is the **tick output**.

AYNIG does not search previous commits nor attempt to reconstruct history.
It only observes the latest state.

Reason: avoid duplication, loops, and temporal ambiguity.

## 7. Takeover

If a runner finds:

```text
aynig-state: working
lease expired
```

it may recover the branch by creating:

```text
aynig-state: stalled
aynig-stalled-run: <run-id>
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

Those belong to frameworks built on top.

Reason: keep the kernel small, deterministic, and universal.

## 9. System guarantees

AYNIG guarantees:

1. A single active executor per branch
2. Auditable execution (everything is a commit)
3. Recovery after crashes
4. HEAD-based determinism
5. Distributed compatibility without external services

## Summary

AYNIG turns Git into:

- an event bus
- a distributed lock
- a state machine

The runner acts as an execution kernel: execute → validate → observe.

All intelligence lives above it.
