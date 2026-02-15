---
title: install
description: Install workflows from another repository.
---

```bash
aynig install <repo> [ref] [subfolder]
```

This clones the source repository, copies its `.aynig/` folder, and reports any overwritten or newly installed files. If `COMMANDS.md` is modified, the installer can invoke `opencode` or `claude` to assist with merge.
