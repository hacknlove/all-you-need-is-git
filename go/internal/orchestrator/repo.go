package orchestrator

import (
	"os"
	"strings"
	"sync"

	"all-you-need-is-git/go/internal/config"
	"all-you-need-is-git/go/internal/gitx"
	"all-you-need-is-git/go/internal/logx"
)

type Repo struct {
	config config.Config
	logger logx.Logger
}

func NewRepo(cfg config.Config) *Repo {
	resolved := logx.ResolveLevel(cfg.LogLevel, cfg.LogLevelSet, "", os.Getenv("AYNIG_LOG_LEVEL"), config.Default().LogLevel)
	return &Repo{config: cfg, logger: logx.New(resolved)}
}

func (r *Repo) Run() error {
	repoRoot, err := gitx.RepoRoot()
	if err != nil {
		return err
	}
	r.config.RepoRoot = repoRoot
	r.logger.Infof("Repository root: %s", repoRoot)

	if r.config.UseRemote == "" {
		commit, readErr := gitx.ReadCommit("HEAD")
		if readErr == nil {
			if remote := firstLowerTrailerValue(commit.Trailers, "aynig-remote"); remote != "" {
				r.config.UseRemote = remote
				r.logger.Infof("Using remote from trailer aynig-remote=%s", remote)
			}
		}
	}

	if r.config.UseRemote != "" {
		r.logger.Infof("Fetching remote branches from %s", r.config.UseRemote)
		if err := gitx.Fetch(repoRoot); err != nil {
			return err
		}
	}

	current, err := gitx.BranchCurrent(repoRoot)
	if err != nil {
		return err
	}
	r.logger.Debugf("Current branch: %s", current)

	var branches []string
	if r.config.UseRemote != "" {
		branches, err = gitx.BranchListRemote(repoRoot)
	} else {
		branches, err = gitx.BranchListLocal(repoRoot)
	}
	if err != nil {
		return err
	}
	r.logger.Debugf("Found %d branches", len(branches))

	branchNames := filterBranches(branches, current, r.config.CurrentBranch)
	if r.config.UseRemote != "" {
		upstream, err := gitx.BranchUpstream(repoRoot, current)
		if err != nil {
			return err
		}
		if upstream != "" {
			if strings.HasPrefix(upstream, r.config.UseRemote+"/") {
				branchNames = filterBranches(branches, upstream, r.config.CurrentBranch)
			} else {
				r.logger.Warnf("Current branch upstream %s does not belong to remote %s", upstream, r.config.UseRemote)
				branchNames = filterBranches(branches, "", r.config.CurrentBranch)
			}
		} else {
			if r.config.CurrentBranch == "only" {
				r.logger.Warnf("Current branch has no upstream for --current-branch=only in remote mode")
			}
			branchNames = filterBranches(branches, "", r.config.CurrentBranch)
		}
	}
	r.logger.Infof("Running %d branches (current-branch=%s)", len(branchNames), r.config.CurrentBranch)

	var wg sync.WaitGroup
	errCh := make(chan error, len(branchNames))
	for _, name := range branchNames {
		branch := NewBranch(r.config, name, name == current)
		wg.Add(1)
		go func(b *Branch) {
			defer wg.Done()
			if err := b.Run(); err != nil {
				errCh <- err
			}
		}(branch)
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}

func filterBranches(all []string, current string, mode string) []string {
	switch mode {
	case "skip":
		out := []string{}
		for _, name := range all {
			if name != current {
				out = append(out, name)
			}
		}
		return out
	case "only":
		if current == "" {
			return []string{}
		}
		return []string{current}
	default:
		return all
	}
}

func firstLowerTrailerValue(trailers map[string][]string, key string) string {
	want := strings.ToLower(strings.TrimSpace(key))
	if want == "" {
		return ""
	}
	for k, values := range trailers {
		if strings.ToLower(strings.TrimSpace(k)) != want {
			continue
		}
		if len(values) == 0 {
			return ""
		}
		return strings.TrimSpace(values[0])
	}
	return ""
}
