package orchestrator

import (
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
	return &Repo{config: cfg, logger: logx.New(cfg.LogLevel)}
}

func (r *Repo) Run() error {
	repoRoot, err := gitx.RepoRoot()
	if err != nil {
		return err
	}
	r.config.RepoRoot = repoRoot
	r.logger.Infof("Repository root: %s", repoRoot)

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
