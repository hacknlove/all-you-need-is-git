---
title: Install Workflow Packs
description: Copy a workflow pack (.aynig/) from another repo.
---

You can install pre-made workflows from other repositories:

```bash
aynig install <repo> [ref] [subfolder]
```

If `<repo>` is in the form `owner/name`, AYNIG expands it to:

```text
https://github.com/owner/name.git
```

### Example: install the ops workflow pack from this repo

```bash
aynig install hacknlove/all-you-need-is-git ops-workflow-pack
```

This copies the source `.aynig/` directory into your current repository.
