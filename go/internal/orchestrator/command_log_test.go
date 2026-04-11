package orchestrator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPrepareCommandLogFile(t *testing.T) {
	worktreePath := t.TempDir()

	file, logPath, err := prepareCommandLogFile(worktreePath, "deadbeef")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}

	expected := filepath.Join(worktreePath, ".aynig", "logs", "deadbeef.log")
	if logPath != expected {
		t.Fatalf("unexpected log path: got %q want %q", logPath, expected)
	}
	if _, err := os.Stat(logPath); err != nil {
		t.Fatalf("expected log file to exist: %v", err)
	}
}
