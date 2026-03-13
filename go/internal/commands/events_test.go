package commands

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestEventsReadsDwpState(t *testing.T) {
	tempDir := t.TempDir()
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
	if output, err := exec.Command("git", "init", ".").CombinedOutput(); err != nil {
		t.Fatalf("git init failed: %v (%s)", err, string(output))
	}
	if output, err := exec.Command("git", "config", "user.email", "test@example.com").CombinedOutput(); err != nil {
		t.Fatalf("git config failed: %v (%s)", err, string(output))
	}
	if output, err := exec.Command("git", "config", "user.name", "Test User").CombinedOutput(); err != nil {
		t.Fatalf("git config failed: %v (%s)", err, string(output))
	}
	if err := os.WriteFile(filepath.Join(tempDir, "file.txt"), []byte("ok"), 0o644); err != nil {
		t.Fatalf("write failed: %v", err)
	}
	if output, err := exec.Command("git", "add", ".").CombinedOutput(); err != nil {
		t.Fatalf("git add failed: %v (%s)", err, string(output))
	}
	if output, err := exec.Command("git", "commit", "-m", "chore: test", "-m", "body\n\ndwp-state: build").CombinedOutput(); err != nil {
		t.Fatalf("git commit failed: %v (%s)", err, string(output))
	}

	stdout := os.Stdout
	pipeR, pipeW, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe failed: %v", err)
	}
	os.Stdout = pipeW
	if err := Events(EventsOptions{History: false, Limit: 1, JSON: true}); err != nil {
		_ = pipeW.Close()
		os.Stdout = stdout
		t.Fatalf("events failed: %v", err)
	}
	_ = pipeW.Close()
	os.Stdout = stdout
	output, err := io.ReadAll(pipeR)
	if err != nil {
		t.Fatalf("read output failed: %v", err)
	}

	var events []Event
	if err := json.Unmarshal(output, &events); err != nil {
		t.Fatalf("json decode failed: %v", err)
	}
	if len(events) != 1 || events[0].State == nil || *events[0].State != "build" {
		t.Fatalf("unexpected events: %+v", events)
	}
}
