---
title: Security & Trust Model
description: What AYNIG assumes and what you must control.
---

AYNIG executes code from your repository (`.aynig/command/<state>`). Treat it like CI.

## What AYNIG assumes

- You trust the **Git branch** you are running on.
- You control who can push to the branch (or to the remote you run against).
- Workflow commands are reviewed like any other code.

## Threat model

If an attacker can push a commit that sets:

```text
aynig-state: <some-state>
```

they can potentially cause a runner to execute `.aynig/command/<some-state>`.

## Recommendations

- Run workflows on dedicated branches.
- Restrict who can push to those branches.
- Keep `.aynig/command/` small and reviewable.
- Prefer running in a sandboxed environment (CI runner / container / VM).
- Avoid storing secrets in the repo; use environment injection.

## Auditing

Because each tick is a commit, you can audit:

- which state was executed
- which prompt/body was provided
- what trailers were set
- what code changes were made

This is one of the core benefits of a Git-native control plane.
