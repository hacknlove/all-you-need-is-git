---
title: Authoring Commands
description: Create executable commands in .aynig/command/<state>.
---

Commands are executable files located at:

```text
.aynig/command/<state>
```

When `aynig-state: <state>` appears in the latest commit trailer, AYNIG executes the matching command.

## Example

Create a command `review`:

```bash
cat > .aynig/command/review <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

echo "Review requested: $AYNIG_BODY"
EOF

chmod +x .aynig/command/review
```

## Tips

- Commands run with the working directory set to the worktree.
- Keep commands idempotent when possible.
- Commands should emit a new commit advancing the workflow by setting a new `aynig-state`.
- Honor `AYNIG_LOG_LEVEL` if your command supports verbosity.
