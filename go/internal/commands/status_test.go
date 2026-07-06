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
	repoDir := newStatusTestRepo(t)
	writeRepoFile(t, repoDir, ".aynig/roles/ops/command/build", "#!/bin/sh\n", 0o755)
	writeRepoFile(t, repoDir, ".aynig/command/build", "#!/bin/sh\n", 0o755)
	commitTree(t, repoDir, "add commands")

	status, path := resolveCommandPath(repoDir, "main", "ops", "build")
	if status != "exists" {
		t.Fatalf("expected exists, got %q", status)
	}
	if path != "main:.aynig/roles/ops/command/build" {
		t.Fatalf("unexpected path: %q", path)
	}
}

func TestResolveCommandPathFallsBackToBase(t *testing.T) {
	repoDir := newStatusTestRepo(t)
	writeRepoFile(t, repoDir, ".aynig/command/build", "#!/bin/sh\n", 0o755)
	commitTree(t, repoDir, "add commands")

	status, path := resolveCommandPath(repoDir, "main", "ops", "build")
	if status != "exists" {
		t.Fatalf("expected exists, got %q", status)
	}
	if path != "main:.aynig/command/build" {
		t.Fatalf("unexpected path: %q", path)
	}
}

func TestResolveCommandPathNonExecutableIsMissing(t *testing.T) {
	repoDir := newStatusTestRepo(t)
	writeRepoFile(t, repoDir, ".aynig/command/build", "#!/bin/sh\n", 0o644)
	commitTree(t, repoDir, "add non-executable command")

	status, path := resolveCommandPath(repoDir, "main", "", "build")
	if status != "missing" {
		t.Fatalf("expected missing, got %q", status)
	}
	if path != "main:.aynig/command/build" {
		t.Fatalf("unexpected path: %q", path)
	}
}

func TestStatusReadsSpecificBranchWithoutCheckout(t *testing.T) {
	repoDir := newStatusTestRepo(t)
	writeRepoFile(t, repoDir, ".aynig/command/build", "#!/bin/sh\n", 0o755)
	commitTree(t, repoDir, "seed")
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
	writeRepoFile(t, repoDir, ".aynig/command/build", "#!/bin/sh\n", 0o755)
	writeRepoFile(t, repoDir, ".aynig/command/review", "#!/bin/sh\n", 0o755)
	commitTree(t, repoDir, "seed")
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
	writeRepoFile(t, repoDir, ".aynig/command/build", "#!/bin/sh\n", 0o755)
	writeRepoFile(t, repoDir, ".aynig/command/review", "#!/bin/sh\n", 0o755)
	commitTree(t, repoDir, "seed")
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

func TestStatusIgnoresCommandsFromEventBranchByDefault(t *testing.T) {
	repoDir := newStatusTestRepo(t)
	commitEmpty(t, repoDir, "seed", "body")
	runGit(t, repoDir, "checkout", "-b", "1-evil")
	writeRepoFile(t, repoDir, ".aynig/command/build", "#!/bin/sh\n", 0o755)
	runGit(t, repoDir, "add", "-A")
	runGit(t, repoDir, "commit", "-m", "feat: evil", "-m", "body\n\ndwp-state: build")
	runGit(t, repoDir, "checkout", "main")

	output := captureStatusOutput(t, repoDir, StatusOptions{Branch: "1-evil"})

	if !strings.Contains(output, "command: missing\n") {
		t.Fatalf("expected command missing on the default commands ref, got %q", output)
	}
	if !strings.Contains(output, "command-path: main:.aynig/command/build\n") {
		t.Fatalf("expected command path on the default branch, got %q", output)
	}
}

func TestStatusCommandsRefSameResolvesFromInspectedBranch(t *testing.T) {
	repoDir := newStatusTestRepo(t)
	commitEmpty(t, repoDir, "seed", "body")
	runGit(t, repoDir, "checkout", "-b", "1-evil")
	writeRepoFile(t, repoDir, ".aynig/command/build", "#!/bin/sh\n", 0o755)
	runGit(t, repoDir, "add", "-A")
	runGit(t, repoDir, "commit", "-m", "feat: evil", "-m", "body\n\ndwp-state: build")
	runGit(t, repoDir, "checkout", "main")

	output := captureStatusOutput(t, repoDir, StatusOptions{Branch: "1-evil", CommandsRef: "same"})

	if !strings.Contains(output, "command: exists\n") {
		t.Fatalf("expected command from the inspected branch, got %q", output)
	}
	if !strings.Contains(output, "command-path: refs/heads/1-evil:.aynig/command/build\n") {
		t.Fatalf("expected command path on the inspected branch, got %q", output)
	}
}

func TestStatusExplicitCommandsRef(t *testing.T) {
	repoDir := newStatusTestRepo(t)
	commitEmpty(t, repoDir, "seed", "body")
	runGit(t, repoDir, "checkout", "-b", "commands")
	writeRepoFile(t, repoDir, ".aynig/command/build", "#!/bin/sh\n", 0o755)
	commitTree(t, repoDir, "add commands")
	runGit(t, repoDir, "checkout", "main")
	createBranchCommit(t, repoDir, "1-bootstrap", "feat: bootstrap", "body\n\ndwp-state: build")

	output := captureStatusOutput(t, repoDir, StatusOptions{Branch: "1-bootstrap", CommandsRef: "commands"})

	if !strings.Contains(output, "command: exists\n") {
		t.Fatalf("expected command from the explicit commands ref, got %q", output)
	}
	if !strings.Contains(output, "command-path: commands:.aynig/command/build\n") {
		t.Fatalf("expected command path on the commands ref, got %q", output)
	}
}

func TestStatusUnknownCommandsRefFails(t *testing.T) {
	repoDir := newStatusTestRepo(t)
	commitEmpty(t, repoDir, "seed", "body")
	chdir(t, repoDir)

	err := Status(StatusOptions{CommandsRef: "does-not-exist"})
	if err == nil {
		t.Fatal("expected an error for an unknown commands ref")
	}
	if !strings.Contains(err.Error(), "cannot resolve commands ref") {
		t.Fatalf("unexpected error: %v", err)
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

func writeRepoFile(t *testing.T, repoDir string, relPath string, content string, mode os.FileMode) {
	t.Helper()
	fullPath := filepath.Join(repoDir, filepath.FromSlash(relPath))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(fullPath, []byte(content), mode); err != nil {
		t.Fatalf("write failed: %v", err)
	}
}

func commitTree(t *testing.T, repoDir string, subject string) {
	t.Helper()
	runGit(t, repoDir, "add", "-A")
	runGit(t, repoDir, "commit", "-m", subject)
}

func chdir(t *testing.T, dir string) {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(cwd)
	})
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
}

func captureStatusOutput(t *testing.T, repoDir string, options StatusOptions) string {
	t.Helper()
	chdir(t, repoDir)

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
