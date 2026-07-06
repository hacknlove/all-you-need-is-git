package commands

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"all-you-need-is-git/go/internal/gitx"
	"all-you-need-is-git/go/internal/orchestrator"
)

func leaseStatusForState(state string, leaseSecondsRaw string, committerDate string) string {
	if state != "working" {
		return "n/a"
	}
	leaseSeconds, parseErr := strconv.Atoi(leaseSecondsRaw)
	if parseErr != nil {
		return "unknown"
	}
	lastCommitTime, timeErr := time.Parse(time.RFC3339, strings.TrimSpace(committerDate))
	if timeErr != nil {
		return "unknown"
	}
	expiresAt := lastCommitTime.Add(time.Duration(leaseSeconds) * time.Second)
	if time.Now().After(expiresAt) {
		return "expired"
	}
	return "active"
}

type StatusOptions struct {
	Role          string
	Branch        string
	BranchPattern string
	CommandsRef   string
}

type branchStatus struct {
	Branch      string
	HeadCommit  string
	State       string
	RunID       string
	LeaseStatus string
	OriginState string
	InDWPState  bool
	Command     string
	CommandPath string
}

func Status(options StatusOptions) error {
	repoRoot, err := gitx.Run("", "rev-parse", "--show-toplevel")
	if err != nil {
		return fmt.Errorf("Error: Not a Git repository. Please run `git init` first.")
	}
	repoRoot = strings.TrimSpace(repoRoot)

	roleName := strings.TrimSpace(options.Role)
	if roleName == "" {
		roleName = strings.TrimSpace(os.Getenv("ROLE"))
	}

	// Same defaulting as `aynig run`: empty means the default branch, "same"
	// (returned here as an empty ref) means each inspected branch.
	commandsRef, _, _, err := orchestrator.ResolveCommandsRef(repoRoot, options.CommandsRef, "")
	if err != nil {
		return err
	}

	branches, err := resolveStatusBranches(options)
	if err != nil {
		return err
	}

	for i, branch := range branches {
		status, err := readBranchStatus(branch, repoRoot, roleName, commandsRef)
		if err != nil {
			return err
		}
		printBranchStatus(status)
		if i < len(branches)-1 {
			fmt.Println()
		}
	}
	return nil
}

func resolveStatusBranches(options StatusOptions) ([]string, error) {
	branch := strings.TrimSpace(options.Branch)
	pattern := strings.TrimSpace(options.BranchPattern)
	if branch != "" && pattern != "" {
		return nil, fmt.Errorf("status accepts either --branch or --branch-pattern, not both")
	}
	if branch != "" {
		if _, err := gitx.Run("", "rev-parse", "--verify", "refs/heads/"+branch); err != nil {
			return nil, fmt.Errorf("branch %q not found", branch)
		}
		return []string{branch}, nil
	}
	if pattern != "" {
		out, err := gitx.Run("", "branch", "--list", "--format=%(refname:short)", pattern)
		if err != nil {
			return nil, err
		}
		branches := splitStatusLines(out)
		if len(branches) == 0 {
			return nil, fmt.Errorf("no branches match pattern %q", pattern)
		}
		sort.Strings(branches)
		return branches, nil
	}
	branch, err := gitx.Run("", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return nil, err
	}
	return []string{strings.TrimSpace(branch)}, nil
}

func readBranchStatus(branch string, repoRoot string, roleName string, commandsRef string) (branchStatus, error) {
	branchRef := localBranchRef(branch)
	headCommit, err := gitx.Run("", "rev-parse", branchRef)
	if err != nil {
		return branchStatus{}, err
	}
	headCommit = strings.TrimSpace(headCommit)

	committerDate, err := gitx.Run("", "log", "-1", "--format=%cI", branchRef)
	if err != nil {
		return branchStatus{}, err
	}
	committerDate = strings.TrimSpace(committerDate)

	fullMessage, err := gitx.Run("", "log", "-1", "--format=%B", branchRef)
	if err != nil {
		return branchStatus{}, err
	}
	trailers, err := parseTrailersFromMessage(fullMessage)
	if err != nil {
		return branchStatus{}, err
	}

	state := trailerValue(trailers, "dwp-state")
	runID := trailerValue(trailers, "dwp-run-id")
	leaseSecondsRaw := trailerValue(trailers, "dwp-lease-seconds")
	originState := trailerValue(trailers, "dwp-origin-state")
	inDWPState := state != ""
	leaseStatus := leaseStatusForState(state, leaseSecondsRaw, committerDate)

	commandStatus := "missing"
	commandPath := ""
	commandState := state
	shouldResolveCommand := true
	if !inDWPState {
		shouldResolveCommand = false
		commandStatus = "not in a DWP state"
	} else if state == "working" && originState != "" {
		commandState = originState
	} else if state == "working" {
		shouldResolveCommand = false
	}
	if shouldResolveCommand && commandState != "" && commandState != "working" {
		ref := commandsRef
		if ref == "" {
			// --commands-ref same: commands come from the inspected branch.
			ref = branchRef
		}
		commandStatus, commandPath = resolveCommandPath(repoRoot, ref, roleName, commandState)
	} else if !shouldResolveCommand {
		commandStatus = "lease"
	}

	return branchStatus{
		Branch:      branch,
		HeadCommit:  headCommit,
		State:       state,
		RunID:       runID,
		LeaseStatus: leaseStatus,
		OriginState: originState,
		InDWPState:  inDWPState,
		Command:     commandStatus,
		CommandPath: commandPath,
	}, nil
}

func printBranchStatus(status branchStatus) {
	fmt.Printf("branch: %s\n", status.Branch)
	fmt.Printf("head: %s\n", status.HeadCommit)
	if !status.InDWPState {
		fmt.Printf("dwp-state: not in a DWP state\n")
		fmt.Printf("command: %s\n", status.Command)
		return
	}
	fmt.Printf("dwp-state: %s\n", status.State)
	if status.State == "working" && status.OriginState != "" {
		fmt.Printf("dwp-origin-state: %s\n", status.OriginState)
	}
	if status.RunID != "" {
		fmt.Printf("dwp-run-id: %s\n", status.RunID)
	} else {
		fmt.Printf("dwp-run-id: n/a\n")
	}
	fmt.Printf("lease: %s\n", status.LeaseStatus)
	fmt.Printf("command: %s\n", status.Command)
	if status.CommandPath != "" {
		fmt.Printf("command-path: %s\n", status.CommandPath)
	}
}

func splitStatusLines(out string) []string {
	lines := []string{}
	for _, line := range strings.Split(out, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			lines = append(lines, trimmed)
		}
	}
	return lines
}

func localBranchRef(branch string) string {
	return "refs/heads/" + branch
}

// resolveCommandPath reports whether the command for commandState exists in
// the tree of ref, mirroring how `aynig run` resolves commands: the
// role-specific command wins when it is an executable blob, otherwise the base
// command is checked. Paths are reported in `<ref>:<path>` notation.
func resolveCommandPath(repoRoot string, ref string, roleName string, commandState string) (string, string) {
	if roleName != "" {
		rolePath := path.Join(".aynig", "roles", roleName, "command", commandState)
		if strings.HasPrefix(rolePath, ".aynig/roles/") && treeEntryIsExecutable(repoRoot, ref, rolePath) {
			return "exists", ref + ":" + rolePath
		}
	}
	basePath := path.Join(".aynig", "command", commandState)
	if !strings.HasPrefix(basePath, ".aynig/command/") {
		return "missing", ""
	}
	if treeEntryIsExecutable(repoRoot, ref, basePath) {
		return "exists", ref + ":" + basePath
	}
	return "missing", ref + ":" + basePath
}

func treeEntryIsExecutable(repoRoot string, ref string, treePath string) bool {
	mode, err := gitx.LsTreeEntryMode(repoRoot, ref, treePath)
	return err == nil && mode == "100755"
}
