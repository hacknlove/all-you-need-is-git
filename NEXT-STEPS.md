# NEXT-STEPS

This document describes a concrete roadmap to address current AYNIG limitations (as identified in docs/review feedback) **without changing the kernel philosophy**.

The goal is to improve usability, operability, and adoption while keeping the contract intact:

- Kernel stays minimal and deterministic
- No hidden control plane
- Git remains the source of truth
- Workflow semantics stay in upper layers

---

## 1) Current State Snapshot

### Core contract and model
- Canonical kernel contract: `CONTRACT.md`
- Mandatory dispatch trailer: `aynig-state: <state>`
- Command path: `.aynig/command/<state>`
- Lease model:
  - lock state: `aynig-state: working`
  - stale recovery state: `aynig-state: stalled`
- Kernel intentionally does not:
  - define workflows
  - retry commands
  - interpret business semantics

### Current implementation (important paths)
- JS runner:
  - `js/AgentsOrchestrator/Command.js`
  - dispatch path: `.aynig/command`
  - env exports:
    - `AYNIG_BODY`
    - `AYNIG_COMMIT_HASH`
    - `AYNIG_TRAILER_<UPPER_KEY>` (hyphens preserved)
- Go runner:
  - `go/internal/orchestrator/command.go`
  - dispatch path: `.aynig/command`
  - env exports equivalent to JS
- CLI docs:
  - main docs index: `docs/src/content/docs/index.md`
  - tradeoffs page: `docs/src/content/docs/why/tradeoffs.md`

---

## Status

Use this section to track progress across phases. Add more checkboxes for any subtask you expect to pause and resume later.

- [x] Phase A (Docs-first ops layer)
- [ ] Phase B (Observability commands)
  - [ ] B1 `status`
  - [ ] B2 `events --json`
- [ ] Phase C (Env var normalization)
- [ ] Phase D (Ops workflow pack)

---

## 2) Problem Statement

Feedback is positive about kernel elegance, but concerns remain in these areas:

1. Lack of out-of-the-box operational patterns (retry/failure/recovery playbooks)
2. Limited observability (no easy run/event view from CLI)
3. Git history noise in long-running workflows
4. Advanced concurrency patterns under-documented (fan-out/fan-in, multi-branch coordination)
5. Env variable ergonomics: hyphenated trailer env names are awkward in shell

We want to reduce friction **without bloating the kernel**.

---

## 3) Guiding Principles (Non-Negotiable)

1. **Do not change kernel responsibilities** in `CONTRACT.md`
2. Additions must be:
   - optional
   - composable
   - Git-native
3. New features should live in:
   - docs
   - helper CLI commands (read-only/operational)
   - optional workflow packs/profiles
4. Avoid introducing external dependencies as mandatory runtime components

---

## 4) Delivery Plan

## Phase A — Docs-first operational layer (low risk, high value)

### A1. Add "Operating AYNIG" documentation section
Create new docs section and pages:

- `docs/src/content/docs/operate/runbooks.md`
  - command failure handling
  - stalled lease handling
  - takeover process
  - human intervention flow
- `docs/src/content/docs/operate/retries.md`
  - retry patterns via trailers/state transitions
  - max-attempt conventions
  - exponential backoff via command scripts
- `docs/src/content/docs/operate/history-hygiene.md`
  - branch strategy for workflow commits
  - checkpointing and squash recommendations
  - lease heartbeat cadence guidance
- `docs/src/content/docs/operate/concurrency-patterns.md`
  - fan-out/fan-in examples
  - correlation trailers conventions
  - multi-repo coordination pattern (without external control plane)

Update sidebar in:
- `docs/astro.config.mjs`

### A2. Define recommended trailer conventions (non-kernel)
Document conventions such as:
- `aynig-attempt`
- `aynig-max-attempts`
- `aynig-next-on-failure`
- `aynig-parent-run`
- `aynig-correlation-id`

These are **conventions only**; kernel should not interpret them.

---

## Phase B — Observability command surface (read-only CLI improvements)

### B1. Add `aynig status`
Purpose:
- show current branch head state at a glance
- indicate if `working` lease appears expired
- indicate actionable state and command existence

Suggested output:
- branch
- head commit
- `aynig-state`
- run id (if present)
- lease status (`active` / `expired` / `n/a`)
- command path resolution (`exists` / `missing`)

### B2. Add `aynig events` (or `aynig inspect`)
Purpose:
- print a compact event stream from recent commits (N commits or HEAD-only mode)
- support machine output (`--json`) for dashboards/scripts

Important:
- keep default behavior conservative (HEAD-focused) to preserve kernel philosophy
- if history mode exists, make it explicitly optional

### B3. Keep all observability read-only
No mutation behavior in these commands.

Implementation touchpoints:
- JS CLI entry: `js/index.js`
- Go CLI entry: `go/cmd/aynig/main.go`
- New command files in:
  - JS: `js/commands/`
  - Go: `go/internal/commands/`

---

## Phase C — Env var ergonomics

Current issue:
- `AYNIG_TRAILER_AYNIG-STATE` is awkward in shell due to hyphen

### C1. Normalize trailer env names
For each trailer key:
- replace hyphens with underscores

Example:
- before: `AYNIG_TRAILER_AYNIG-STATE`
- after: `AYNIG_TRAILER_AYNIG_STATE`

Apply in both:
- `js/AgentsOrchestrator/Command.js`
- `go/internal/orchestrator/command.go`

### C2. Document the new naming
Update docs to reflect the underscore format and remove references to hyphenated names.

---

## Phase D — Optional ops workflow pack (outside kernel)

Create an optional workflow pack repository or folder template containing:
- retry command patterns
- failure notification command templates
- stalled recovery command template
- sample `COMMANDS.md` state map

Installable via:
- `aynig install <repo> [ref] [subfolder]`

This provides defaults without forcing policy into kernel.

---

## 5) Acceptance Criteria

## Kernel integrity
- `CONTRACT.md` remains valid and unchanged in meaning
- No new kernel-level semantic interpretation of workflow trailers

## Docs completeness
- New operations docs explain:
  - failures
  - retries
  - stalled takeover
  - history hygiene
  - concurrency patterns
- Sidebar includes discoverable navigation

## CLI improvements
- `aynig status` works in JS and Go builds
- `aynig events --json` (or chosen equivalent) works consistently
- outputs are script-friendly and stable

## Env ergonomics
- both raw and normalized trailer env vars are exported
- existing workflows using raw names continue to work

## Testing
- Unit tests added for:
  - status/event parsing logic
  - env alias generation
- Smoke tests updated for command presence and basic output

---

## 6) Suggested Task Breakdown (Issue-Ready)

1. Docs: add operating runbooks and patterns pages
2. Docs: add trailer conventions reference page
3. CLI(JS): implement `status` command
4. CLI(Go): implement `status` command
5. CLI(JS): implement `events --json`
6. CLI(Go): implement `events --json`
7. Runner(JS): add normalized trailer env aliases
8. Runner(Go): add normalized trailer env aliases
9. Tests(JS): status/events/env alias coverage
10. Tests(Go): status/events/env alias coverage
11. Optional: publish ops workflow pack template

---

## 7) Verification Commands

### JS
- `cd js && npm test`

### Go
- `cd go && go test ./...`

### Docs
- `cd docs && npm install && npm run build`

---

## 8) Risks and Mitigations

### Risk: Scope creep into orchestration platform
Mitigation:
- enforce "read-only + optional" rule for additions
- keep policies in docs/pack conventions, not kernel code

### Risk: Divergence between JS and Go behavior
Mitigation:
- define output contract for `status/events`
- add parity tests and cross-implementation checklist

### Risk: docs promise features not implemented
Mitigation:
- merge docs and feature branches together or behind clear "experimental" labels

---

## 9) Out of Scope (for now)

- Hosted UI/dashboard service
- Mandatory external lock/state service
- Built-in retry engine in kernel
- Automatic fan-in coordinator in kernel
- Breaking changes to existing env variable names

---

## 10) Definition of Done

This initiative is done when:
1. New “Operate AYNIG” docs are published and discoverable
2. `status` and `events` read-only commands exist in JS and Go
3. Trailer env aliasing is implemented in JS and Go without breaking compatibility
4. Tests pass in JS and Go
5. No kernel contract semantics were expanded
