package orchestrator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPrepareCommandLogFile(t *testing.T) {
	worktreePath := t.TempDir()

	stdoutLogPath, stderrLogPath, err := prepareCommandLogPaths(worktreePath, "deadbeef")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedStdout := filepath.Join(worktreePath, ".aynig", "logs", "deadbeef.stdout.log")
	if stdoutLogPath != expectedStdout {
		t.Fatalf("unexpected stdout log path: got %q want %q", stdoutLogPath, expectedStdout)
	}

	expectedStderr := filepath.Join(worktreePath, ".aynig", "logs", "deadbeef.stderr.log")
	if stderrLogPath != expectedStderr {
		t.Fatalf("unexpected stderr log path: got %q want %q", stderrLogPath, expectedStderr)
	}

	logsDir := filepath.Join(worktreePath, ".aynig", "logs")
	if info, err := os.Stat(logsDir); err != nil {
		t.Fatalf("expected logs directory to exist: %v", err)
	} else if !info.IsDir() {
		t.Fatalf("expected logs directory, got file")
	}
}
