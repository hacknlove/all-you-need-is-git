package statex

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildCommitMessage(t *testing.T) {
	message := BuildCommitMessage("chore: working", "keep lease alive", []Trailer{
		{Key: "dwp-state", Value: "working"},
		{Key: "dwp-run-id", Value: "abc"},
	})

	expected := []string{
		"chore: working",
		"",
		"keep lease alive",
		"",
		"dwp-state: working",
		"dwp-run-id: abc",
		"",
	}
	if message != strings.Join(expected, "\n") {
		t.Fatalf("unexpected message:\n%s", message)
	}
}

func TestVerifyHeadStateTrailer(t *testing.T) {
	tests := []struct {
		name        string
		message     string
		wantErrText string
	}{
		{
			name: "valid state trailer",
			message: strings.Join([]string{
				"chore: build",
				"",
				"run build",
				"",
				"dwp-state: build",
			}, "\n"),
		},
		{
			name: "multiple state trailers last wins",
			message: strings.Join([]string{
				"chore: review",
				"",
				"review output",
				"",
				"dwp-state: build",
				"dwp-state: review",
			}, "\n"),
		},
		{
			name: "missing state trailer",
			message: strings.Join([]string{
				"chore: build",
				"",
				"run build",
				"",
				"dwp-run-id: abc",
			}, "\n"),
			wantErrText: "Invalid trailer block: missing dwp-state trailer",
		},
		{
			name: "empty state trailer",
			message: strings.Join([]string{
				"chore: build",
				"",
				"run build",
				"",
				"dwp-state:   ",
			}, "\n"),
			wantErrText: "Invalid trailer block: empty dwp-state trailer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoDir := initGitRepo(t)
			commitMessage(t, repoDir, tt.message)

			err := VerifyHeadStateTrailer(repoDir)
			if tt.wantErrText == "" {
				if err != nil {
					t.Fatalf("VerifyHeadStateTrailer() error = %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("VerifyHeadStateTrailer() error = nil, want %q", tt.wantErrText)
			}
			if err.Error() != tt.wantErrText {
				t.Fatalf("VerifyHeadStateTrailer() error = %q, want %q", err.Error(), tt.wantErrText)
			}
		})
	}
}

func initGitRepo(t *testing.T) string {
	t.Helper()

	repoDir := t.TempDir()
	runGit(t, repoDir, "init", ".")
	runGit(t, repoDir, "config", "user.email", "test@example.com")
	runGit(t, repoDir, "config", "user.name", "Test User")
	if err := os.WriteFile(filepath.Join(repoDir, "file.txt"), []byte("ok\n"), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}
	runGit(t, repoDir, "add", ".")
	return repoDir
}

func commitMessage(t *testing.T, repoDir string, message string) {
	t.Helper()
	runGit(t, repoDir, "commit", "-m", message)
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s failed: %v (%s)", strings.Join(args, " "), err, string(output))
	}
}
