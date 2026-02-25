---
title: Runbooks
description: Operational guidance for common situations.
---

## Command failures

If a command fails, AYNIG does not retry automatically. Your workflow should decide what to do next.

Common patterns:

- Emit a new commit with a failure state (e.g. `aynig-state: failed`)
- Increment an attempt trailer (`aynig-attempt:`) and re-emit the same state
- Notify a human (out of band) and stop advancing

## Working lease stalled

If a branch is stuck at:

```text
aynig-state: working
```

and the lease appears expired, a runner may take over and mark it `stalled` (see Kernel Contract). Keep recovery behavior explicit in your workflow.
