package commands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveCommandPathPrefersRole(t *testing.T) {
	repoRoot := t.TempDir()
	rolePath := filepath.Join(repoRoot, ".dwp", "roles", "ops", "command", "build")
	if err := os.MkdirAll(filepath.Dir(rolePath), 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(rolePath, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	status, path := resolveCommandPath(repoRoot, "ops", "build")
	if status != "exists" {
		t.Fatalf("expected exists, got %q", status)
	}
	if path != rolePath {
		t.Fatalf("unexpected path: %q", path)
	}
}

func TestResolveCommandPathFallsBackToBase(t *testing.T) {
	repoRoot := t.TempDir()
	basePath := filepath.Join(repoRoot, ".dwp", "command", "build")
	if err := os.MkdirAll(filepath.Dir(basePath), 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(basePath, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	status, path := resolveCommandPath(repoRoot, "ops", "build")
	if status != "exists" {
		t.Fatalf("expected exists, got %q", status)
	}
	if path != basePath {
		t.Fatalf("unexpected path: %q", path)
	}
}
