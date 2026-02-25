---
title: Leases and Liveness
description: How AYNIG prevents two runners from executing the same branch.
---

AYNIG implements a distributed mutual exclusion mechanism (“lease”) using Git commits.

Before executing, the runner creates a `working` commit to claim the branch.

Reserved trailers include:

```text
aynig-state: working
aynig-origin-state: <state>
aynig-run-id: <uuid>
aynig-runner-id: <host-id>
aynig-lease-seconds: <ttl>
```

Liveness is tracked via the **committer timestamp** of `HEAD`. If the lease expires, another runner may take over.

While running, commands must keep `aynig-state: working` and renew the lease by pushing commits. Liveness is based only on the committer timestamp of `HEAD`.

If a lease expires, a runner may mark the branch as stalled:

```
aynig-state: stalled
```

For the full specification, see [Kernel Contract](/contract/).
