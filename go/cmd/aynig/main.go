package main

import (
	"flag"
	"fmt"
	"os"

	"all-you-need-is-git/go/internal/commands"
	"all-you-need-is-git/go/internal/config"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "version", "-v", "--version":
		fmt.Println(version)
	case "run":
		runCmd(os.Args[2:])
	case "set-working":
		setWorkingCmd(os.Args[2:])
	case "set-state":
		setStateCmd(os.Args[2:])
	case "status":
		if err := commands.Status(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "init":
		if err := commands.Init(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "install":
		installCmd(os.Args[2:])
	case "events":
		eventsCmd(os.Args[2:])
	case "update":
		if err := commands.Update(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "-h", "--help", "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func runCmd(args []string) {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	fs.Usage = func() {
		out := fs.Output()
		fmt.Fprintln(out, "Usage of run:")
		fmt.Fprintln(out, "  --role <name>")
		fmt.Fprintln(out, "        Use role-specific commands from .aynig/roles/<name>/command when available")
		fmt.Fprintln(out, "  --current-branch <mode>")
		fmt.Fprintln(out, "        How to handle the current branch: skip (default), include, only (default \"skip\")")
		fmt.Fprintln(out, "  --log-level <level>")
		fmt.Fprintln(out, "        Log verbosity: debug, info, warn, error (default \"error\")")
		fmt.Fprintln(out, "        Precedence: --log-level > aynig-log-level trailer > AYNIG_LOG_LEVEL")
		fmt.Fprintln(out, "  --aynig-remote <name>")
		fmt.Fprintln(out, "        Use remote branches instead of local (specify remote name, e.g., origin)")
		fmt.Fprintln(out, "  -w, --worktree <path>")
		fmt.Fprintln(out, "        Specify custom worktree directory (default: .worktrees) (default \".worktrees\")")
	}
	worktree := fs.String("worktree", config.Default().WorkTree, "Specify custom worktree directory (default: .worktrees)")
	fs.StringVar(worktree, "w", config.Default().WorkTree, "Specify custom worktree directory (default: .worktrees)")
	useRemote := fs.String("aynig-remote", "", "Use remote branches instead of local (specify remote name, e.g., origin)")
	currentBranch := fs.String("current-branch", config.Default().CurrentBranch, "How to handle the current branch: skip (default), include, only")
	role := fs.String("role", "", "Use role-specific commands from .aynig/roles/<name>/command when available")
	logLevel := &stringFlag{value: config.Default().LogLevel}
	fs.Var(logLevel, "log-level", "Log verbosity: debug, info, warn, error")
	fs.Parse(args)

	cfg := config.Default()
	cfg.WorkTree = *worktree
	cfg.UseRemote = *useRemote
	cfg.CurrentBranch = *currentBranch
	cfg.Role = *role
	cfg.LogLevel = logLevel.value
	cfg.LogLevelSet = logLevel.set

	if err := commands.Run(cfg); err != nil {
		fmt.Fprintln(os.Stderr, "Error trying to set up the Repository:", err)
		os.Exit(1)
	}
}

type stringFlag struct {
	value string
	set   bool
}

func (s *stringFlag) String() string {
	return s.value
}

func (s *stringFlag) Set(value string) error {
	s.value = value
	s.set = true
	return nil
}

func installCmd(args []string) {
	fs := flag.NewFlagSet("install", flag.ExitOnError)
	fs.Parse(args)
	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "install requires <repo> [ref] [subfolder]")
		printUsage()
		os.Exit(1)
	}

	repo := fs.Arg(0)
	ref := ""
	subfolder := ""
	if fs.NArg() > 1 {
		ref = fs.Arg(1)
	}
	if fs.NArg() > 2 {
		subfolder = fs.Arg(2)
	}

	if err := commands.Install(repo, ref, subfolder); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func eventsCmd(args []string) {
	fs := flag.NewFlagSet("events", flag.ExitOnError)
	history := fs.Bool("history", false, "Allow scanning recent history")
	limit := fs.Int("limit", 10, "Number of commits to inspect (requires --history)")
	fs.IntVar(limit, "n", 10, "Number of commits to inspect (requires --history)")
	jsonOut := fs.Bool("json", false, "Output JSON")
	fs.Parse(args)

	if err := commands.Events(commands.EventsOptions{
		History: *history,
		Limit:   *limit,
		JSON:    *jsonOut,
	}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func setWorkingCmd(args []string) {
	fs := flag.NewFlagSet("set-working", flag.ExitOnError)
	subject := fs.String("subject", "", "Commit title")
	prompt := fs.String("prompt", "", "Commit prompt/body")
	promptFile := fs.String("prompt-file", "", "Path to file used as prompt/body")
	promptStdin := fs.Bool("prompt-stdin", false, "Read prompt/body from stdin")
	remote := fs.String("aynig-remote", "", "Remote name to push after commit")
	var trailers trailerListFlag
	fs.Var(&trailers, "trailer", "Additional trailer in key:value format (repeatable)")
	fs.Parse(args)

	if err := commands.SetWorking(commands.SetWorkingOptions{
		Subject:     *subject,
		Prompt:      *prompt,
		PromptFile:  *promptFile,
		PromptStdin: *promptStdin,
		AynigRemote: *remote,
		Trailers:    trailers,
	}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func setStateCmd(args []string) {
	fs := flag.NewFlagSet("set-state", flag.ExitOnError)
	state := fs.String("aynig-state", "", "Next state to set (required, must not be working)")
	subject := fs.String("subject", "", "Commit title")
	prompt := fs.String("prompt", "", "Commit prompt/body")
	promptFile := fs.String("prompt-file", "", "Path to file used as prompt/body")
	promptStdin := fs.Bool("prompt-stdin", false, "Read prompt/body from stdin")
	remote := fs.String("aynig-remote", "", "Remote name to push after commit")
	var trailers trailerListFlag
	fs.Var(&trailers, "trailer", "Additional trailer in key:value format (repeatable)")
	fs.Parse(args)

	if err := commands.SetState(commands.SetStateOptions{
		State:       *state,
		Subject:     *subject,
		Prompt:      *prompt,
		PromptFile:  *promptFile,
		PromptStdin: *promptStdin,
		AynigRemote: *remote,
		Trailers:    trailers,
	}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type trailerListFlag []string

func (t *trailerListFlag) String() string {
	if t == nil {
		return ""
	}
	return fmt.Sprintf("%v", []string(*t))
}

func (t *trailerListFlag) Set(value string) error {
	*t = append(*t, value)
	return nil
}

func printUsage() {
	fmt.Println("Usage: aynig <command> [options]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  run       Run AYNIG for the current repository")
	fmt.Println("            --aynig-remote can also be persisted in commits via aynig-remote trailer")
	fmt.Println("  set-working  Create a working lease commit")
	fmt.Println("  set-state    Create a commit with a non-working aynig-state")
	fmt.Println("  status    Show current AYNIG state")
	fmt.Println("  init      Initialize AYNIG in the current repository")
	fmt.Println("  install   Install AYNIG workflows from another repository")
	fmt.Println("  events    Show recent AYNIG events")
	fmt.Println("  update    Download and install the latest AYNIG release")
	fmt.Println("  version   Print the current AYNIG version")
}
