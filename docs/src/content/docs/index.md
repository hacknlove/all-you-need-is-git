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

AYNIG uses Git as the source of truth for a workflow runner.
You write commands that the runner selects and executes them based on the trailers in the latest commit.

If you prefer small, composable tools and want to shape your own process instead of adopting a monolithic platform, AYNIG follows that philosophy.

It’s aimed at teams and individuals who favor a Unix-style model: simple primitives, composition, a few sensible conventions, and full control over their workflow.

## What AYNIG is

- A Git-native execution kernel for agentic workflows
- A runner that reads `HEAD`, dispatches to a command, and validates the new `HEAD`
- A system that uses commits as the control plane

## Why use AYNIG

- Zero external control plane: no SaaS, no webhooks, no hidden state
- Fully auditable workflows: every step is a commit you can inspect
- Works with your existing Git practices: branching, history, and collaboration
- Agent-agnostic and composable: bring your own tools and scripts

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
