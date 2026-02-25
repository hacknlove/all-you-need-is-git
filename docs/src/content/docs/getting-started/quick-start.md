---
title: Quick Start
description: Install the CLI, initialize a repo, and run your first command.
---

## 1) Install the CLI

```bash
curl -fsSL https://aynig.org/install.sh | bash
```

If your `PATH` does not include `~/.local/bin`, add it to your shell profile.

## 2) Initialize the repository

Inside the repository where you want to use AYNIG:

```bash
aynig init
```

This creates:

- `.aynig/` with starter files
- `.worktrees/` for ephemeral execution worktrees
- `.worktrees/` and `.aynig/logs/` entries in `.gitignore`

## 3) Add a command

Create an executable command at `.aynig/command/build`:

```bash
mkdir -p .aynig/command
cat > .aynig/command/build <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

echo "Build requested: $AYNIG_BODY"
EOF
chmod +x .aynig/command/build
```

## 4) Create a commit with a state

Create a commit whose **message** includes an `aynig-state:` trailer:

```text
chore: request a build

Build the project.

aynig-state: build
```

## 5) Run AYNIG

```bash
aynig run
```

AYNIG will read `HEAD`, resolve the command for `aynig-state: build`, execute it, and then the command should produce a new commit advancing the state.
