---
title: Glossary
description: Shared vocabulary used throughout AYNIG docs.
---

- **state**: the value of the `aynig-state:` trailer used to select a command.
- **trailer**: key/value metadata at the end of a commit message (e.g. `aynig-state: build`).
- **command**: an executable file at `.aynig/command/<state>`.
- **tick**: one workflow step: dispatch → execute command → validate new `HEAD`.
- **worktree**: an ephemeral Git worktree created by the runner to isolate execution.
- **lease**: the distributed lock mechanism implemented via `aynig-state: working` commits.
- **origin state**: the state the runner was about to execute before claiming a `working` lease.
