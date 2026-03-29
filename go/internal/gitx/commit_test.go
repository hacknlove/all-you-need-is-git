package gitx

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestReadCommitParsesTrailerOnlyBody(t *testing.T) {
	repoDir := t.TempDir()
	runGit(t, repoDir, "init", ".")
	runGit(t, repoDir, "config", "user.email", "test@example.com")
	runGit(t, repoDir, "config", "user.name", "Test User")
	if err := os.WriteFile(filepath.Join(repoDir, "file.txt"), []byte("ok\n"), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}
	runGit(t, repoDir, "add", ".")
	runGit(t, repoDir, "commit", "-m", "chore: set bar", "-m", "dwp-state: bar")

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}
	if err := os.Chdir(repoDir); err != nil {
		t.Fatalf("chdir failed: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(wd)
	})

	commit, err := ReadCommit("HEAD")
	if err != nil {
		t.Fatalf("ReadCommit failed: %v", err)
	}
	if got := commit.Trailers["dwp-state"]; len(got) != 1 || got[0] != "bar" {
		t.Fatalf("unexpected trailers: %#v", commit.Trailers)
	}
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v (%s)", args, err, string(output))
	}
}
