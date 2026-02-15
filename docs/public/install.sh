#!/usr/bin/env sh
set -e

BASE_URL="https://aynig.org/releases"
VERSION="${VERSION:-}"

download() {
  url="$1"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url"
    return
  fi
  if command -v wget >/dev/null 2>&1; then
    wget -qO- "$url"
    return
  fi
  echo "Error: curl or wget is required." >&2
  exit 1
}

if [ -z "$VERSION" ]; then
  VERSION="$(download "$BASE_URL/latest.txt" 2>/dev/null || true)"
fi

if [ -z "$VERSION" ]; then
  echo "Error: VERSION is required. Example: VERSION=v0.1.0 $0" >&2
  exit 1
fi

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$OS" in
  linux|darwin) : ;;
  msys*|mingw*|cygwin*) OS="windows" ;;
  *) echo "Unsupported OS: $OS" >&2; exit 1 ;;
esac

case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

ARCHIVE="aynig_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="$BASE_URL/$VERSION/$ARCHIVE"

tmpdir="$(mktemp -d)"
cleanup() { rm -rf "$tmpdir"; }
trap cleanup EXIT

echo "Downloading $URL"
download "$URL" > "$tmpdir/$ARCHIVE"

tar -xzf "$tmpdir/$ARCHIVE" -C "$tmpdir"

BIN_NAME="aynig"
if [ "$OS" = "windows" ]; then
  BIN_NAME="aynig.exe"
fi

BIN_PATH="$(find "$tmpdir" -type f -name "$BIN_NAME" 2>/dev/null | head -n 1)"
if [ -z "$BIN_PATH" ]; then
  echo "Error: aynig binary not found in archive." >&2
  exit 1
fi

INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
  INSTALL_DIR="$HOME/.local/bin"
fi

mkdir -p "$INSTALL_DIR"
cp "$BIN_PATH" "$INSTALL_DIR/aynig"
chmod +x "$INSTALL_DIR/aynig"

echo "Installed aynig to $INSTALL_DIR/aynig"
echo "Make sure $INSTALL_DIR is on your PATH."
