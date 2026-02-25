---
title: Environment Variables
description: What AYNIG exposes to commands.
---

AYNIG exposes commit metadata to commands via environment variables.

Common variables:

- `AYNIG_BODY` — the commit message body (prompt)
- `AYNIG_COMMIT_HASH` — the triggering commit hash
- `AYNIG_LOG_LEVEL` — resolved log level for the run/branch

Precedence: `--log-level` > `aynig-log-level` trailer > `AYNIG_LOG_LEVEL` env.

Trailers are also exposed as environment variables:

- trailer key `foo: bar` → `AYNIG_TRAILER_FOO=bar`
- trailer key `baz: qux` → `AYNIG_TRAILER_BAZ=qux`

> Implementation note: keys are uppercased and normalized for shells.
