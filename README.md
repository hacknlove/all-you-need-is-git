# AYNIG (developer README)

## What is AYNIG?
AYNIG (All You Need Is Git) is a coordination protocol for software development teams with one or more humans and one or more agents, so everyone knows whose turn it is and what to do next.

Imagine the following workflow:
1. An architect prompts an agent to generate a plan, and they iterate together until the architect is satisfied.
2. A developer prompts an agent to implement one step of the plan, and they iterate using other agents for code review and testing until the developer is satisfied. (Sometimes they send the plan back to the architect for revision, which is also part of the workflow.)
3. A different developer reviews the commit with the help of more agents, sending it back to the original developer for fixes if needed, until the reviewer is satisfied.
4. A QA engineer performs further testing with the help of different agents, sending it back to the developer for fixes if needed, until the QA engineer is satisfied.

How would you implement it?

You could have humans signal each other with messages in a direct messaging platform or by moving cards on a kanban board (or both), but this only captures the human side of the workflow.

Then humans also need to coordinate with agents through a different interface, likely manual prompting and copying outputs between tools. Agents might read the messages or the kanban board, but everything is flaky, error-prone, hard to audit, and heavily reliant on human discipline and constant supervision. That makes it hard to automate or scale.

AYNIG takes a different approach: it uses Git commits as the single source of truth for the workflow. Humans and agents interact through Git, using commit trailers to signal whose turn it is and what to do next. AYNIG runners read the latest commit, dispatch the appropriate command, and validate the new commit, ensuring that the workflow progresses smoothly and reliably.

Humans can still interact with each other through their preferred channels, and manually with agents when needed, but AYNIG provides them with a clear, auditable, and robust protocol to handle the main parts of the workflow.

> **WORK IN PROGRESS:** This project is under active development. APIs, commands, and documentation may change without notice.

This repository contains the AYNIG implementations and the documentation site.

- User docs: https://aynig.org
- Kernel contract: `CONTRACT.md`

## Repository layout

- `go/` — Go implementation
- `js/` — Node.js implementation (CLI entrypoint)
- `docs/` — Documentation site (Astro + Starlight)
- `slides/` — Slidev presentation workspace (index + scoped decks, Cloudflare Pages-ready build output)
- `ops-workflow-pack/` — optional workflow pack

## Development

### Go

```bash
cd go
go test ./...
```

### JS

```bash
cd js
npm ci
npm test
```

Run the JS CLI locally:

```bash
cd js
npm run dev
```

### Docs

```bash
cd docs
npm ci
npm run dev
```

Build docs:

```bash
cd docs
npm run build
```

### Slides

```bash
cd slides
npm ci
npm run dev
```

Build slides (index + scoped decks for deployment):

```bash
cd slides
npm run build
```

## Release

This repo uses GoReleaser (`.goreleaser.yaml`).

See GitHub Releases and workflows for the current release process.

## Contributing

See `CONTRIBUTING.md`.

## License

Apache-2.0 — see `LICENSE`.
