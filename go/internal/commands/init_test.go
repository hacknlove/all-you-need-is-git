package commands

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestInitCreatesCommandDirAndClean(t *testing.T) {
	tempDir := t.TempDir()

	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init failed: %v (%s)", err, string(output))
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(cwd)
	})
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	if err := Init(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	commandDir := filepath.Join(tempDir, ".aynig", "command")
	cleanCommand := filepath.Join(commandDir, "clean")
	commandsDoc := filepath.Join(tempDir, ".aynig", "COMMANDS.md")

	if _, err := os.Stat(commandDir); err != nil {
		t.Fatalf("command dir missing: %v", err)
	}
	if _, err := os.Stat(cleanCommand); err != nil {
		t.Fatalf("clean command missing: %v", err)
	}
	if _, err := os.Stat(commandsDoc); err != nil {
		t.Fatalf("COMMANDS.md missing: %v", err)
	}
}
