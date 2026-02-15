package config

type Config struct {
	UseRemote     string
	WorkTree      string
	CurrentBranch string
	LeaseSeconds  int
	RepoRoot      string
}

func Default() Config {
	return Config{
		UseRemote:     "",
		WorkTree:      ".worktrees",
		CurrentBranch: "skip",
		LeaseSeconds:  300,
	}
}
