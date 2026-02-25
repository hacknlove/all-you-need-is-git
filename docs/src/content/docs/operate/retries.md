---
title: Retries
description: Retry conventions that keep the kernel minimal.
---

AYNIG does not implement retries. Retrying is a **workflow policy**.

A simple convention:

- Use trailers like:
  - `aynig-attempt: 1`
  - `aynig-max-attempts: 3`
- If a command fails, create a new commit with the same `aynig-state` and increment the attempt.
- Stop retrying when attempt reaches max.

These are conventions only; the kernel does not interpret them.
