# AGENTS.md

## PRIME DIRECTIVE
- MUST treat Git commits as the workflow control plane; `aynig-state` trailers are authoritative.
- MUST NOT guess missing behavior; if not specified here or in docs/CONTRACT.md, mark UNKNOWN.
- MUST read `HEAD` only for state decisions; history scanning is optional and outside the kernel.
- MUST map state to `.aynig/command/<state>` unless state is reserved (`working`).
- MUST advance state by creating a new commit with a new `aynig-state` (runner does not).
- MUST keep commands idempotent when possible; safe re-run is a design goal.
- MUST respect lease semantics: `working` is reserved and used for mutual exclusion.
- MUST treat `.worktrees/` as runner-owned and ephemeral; do not manage worktrees manually.

## QUICK START
- Install CLI: `curl -fsSL https://aynig.org/install.sh | bash`
- Init repo: `aynig init` (creates `.aynig/`, `.worktrees/`, `.aynig/CONTRACT.md`, updates `.gitignore`)
- Add command: create executable `.aynig/command/<state>`
- Create commit with trailer `aynig-state: <state>`
- Run one tick: `aynig run` (use `--current-branch only` if testing on current branch)

## REPO LAYOUT
| Path | Purpose | Owned-by | Notes |
| --- | --- | --- | --- |
| `CONTRACT.md` | Kernel contract (canonical) | human | Do not diverge from docs/contract. |
| `.aynig/` | Workflow assets | human/agent | Created by `aynig init`. |
| `.aynig/command/` | Command scripts, one per state | human/agent | Executable files; name == state. |
| `.aynig/COMMANDS.md` | List of available states | human | Optional; used as reference. |
| `.aynig/CONTRACT.md` | Contract copy inside repo | human | Created by `aynig init`; keep in sync with root. |
| `.aynig/logs/` | Command stdout/stderr logs | runner | Files: `<commit-hash>.log`. |
| `.worktrees/` | Ephemeral worktrees | runner | Created/cleaned by runner. |
| `ops-workflow-pack/` | Optional workflow pack | human/agent | Provides `stalled` recovery command. |
| `docs/` | Documentation site | human | Source for behavior in this file. |

## STATE MACHINE
Definitions:
- **state**: value of `aynig-state: <state>` in trailing commit trailers.
- **dispatch**: map `state` → `.aynig/command/<state>`.
- **tick**: read `HEAD` → dispatch → execute → validate new `HEAD`.

Resolution algorithm (canonical):
1) Runner selects branches (local by default; remote if `--aynig-remote` or `aynig-remote` trailer is set).
2) For each branch, read `HEAD` commit trailers.
3) If `aynig-state` exists, resolve command at `.aynig/command/<state>`.
4) Execute command in a dedicated worktree.
5) Valid completion: new `HEAD` has `aynig-state: <state>` and `state != working`.

Mapping table:
| `aynig-state` value | Command path | Expected inputs | Expected outputs | Next state rules |
| --- | --- | --- | --- | --- |
| `<any-non-working>` | `.aynig/command/<state>` | `AYNIG_BODY`, `AYNIG_COMMIT_HASH`, `AYNIG_TRAILER_*`, `AYNIG_LOG_LEVEL` | New commit with `aynig-state: <next>` | Command MUST decide next state and commit it. |
| `working` (reserved) | NONE | Lease trailers | Lease commit pushed to branch | Runner uses this to claim work; commands MUST NOT use as terminal state. |
| `stalled` (optional) | `.aynig/command/stalled` if installed | Same as above | New commit with next state | Used for recovery after lease expiry; behavior is workflow-defined. |

## COMMIT TRAILERS
Syntax: trailer lines in commit message footer, `key: value` (exact; no `key=value`).

| Trailer key | Allowed values/pattern | Meaning | Example |
| --- | --- | --- | --- |
| `aynig-state` | string, REQUIRED, must appear once | Dispatch key / current state | `aynig-state: build` |
| `aynig-remote` | remote name | Remote to use for `aynig run` | `aynig-remote: origin` |
| `aynig-log-level` | `debug|info|warn|error` | Log level override | `aynig-log-level: info` |
| `aynig-origin-state` | string | Original state before `working` | `aynig-origin-state: build` |
| `aynig-run-id` | uuid string | Lease/run identifier | `aynig-run-id: 123e4567...` |
| `aynig-runner-id` | host id string | Runner identity | `aynig-runner-id: host-01` |
| `aynig-lease-seconds` | integer seconds | Lease TTL | `aynig-lease-seconds: 600` |
| `aynig-stalled-run` | run id | Run being marked stalled | `aynig-stalled-run: <run-id>` |
| `aynig-attempt` | integer | Retry attempt (convention) | `aynig-attempt: 2` |
| `aynig-max-attempts` | integer | Max retries (convention) | `aynig-max-attempts: 3` |
| `aynig-next-on-failure` | state string | Failure fallback (convention) | `aynig-next-on-failure: ops-triage` |
| `aynig-parent-run` | string | Parent/coordinator run (convention) | `aynig-parent-run: run-123` |
| `aynig-correlation-id` | string | Cross-branch correlation (convention) | `aynig-correlation-id: 2026-02-16-build-42` |
| `aynig-checkpoint` | string | Checkpoint marker (convention) | `aynig-checkpoint: package-published` |
| `aynig-note` | string | Human note (convention) | `aynig-note: manual rollback completed` |

## RUNNER CONTRACT
- Reads `HEAD` commit trailers; ignores title.
- Source of truth is the branch (local or remote in remote mode).
- Creates a `working` commit to claim branch; pushes; aborts if push fails.
- Liveness is based on `HEAD` committer timestamp and `aynig-lease-seconds`.
- Workdir is a dedicated worktree under `.worktrees/`.
- Environment variables:
  - `AYNIG_BODY` (commit body)
  - `AYNIG_COMMIT_HASH` (triggering commit)
  - `AYNIG_LOG_LEVEL` (resolved log level)
  - `AYNIG_TRAILER_*` (uppercase, normalized keys)
- Logging: stdout/stderr to `.aynig/logs/<commit-hash>.log`.
- Exit code semantics: UNKNOWN (not specified in docs).
- Retries: NOT implemented by kernel; workflow code must decide.

## PLAYBOOKS
### stalled (state didn’t advance / lease expired)
- Trigger: `HEAD` is `aynig-state: working` and lease appears expired.
- Checks:
  - `aynig-run-id`, `aynig-lease-seconds`, committer timestamp of `HEAD`.
  - `COMMANDS.md` for available states.
- Actions:
  - Create commit `aynig-state: stalled` and `aynig-stalled-run: <run-id>`.
  - Optionally run `aynig run` to let `.aynig/command/stalled` handle recovery.
- Safety:
  - Only act after lease expiry; avoid double execution.
  - Preserve `aynig-run-id` in the stalled marker.

### retry (re-run same state)
- Trigger: command failed and workflow policy allows retry.
- Checks:
  - Current `aynig-attempt`, `aynig-max-attempts`.
  - Ensure `.aynig/command/<state>` exists and is executable.
- Actions:
  - Create new commit with same `aynig-state` and increment `aynig-attempt`.
  - Keep prompt/body as needed (use `--prompt` or `--prompt-file`).
- Safety:
  - Stop when attempt reaches max.
  - Avoid retry loops without increasing attempt.

### fail (terminal failure)
- Trigger: command failure and no retry remaining.
- Checks:
  - Decide terminal state (example: `failed`).
- Actions:
  - Create commit with `aynig-state: failed` and optional `aynig-note`.
- Safety:
  - Do not use `working` as a terminal state.

### resume/continue
- Trigger: you have fixed the cause and want to continue.
- Checks:
  - `aynig status` to confirm current state and command existence.
- Actions:
  - If state is correct, run `aynig run`.
  - Otherwise, create commit with intended next `aynig-state`.
- Safety:
  - Ensure you are on the intended branch (or remote via `aynig-remote`).

## MINIMAL EXAMPLES
1) Request a build (commit message body):
```
chore: request build

Build this project.

aynig-state: build
```

2) Minimal command stub:
```bash
#!/usr/bin/env bash
set -euo pipefail
echo "Review requested: $AYNIG_BODY"
```

3) Lease marker (reserved):
```
aynig-state: working
aynig-origin-state: build
aynig-run-id: 123e4567-e89b-12d3-a456-426614174000
aynig-runner-id: host-01
aynig-lease-seconds: 600
```

4) Stalled marker:
```
aynig-state: stalled
aynig-stalled-run: 123e4567-e89b-12d3-a456-426614174000
```

5) Retry marker:
```
aynig-state: build
aynig-attempt: 2
aynig-max-attempts: 3
```
