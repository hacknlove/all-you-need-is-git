#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
binary_path="${HOME}/.local/bin/aynig"
version="$(git -C "${repo_root}" describe --tags --always --dirty 2>/dev/null || echo "dev")"

echo "Building aynig (${version})..."
go -C "${repo_root}/go" build -ldflags "-X main.version=${version}" -o "${repo_root}/go/aynig" ./cmd/aynig

echo "Installing to ${binary_path}..."
install -m 0755 "${repo_root}/go/aynig" "${binary_path}"

echo "Installed version:"
"${binary_path}" version || true
