#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
GO_DIR="${REPO_ROOT}/go"

VERSION="${VERSION:-}"
if [[ -z "${VERSION}" ]]; then
  if command -v git >/dev/null 2>&1; then
    VERSION="$(git -C "${REPO_ROOT}" describe --tags --always)"
  else
    VERSION="dev"
  fi
fi

OUT_DIR="${REPO_ROOT}/docs/public/releases/${VERSION}"
mkdir -p "${OUT_DIR}"

TARGETS=(
  "linux/amd64"
  "linux/arm64"
  "darwin/amd64"
  "darwin/arm64"
  "windows/amd64"
)

for target in "${TARGETS[@]}"; do
  IFS="/" read -r GOOS GOARCH <<< "${target}"
  EXT=""
  if [[ "${GOOS}" == "windows" ]]; then
    EXT=".exe"
  fi
  OUTPUT_NAME="aynig_${GOOS}_${GOARCH}${EXT}"
  echo "Building ${OUTPUT_NAME}"
  CGO_ENABLED=0 GOOS="${GOOS}" GOARCH="${GOARCH}" \
    go -C "${GO_DIR}" build -trimpath -o "${OUT_DIR}/${OUTPUT_NAME}" ./cmd/aynig
done

echo "Release artifacts written to ${OUT_DIR}"
