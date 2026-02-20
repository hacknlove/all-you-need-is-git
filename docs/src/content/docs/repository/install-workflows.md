---
title: Install Workflow Packs
description: Copy commands from another repository.
---

```bash
aynig install <repo> [ref] [subfolder]
```

If `<repo>` is in the form `owner/name`, AYNIG expands it to `https://github.com/owner/name.git`.

Examples:

```bash
aynig install https://github.com/org/workflows
aynig install https://github.com/org/workflows main
aynig install https://github.com/org/workflows main path/to/.aynig
```

Tip: keep a `COMMANDS.md` file in your repo to document available commands and expected transitions.

## Ops workflow pack (optional)

This repository ships an optional ops workflow pack under `ops-workflow-pack/`.

```bash
aynig install hacknlove/all-you-need-is-git ops-workflow-pack
```
