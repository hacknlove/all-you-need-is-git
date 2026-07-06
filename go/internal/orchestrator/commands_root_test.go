package orchestrator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"all-you-need-is-git/go/internal/config"
)

func TestPrepareCommandsRootSameReturnsEmpty(t *testing.T) {
	repo := newCommandsRootRepo(t, config.Config{CommandsRef: "same"})

	got, err := repo.prepareCommandsRoot()
	if err != nil {
		t.Fatalf("prepareCommandsRoot failed: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty commands root for 'same', got %q", got)
	}
}

func TestPrepareCommandsRootDefaultsToMaster(t *testing.T) {
	repoDir := initCommandsRootTestRepo(t, "master")
	commitCommand(t, repoDir, "build", "#!/bin/sh\necho master\n")
	runGit(t, repoDir, "checkout", "-b", "evil")
	commitCommand(t, repoDir, "build", "#!/bin/sh\necho evil\n")

	repo := newCommandsRootRepo(t, config.Config{RepoRoot: repoDir, WorkTree: ".worktrees"})
	got, err := repo.prepareCommandsRoot()
	if err != nil {
		t.Fatalf("prepareCommandsRoot failed: %v", err)
	}
	if got == "" {
		t.Fatal("expected a commands root path")
	}
	content := readFile(t, filepath.Join(got, ".aynig", "command", "build"))
	if !strings.Contains(content, "echo master") {
		t.Fatalf("expected commands from master, got: %q", content)
	}
}

func TestPrepareCommandsRootExplicitRef(t *testing.T) {
	repoDir := initCommandsRootTestRepo(t, "master")
	commitCommand(t, repoDir, "build", "#!/bin/sh\necho master\n")
	runGit(t, repoDir, "checkout", "-b", "other")
	commitCommand(t, repoDir, "build", "#!/bin/sh\necho other\n")
	runGit(t, repoDir, "checkout", "master")

	repo := newCommandsRootRepo(t, config.Config{RepoRoot: repoDir, WorkTree: ".worktrees", CommandsRef: "other"})
	got, err := repo.prepareCommandsRoot()
	if err != nil {
		t.Fatalf("prepareCommandsRoot failed: %v", err)
	}
	content := readFile(t, filepath.Join(got, ".aynig", "command", "build"))
	if !strings.Contains(content, "echo other") {
		t.Fatalf("expected commands from 'other', got: %q", content)
	}
}

func TestPrepareCommandsRootReusesExistingCheckout(t *testing.T) {
	repoDir := initCommandsRootTestRepo(t, "master")
	commitCommand(t, repoDir, "build", "#!/bin/sh\necho master\n")

	repo := newCommandsRootRepo(t, config.Config{RepoRoot: repoDir, WorkTree: ".worktrees"})
	first, err := repo.prepareCommandsRoot()
	if err != nil {
		t.Fatalf("first prepareCommandsRoot failed: %v", err)
	}
	second, err := repo.prepareCommandsRoot()
	if err != nil {
		t.Fatalf("second prepareCommandsRoot failed: %v", err)
	}
	if first != second {
		t.Fatalf("expected the same checkout path, got %q and %q", first, second)
	}
}

func TestPrepareCommandsRootUnknownRefFails(t *testing.T) {
	repoDir := initCommandsRootTestRepo(t, "master")

	repo := newCommandsRootRepo(t, config.Config{RepoRoot: repoDir, WorkTree: ".worktrees", CommandsRef: "does-not-exist"})
	if _, err := repo.prepareCommandsRoot(); err == nil {
		t.Fatal("expected an error for an unknown commands ref")
	}
}

func TestDefaultCommandsRefErrorsWithoutMasterOrMain(t *testing.T) {
	repoDir := initCommandsRootTestRepo(t, "trunk")

	repo := newCommandsRootRepo(t, config.Config{RepoRoot: repoDir, WorkTree: ".worktrees"})
	_, err := repo.prepareCommandsRoot()
	if err == nil {
		t.Fatal("expected an error when neither master nor main exists")
	}
	if !strings.Contains(err.Error(), "--commands-ref") {
		t.Fatalf("expected the error to mention --commands-ref, got: %v", err)
	}
}

func newCommandsRootRepo(t *testing.T, cfg config.Config) *Repo {
	t.Helper()
	if cfg.LogLevel == "" {
		cfg.LogLevel = "error"
	}
	return NewRepo(cfg)
}

func initCommandsRootTestRepo(t *testing.T, defaultBranch string) string {
	t.Helper()
	repoDir := t.TempDir()
	runGit(t, repoDir, "init", "-b", defaultBranch, ".")
	runGit(t, repoDir, "config", "user.email", "test@example.com")
	runGit(t, repoDir, "config", "user.name", "Test User")
	runGit(t, repoDir, "commit", "--allow-empty", "-m", "seed")
	return repoDir
}

func commitCommand(t *testing.T, repoDir string, name string, script string) {
	t.Helper()
	commandDir := filepath.Join(repoDir, ".aynig", "command")
	if err := os.MkdirAll(commandDir, 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(commandDir, name), []byte(script), 0o755); err != nil {
		t.Fatalf("write command failed: %v", err)
	}
	runGit(t, repoDir, "add", ".aynig")
	runGit(t, repoDir, "commit", "-m", "add command "+name)
}
