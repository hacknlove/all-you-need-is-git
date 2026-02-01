# AYNIG

**Agentic Yet Native in Git**  
*Because all you need is Git.*

AYNIG is a Git-native orchestration tool for agentic workflows.  
No SaaS. No dependencies. No control plane outside your repo.

If that sounds wrong to you, this tool is probably not for you.

![hero](/hero.jpg)

---

## What is this?

**AYNIG** is a Git-native runner for agentic workflows.

It treats **Git commits as the control plane**:

- commit messages carry **state and intent**
- agents *(and humans)* **take turns** by creating and responding to commits

There is no external state, no webhooks, no SaaS backend.  
Everything is transparent, inspectable, and agnostic to the agents or tools you use.

---

## Why?

Because **simple, transparent, open, and composable** beats  
**complex, opaque, closed, and “complete”**.

Other agentic frameworks lock you into their platform and, even worse, into **their way of thinking about state and orchestration**.

You are forced to:

- shape your workflows around their abstractions,
- accept dark magic and hidden state,
- and depend on a single point of failure.

But Git already gives us:

- ordering,
- history,
- branching,
- collaboration,
- and auditability.

Humans have been using Git to communicate and coordinate work for decades, we just need a way to let agents join the party.

With AYNIG, you:

- fetch refs,
- inspect commit messages for `aynig:` trailers,
- run the corresponding agents,
- and let agents respond by creating new commits using the same syntax.

The result is an **auditable, replicable, distributed, and fully customizable workflow**, out of the box.

You don't need webhooks, APIs, and keys, and certainly you don't need to be tied to a certain project management system, control version platform, LLM provider, or anything else. 

---

## How does it work?

AYNIG scans branches and inspects the **latest commit**.

If it finds a commit with an `aynig:` trailer, like this:

```
chore: whatever title for humans

Here goes the prompt to be processed by the agent.

aynig: some-state
foo: bar
baz: qux
```

AYNIG will:

1. Detect the state (`some-state`)
2. Execute the corresponding script:

```
.aynig/some-state
```

3. Expose the commit contents to the script via environment variables:

* `AYNIG_BODY`
* `AYNIG_FOO`
* `AYNIG_BAZ`
* `AYNIG_WORKTREE_PATH`
* `AYNIG_COMMIT_SHA`
* …and others

The agent is then responsible for:

* doing the work, and
* creating a **new commit** that
* advances the workflow by setting a new `aynig:` state trailer.

After that, the agent may:

* call `aynig run` to continue immediately,
* or rely on the human calling it, on cron, or any other mechanism.

AYNIG does not care how you keep it running.

---

## Worktrees and isolation

AYNIG runs each agent execution inside a **dedicated Git worktree**.

This is not an implementation detail — it is a core design choice.

For every actionable commit, AYNIG:

* creates a worktree checked out on the target branch,
* runs the corresponding workflow script inside that worktree,
* and removes it once execution finishes.

This guarantees that:

* agents always operate on a **clean and reproducible snapshot**,
* the main working directory is never modified,
* concurrent executions do not interfere with each other,
* and crashes or partial runs do not leave the repository in an inconsistent state.

The worktree path is exposed to scripts via:

* `AYNIG_WORKTREE_PATH`

Worktree creation, lifecycle, and cleanup are handled entirely by AYNIG.
Agents should **not** create or manage worktrees themselves.

The default worktree location is:

```
.worktrees/
```

and it is added to `.gitignore` automatically.

---

### Humans and worktrees

Worktrees are not just an execution detail — they are also the **human interface**.

If a human wants to inspect, debug, or modify what an agent is doing, they can simply enter the corresponding worktree:


That worktree contains the exact files, state, and history the agent is working with.

Humans and agents operate symmetrically:

* both see the same files,
* both can edit,
* both can commit.

AYNIG does not introduce a separate UI or abstraction layer.
The worktree *is* the interface.

If a human commits from the worktree, that commit becomes part of the workflow and is handled exactly like any agent-produced commit.

---

## Install

```bash
curl -fsSL https://aynig.com/install.sh | bash
```

---

## Setup

Inside the repository where you want to use AYNIG, run:

```bash
aynig init
```

This will:

* create `.aynig/` and `.worktrees/`,
* and add `.worktrees/` to your `.gitignore`.

---

## Workflows

You define workflows by creating executable scripts inside `.aynig/`.

Each script corresponds to a state name and is invoked when that state appears in an `aynig:` commit trailer.

You can also install pre-made workflows from other repositories:

```bash
aynig install <git-repo-url> <branch-or-tag> [subfolder]
```

AYNIG will copy the `.aynig/` scripts from the source repository into your own, allowing workflows to be **shared, versioned, and forked** like any other code.

It is recommended to create a COMMANDS.md file to document the available workflows in your repository, in case any agent needs to decide the next step based on the available options.

for instance:
```markdown
build: Build the project
review: Review the changes
test: Run tests
deploy: Deploy to production
cleanup: Clean the worktree
```
---

### Script example

```bash
#!/bin/bash

cd "$AYNIG_WORKTREE_PATH"

PROMPT=$(cat <<EOF
Build this spec:
$AYNIG_BODY
EOF
)

claude -p "$PROMPT" --session-id "$AYNIG_COMMIT_SHA"
```

---

## CLI syntax

```bash
Usage: aynig [options] [command]

Global Options:
  -v, --version        output the version number
  -h, --help           display help for command

Commands:
  init                 Initialize AYNIG in the current repository

  run [options]        Run AYNIG for the current repository
    Options:
      -w, --worktree <path>   Specify custom worktree directory (default: .worktrees)
      --use-remote <name>     Specify which remote to use (default: origin)

  install <repo> [ref] [subfolder]
                       Install AYNIG workflows from another repository

  help [command]       display help for command
```