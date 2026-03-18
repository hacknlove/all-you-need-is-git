---
title: Hello World Workflow
description: A minimal, real workflow you can copy-paste.
---

This guide builds the smallest end-to-end workflow:

- You create a commit with `dwp-state: build`
- AYNIG runs `.dwp/command/build`
- The command emits a new commit advancing the state

## 1) Initialize the repo

```bash
aynig init
```

## 2) Create a command

```bash
cat > .dwp/command/build <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

echo "Build requested: ${AYNIG_BODY}" > build.out

git add build.out
git commit -m "build: done" -m $'Build completed.\n\ndwp-state: done'
EOF

chmod +x .dwp/command/build
```

## 3) Request a build (commit protocol)

Create a commit with a trailer:

```text
chore: request build

Build this project.

dwp-state: build
```

## 4) Run

```bash
aynig run
```

## 5) Verify

- `dwp-state` should now be `done`
- `build.out` should exist

## Notes

- AYNIG does not decide the next state — the command does.
- Use `COMMANDS.md` in your repo to document which states exist.
