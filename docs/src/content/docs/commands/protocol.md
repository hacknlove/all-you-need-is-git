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

`aynig-state` must appear exactly once in the trailing trailer block.

Optional trailers:

```
aynig-log-level: <debug|info|warn|error>
```

Example:

```text
chore: request build

Build the project and report warnings.

aynig-state: build
priority: high
```
