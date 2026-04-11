package commands

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestLeaseStatusForStateNonWorkingIsNA(t *testing.T) {
	if got := leaseStatusForState("", "", ""); got != "n/a" {
		t.Fatalf("expected n/a for empty state, got %q", got)
	}
}

func TestResolveCommandPathPrefersRole(t *testing.T) {
	repoRoot := t.TempDir()
	rolePath := filepath.Join(repoRoot, ".aynig", "roles", "ops", "command", "build")
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
	basePath := filepath.Join(repoRoot, ".aynig", "command", "build")
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

func TestStatusReadsSpecificBranchWithoutCheckout(t *testing.T) {
	repoDir := newStatusTestRepo(t)
	writeExecutable(t, filepath.Join(repoDir, ".aynig", "command", "build"))
	commitEmpty(t, repoDir, "seed", "body")
	createBranchCommit(t, repoDir, "1-bootstrap", "feat: bootstrap", "body\n\ndwp-state: build\ndwp-run-id: run-123")

	output := captureStatusOutput(t, repoDir, StatusOptions{Branch: "1-bootstrap"})

	if !strings.Contains(output, "branch: 1-bootstrap\n") {
		t.Fatalf("expected branch in output, got %q", output)
	}
	if !strings.Contains(output, "dwp-state: build\n") {
		t.Fatalf("expected branch dwp-state in output, got %q", output)
	}
	if !strings.Contains(output, "dwp-run-id: run-123\n") {
		t.Fatalf("expected branch run id in output, got %q", output)
	}
	if !strings.Contains(output, "command: exists\n") {
		t.Fatalf("expected command resolution in output, got %q", output)
	}
	current := currentBranch(t, repoDir)
	if current != "main" {
		t.Fatalf("status should not checkout the target branch, current branch is %q", current)
	}
}

func TestStatusReadsBranchPattern(t *testing.T) {
	repoDir := newStatusTestRepo(t)
	writeExecutable(t, filepath.Join(repoDir, ".aynig", "command", "build"))
	writeExecutable(t, filepath.Join(repoDir, ".aynig", "command", "review"))
	commitEmpty(t, repoDir, "seed", "body")
	createBranchCommit(t, repoDir, "1-bootstrap", "feat: bootstrap", "body\n\ndwp-state: build")
	createBranchCommit(t, repoDir, "1-probing", "feat: probing", "body\n\ndwp-state: review")
	createBranchCommit(t, repoDir, "other", "feat: other", "body\n\ndwp-state: build")

	output := captureStatusOutput(t, repoDir, StatusOptions{BranchPattern: "1-*"})

	if !strings.Contains(output, "branch: 1-bootstrap\n") || !strings.Contains(output, "branch: 1-probing\n") {
		t.Fatalf("expected matching branches in output, got %q", output)
	}
	if strings.Contains(output, "branch: other\n") {
		t.Fatalf("unexpected non-matching branch in output, got %q", output)
	}
	if strings.Count(output, "branch: ") != 2 {
		t.Fatalf("expected two branch blocks, got %q", output)
	}
	if !strings.Contains(output, "\n\nbranch: 1-probing\n") {
		t.Fatalf("expected blank line between branch blocks, got %q", output)
	}
	current := currentBranch(t, repoDir)
	if current != "main" {
		t.Fatalf("status should not checkout matching branches, current branch is %q", current)
	}
}

func TestStatusPrefersLocalBranchRefOverMatchingTag(t *testing.T) {
	repoDir := newStatusTestRepo(t)
	writeExecutable(t, filepath.Join(repoDir, ".aynig", "command", "build"))
	writeExecutable(t, filepath.Join(repoDir, ".aynig", "command", "review"))
	commitEmpty(t, repoDir, "seed", "body")
	createBranchCommit(t, repoDir, "foo", "feat: branch foo", "body\n\ndwp-state: build")
	runGit(t, repoDir, "tag", "foo", "main")

	output := captureStatusOutput(t, repoDir, StatusOptions{Branch: "foo"})

	if !strings.Contains(output, "branch: foo\n") {
		t.Fatalf("expected branch in output, got %q", output)
	}
	if !strings.Contains(output, "dwp-state: build\n") {
		t.Fatalf("expected branch state to win over tag target, got %q", output)
	}
	if strings.Contains(output, "dwp-state: review\n") {
		t.Fatalf("expected status to ignore matching tag ref, got %q", output)
	}
	if !strings.Contains(output, "command: exists\n") {
		t.Fatalf("expected branch command resolution, got %q", output)
	}
}

func newStatusTestRepo(t *testing.T) string {
	t.Helper()
	repoDir := t.TempDir()
	runGit(t, repoDir, "init", ".")
	runGit(t, repoDir, "config", "user.email", "test@example.com")
	runGit(t, repoDir, "config", "user.name", "Test User")
	runGit(t, repoDir, "checkout", "-b", "main")
	return repoDir
}

func createBranchCommit(t *testing.T, repoDir string, branch string, subject string, body string) {
	t.Helper()
	runGit(t, repoDir, "checkout", "-b", branch)
	commitEmpty(t, repoDir, subject, body)
	runGit(t, repoDir, "checkout", "main")
}

func commitEmpty(t *testing.T, repoDir string, subject string, body string) {
	t.Helper()
	runGit(t, repoDir, "commit", "--allow-empty", "-m", subject, "-m", body)
}

func writeExecutable(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("write failed: %v", err)
	}
}

func captureStatusOutput(t *testing.T, repoDir string, options StatusOptions) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(cwd)
	})
	if err := os.Chdir(repoDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}

	stdout := os.Stdout
	pipeR, pipeW, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe failed: %v", err)
	}
	os.Stdout = pipeW
	if err := Status(options); err != nil {
		_ = pipeW.Close()
		os.Stdout = stdout
		t.Fatalf("status failed: %v", err)
	}
	_ = pipeW.Close()
	os.Stdout = stdout
	output, err := io.ReadAll(pipeR)
	if err != nil {
		t.Fatalf("read output failed: %v", err)
	}
	return string(output)
}

func currentBranch(t *testing.T, repoDir string) string {
	t.Helper()
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = repoDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git branch --show-current failed: %v (%s)", err, string(output))
	}
	return strings.TrimSpace(string(output))
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v (%s)", strings.Join(args, " "), err, string(output))
	}
}
