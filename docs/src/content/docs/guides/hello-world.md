---
title: Hello World Workflow
description: A minimal, real workflow you can copy-paste.
---

This guide builds the smallest end-to-end workflow:

- You create a commit with `aynig-state: build`
- AYNIG runs `.aynig/command/build`
- The command emits a new commit advancing the state

## 1) Initialize the repo

```bash
aynig init
```

## 2) Create a command

```bash
cat > .aynig/command/build <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

echo "Build requested: ${AYNIG_BODY}" > build.out

git add build.out
git commit -m "build: done" -m $'Build completed.\n\naynig-state: done'
EOF

chmod +x .aynig/command/build
```

## 3) Request a build (commit protocol)

Create a commit with a trailer:

```text
chore: request build

Build this project.

aynig-state: build
```

## 4) Run

```bash
aynig run
```

## 5) Verify

- `aynig-state` should now be `done`
- `build.out` should exist

## Notes

- The kernel does not decide the next state — the command does.
- Use `COMMANDS.md` in your repo to document which states exist.
