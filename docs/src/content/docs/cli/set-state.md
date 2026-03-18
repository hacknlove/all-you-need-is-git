---
title: CLI — set-state
description: Create a commit with a new non-working dwp-state.
---

```bash
aynig set-state --dwp-state <state> [options]
```

`aynig set-state` creates a state commit for any non-`working` DWP state.
Use it when a DWP run has moved into a new state such as `review`, `failed`, or
`triage` and you want the runner to take action on the next cycle.

The command writes `dwp-state: <state>` automatically, validates trailer
format, and can push the current branch when a remote is configured.

## Behavior

- `--dwp-state` is required.
- `working` is not allowed. Use `aynig set-working` for that state.
- The default commit title is `chore: set <state>` when `--subject` is not provided.
- The body can be provided with `--prompt`, `--prompt-file`, or `--prompt-stdin`.
- Trailers from `HEAD` are not copied by default, because this is treated as a completely new state.
- The push remote is resolved from `--dwp-remote` first, then from the `dwp-source` trailer on `HEAD`.
- When a remote is resolved, the command pushes the current branch after the commit.

## Options

### `--dwp-state <state>`

Required. Sets the new non-`working` DWP state and writes it as the
`dwp-state` trailer.

Example:

```bash
aynig set-state --dwp-state review
```

### `--subject <text>`

Sets the commit title. If omitted, AYNIG uses `chore: set <state>`.

Example:

```bash
aynig set-state --dwp-state review --subject "review: ready for feedback"
```

### `--prompt <text>`

Sets the commit body directly on the command line.

Example:

```bash
aynig set-state --dwp-state failed --prompt "Tests failed in CI; investigating flaky snapshot output."
```

### `--prompt-file <path>`

Reads the commit body from a file.

Example:

```bash
aynig set-state --dwp-state failed --prompt-file ./failure.txt
```

### `--prompt-stdin`

Reads the commit body from standard input. This is useful when generating the
body from another command or when writing a longer message in a shell pipeline.

Example:

```bash
printf "Blocked on API credentials from staging environment.\n" | aynig set-state --dwp-state triage --prompt-stdin
```

### `--dwp-remote <name>`

Uses the named git remote for the post-commit push and records it in a DWP trailer as
`dwp-source: git:<name>`.

Example:

```bash
aynig set-state --dwp-state review --dwp-remote origin
```

### `--trailer <key:value>`

Adds an extra trailer to the commit. Repeat the flag to add more than one.
`dwp-state` cannot be set with `--trailer`; it is managed by `--dwp-state`.

Example:

```bash
aynig set-state --dwp-state triage --trailer "dwp-issue: 42"
```

Repeated example:

```bash
aynig set-state --dwp-state review --trailer "dwp-issue: 42" --trailer "dwp-note: waiting on design sign-off"
```

## Examples

Set a simple state with the default title:

```bash
aynig set-state --dwp-state review
```

Set both the title and the body explicitly:

```bash
aynig set-state --dwp-state review --subject "review: ready for feedback" --prompt "Core flow is implemented and ready for another pass."
```

Read a longer body from a file:

```bash
aynig set-state --dwp-state failed --prompt-file ./failure.txt
```

Pipe a generated body through stdin:

```bash
git status --short | aynig set-state --dwp-state triage --prompt-stdin
```

Push through a specific remote and add extra trailers:

```bash
aynig set-state --dwp-state triage --dwp-remote origin --trailer "dwp-issue: 42" --trailer "dwp-note: waiting on upstream fix"
```
