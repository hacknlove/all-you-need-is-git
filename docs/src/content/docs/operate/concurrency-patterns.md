---
title: Concurrency Patterns
description: Fan-out/fan-in and multi-branch coordination patterns.
---

AYNIG provides mutual exclusion per branch. Concurrency patterns are implemented at the workflow level.

Examples (conceptual):

- Fan-out: create N branches each with its own `aynig-state`
- Fan-in: a coordinator branch waits for results and merges or references outputs

Use correlation trailers (e.g. `aynig-correlation-id`) to relate work across branches.
