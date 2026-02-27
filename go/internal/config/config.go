package config

type Config struct {
	UseRemote     string
	WorkTree      string
	CurrentBranch string
	Role          string
	LogLevel      string
	LogLevelSet   bool
	LeaseSeconds  int
	RepoRoot      string
}

func Default() Config {
	return Config{
		UseRemote:     "",
		WorkTree:      ".worktrees",
		CurrentBranch: "skip",
		Role:          "",
		LogLevel:      "error",
		LogLevelSet:   false,
		LeaseSeconds:  300,
	}
}
