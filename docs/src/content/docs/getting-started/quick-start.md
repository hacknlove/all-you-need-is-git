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

Create an executable command at `.aynig/command/human-turn` with something like this:

```bash
#!/bin/bash

case "$(uname)" in
  Darwin)
    osascript -e "display notification \"$AYNIG_BODY\" with title \"AYNIG\""
    afplay /System/Library/Sounds/Glass.aiff >/dev/null 2>&1 &
    ;;
  Linux)
    notify-send "AYNIG" "$AYNIG_BODY"
    paplay /usr/share/sounds/freedesktop/stereo/complete.oga >/dev/null 2>&1 &
    ;;
esac

aynig set-working "human-turn" --lease-seconds 3600
```

## 4) Create a commit with a state

Create a commit whose **message** includes an `aynig-state:` trailer:

```text
chore: Implementation done

Please review and merge.

aynig-state: human-turn
```

## 5) Run AYNIG

```bash
aynig run
```

AYNIG will read `HEAD`, resolve the command for `aynig-state: human-turn`, execute it, and then the command should produce a new commit advancing the state.

## What a command can do

Commands can do anything you can do in a script, including prompting agents.