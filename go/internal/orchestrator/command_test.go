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

func TestCommandEnvIncludesLogPath(t *testing.T) {
	cmd := NewCommand(CommandParams{
		Config: config.Config{Role: "reviewer"},
		Trailers: map[string][]string{
			"dwp-state": {"review"},
			"foo-bar":   {"baz", "qux"},
		},
		Body:     "prompt body",
		LogLevel: "debug",
	})

	logPath := filepath.Join("/tmp", ".dwp", "logs", "deadbeef.log")
	env := cmd.commandEnv("deadbeef", logPath)

	wantEntries := []string{
		"BODY=prompt body",
		"AYNIG_BODY=prompt body",
		"COMMIT_HASH=deadbeef",
		"AYNIG_COMMIT_HASH=deadbeef",
		"LOG_PATH=" + logPath,
		"AYNIG_LOG_PATH=" + logPath,
		"LOG_LEVEL=debug",
		"AYNIG_LOG_LEVEL=debug",
		"ROLE=reviewer",
		"AYNIG_ROLE=reviewer",
		"DWP_STATE=review",
		"AYNIG_TRAILER_DWP_STATE=review",
		"FOO_BAR=baz,qux",
		"AYNIG_TRAILER_FOO_BAR=baz,qux",
	}
	for _, want := range wantEntries {
		if !containsEnvEntry(env, want) {
			t.Fatalf("missing env entry %q in %v", want, env)
		}
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
