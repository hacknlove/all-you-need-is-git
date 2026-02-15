# Releasing AYNIG

This repo uses GoReleaser to build and publish cross‑platform binaries into the docs public folder.

## Prerequisites

- Go 1.22+
- GoReleaser installed (`goreleaser --version`)

Install GoReleaser:

```bash
go install github.com/goreleaser/goreleaser@latest
```

## Release workflow

1) Choose a version and create a tag:

```bash
git tag vX.Y.Z
```

2) Run GoReleaser in snapshot mode locally (optional sanity check):

```bash
goreleaser release --snapshot --clean
```

3) Run the real release (writes artifacts to docs):

```bash
goreleaser release --clean
```

Artifacts will be written to:

```
docs/public/releases/<version>/
```

4) Update the latest pointer:

```bash
echo "vX.Y.Z" > docs/public/releases/latest.txt
```

5) Commit and push the artifacts and `latest.txt`.

## Install script

The docs site serves `docs/public/install.sh`. It downloads the proper archive from:

```
https://aynig.org/releases/<version>/
```

Ensure `latest.txt` points at your newest release so `curl https://aynig.org/install.sh | bash` uses the right version.
