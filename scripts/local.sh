#!/bin/bash
set -e

log() {
  printf '[local.sh] %s\n' "$1"
}

ROOT_DIR="$(dirname "$0")/.."
TARGET_DIR="$HOME/.local/bin"
TARGET_BIN="$TARGET_DIR/aynig"
BUILD_TIMESTAMP_UNIX="$(date +%s)"

log "Using repo root: $ROOT_DIR"
cd "$(dirname "$0")/.."
log "Ensuring target directory exists: $TARGET_DIR"
mkdir -p "$TARGET_DIR"

log "Building Go CLI into: $TARGET_BIN"
log "Embedding build timestamp: $BUILD_TIMESTAMP_UNIX"
cd go && go build -ldflags "-X main.buildTimestampUnix=$BUILD_TIMESTAMP_UNIX" -o "$TARGET_BIN" ./cmd/aynig

if [ -x "$TARGET_BIN" ]; then
  log "Build completed successfully: $TARGET_BIN"
else
  log "Build finished but binary is missing or not executable: $TARGET_BIN"
  exit 1
fi

log "Checking built binary version"
"$TARGET_BIN" version
log "Binary verification completed successfully"
