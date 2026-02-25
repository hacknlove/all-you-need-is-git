# Contributing

Thanks for your interest in contributing to **AYNIG**.

## Scope

This repository contains:

- **Go** implementation: `go/`
- **JS** implementation: `js/`
- **Docs site** (Astro + Starlight): `docs/`
- Optional workflow pack: `ops-workflow-pack/`

## Development setup

### Prereqs

- Git
- Node.js (for `js/` and `docs/`)
- Go 1.22+ (for `go/`)

### Install dependencies

```bash
# JS runner/CLI
cd js
npm ci

# Docs
cd ../docs
npm ci
```

### Run tests

```bash
# JS
cd js
npm test

# Go
cd ../go
go test ./...

# Docs build
cd ../docs
npm run build
```

### Run locally

```bash
# JS CLI (watch)
cd js
npm run dev

# Docs (dev server)
cd ../docs
npm run dev
```

## Documentation contributions

User-facing docs live in `docs/src/content/docs/`.

- If you add a new page, also add it to the sidebar in `docs/astro.config.mjs`.
- Keep docs focused on **how to use AYNIG**, not repo-internal details.

## PR guidelines

- Prefer small PRs with clear scope.
- Add/adjust tests when behavior changes.
- Keep the kernel philosophy intact (see `CONTRACT.md`).

## License

By contributing, you agree that your contributions will be licensed under the repository license (see `LICENSE`).
