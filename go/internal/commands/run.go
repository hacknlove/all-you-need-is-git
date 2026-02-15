package commands

import (
	"all-you-need-is-git/go/internal/config"
	"all-you-need-is-git/go/internal/orchestrator"
)

func Run(cfg config.Config) error {
	repo := orchestrator.NewRepo(cfg)
	return repo.Run()
}
