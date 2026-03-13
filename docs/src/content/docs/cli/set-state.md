---
title: set-state
description: Create a commit with a new non-working dwp-state.
---

```bash
aynig set-state --dwp-state <state> [options]
```

Creates a commit with `dwp-state: <state>` and validates trailer parsing.

Rules:

- `--dwp-state` is required
- `working` is not allowed (use `aynig set-working`)
- does not copy trailers by default
- resolves remote from `--dwp-remote` or `dwp-source` trailer and pushes when remote is set

Examples:

```bash
aynig set-state --dwp-state review
aynig set-state --dwp-state failed --prompt-file ./failure.txt
aynig set-state --dwp-state triage --trailer "dwp-issue: 42"
```
