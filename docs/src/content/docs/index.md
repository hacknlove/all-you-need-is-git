---
title: Overview
description: Use AYNIG to run workflows from Git commits.
hero:
  title: AYNIG
  tagline: Agentic Yet Native in Git
  image:
    file: ./hero.jpg
    alt: AYNIG hero illustration
---

AYNIG lets you drive workflows from Git commits. You define commands in `.aynig/command/`, and the runner executes them based on the `aynig-state` trailer in the latest commit.

## What AYNIG is

- A Git-native execution kernel for agentic workflows
- A runner that reads `HEAD`, dispatches to a command, and validates the new `HEAD`
- A system that uses commits as the control plane

## Where to start

- New to AYNIG? Start with [Quick Start](/getting-started/quick-start/)
- Installing the CLI? See [Install CLI](/installation/cli/)
- Writing commands? See [Authoring Commands](/commands/authoring/)
