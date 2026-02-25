# AYNIG (developer README)

This repository contains the AYNIG implementations and its documentation site.

- User docs: https://aynig.org
- Kernel contract: `CONTRACT.md`

## Repository layout

- `go/` — Go implementation
- `js/` — Node.js implementation (CLI entrypoint)
- `docs/` — Documentation site (Astro + Starlight)
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

## Release

This repo uses GoReleaser (`.goreleaser.yaml`).

See GitHub Releases and workflows for the current release process.

## Contributing

See `CONTRIBUTING.md`.

## License

Apache-2.0 — see `LICENSE`.
