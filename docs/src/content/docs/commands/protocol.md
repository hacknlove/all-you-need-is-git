---
title: Commit Protocol
description: Structure commit messages for AYNIG.
---

AYNIG reads the latest commit and interprets:

- Title: human-only (ignored by the runner)
- Body: prompt delivered to the command
- Trailers: structured metadata

The mandatory trailer is:

```
aynig-state: <state>
```

Example:

```text
chore: request build

Build the project and report warnings.

aynig-state: build
priority: high
```
