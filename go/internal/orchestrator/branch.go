package orchestrator

import (
	"strings"

	"all-you-need-is-git/go/internal/config"
	"all-you-need-is-git/go/internal/gitx"
)

type Branch struct {
	config          config.Config
	branchName      string
	isCurrentBranch bool
}

func NewBranch(cfg config.Config, branchName string, isCurrentBranch bool) *Branch {
	return &Branch{
		config:          cfg,
		branchName:      branchName,
		isCurrentBranch: isCurrentBranch,
	}
}

func (b *Branch) Run() error {
	if b.config.UseRemote != "" && !strings.HasPrefix(b.branchName, b.config.UseRemote+"/") {
		return nil
	}

	commit, err := gitx.ReadCommit(b.branchName)
	if err != nil {
		return err
	}

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
