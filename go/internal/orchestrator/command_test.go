package orchestrator

import (
	"path/filepath"
	"strings"
	"testing"

	"all-you-need-is-git/go/internal/config"
)

func TestResolveStateTrailer(t *testing.T) {
	tests := []struct {
		name       string
		trailers   map[string][]string
		wantState  string
		wantReason string
	}{
		{
			name:       "missing state",
			trailers:   map[string][]string{},
			wantState:  "",
			wantReason: "",
		},
		{
			name:       "single state",
			trailers:   map[string][]string{"dwp-state": []string{" Build "}},
			wantState:  "build",
			wantReason: "",
		},
		{
			name:       "multiple state trailers (last wins)",
			trailers:   map[string][]string{"dwp-state": []string{"build", "review"}},
			wantState:  "review",
			wantReason: "",
		},
		{
			name:       "empty state trailer",
			trailers:   map[string][]string{"dwp-state": []string{"  "}},
			wantState:  "",
			wantReason: "empty dwp-state trailer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotState, gotReason := resolveStateTrailer(tt.trailers)
			if gotState != tt.wantState {
				t.Fatalf("state mismatch: got %q want %q", gotState, tt.wantState)
			}
			if gotReason != tt.wantReason {
				t.Fatalf("reason mismatch: got %q want %q", gotReason, tt.wantReason)
			}
		})
	}
}

func TestGetWorkspaceUsesRepoRootForCurrentBranch(t *testing.T) {
	repoRoot := t.TempDir()
	cmd := NewCommand(CommandParams{
		Config:          config.Config{RepoRoot: repoRoot},
		BranchName:      "main",
		IsCurrentBranch: true,
	})

	got, err := cmd.getWorkspace()
	if err != nil {
		t.Fatalf("getWorkspace returned error: %v", err)
	}
	if got != repoRoot {
		t.Fatalf("workspace mismatch: got %q want %q", got, repoRoot)
	}
}

func TestCommandEnvIncludesLogPaths(t *testing.T) {
	cmd := NewCommand(CommandParams{
		Config: config.Config{Role: "reviewer"},
		Trailers: map[string][]string{
			"dwp-state": {"review"},
			"foo-bar":   {"baz", "qux"},
		},
		Body:     "prompt body",
		LogLevel: "debug",
	})

	stdoutLogPath := filepath.Join("/tmp", ".aynig", "logs", "deadbeef.stdout.log")
	stderrLogPath := filepath.Join("/tmp", ".aynig", "logs", "deadbeef.stderr.log")
	env := cmd.commandEnv("deadbeef", stdoutLogPath, stderrLogPath, "/tmp/worktree")

	wantEntries := []string{
		"BODY=prompt body",
		"COMMIT_HASH=deadbeef",
		"STDOUT_LOG_PATH=" + stdoutLogPath,
		"STDERR_LOG_PATH=" + stderrLogPath,
		"WORKTREE_PATH=/tmp/worktree",
		"LOG_LEVEL=debug",
		"ROLE=reviewer",
		"DWP_STATE=review",
		"FOO_BAR=baz,qux",
	}
	for _, want := range wantEntries {
		if !containsEnvEntry(env, want) {
			t.Fatalf("missing env entry %q in %v", want, env)
		}
	}
}

func TestCommandEnvDoesNotOverrideInheritedOrReservedEnvNames(t *testing.T) {
	t.Setenv("PATH", "/tmp/original-path")
	t.Setenv("HOME", "/tmp/original-home")

	cmd := NewCommand(CommandParams{
		Config: config.Config{Role: "reviewer"},
		Trailers: map[string][]string{
			"path":            {"/tmp/evil-bin"},
			"home":            {"/tmp/evil-home"},
			"role":            {"other"},
			"stdout-log-path": {"/tmp/other.stdout.log"},
			"stderr-log-path": {"/tmp/other.stderr.log"},
			"dwp-state":       {"review"},
		},
		Body:     "prompt body",
		LogLevel: "debug",
	})

	env := cmd.commandEnv("deadbeef", "/tmp/stdout.log", "/tmp/stderr.log", "/tmp/worktree")

	if !containsEnvEntry(env, "PATH=/tmp/original-path") {
		t.Fatalf("expected inherited PATH to be preserved: %v", env)
	}
	if !containsEnvEntry(env, "HOME=/tmp/original-home") {
		t.Fatalf("expected inherited HOME to be preserved: %v", env)
	}
	if containsEnvEntry(env, "PATH=/tmp/evil-bin") {
		t.Fatalf("unexpected trailer override for PATH: %v", env)
	}
	if containsEnvEntry(env, "HOME=/tmp/evil-home") {
		t.Fatalf("unexpected trailer override for HOME: %v", env)
	}
	if containsEnvEntry(env, "ROLE=other") {
		t.Fatalf("unexpected trailer override for ROLE: %v", env)
	}
	if containsEnvEntry(env, "STDOUT_LOG_PATH=/tmp/other.stdout.log") {
		t.Fatalf("unexpected trailer override for STDOUT_LOG_PATH: %v", env)
	}
	if containsEnvEntry(env, "STDERR_LOG_PATH=/tmp/other.stderr.log") {
		t.Fatalf("unexpected trailer override for STDERR_LOG_PATH: %v", env)
	}
}

func TestCommandEnvReservesRoleEvenWhenUnset(t *testing.T) {
	cmd := NewCommand(CommandParams{
		Config: config.Config{},
		Trailers: map[string][]string{
			"role":      {"reviewer"},
			"dwp-state": {"review"},
		},
		Body:     "prompt body",
		LogLevel: "debug",
	})

	env := cmd.commandEnv("deadbeef", "/tmp/stdout.log", "/tmp/stderr.log", "/tmp/worktree")

	if containsEnvEntry(env, "ROLE=reviewer") {
		t.Fatalf("unexpected trailer-created ROLE when no role was configured: %v", env)
	}
	if !containsEnvEntry(env, "DWP_STATE=review") {
		t.Fatalf("expected non-reserved trailer env to remain exported: %v", env)
	}
}

func containsEnvEntry(env []string, want string) bool {
	for _, entry := range env {
		if strings.TrimSpace(entry) == want {
			return true
		}
	}
	return false
}
