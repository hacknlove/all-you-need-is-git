package commands

import (
	"os"
	"os/exec"
	"testing"
)

func TestSetStateRequiresState(t *testing.T) {
	err := SetState(SetStateOptions{})
	if err == nil || err.Error() != "Missing required flag: --dwp-state" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSetStateRejectsWorking(t *testing.T) {
	err := SetState(SetStateOptions{State: "working"})
	if err == nil || err.Error() != "Invalid dwp-state: working (use aynig set-working)" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSetStateKeepsExistingDwpTrailersWhenRequested(t *testing.T) {
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
	if output, err := exec.Command("git", "commit", "--allow-empty", "-m", "seed", "-m", "body\n\ndwp-state: triage\ndwp-note: first\ndwp-note: second\ndwp-origin-state: queued").CombinedOutput(); err != nil {
		t.Fatalf("git commit failed: %v (%s)", err, string(output))
	}

	if err := SetState(SetStateOptions{State: "review", Prompt: "Mantener trailers", KeepTrailers: true}); err != nil {
		t.Fatalf("SetState failed: %v", err)
	}

	output, err := exec.Command("git", "log", "-1", "--format=%B").CombinedOutput()
	if err != nil {
		t.Fatalf("git log failed: %v (%s)", err, string(output))
	}
	want := "chore: set review\n\nMantener trailers\n\ndwp-state: review\ndwp-note: first\ndwp-note: second\ndwp-origin-state: queued\n\n"
	if string(output) != want {
		t.Fatalf("unexpected commit message:\n%s", string(output))
	}
}
