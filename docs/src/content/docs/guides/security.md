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
dwp-state: <some-state>
```

they can cause a runner to execute the command registered for that state.

By default, commands are resolved from a **trusted commands ref** (the default
branch, or whatever `--commands-ref` selects) — never from the event branch.
Pushing a branch that redefines `.aynig/command/<state>` does not change what
a runner executes; changing a command requires landing it on the trusted ref,
where branch protection and code review apply.

Two limits to keep in mind:

- Whoever can write to the trusted ref controls what runners execute. Protect
  it accordingly.
- The commit `BODY` and trailers still come from the event branch. They are
  attacker-controlled *input* to your commands — including prompts delivered
  to agents — even though the command *code* is trusted.

`--commands-ref same` disables this protection and resolves commands from each
event branch. Only use it in fully trusted environments.

## Recommendations

- Keep the default trusted commands ref; pin `--commands-ref` to a commit for
  the strongest guarantee.
- Restrict who can push to the trusted ref and to the branches runners watch.
- Keep `.aynig/command/` small and reviewable.
- Treat `BODY` and trailers as untrusted input in your commands.
- Prefer running in a sandboxed environment (CI runner / container / VM).
- Avoid storing secrets in the repo; use environment injection.

## Auditing

Because each step is a commit, you can audit:

- which state was executed
- which prompt/body was provided
- what trailers were set
- what code changes were made

This is one of the core benefits of a Git-native workflow.
