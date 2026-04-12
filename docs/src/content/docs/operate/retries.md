---
title: Retries
description: Retry conventions that keep AYNIG minimal.
---

AYNIG does not implement retries. Retrying is a **workflow policy**.

A simple convention:

- Use trailers like:
  - `dwp-attempt: 1`
  - `dwp-max-attempts: 3`
- If a command wants a retry, emit `SET_STATE {...}` with the same `dwp-state` and increment the attempt trailer.
- Stop retrying when attempt reaches max.

These are conventions only; AYNIG does not interpret them.
