---
title: History Hygiene
description: Keep long-running workflows readable and reliable.
---

AYNIG workflows can produce many small commits. Use a strategy that keeps operational history useful without losing traceability.

## Branch strategy

- Run automation on a dedicated workflow branch per run or per pipeline.
- Merge results into main branches using squash or rebase when appropriate.
- Keep the run branch as an audit trail when needed.

## Checkpointing

Create a checkpoint commit when you finish a major phase.

Recommended trailers:

```
aynig-checkpoint: build-complete
```

This makes it easy to find stable points in the history.

## Squash recommendations

If the run branch is not needed for long-term audit, squash the operational commits and retain:

- Final state transition
- Final artifacts or metadata
- A link to logs or external outputs

## Lease heartbeat cadence

Choose a cadence that matches your runtime and risk tolerance:

- Short tasks: 1 to 5 minutes
- Long tasks: 5 to 15 minutes

Ensure your commands produce a commit often enough to keep the lease active.
