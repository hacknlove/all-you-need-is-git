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

func TestParseSetStateLineReadsKeepTrailers(t *testing.T) {
	line := "SET_STATE {\"state\":\"review\",\"keep_trailers\":true}\n"
	result, ok, err := parseSetStateLine(line)
	if err != nil {
		t.Fatalf("parseSetStateLine returned error: %v", err)
	}
	if !ok {
		t.Fatalf("expected SET_STATE line to be recognized")
	}
	if !result.KeepTrailers {
		t.Fatalf("expected keep_trailers to be true")
	}
}

func TestSuperviseAppliesLastValidSetState(t *testing.T) {
	repoDir := initSupervisorTestRepo(t)
	writeWorkingState(t, repoDir, "run-123", []statex.Trailer{})

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

func TestSuperviseKeepsWorkflowTrailersWhenRequested(t *testing.T) {
	repoDir := initSupervisorTestRepo(t)
	writeWorkingState(t, repoDir, "run-123", []statex.Trailer{
		{Key: "dwp-attempt", Value: "2"},
		{Key: "dwp-max-attempts", Value: "4"},
		{Key: "dwp-note", Value: "old-note"},
		{Key: "dwp-issue", Value: "42"},
	})

	commandPath := writeSupervisorCommand(t, repoDir, "keep.sh", `#!/bin/sh
printf 'SET_STATE {"state":"review","keep_trailers":true,"trailers":[{"key":"dwp-note","value":"new-note"}]}\n'
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

	commit, err := gitx.ReadCommitInDir(repoDir, "HEAD")
	if err != nil {
		t.Fatalf("ReadCommitInDir failed: %v", err)
	}
	if state := firstTrailer(commit.Trailers["dwp-state"], ""); state != "review" {
		t.Fatalf("expected review state, got %q", state)
	}
	if attempt := firstTrailer(commit.Trailers["dwp-attempt"], ""); attempt != "2" {
		t.Fatalf("expected dwp-attempt trailer to be preserved, got %q", attempt)
	}
	if maxAttempts := firstTrailer(commit.Trailers["dwp-max-attempts"], ""); maxAttempts != "4" {
		t.Fatalf("expected dwp-max-attempts trailer to be preserved, got %q", maxAttempts)
	}
	if issue := firstTrailer(commit.Trailers["dwp-issue"], ""); issue != "42" {
		t.Fatalf("expected dwp-issue trailer to be preserved, got %q", issue)
	}
	notes := commit.Trailers["dwp-note"]
	if len(notes) != 2 || notes[0] != "old-note" || notes[1] != "new-note" {
		t.Fatalf("expected old and new dwp-note trailers, got %#v", notes)
	}
	if originState := firstTrailer(commit.Trailers["dwp-origin-state"], ""); originState != "" {
		t.Fatalf("expected working-only dwp-origin-state to be dropped, got %q", originState)
	}
	if runID := firstTrailer(commit.Trailers["dwp-run-id"], ""); runID != "" {
		t.Fatalf("expected working-only dwp-run-id to be dropped, got %q", runID)
	}
	if lease := firstTrailer(commit.Trailers["dwp-lease-seconds"], ""); lease != "" {
		t.Fatalf("expected working-only dwp-lease-seconds to be dropped, got %q", lease)
	}
}

func TestSuperviseRefreshesWorkingWithoutSetState(t *testing.T) {
	repoDir := initSupervisorTestRepo(t)
	writeWorkingState(t, repoDir, "run-123", []statex.Trailer{
		{Key: "dwp-note", Value: "carry-me"},
		{Key: "dwp-attempt", Value: "2"},
	})
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
	if before == after {
		t.Fatalf("expected HEAD to advance, before=%s after=%s", before, after)
	}

	commit, err := gitx.ReadCommitInDir(repoDir, "HEAD")
	if err != nil {
		t.Fatalf("ReadCommitInDir failed: %v", err)
	}
	if state := firstTrailer(commit.Trailers["dwp-state"], ""); state != "working" {
		t.Fatalf("expected HEAD to stay in working, got %q", state)
	}
	if runID := firstTrailer(commit.Trailers["dwp-run-id"], ""); runID != "run-123" {
		t.Fatalf("expected run id to be preserved, got %q", runID)
	}
	if note := firstTrailer(commit.Trailers["dwp-note"], ""); note != "carry-me" {
		t.Fatalf("expected dwp-note trailer to be preserved, got %q", note)
	}
	if attempt := firstTrailer(commit.Trailers["dwp-attempt"], ""); attempt != "2" {
		t.Fatalf("expected dwp-attempt trailer to be preserved, got %q", attempt)
	}
	if !strings.Contains(commit.Body, "command ended without changing the state") {
		t.Fatalf("unexpected working body: %q", commit.Body)
	}

	stderrLog := readFile(t, stderrLogPath)
	if !strings.Contains(stderrLog, "no valid SET_STATE line found; refreshing working state") {
		t.Fatalf("stderr log missing supervisor note:\n%s", stderrLog)
	}
}

func TestSuperviseMarksStalledWhenCommandFails(t *testing.T) {
	repoDir := initSupervisorTestRepo(t)
	writeWorkingState(t, repoDir, "run-123", []statex.Trailer{})

	commandPath := writeSupervisorCommand(t, repoDir, "fails.sh", `#!/bin/sh
printf 'line-1\n'
printf 'line-2\n'
printf 'line-3\n'
printf 'line-4\n'
printf 'line-5\n'
printf 'line-6\n'
printf 'SET_STATE {"state":"done","body":"too early"}\n'
printf 'err-1\n' >&2
printf 'err-2\n' >&2
printf 'err-3\n' >&2
printf 'err-4\n' >&2
printf 'err-5\n' >&2
printf 'err-6\n' >&2
exit 7
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

	commit, err := gitx.ReadCommitInDir(repoDir, "HEAD")
	if err != nil {
		t.Fatalf("ReadCommitInDir failed: %v", err)
	}
	if state := firstTrailer(commit.Trailers["dwp-state"], ""); state != "stalled" {
		t.Fatalf("expected stalled state, got %q", state)
	}
	if stalledRun := firstTrailer(commit.Trailers["dwp-stalled-run"], ""); stalledRun != "run-123" {
		t.Fatalf("expected stalled run trailer, got %q", stalledRun)
	}
	if !strings.Contains(commit.Body, "command exited with code 7") {
		t.Fatalf("expected exit code in stalled body, got %q", commit.Body)
	}
	if !strings.Contains(commit.Body, "line-3\nline-4\nline-5\nline-6\nSET_STATE {\"state\":\"done\",\"body\":\"too early\"}") {
		t.Fatalf("expected stdout tail in stalled body, got %q", commit.Body)
	}
	if strings.Contains(commit.Body, "line-1") {
		t.Fatalf("expected stalled body to include only the stdout tail, got %q", commit.Body)
	}
	if !strings.Contains(commit.Body, "err-3\nerr-4\nerr-5\nerr-6\naynig __supervise: command exited with error: exit status 7") {
		t.Fatalf("expected stderr tail in stalled body, got %q", commit.Body)
	}
}

func TestSuperviseSkipsStateUpdateWhenRunIDChanges(t *testing.T) {
	repoDir := initSupervisorTestRepo(t)
	writeWorkingState(t, repoDir, "run-123", []statex.Trailer{})
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

func writeWorkingState(t *testing.T, repoDir string, runID string, extraTrailers []statex.Trailer) {
	t.Helper()
	trailers := []statex.Trailer{
		{Key: "dwp-state", Value: "working"},
		{Key: "dwp-origin-state", Value: "review"},
		{Key: "dwp-run-id", Value: runID},
		{Key: "dwp-lease-seconds", Value: "300"},
	}
	trailers = append(trailers, extraTrailers...)
	err := statex.CommitState(repoDir, "chore: working", "lease", trailers)
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
