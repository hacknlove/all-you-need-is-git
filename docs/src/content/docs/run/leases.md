---
title: Leases and Liveness
description: How AYNIG prevents two runners from executing the same branch.
---

AYNIG prevents two runners from executing the same branch by using a Git-based lease.

Before executing, the runner creates a `working` commit to claim the branch.

Reserved trailers include:

```text
dwp-state: working
dwp-origin-state: <state>
dwp-run-id: <uuid>
dwp-runner-id: <host-id>
dwp-lease-seconds: <ttl>
```

Liveness is tracked via the **committer timestamp** of `HEAD`. If the lease expires, another runner may take over.

While running, commands must keep `dwp-state: working` and renew the lease by pushing commits. Liveness is based only on the committer timestamp of `HEAD`.

If a lease expires, a runner may mark the branch as stalled:

```
dwp-state: stalled
```

For the full specification, see [Runner Contract](/contract/).
