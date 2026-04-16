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

`dwp-state` must appear in the trailer block. If multiple are present, last wins.

The `<state>` value is the dispatch key of the command to execute.

AYNIG:

1. reads `HEAD`
2. extracts trailers
3. resolves the command
4. executes
5. watches command stdout for `SET_STATE {...}` lines
6. materializes the result by writing a new `HEAD`

AYNIG never interprets business semantics.

## 2. Command selection

`dwp-state: <state>` → executable command.

If a role is specified (`--role <name>` or `ROLE`), AYNIG first looks for
`.aynig/roles/<role>/command/<state>` and falls back to `.aynig/command/<state>`.

AYNIG does not define what a state means; it only uses it as a selector.
Semantics belong to upper layers (frameworks, policies, profiles).

## 3. Execution

The command receives:

- body (prompt)
- trailers
- commit hash
- runner configuration

Metadata is delivered as environment variables. Common variables are `BODY`,
`COMMIT_HASH`, `WORKTREE_PATH`, `STDOUT_LOG_PATH`, `STDERR_LOG_PATH`,
`LOG_LEVEL`, and `ROLE`. Commit trailers are also exposed as uppercase
environment variables with dashes converted to underscores, unless that would
overwrite an existing or reserved variable.

AYNIG:

- does not modify the repository during execution
- does not infer the next state
- does not interpret business semantics

Command stdout and stderr are logged separately under `.aynig/logs/` as
`<commit-hash>.stdout.log` and `<commit-hash>.stderr.log`.

The command declares the next state by emitting a line on stdout:

```text
SET_STATE {"state":"review","subject":"review: ready","body":"..."}
```

AYNIG watches stdout, keeps the last valid `SET_STATE` line it sees, and
creates the final commit after the command exits successfully.

The payload may include `"keep_trailers": true` to preserve existing
non-reserved `dwp-*` trailers from the current `working` commit.

If the command exits non-zero, AYNIG ignores any observed `SET_STATE` line and
marks the branch as `stalled` with diagnostic context from stdout/stderr.

If the command exits zero without emitting a valid `SET_STATE`, AYNIG writes a
fresh `working` commit with the same trailers to keep the lease alive while
waiting for any spawned follow-up process.

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
now > committer_date + lease-seconds
```

Reason:

- prevent permanent blocking
- avoid dependence on local processes
- tolerate machine crashes

History is never scanned.

## 6. Valid completion

A tick is valid when, after execution:

- the command emitted a valid `SET_STATE {...}` line on stdout
- the command exited successfully
- `HEAD` contains `dwp-state: <state>`
- `state != working`

That commit is the **tick output**.

AYNIG does not search previous commits nor attempt to reconstruct history. It
only applies the last valid `SET_STATE` observed in the current run and then
observes the latest state.

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
dwp-origin-state: <state>
```

and continue evaluation.

Reason: self-healing system without external coordination or mandatory human intervention.

## 8. What AYNIG does not do

AYNIG does not:

- retry commands
- define workflows
- interpret states
- scan history
- parse stderr for state transitions
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
