---
title: Ops Workflow Pack
description: Optional operational commands you can install into any repo.
---

This repository includes an optional workflow pack in `ops-workflow-pack/`.

## Install

```bash
aynig install hacknlove/all-you-need-is-git ops-workflow-pack
```

This copies a `.aynig/` directory into your current repository.

## Included states

See the pack’s `COMMANDS.md` after installing.

Currently included:

- `stalled` — recover from a stalled lease

## How `stalled` is intended to be used

When a runner crashes or a lease expires, a branch may be left in:

```text
aynig-state: working
```

(or a recovery marker such as `stalled`, depending on your policy).

The `stalled` command provides a human/agent playbook:

- inspect the job start
- inspect `working` commits and changes
- decide whether to retry, continue, switch state, or fail

It is intentionally minimal and designed to be adapted.
