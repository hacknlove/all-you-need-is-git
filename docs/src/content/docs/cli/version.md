---
title: CLI — version
description: Print the current AYNIG version.
---

```bash
aynig version
```

`aynig version` prints the currently installed AYNIG version string.

## Behavior

- Prints a single version value to stdout.
- Does not inspect the repository or modify any files.
- The CLI also accepts `aynig -v` and `aynig --version`.

## Options

This command has no flags.

## Examples

Print the installed version:

```bash
aynig version
```

Use the short alias:

```bash
aynig -v
```

Use the long alias:

```bash
aynig --version
```
