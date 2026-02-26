---
title: set-working
description: Create or refresh a working lease commit.
---

```bash
aynig set-working [options]
```

Creates a `working` state commit with lease trailers and validates trailer parsing.

Defaults:

- copies only `aynig-*` trailers from `HEAD`
- keeps `aynig-run-id` if present (generates one if missing)
- writes `aynig-state: working`, `aynig-runner-id`, and `aynig-lease-seconds`
- resolves remote from `--aynig-remote` or `aynig-remote` trailer and pushes when remote is set

Examples:

```bash
aynig set-working
aynig set-working --subject "chore: heartbeat" --prompt "Lease refresh"
aynig set-working --prompt-file ./prompt.txt --trailer "aynig-note: retry"
```
