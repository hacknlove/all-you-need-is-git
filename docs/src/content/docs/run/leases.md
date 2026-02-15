---
title: Leases and Liveness
description: How AYNIG avoids double execution.
---

AYNIG uses Git commits as a distributed lock. Before executing, the runner writes:

```
aynig-state: working
```

If the push fails, another runner won the lease.

While running, commands must keep `aynig-state: working` and renew the lease by pushing commits. Liveness is based on the committer timestamp of `HEAD`.

If a lease expires, a runner may mark the branch as stalled:

```
aynig-state: stalled
```

For the full specification, see [Kernel Contract](/contract/).
