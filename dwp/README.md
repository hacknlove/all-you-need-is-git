# DWP — Distributed Workflow Protocol (draft)

Location: this protocol lives under `dwp/` in the monorepo.

This document captures the **protocol layer** extracted from the current AYNIG implementation.

Goal: separate **protocol (DWP)** from **implementation (AYNIG runner/CLI)** so other products can adopt the same coordination mechanism.

Status: **draft** (decisions + open questions).

---

## 1) Naming

- Protocol name: **DWP** (Distributed Workflow Protocol)
- Transport/binding (for the current implementation): **DWP/GC** (DWP over Git commits)

Rationale: like HTTP over TCP/IP, DWP is transport-agnostic; Git commits are one binding.

---

## 2) Core concept (what DWP is for)

DWP is a protocol for coordinating distributed workers (humans and/or agents) as a workflow.

A workflow is advanced via **state events**. Only the latest event is required for dispatch (HEAD-based / latest-event based).

---

## 3) Message / metadata key

- Canonical state trailer key: `dwp-state: <state>`
- The state key MUST appear exactly once in the trailing metadata block.

Back-compat:
- **None planned.** AYNIG will migrate to `dwp-state:` as the canonical key.

---

## 4) Reserved states (standard)

DWP reserves a small set of state names.

- Some are **core** and have protocol-level semantics.
- Others are **reserved for future standard semantics** (like reserved words in a programming language).

Workflows **MUST NOT** redefine reserved states with conflicting meaning.

### 4.1 Core reserved states (normative semantics)

- `working`
  - indicates a runner has claimed an execution lease for a branch/job
- `stalled`
  - indicates recovery is required (previous lease expired or runner crashed)

### 4.2 Reserved states (no normative semantics yet)

These are reserved to enable future interoperability (dashboards, tooling) without breaking workflows.

- `done` — terminal success candidate
- `failed` — terminal failure candidate
- `canceled` — terminal cancellation candidate
- `blocked` — waiting on external input/decision (human/system)
- `waiting` — temporary wait (rate limit, window, dependency)
- `noop` — no work / no-op

---

## 5) Reserved trailers (standard)

When `dwp-state: working` is used, the following trailers are reserved:

- `dwp-origin-state: <state>`
- `dwp-run-id: <uuid>`
- `dwp-runner-id: <host-id>`
- `dwp-lease-seconds: <ttl>`

Optional trailers (standardized):

- `dwp-source: <locator>`
- `dwp-log-level: <debug|info|warn|error>`

`dwp-log-level` is a **hint** for runner verbosity. It MUST NOT change workflow semantics.
Runners MAY ignore it.

---

## 6) Execution / dispatch model

- Dispatch is based on the **latest event** (e.g., `HEAD` for DWP/GC).
- `dwp-state: <state>` selects the command/handler for that state.
- The handler is responsible for:
  - doing work
  - emitting a new event with a new `dwp-state`

The protocol does NOT define workflow semantics beyond the reserved states.

---

## 7) Lease / mutual exclusion

DWP standardizes a distributed mutual exclusion mechanism via `working` state:

- Runner claims by emitting `dwp-state: working` with reserved trailers.
- Lease liveness is determined by the timestamp of the latest event.
- Takeover is allowed when lease is expired (`now > last_event_time + dwp-lease-seconds`).

---

## 8) Actions pending

### Spec structure

- Split spec into:
  - `DWP.md` (core)
  - `DWP-GC.md` (Git commit binding)

### Rename and docs alignment

- Decide canonical naming for all reserved keys (`dwp-*`) and document exact rules:
  - metadata block detection ("trailing trailer block")
  - whitespace/case normalization
  - multiple occurrences handling

### Implementation plan (AYNIG)

- Update AYNIG implementation to use DWP keys:
  - parse `dwp-state` (canonical)
  - emit `dwp-*` trailers for the working lease
- Update docs pages that currently mention `aynig-state` to mention DWP + DWP/GC.

### Compatibility

- **No backwards compatibility planned** (pre-1.0; no external users yet).

---

## 9) Open questions

- Should DWP standardize terminal semantics for `done`/`failed`/`canceled`, or keep them reserved-only?
- Should DWP reserve/standardize additional lease-related trailers beyond the current set?
- Are we happy with strict case-sensitivity for trailer keys long-term?

Resolved (for DWP/GC):
- Trailer parsing: **Git trailers as the standard**.
- Duplicate keys: **last wins**.
- Case handling: **case-sensitive**.

