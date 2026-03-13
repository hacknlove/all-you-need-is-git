---
title: CLI — status
description: Show the current AYNIG state for the repo.
---

```bash
aynig status [--role <name>]
```

Prints the current branch, `HEAD` hash, current `dwp-state`, lease information (if present), and whether the command for the current state exists.

When `--role` is provided (or `AYNIG_ROLE` is set), it checks `.dwp/roles/<role>/command` first.

This command is read-only.
