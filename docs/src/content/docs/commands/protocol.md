---
title: Commit Protocol
description: How AYNIG reads intent from commit messages.
---

AYNIG chooses the command from the latest commit (`HEAD`) by reading trailers in the commit message.

Minimum required trailer:

```text
dwp-state: <state>
```

`dwp-state` must appear in the trailing trailer block. If multiple appear, last wins.

Optional trailers:

```text
dwp-source: git:<remote-name>
dwp-log-level: <debug|info|warn|error>
```

Recommended structure:

```text
<subject>

<prompt/body>

dwp-state: <state>
<key>: <value>
```

Only the workflow command decides the next state by creating a new commit with a new `dwp-state` trailer.
