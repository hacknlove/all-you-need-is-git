---
title: Authoring Commands
description: Write executable scripts that AYNIG dispatches.
---

Commands are executable files in `.aynig/command/<state>`.

Example:

```bash
#!/usr/bin/env bash
set -euo pipefail

echo "Review requested: $AYNIG_BODY"
```

Make the file executable:

```bash
chmod +x .aynig/command/review
```

Tips:

- Commands run with the working directory set to the worktree
- Keep commands idempotent when possible
- Emit a new commit with the next `aynig-state`
