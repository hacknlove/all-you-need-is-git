---
title: CLI — install
description: Install a workflow pack from another repository.
---

```bash
aynig install <repo> [ref] [subfolder]
```

Copies the source repository’s `.aynig/` directory into the current repository.

If `<repo>` is in the form `owner/name`, AYNIG expands it to `https://github.com/owner/name.git`.

Example:

```bash
aynig install hacknlove/all-you-need-is-git ops-workflow-pack
```
