---
title: Authoring Commands
description: Create executable commands in .dwp/command/<state>.
---

Commands are executable files located at:

```text
.dwp/command/<state>
```

When `dwp-state: <state>` appears in the latest commit trailer, AYNIG executes the matching command.

## Example

Create a command `review`:

```bash
cat > .dwp/command/review <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

echo "Review requested: $AYNIG_BODY"
EOF

chmod +x .dwp/command/review
```

## Tips

- Commands run with the working directory set to the worktree.
- Keep commands idempotent when possible.
- Commands should emit a new commit advancing the workflow by setting a new `dwp-state`.
- Honor `AYNIG_LOG_LEVEL` if your command supports verbosity.
