---
title: CLI — events
description: Show recent DWP events.
---

```bash
aynig events [options]
```

`aynig events` reads recent commits and summarizes DWP-related trailer data as
an event stream.

It is useful for quickly checking recent state transitions without reading full
commit messages.

## Behavior

- Requires the current directory to be a git repository.
- Reads only the latest commit by default.
- Reads more than one commit only when `--history` is enabled.
- Uses a default history limit of `10` when `--history` is set without `-n`.
- Treats invalid or missing limits as `1`.
- Prints a compact text format by default.
- Prints structured JSON when `--json` is provided.

## Options

### `--history`

Enables scanning recent commit history instead of showing only the latest commit.

Example:

```bash
aynig events --history
```

### `-n, --limit <number>`

Sets how many commits to inspect when `--history` is enabled.

Example:

```bash
aynig events --history --limit 25
```

### `--json`

Outputs the event list as formatted JSON.

Example:

```bash
aynig events --history --json
```

## Output

Default text output looks like this:

```text
<commit> <date> state=<state> run=<run-id> origin=<origin-state> <subject>
```

The `origin=<origin-state>` segment appears only when `dwp-origin-state` is present.

## Examples

Show the latest DWP event:

```bash
aynig events
```

Show the last ten events in history mode:

```bash
aynig events --history
```

Show the last twenty-five events:

```bash
aynig events --history --limit 25
```

Return event data as JSON:

```bash
aynig events --history --json
```
