package orchestrator

import (
	"strings"

	"all-you-need-is-git/go/internal/config"
	"all-you-need-is-git/go/internal/gitx"
	"all-you-need-is-git/go/internal/logx"
)

type Branch struct {
	config          config.Config
	branchName      string
	isCurrentBranch bool
	logger          logx.Logger
}

func NewBranch(cfg config.Config, branchName string, isCurrentBranch bool) *Branch {
	return &Branch{
		config:          cfg,
		branchName:      branchName,
		isCurrentBranch: isCurrentBranch,
		logger:          logx.New(cfg.LogLevel),
	}
}

func (b *Branch) Run() error {
	if b.config.UseRemote != "" && !strings.HasPrefix(b.branchName, b.config.UseRemote+"/") {
		b.logger.Debugf("Skipping branch %s (not on remote %s)", b.branchName, b.config.UseRemote)
		return nil
	}

	commit, err := gitx.ReadCommit(b.branchName)
	if err != nil {
		return err
	}
	b.logger.Debugf("Inspecting branch %s", b.branchName)

	cmd := NewCommand(CommandParams{
		Config:          b.config,
		BranchName:      b.branchName,
		IsCurrentBranch: b.isCurrentBranch,
		Trailers:        commit.Trailers,
		Body:            commit.Body,
		CommitDate:      commit.Date,
	})
	return cmd.Run()
}
