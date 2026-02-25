---
title: Commit Protocol
description: How AYNIG reads intent from commit messages.
---

AYNIG dispatches from the latest commit (`HEAD`) by reading trailers in the commit message.

Minimum required trailer:

```text
aynig-state: <state>
```

`aynig-state` must appear exactly once in the trailing trailer block.

Optional trailers:

```text
aynig-remote: <remote-name>
aynig-log-level: <debug|info|warn|error>
```

Recommended structure:

```text
<subject>

<prompt/body>

aynig-state: <state>
<key>: <value>
```

Only the workflow command decides the next state by creating a new commit with a new `aynig-state` trailer.
