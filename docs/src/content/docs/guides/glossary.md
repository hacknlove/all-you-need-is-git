---
title: Glossary
description: Shared vocabulary used throughout AYNIG docs.
---

- **state**: the value of the `dwp-state:` trailer used to select a command.
- **trailer**: key/value metadata at the end of a commit message (e.g. `dwp-state: build`).
- **command**: an executable file at `.dwp/command/<state>`.
- **step**: one workflow step: select command → execute → check new `HEAD`.
- **worktree**: an ephemeral Git worktree created by the runner to isolate execution.
- **lease**: the distributed lock mechanism implemented via `dwp-state: working` commits.
- **origin state**: the state the runner was about to execute before claiming a `working` lease.
