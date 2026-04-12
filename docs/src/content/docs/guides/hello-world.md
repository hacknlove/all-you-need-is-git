---
title: Hello World Workflow
description: A minimal, real workflow you can copy-paste.
---

This guide builds the smallest end-to-end workflow:

- You create a commit with `dwp-state: build`
- AYNIG runs `.aynig/command/build`
- The command emits `SET_STATE {...}` and AYNIG writes the next state commit

## 1) Initialize the repo

```bash
aynig init
```

## 2) Create a command

```bash
cat > .aynig/command/build <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

echo "Build requested: ${BODY}" > build.out

git add build.out
printf '%s\n' 'SET_STATE {"state":"done","subject":"build: done","body":"Build completed."}'
EOF

chmod +x .aynig/command/build
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

- AYNIG does not infer the next state — the command declares it with `SET_STATE {...}`.
- Use `COMMANDS.md` in your repo to document which states exist.
