---
title: Runner Contract
description: The minimum guarantees of the AYNIG runner.
---

# AYNIG Runner Contract (v0)

This page mirrors the canonical contract at the repository root in `CONTRACT.md`.

AYNIG is a small Git-based execution kernel: workflows define meaning; the
runner provides dispatch, locking, logging, and state materialization.

The source of truth is always the latest commit (`HEAD`) on a branch.

## 1. State Events

AYNIG treats each commit as a workflow event:

- title: human-readable only
- body: prompt delivered to the command
- trailers: structured metadata

The required trailer is:

```text
dwp-state: <state>
```

`dwp-state` must be in the trailer block. If it appears more than once, the
last value wins.

For each selected branch, AYNIG:

1. reads `HEAD`
2. extracts trailers
3. resolves the command for `dwp-state`
4. creates a `working` lease commit
5. runs the command under supervision
6. materializes the command result as a new `HEAD`

AYNIG never interprets workflow-specific state semantics.

## 2. Command Dispatch

`dwp-state: <state>` maps to an executable command:

```text
.aynig/command/<state>
```

When `--role <name>` or `ROLE` is set, AYNIG tries this path first:

```text
.aynig/roles/<role>/command/<state>
```

If no role-specific command exists, AYNIG falls back to `.aynig/command/<state>`.

## 3. Command Inputs And Logs

Commands run with the working directory set to the selected worktree.

Common environment variables:

- `BODY`
- `COMMIT_HASH`
- `WORKTREE_PATH`
- `STDOUT_LOG_PATH`
- `STDERR_LOG_PATH`
- `LOG_LEVEL`
- `ROLE`

Commit trailers are also exposed as uppercase environment variables with dashes
converted to underscores, unless that would overwrite an existing or reserved
variable.

Command stdout and stderr are logged separately:

```text
.aynig/logs/<commit-hash>.stdout.log
.aynig/logs/<commit-hash>.stderr.log
```

AYNIG watches stdout for state protocol lines. It does not parse stderr for
state transitions.

## 4. Command Output Protocol

A command declares the next state by emitting a single-line JSON payload on
stdout:

```text
SET_STATE {"state":"review","subject":"review: ready","body":"..."}
```

Rules:

- the line must start with `SET_STATE `
- the payload must be valid JSON on one line
- `state` is required and cannot be `working`
- the last valid `SET_STATE` line wins
- the runner applies it only if the command exits with code `0`

Optional payload fields:

- `subject`: final commit title; defaults to `chore: set <state>`
- `body`: final commit body
- `keep_trailers: true`: copies non-reserved `dwp-*` trailers from `working`
- `trailers`: extra final commit trailers

`trailers` must not include runner-managed keys:

```text
dwp-state
dwp-source
dwp-origin-state
dwp-run-id
dwp-runner-id
dwp-lease-seconds
dwp-stalled-run
```

## 5. Working Lease

Before running a command, AYNIG claims the branch by creating:

```text
dwp-state: working
dwp-origin-state: <state>
dwp-run-id: <uuid>
dwp-runner-id: <host-id>
dwp-lease-seconds: <ttl>
```

When remote mode is active, it also records:

```text
dwp-source: git:<remote>
```

The push of this `working` commit is the compare-and-swap: if the push fails,
another runner advanced the branch first.

Liveness is based on the committer timestamp of `HEAD`:

```text
HEAD is working
and now > committer_date + dwp-lease-seconds
```

When that is true, another runner may recover the branch by creating:

```text
dwp-state: stalled
dwp-stalled-run: <run-id>
dwp-origin-state: <state>
```

## 6. Completion Outcomes

Successful completion:

- command exits `0`
- command emitted a valid `SET_STATE`
- runner writes a non-`working` final state commit

Command failure:

- command exits non-zero
- runner ignores any observed `SET_STATE`
- runner writes `dwp-state: stalled` with diagnostic stdout/stderr context

No declared state:

- command exits `0`
- no valid `SET_STATE` was emitted
- runner writes a fresh `working` commit with the same trailers

That final `working` commit keeps the lease alive in case the command spawned
another process that will eventually advance the branch.

## 7. What AYNIG Does Not Do

AYNIG does not:

- define workflows
- interpret non-reserved states
- retry commands
- scan history to infer state
- parse stderr for transitions
- decide merges
- resolve semantic conflicts
- guarantee task success

Those are workflow or policy concerns.

## 8. Guarantees

AYNIG guarantees:

1. one active executor per branch
2. auditable state transitions through commits
3. recovery from abandoned `working` leases
4. HEAD-based determinism
5. distributed execution without external coordination services
