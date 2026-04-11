---
title: Environment Variables
description: What AYNIG exposes to commands.
---

AYNIG exposes commit metadata to commands via environment variables.

Common variables:

- `BODY` — the commit message body (prompt)
- `COMMIT_HASH` — the triggering commit hash
- `LOG_PATH` — absolute path to the command log file for this run
- `LOG_LEVEL` — resolved log level for the run/branch
- `ROLE` — selects `.dwp/roles/<role>/command` when set

Precedence: `--log-level` > `dwp-log-level` trailer > `LOG_LEVEL` env.

Trailers are also exposed as environment variables:

- trailer key `foo: bar` → `FOO=bar`
- trailer key `baz: qux` → `BAZ=qux`

> Implementation note: keys are uppercased and normalized for shells.
