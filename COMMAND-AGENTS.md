# AGENTS.md

## PRIME DIRECTIVE
- MUST treat Git commits as the workflow control plane; `dwp-state` trailers are authoritative.
- MUST NOT guess missing behavior; if not specified here or in docs/CONTRACT.md, mark UNKNOWN.
- MUST read `HEAD` only for state decisions; history scanning is optional.
- MUST map state to `.aynig/command/<state>` unless state is reserved (`working`).
- MUST advance state by emitting `SET_STATE {...}` on stdout; the runner materializes the commit.
- MUST keep commands idempotent when possible; safe re-run is a design goal.
- MUST respect lease semantics: `working` is reserved and used for mutual exclusion.
- MUST treat `.worktrees/` as runner-owned and ephemeral; do not manage worktrees manually.

## WHAT YOU ARE
- You are a command agent executed by `aynig run`.
- Your job: take action for the current `dwp-state`, then emit a `SET_STATE {...}` line describing the next state.

## WHERE TO LOOK
- Available states: `.aynig/COMMANDS.md` (if present) and executable scripts in `.aynig/command/`.
- Contract: `.aynig/CONTRACT.md`.

## RELEVANT CLI
- `aynig set-working` renews the lease (`working`) while the command runs.
- `aynig set-state` is a manual escape hatch for humans; normal command completion goes through `SET_STATE {...}`.

## STATE DISPATCH (SUMMARY)
- `dwp-state: <state>` maps to `.aynig/command/<state>`.
- `working` is reserved for leases; never use it as a terminal state.

## INPUTS YOU RECEIVE
- Environment variables (if provided by runner): `BODY`, `COMMANDS_PATH`, `COMMIT_HASH`, `LOG_LEVEL`, `ROLE`, `STDOUT_LOG_PATH`, `STDERR_LOG_PATH`, `WORKTREE_PATH`.
- `COMMANDS_PATH` is the `.aynig` directory your command was resolved from. Commands run from a trusted commands ref (default branch unless the runner overrides it), so relative `.aynig/...` paths under the worktree may differ from your own files; use `COMMANDS_PATH` to reference them.
- Commit trailers are also exposed as uppercase env vars, with dashes converted to underscores. Example: `dwp-state` -> `DWP_STATE`.

## OUTPUT PROTOCOL
- The runner watches stdout for lines that begin with `SET_STATE `.
- The text after `SET_STATE ` must be a single-line JSON object.
- The last valid `SET_STATE` line wins.
- Use `"keep_trailers": true` when you want AYNIG to preserve existing non-reserved `dwp-*` workflow trailers.
- The runner only applies `SET_STATE` if your process exits with code `0`.
- If your process exits non-zero, AYNIG moves the branch to `stalled`.
- If your process exits `0` without any valid `SET_STATE`, AYNIG refreshes `working` and keeps waiting.
- The runner does not parse stderr for state transitions.

## COMMIT TRAILER RULES
- Trailer lines use `key: value` in the commit footer.
- `dwp-state` must appear in the trailer block; if multiple are present, last wins.

## LEASE SAFETY
- If `HEAD` is `dwp-state: working`, do not proceed unless the lease is confirmed expired per workflow policy.
