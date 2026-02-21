---
title: Kernel Contract
description: The minimum guarantees of the AYNIG runner.
---

The AYNIG kernel is a small, deterministic runner that:

- Interprets each commit as a state event
- Dispatches using the `aynig-state` trailer
- Executes commands and validates the new `HEAD`

It does not define workflows, interpret business semantics, or retry commands. All intelligence lives above the kernel.

For the full specification, see `.aynig/CONTRACT.md` in your repository (created by `aynig init`).
