---
title: Tradeoffs
description: Honest limitations and what AYNIG does not do.
---

AYNIG keeps the kernel small and deterministic. That means some things are intentionally out of scope.

## Tradeoffs

- No hosted UI or dashboards
- No built-in workflow policies or retries
- You own the scripts, tools, and operational glue

## What you still need

- A way to run `aynig run` (human, cron, or CI)
- Command scripts that handle your domain logic
- Your own safety checks and guardrails
