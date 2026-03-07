---
title: History Hygiene
description: Keep long-running workflows auditable without making branches unusable.
---

Because each step is a commit, long-running workflows can produce many commits.

Common approaches:

- Run workflows on dedicated branches.
- Periodically checkpoint and squash (human decision).
- Keep `working`/lease noise out of main development branches.

AYNIG stays HEAD-based; history scanning is optional and belongs to higher layers.
