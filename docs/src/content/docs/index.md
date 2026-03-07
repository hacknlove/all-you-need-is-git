---
title: Overview
description: Use AYNIG to run workflows from Git commits.
hero:
  title: AYNIG
  tagline: All you need is Git
  image:
    file: ./hero.jpg
    alt: AYNIG hero illustration
---

:::caution[WORK IN PROGRESS]
This project is under active development. APIs, commands, and documentation may change without notice.
:::

## What AYNIG does, in a nutshell.

1. Reads the state from the latest commit
2. Runs the command or agent that handles that state
3. Provides it with the necessary context and tools to do its job and update the state

## What's the point of AYNIG?

AYNIG is a coordination protocol and workflow system for software development teams with one or more humans and one or more agents, so everyone knows whose turn it is and what to do next.

### workflow example:

1. An architect prompts an agent to generate a plan, and they iterate together until the architect is satisfied.
2. A developer prompts an agent to implement one step of the plan, and they iterate using other agents for code review and testing until the developer is satisfied. (Sometimes they send the plan back to the architect for revision, which is also part of the workflow.)
3. A different developer reviews the commit with the help of more agents, sending it back to the original developer for fixes if needed, until the reviewer is satisfied.
4. A QA engineer performs further testing with the help of different agents, sending it back to the developer for fixes if needed, until the QA engineer is satisfied.

### How would you implement it?

You could have humans signal each other with messages in a direct messaging platform or by moving cards on a kanban board (or both), but this only captures the human side of the workflow.
Then humans also need to coordinate with agents through a different interface, likely manual prompting and copying outputs between tools. Agents might read the messages or the kanban board, but everything is flaky, error-prone, hard to audit, and heavily reliant on human discipline and constant supervision. That makes it hard to automate or scale.

### AYNIG's approach

AYNIG takes a different approach: it uses Git commits as the single source of truth for the workflow. Humans and agents interact through Git, using commit trailers to signal whose turn it is and what to do next. AYNIG runners read the latest commit, run the appropriate command, and check the new commit, ensuring that the workflow progresses smoothly and reliably.
Humans can still interact with each other through their preferred channels, and manually with agents when needed, but AYNIG provides them with a clear, auditable, and robust protocol to handle the main parts of the workflow.


## Why use AYNIG

- No external control system: no SaaS, no webhooks, no hidden state
- Fully auditable workflows: every step is a commit you can inspect
- Works with your existing Git practices: branching, history, and collaboration
- Tool-agnostic and composable: bring your own scripts

## Where to start

- New to AYNIG? Start with [Quick Start](/getting-started/quick-start/)
- Installing the CLI? See [Install CLI](/installation/cli/)
- Writing commands? See [Authoring Commands](/commands/authoring/)

## Helpful guides

- [Hello World Workflow](/guides/hello-world/)
- [FAQ / Troubleshooting](/guides/faq/)
- [Security & Trust Model](/guides/security/)
- [Ops Workflow Pack](/guides/ops-workflow-pack/)
- [Glossary](/guides/glossary/)
