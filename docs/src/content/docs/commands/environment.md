---
title: Environment Variables
description: Data available to commands during execution.
---

AYNIG exposes commit data to commands via environment variables.

Common variables:

- `AYNIG_BODY` - commit body (prompt)
- `AYNIG_COMMIT_HASH` - commit hash being processed

Each trailer key becomes an uppercase variable prefixed with `AYNIG_TRAILER_`. Hyphens are replaced with underscores.

```
aynig-state: build
priority: high
```

becomes:

- `AYNIG_TRAILER_AYNIG_STATE=build`
- `AYNIG_TRAILER_PRIORITY=high`
