# AGENTS.md

## PRIME DIRECTIVE
- MUST treat Git commits as the workflow control plane; `dwp-state` trailers are authoritative.
- MUST NOT guess missing behavior; if not specified here or in docs/CONTRACT.md, mark UNKNOWN.
- MUST read `HEAD` only for state decisions; history scanning is optional.
- MUST map state to `.dwp/command/<state>` unless state is reserved (`working`).
- MUST advance state by creating a new commit with a new `dwp-state` (runner does not).
- MUST keep commands idempotent when possible; safe re-run is a design goal.
- MUST respect lease semantics: `working` is reserved and used for mutual exclusion.
- MUST treat `.worktrees/` as runner-owned and ephemeral; do not manage worktrees manually.

## WHAT YOU ARE
- You are a command agent executed by `aynig run`.
- Your job: take action for the current `dwp-state`, then create a new commit with the next `dwp-state` trailer.

## WHERE TO LOOK
- Available states: `.dwp/COMMANDS.md` (if present) and executable scripts in `.dwp/command/`.
- Contract: `.dwp/CONTRACT.md`.

## RELEVANT CLI
- `aynig set-working` renews the lease (`working`) while the command runs.
- `aynig set-state` writes the final state commit when the command completes.

## STATE DISPATCH (SUMMARY)
- `dwp-state: <state>` maps to `.dwp/command/<state>`.
- `working` is reserved for leases; never use it as a terminal state.

## INPUTS YOU RECEIVE
- Environment variables (if provided by runner): `AYNIG_BODY`, `AYNIG_COMMIT_HASH`, `AYNIG_LOG_LEVEL`, `AYNIG_TRAILER_*`.

## COMMIT TRAILER RULES
- Trailer lines use `key: value` in the commit footer.
- `dwp-state` must appear in the trailer block; if multiple are present, last wins.

## LEASE SAFETY
- If `HEAD` is `dwp-state: working`, do not proceed unless the lease is confirmed expired per workflow policy.
