---
title: install
description: Install workflows from another repository.
---

```bash
aynig install <repo> [ref] [subfolder]
```

Copies the `.aynig/` folder from the specified repo ref and subfolder and reports any overwritten or newly installed files.

If the source repository contains `CONTRACT.md`, it is also copied into `.aynig/CONTRACT.md`.

If `<repo>` is in the form `owner/name`, AYNIG expands it to `https://github.com/owner/name.git`.

The optional ops workflow pack in this repository can be installed with a subfolder path:

```bash
aynig install hacknlove/all-you-need-is-git ops-workflow-pack
```
