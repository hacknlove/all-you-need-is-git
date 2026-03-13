---
title: set-working
description: Create or refresh a working lease commit.
---

```bash
aynig set-working [options]
```

Creates a `working` state commit with lease trailers and validates trailer parsing.

Defaults:

- copies only `dwp-*` trailers from `HEAD`
- keeps `dwp-run-id` if present (generates one if missing)
- writes `dwp-state: working`, `dwp-runner-id`, and `dwp-lease-seconds`
- resolves remote from `--dwp-remote` or `dwp-source` trailer and pushes when remote is set

Examples:

```bash
aynig set-working
aynig set-working --subject "chore: heartbeat" --prompt "Lease refresh"
aynig set-working --prompt-file ./prompt.txt --trailer "dwp-note: retry"
aynig set-working --lease-seconds 600
```
