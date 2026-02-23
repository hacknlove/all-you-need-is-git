package orchestrator

import (
	"os"
	"strings"

	"all-you-need-is-git/go/internal/config"
	"all-you-need-is-git/go/internal/gitx"
	"all-you-need-is-git/go/internal/logx"
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
	buffer := logx.NewBufferedLogger()
	buffer.Debugf("Inspecting branch %s", b.branchName)
	baseLevel := logx.ResolveLevel(b.config.LogLevel, b.config.LogLevelSet, "", os.Getenv("AYNIG_LOG_LEVEL"), config.Default().LogLevel)
	baseLogger := logx.New(baseLevel)
	if b.config.UseRemote != "" && !strings.HasPrefix(b.branchName, b.config.UseRemote+"/") {
		buffer.Debugf("Skipping branch %s (not on remote %s)", b.branchName, b.config.UseRemote)
		buffer.Flush(baseLogger)
		return nil
	}

	commit, err := gitx.ReadCommit(b.branchName)
	if err != nil {
		buffer.Flush(baseLogger)
		return err
	}
	trailerLevel := trailerValue(commit.Trailers, "aynig-log-level")
	resolvedLevel := logx.ResolveLevel(b.config.LogLevel, b.config.LogLevelSet, trailerLevel, os.Getenv("AYNIG_LOG_LEVEL"), config.Default().LogLevel)
	branchLogger := logx.New(resolvedLevel)
	buffer.Flush(branchLogger)

	cmd := NewCommand(CommandParams{
		Config:          b.config,
		BranchName:      b.branchName,
		IsCurrentBranch: b.isCurrentBranch,
		Trailers:        commit.Trailers,
		Body:            commit.Body,
		CommitDate:      commit.Date,
		LogLevel:        resolvedLevel,
	})
	return cmd.Run()
}

func trailerValue(trailers map[string][]string, key string) string {
	values, ok := trailers[key]
	if !ok || len(values) == 0 {
		return ""
	}
	value := strings.TrimSpace(values[0])
	if value == "" {
		return ""
	}
	return value
}
