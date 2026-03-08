# DWP/GC — Git Commit Binding (draft)

This document defines **DWP over Git commits**.

DWP/GC uses:
- **commit messages** as the event payload
- **Git trailers** as the metadata block
- **branch HEAD** as the latest event used for dispatch

Status: draft.

---

## 1) Trailer parsing

Normative rule: **Git trailers are the standard**.

- A trailer is a `Key: Value` line as defined by Git trailer conventions.
- The trailer block is the **trailing metadata block** of the commit message.
- Parsers SHOULD behave like `git interpret-trailers`.

## 2) Canonical state key

- Canonical key: `dwp-state:`
- The state value selects the handler/command for dispatch.

### Occurrence rule

- `dwp-state` MUST appear in the trailing trailer block.
- If multiple `dwp-state` trailers are present, **last one wins**.

## 3) Case sensitivity

- Trailer keys are **case-sensitive**.
- Therefore, `dwp-state` is distinct from `Dwp-State`.

Rationale: keep the protocol strict and avoid implicit normalization.

## 4) Reserved keys and states

- DWP reserved states apply (`working`, `stalled`, etc.).
- Reserved trailer keys for lease apply when `dwp-state: working` is used:
  - `dwp-origin-state`
  - `dwp-run-id`
  - `dwp-runner-id`
  - `dwp-lease-seconds`

## 5) Dispatch semantics

- The latest event is `HEAD` of the branch.
- The runner reads `HEAD`, parses the trailer block, and extracts `dwp-state`.
- The runner executes the handler for that state.
- Only the handler advances the workflow by emitting a new commit.

## 6) Error handling (binding-level)

Recommended behavior (non-normative):
- If the trailer block cannot be parsed, treat the commit as non-actionable.
- If `dwp-state` is missing, treat the commit as non-actionable.
- If the handler is missing, surface an error but do not mutate history.

---

## Appendix: Reference implementation hints

- Consider implementing parsing by shelling out to `git interpret-trailers` when available,
  or using a library that matches Git trailer rules.

