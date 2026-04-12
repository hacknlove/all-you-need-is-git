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

## Output protocol

The command declares the next state by writing a single-line JSON payload to stdout:

```text
SET_STATE {"state":"review","subject":"review: ready","body":"Line 1\nLine 2"}
```

Rules:

- The runner only interprets stdout for this protocol.
- The prefix must be exactly `SET_STATE ` at the beginning of the line.
- The payload must be valid JSON on a single line.
- If multiple valid `SET_STATE` lines are emitted, the last one wins.
- The runner creates the final commit only if the command exits with code `0`.
- If the command exits non-zero, the runner ignores `SET_STATE` and marks the branch as `stalled`.
- If the command exits `0` without any valid `SET_STATE`, the runner emits a fresh `working` commit and keeps waiting.
