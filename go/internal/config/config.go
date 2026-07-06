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
	// CommandsRef selects the committish that .aynig commands are resolved
	// from: empty means the default branch, "same" means the event branch.
	CommandsRef string
	// CommandsRoot is the materialized checkout of CommandsRef; empty means
	// commands resolve from each event branch's worktree.
	CommandsRoot string
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
		CommandsRef:   "",
	}
}
