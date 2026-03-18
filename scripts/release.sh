#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

VERSION="${1:-}"
if [[ -z "${VERSION}" ]]; then
  echo "Error: VERSION is required. Usage: $0 v0.0.1" >&2
  exit 1
fi

if [[ ! "${VERSION}" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Error: VERSION must match vX.Y.Z format. Got: ${VERSION}" >&2
  exit 1
fi

cd "${REPO_ROOT}"

echo "=== Releasing ${VERSION} ==="

if gh release view "${VERSION}" >/dev/null 2>&1; then
  echo "Release ${VERSION} exists. Deleting..."
  gh release delete "${VERSION}" --yes
  git tag -d "${VERSION}" 2>/dev/null || true
  git push origin :refs/tags/"${VERSION}" 2>/dev/null || true
fi

if git rev-parse "${VERSION}" >/dev/null 2>&1; then
  echo "Tag ${VERSION} exists locally. Deleting..."
  git tag -d "${VERSION}"
  git push origin :refs/tags/"${VERSION}" 2>/dev/null || true
fi

echo "Creating tag ${VERSION}..."
git tag "${VERSION}"

echo "Running goreleaser..."
GITHUB_TOKEN=$(gh auth token) goreleaser release --clean --skip=validate

RELEASES_DIR="docs/public/releases/${VERSION}"
echo "Copying artifacts to ${RELEASES_DIR}..."
rm -rf "${RELEASES_DIR}"
mkdir -p "${RELEASES_DIR}"
cp -r .dist/"${VERSION}"/* "${RELEASES_DIR}/"

echo "Committing and pushing artifacts..."
git add "docs/public/releases/${VERSION}"
git commit -m "chore: add v${VERSION#v} release artifacts" || true
git push origin master

echo "=== Release ${VERSION} complete ==="
