---
title: Workflow Design
description: Design guidance for state machines on top of AYNIG.
---

Design workflows as explicit state machines:

- Each state has a single command: `.aynig/command/<state>`
- A command does work and emits the next state by creating a commit
- Keep policies (retries, approvals, gates) in workflow code, not in the kernel

Prefer small, idempotent commands and stable trailer conventions.
