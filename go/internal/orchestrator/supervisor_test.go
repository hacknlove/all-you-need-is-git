package orchestrator

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"all-you-need-is-git/go/internal/gitx"
	"all-you-need-is-git/go/internal/statex"
)

func TestParseSetStateLine(t *testing.T) {
	line := "SET_STATE {\"state\":\"review\",\"subject\":\"ready\",\"body\":\"line1\\nline2\"}\n"
	result, ok, err := parseSetStateLine(line)
	if err != nil {
		t.Fatalf("parseSetStateLine returned error: %v", err)
	}
	if !ok {
		t.Fatalf("expected SET_STATE line to be recognized")
	}
	if result.State != "review" {
		t.Fatalf("unexpected state: %q", result.State)
	}
	if result.Subject != "ready" {
		t.Fatalf("unexpected subject: %q", result.Subject)
	}
	if result.Body != "line1\nline2" {
		t.Fatalf("unexpected body: %q", result.Body)
	}
}

func TestSuperviseAppliesLastValidSetState(t *testing.T) {
	repoDir := initSupervisorTestRepo(t)
	writeWorkingState(t, repoDir, "run-123")

	commandPath := writeSupervisorCommand(t, repoDir, "state-command.sh", `#!/bin/sh
printf 'noise before\n'
printf 'SET_STATE {"state":"review","body":"first"}\n'
printf 'SET_STATE {"state":"complete","subject":"ship","body":"line1\\nline2","trailers":[{"key":"dwp-note","value":"ready"}]}\n'
printf 'warn on stderr\n' >&2
`)
	stdoutLogPath := filepath.Join(repoDir, ".aynig", "logs", "stdout.log")
	stderrLogPath := filepath.Join(repoDir, ".aynig", "logs", "stderr.log")

	if err := Supervise(SuperviseOptions{
		WorktreePath:  repoDir,
		CommandPath:   commandPath,
		RunID:         "run-123",
		StdoutLogPath: stdoutLogPath,
		StderrLogPath: stderrLogPath,
	}); err != nil {
		t.Fatalf("Supervise failed: %v", err)
	}

	fullMessage := runGitOutput(t, repoDir, "log", "-1", "--format=%B")
	wantMessage := statex.BuildCommitMessage("ship", "line1\nline2", []statex.Trailer{
		{Key: "dwp-state", Value: "complete"},
		{Key: "dwp-note", Value: "ready"},
	})
	if strings.TrimRight(fullMessage, "\n") != strings.TrimRight(wantMessage, "\n") {
		t.Fatalf("unexpected commit message:\n%s", fullMessage)
	}

	stdoutLog := readFile(t, stdoutLogPath)
	if !strings.Contains(stdoutLog, `SET_STATE {"state":"complete","subject":"ship","body":"line1\nline2","trailers":[{"key":"dwp-note","value":"ready"}]}`) {
		t.Fatalf("stdout log missing final SET_STATE line:\n%s", stdoutLog)
	}

	stderrLog := readFile(t, stderrLogPath)
	if !strings.Contains(stderrLog, "warn on stderr") {
		t.Fatalf("stderr log missing command output:\n%s", stderrLog)
	}
}

func TestSuperviseLeavesWorkingWithoutSetState(t *testing.T) {
	repoDir := initSupervisorTestRepo(t)
	writeWorkingState(t, repoDir, "run-123")
	before := runGitOutput(t, repoDir, "rev-parse", "HEAD")

	commandPath := writeSupervisorCommand(t, repoDir, "no-state.sh", `#!/bin/sh
printf 'still working\n'
printf 'nothing to see here\n' >&2
`)
	stdoutLogPath := filepath.Join(repoDir, ".aynig", "logs", "stdout.log")
	stderrLogPath := filepath.Join(repoDir, ".aynig", "logs", "stderr.log")

	if err := Supervise(SuperviseOptions{
		WorktreePath:  repoDir,
		CommandPath:   commandPath,
		RunID:         "run-123",
		StdoutLogPath: stdoutLogPath,
		StderrLogPath: stderrLogPath,
	}); err != nil {
		t.Fatalf("Supervise failed: %v", err)
	}

	after := runGitOutput(t, repoDir, "rev-parse", "HEAD")
	if before != after {
		t.Fatalf("expected HEAD to remain unchanged, before=%s after=%s", before, after)
	}

	stderrLog := readFile(t, stderrLogPath)
	if !strings.Contains(stderrLog, "no valid SET_STATE line found") {
		t.Fatalf("stderr log missing supervisor note:\n%s", stderrLog)
	}
}

func TestSuperviseSkipsStateUpdateWhenRunIDChanges(t *testing.T) {
	repoDir := initSupervisorTestRepo(t)
	writeWorkingState(t, repoDir, "run-123")
	before := runGitOutput(t, repoDir, "rev-parse", "HEAD")

	commandPath := writeSupervisorCommand(t, repoDir, "mismatch.sh", `#!/bin/sh
printf 'SET_STATE {"state":"review","body":"done"}\n'
`)
	stdoutLogPath := filepath.Join(repoDir, ".aynig", "logs", "stdout.log")
	stderrLogPath := filepath.Join(repoDir, ".aynig", "logs", "stderr.log")

	if err := Supervise(SuperviseOptions{
		WorktreePath:  repoDir,
		CommandPath:   commandPath,
		RunID:         "run-999",
		StdoutLogPath: stdoutLogPath,
		StderrLogPath: stderrLogPath,
	}); err != nil {
		t.Fatalf("Supervise failed: %v", err)
	}

	after := runGitOutput(t, repoDir, "rev-parse", "HEAD")
	if before != after {
		t.Fatalf("expected HEAD to remain unchanged, before=%s after=%s", before, after)
	}

	commit, err := gitx.ReadCommitInDir(repoDir, "HEAD")
	if err != nil {
		t.Fatalf("ReadCommitInDir failed: %v", err)
	}
	if state := firstTrailer(commit.Trailers["dwp-state"], ""); state != "working" {
		t.Fatalf("expected HEAD to stay in working, got %q", state)
	}
}

func initSupervisorTestRepo(t *testing.T) string {
	t.Helper()
	repoDir := t.TempDir()
	runGit(t, repoDir, "init", ".")
	runGit(t, repoDir, "config", "user.email", "test@example.com")
	runGit(t, repoDir, "config", "user.name", "Test User")
	runGit(t, repoDir, "commit", "--allow-empty", "-m", "seed")
	return repoDir
}

func writeWorkingState(t *testing.T, repoDir string, runID string) {
	t.Helper()
	err := statex.CommitState(repoDir, "chore: working", "lease", []statex.Trailer{
		{Key: "dwp-state", Value: "working"},
		{Key: "dwp-origin-state", Value: "review"},
		{Key: "dwp-run-id", Value: runID},
		{Key: "dwp-lease-seconds", Value: "300"},
	})
	if err != nil {
		t.Fatalf("CommitState failed: %v", err)
	}
}

func writeSupervisorCommand(t *testing.T, repoDir string, name string, script string) string {
	t.Helper()
	commandPath := filepath.Join(repoDir, name)
	if err := os.WriteFile(commandPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write command failed: %v", err)
	}
	return commandPath
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v (%s)", args, err, string(output))
	}
}

func runGitOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v (%s)", args, err, string(output))
	}
	return string(output)
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file failed: %v", err)
	}
	return string(content)
}
