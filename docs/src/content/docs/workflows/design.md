---
title: Workflow Design
description: Patterns for building reliable state machines.
---

## Think in states

Design workflows as explicit states, each mapped to a command in `.aynig/command/<state>`.

## Make transitions explicit

Each command should end by creating a new commit with the next `aynig-state`. This keeps the workflow auditable and deterministic.

## Favor idempotency

Commands may be re-run after failures or handoffs. Make operations safe to run multiple times when possible.

## Human and agent handoffs

Humans can enter the worktree, inspect changes, and commit the next state. Treat humans and agents symmetrically.
