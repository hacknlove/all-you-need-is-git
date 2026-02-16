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
	worktree := fs.String("worktree", config.Default().WorkTree, "Specify custom worktree directory (default: .worktrees)")
	fs.StringVar(worktree, "w", config.Default().WorkTree, "Specify custom worktree directory (default: .worktrees)")
	useRemote := fs.String("use-remote", "", "Use remote branches instead of local (specify remote name, e.g., origin)")
	currentBranch := fs.String("current-branch", config.Default().CurrentBranch, "How to handle the current branch: skip (default), include, only")
	fs.Parse(args)

	cfg := config.Default()
	cfg.WorkTree = *worktree
	cfg.UseRemote = *useRemote
	cfg.CurrentBranch = *currentBranch

	if err := commands.Run(cfg); err != nil {
		fmt.Fprintln(os.Stderr, "Error trying to set up the Repository:", err)
		os.Exit(1)
	}
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

func printUsage() {
	fmt.Println("Usage: aynig <command> [options]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  run       Run AYNIG for the current repository")
	fmt.Println("  status    Show current AYNIG state")
	fmt.Println("  init      Initialize AYNIG in the current repository")
	fmt.Println("  install   Install AYNIG workflows from another repository")
	fmt.Println("  events    Show recent AYNIG events")
	fmt.Println("  update    Download and install the latest AYNIG release")
	fmt.Println("  version   Print the current AYNIG version")
}
