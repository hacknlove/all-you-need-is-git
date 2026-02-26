---
title: set-state
description: Create a commit with a new non-working aynig-state.
---

```bash
aynig set-state --aynig-state <state> [options]
```

Creates a commit with `aynig-state: <state>` and validates trailer parsing.

Rules:

- `--aynig-state` is required
- `working` is not allowed (use `aynig set-working`)
- does not copy trailers by default
- resolves remote from `--aynig-remote` or `aynig-remote` trailer and pushes when remote is set

Examples:

```bash
aynig set-state --aynig-state review
aynig set-state --aynig-state failed --prompt-file ./failure.txt
aynig set-state --aynig-state triage --trailer "aynig-issue: 42"
```
