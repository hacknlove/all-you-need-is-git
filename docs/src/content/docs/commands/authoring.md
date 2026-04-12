---
title: Authoring Commands
description: Create executable commands in .aynig/command/<state>.
---

Commands are executable files located at:

```text
.aynig/command/<state>
```

When `dwp-state: <state>` appears in the latest commit trailer, AYNIG executes the matching command.

## Example

Create a command `review`:

```bash
cat > .aynig/command/review <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

echo "Review requested: $BODY" >&2
printf '%s\n' 'SET_STATE {"state":"done","subject":"review: done","body":"Review completed."}'
EOF

chmod +x .aynig/command/review
```

## Tips

- Commands run with the working directory set to the worktree.
- Keep commands idempotent when possible.
- Commands should emit a `SET_STATE {...}` line on stdout instead of creating the final commit directly.
- The runner watches stdout for lines that begin with `SET_STATE ` and applies the last valid one after the command exits.
- stderr is not parsed for state transitions.
- Honor `LOG_LEVEL` if your command supports verbosity.
