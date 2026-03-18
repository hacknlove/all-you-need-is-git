---
title: CLI — update
description: Download and install the latest AYNIG release.
---

```bash
aynig update
```

`aynig update` downloads and installs the latest AYNIG release by running the
official install script.

## Behavior

- Runs `curl -fsSL https://aynig.org/install.sh | bash` through `sh -c`.
- Streams stdout, stderr, and stdin through the current terminal session.
- Returns an error if the install script fails.
- Always targets the latest available release.

## Options

This command has no flags.

## Examples

Update to the latest AYNIG release:

```bash
aynig update
```

Run the update command in automation where the latest release is acceptable:

```bash
aynig update
```
