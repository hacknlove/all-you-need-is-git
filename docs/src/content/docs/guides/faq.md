---
title: FAQ / Troubleshooting
description: Common problems and quick fixes.
---

## `aynig run` does nothing

Common causes:

- There is no `aynig-state:` trailer in `HEAD`.
- The runner is scanning branches that have no actionable state.
- Your command for the state does not exist.

Checks:

```bash
aynig status
```

Confirm:

- `aynig-state` is present
- `command: exists`

## "command missing"

Your state `foo` must have an executable file:

```text
.aynig/command/foo
```

Make it executable:

```bash
chmod +x .aynig/command/foo
```

## "permission denied" running a command

- Make sure the command is executable.
- Ensure the script has a valid shebang:

```bash
#!/usr/bin/env bash
```

## Worktrees keep accumulating

- The default worktree directory is `.worktrees/`.
- It is meant to be ephemeral and is ignored by default.

If you see leftovers after crashes, clean them:

```bash
rm -rf .worktrees/
```

Then run again.

## "working" commits appear in history

This is the lease mechanism. The runner uses `aynig-state: working` to claim a branch.

See: Run and Operate → Leases and Liveness.

## My workflow is stuck in `working`

A runner may have crashed. If the lease is expired, another runner can take over.

- Inspect the latest commit trailers (`aynig-run-id`, `aynig-lease-seconds`).
- Decide whether to continue or mark it stalled.

See: Guides → Ops Workflow Pack (stalled recovery).
