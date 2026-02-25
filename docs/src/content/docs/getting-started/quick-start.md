---
title: Quick Start
description: Run your first AYNIG workflow in minutes.
---

## 1. Install the CLI

```bash
curl -fsSL https://aynig.org/install.sh | bash
```

If your PATH does not include `~/.local/bin`, add it to your shell profile.

## 2. Initialize the repository

```bash
aynig init
```

This creates:

- `.aynig/` with starter files
- `.worktrees/` for ephemeral worktrees
- `.worktrees/` and `.aynig/logs/` entries in `.gitignore`

## 3. Add a command

Create a command directory:

```bash
mkdir -p .aynig/command
```

Create an executable file at `.aynig/command/build` with this content:

```bash
#!/usr/bin/env bash
set -euo pipefail

echo "Build requested: $AYNIG_BODY"
```

Make it executable:

```bash
chmod +x .aynig/command/build
```

## 4. Create a commit with a state

```text
chore: request a build

Build the project.

aynig-state: build
```

## 5. Run AYNIG

```bash
aynig run
```

Next steps:

- [Install CLI](/installation/cli/)
- [Repository setup](/repository/init/)
- [Authoring commands](/commands/authoring/)
