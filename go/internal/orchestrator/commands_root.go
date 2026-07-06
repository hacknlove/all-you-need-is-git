package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"all-you-need-is-git/go/internal/gitx"
)

// ResolveCommandsRef resolves the configured commands ref shared by run and
// status: empty selects the default branch, and in remote mode a plain branch
// name prefers the remote-tracking ref. It returns same=true (with empty
// ref/hash) for the "same" sentinel, meaning commands resolve from each event
// branch.
func ResolveCommandsRef(repoRoot string, commandsRef string, useRemote string) (string, string, bool, error) {
	ref := strings.TrimSpace(commandsRef)
	if ref == "same" {
		return "", "", true, nil
	}
	if ref == "" {
		defaultRef, err := defaultCommandsRef(repoRoot, useRemote)
		if err != nil {
			return "", "", false, err
		}
		ref = defaultRef
	} else if useRemote != "" && !strings.HasPrefix(ref, useRemote+"/") {
		// In remote mode a plain branch name refers to the remote-tracking
		// ref when one exists, so a stale local branch cannot shadow it.
		if _, err := gitx.RevParseVerify(repoRoot, useRemote+"/"+ref); err == nil {
			ref = useRemote + "/" + ref
		}
	}

	hash, err := gitx.RevParseVerify(repoRoot, ref)
	if err != nil {
		return "", "", false, fmt.Errorf("cannot resolve commands ref %q: %w", ref, err)
	}
	return ref, hash, false, nil
}

// prepareCommandsRoot materializes the checkout that .aynig commands are
// resolved from, pinned at the configured commands ref. It returns "" when
// commands should resolve from each event branch instead (--commands-ref same).
func (r *Repo) prepareCommandsRoot() (string, error) {
	ref, hash, same, err := ResolveCommandsRef(r.config.RepoRoot, r.config.CommandsRef, r.config.UseRemote)
	if err != nil {
		return "", err
	}
	if same {
		r.logger.Warnf("Commands resolve from each event branch (--commands-ref same); any branch can redefine them")
		return "", nil
	}

	baseDir := r.config.WorkTree
	if baseDir == "" {
		baseDir = "."
	}
	worktreePath := filepath.Join(r.config.RepoRoot, baseDir, "commands-"+hash[:12])
	if _, statErr := os.Stat(worktreePath); statErr == nil {
		r.logger.Debugf("Reusing commands checkout at %s", worktreePath)
	} else if err := gitx.WorktreeAdd(r.config.RepoRoot, "--detach", worktreePath, hash); err != nil {
		if _, statErr := os.Stat(worktreePath); statErr != nil {
			return "", fmt.Errorf("cannot create commands checkout for %q: %w", ref, err)
		}
		// Another runner created it concurrently; the checkout is pinned at
		// the same commit, so it is safe to share.
	}
	r.logger.Infof("Commands resolve from %s (%s)", ref, hash)
	return worktreePath, nil
}

func defaultCommandsRef(repoRoot string, useRemote string) (string, error) {
	if useRemote != "" {
		head, err := gitx.RemoteHead(repoRoot, useRemote)
		if err == nil && strings.TrimSpace(head) != "" {
			return strings.TrimSpace(head), nil
		}
		return "", fmt.Errorf("cannot resolve the default branch of remote %s; run 'git remote set-head %s --auto' or pass --commands-ref", useRemote, useRemote)
	}
	for _, name := range []string{"master", "main"} {
		if _, err := gitx.RevParseVerify(repoRoot, "refs/heads/"+name); err == nil {
			return name, nil
		}
	}
	return "", fmt.Errorf("cannot determine the default branch; pass --commands-ref <ref> (or --commands-ref same to resolve commands from each event branch)")
}
