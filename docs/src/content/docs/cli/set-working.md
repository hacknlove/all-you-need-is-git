---
title: CLI — set-working
description: Create or refresh a working lease commit.
---

```bash
aynig set-working [options]
```

`aynig set-working` creates a `working` state commit for the current branch.
Use it to claim or refresh a lease while work is in progress.

The command writes the working-state trailers automatically, preserves useful
run metadata, and can push the branch when a remote is configured.

## Behavior

- The default commit title is `chore: working` when `--subject` is not provided.
- The default commit body is `Lease heartbeat` when no prompt source is given.
- Only one prompt source may be used: `--prompt`, `--prompt-file`, or
  `--prompt-stdin`.
- AYNIG requires an existing `dwp-state` or `dwp-origin-state` trailer on
  `HEAD` so it knows what state the lease belongs to.
- The command writes `dwp-state: working`, `dwp-origin-state`, `dwp-run-id`,
  `dwp-runner-id`, and `dwp-lease-seconds`.
- `dwp-run-id` is kept from `HEAD` when present, or generated when missing.
- `dwp-lease-seconds` comes from `--lease-seconds`, then from `HEAD`, then from
  the default of `300` seconds.
- Existing `dwp-*` trailers are copied from `HEAD` except for the reserved
  working-state trailers that AYNIG manages itself.
- The push remote is resolved from `--dwp-remote` first, then from the
  `dwp-source` trailer on `HEAD`.
- When a remote is resolved, the command pushes the current branch after the commit.

## Options

### `--subject <text>`

Sets the commit title. If omitted, AYNIG uses `chore: working`.

Example:

```bash
aynig set-working --subject "chore: heartbeat"
```

### `--prompt <text>`

Sets the commit body directly on the command line.

Example:

```bash
aynig set-working --prompt "Lease refresh while integration tests run."
```

### `--prompt-file <path>`

Reads the commit body from a file.

Example:

```bash
aynig set-working --prompt-file ./prompt.txt
```

### `--prompt-stdin`

Reads the commit body from standard input.

Example:

```bash
printf "Waiting on review feedback.\n" | aynig set-working --prompt-stdin
```

### `--lease-seconds <seconds>`

Overrides the lease duration written to `dwp-lease-seconds`.

Example:

```bash
aynig set-working --lease-seconds 600
```

### `--dwp-remote <name>`

Uses the named git remote for the post-commit push and records it as
`dwp-source: git:<name>`.

Example:

```bash
aynig set-working --dwp-remote origin
```

### `--trailer <key:value>`

Adds an extra trailer to the commit. Repeat the flag to add more than one.

Example:

```bash
aynig set-working --trailer "dwp-note: retry"
```

Repeated example:

```bash
aynig set-working --trailer "dwp-note: retry" --trailer "dwp-ticket: 42"
```

## Examples

Create a basic working lease commit:

```bash
aynig set-working
```

Set both the title and the body explicitly:

```bash
aynig set-working --subject "chore: heartbeat" --prompt "Lease refresh while the migration is running."
```

Read the body from a file:

```bash
aynig set-working --prompt-file ./prompt.txt
```

Pipe the body through stdin and extend the lease:

```bash
printf "Still debugging the failing worker.\n" | aynig set-working --prompt-stdin --lease-seconds 900
```

Push through a specific remote and add extra trailers:

```bash
aynig set-working --dwp-remote origin --trailer "dwp-note: retry" --trailer "dwp-ticket: 42"
```
