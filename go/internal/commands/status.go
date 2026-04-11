package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"all-you-need-is-git/go/internal/gitx"
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

	branches, err := resolveStatusBranches(options)
	if err != nil {
		return err
	}

	for i, branch := range branches {
		status, err := readBranchStatus(branch, repoRoot, roleName)
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

func readBranchStatus(branch string, repoRoot string, roleName string) (branchStatus, error) {
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
		commandStatus, commandPath = resolveCommandPath(repoRoot, roleName, commandState)
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

func resolveCommandPath(repoRoot string, roleName string, commandState string) (string, string) {
	commandStatus := "missing"
	commandPath := ""
	if roleName != "" {
		rolePath := filepath.Join(repoRoot, ".dwp", "roles", filepath.FromSlash(roleName), "command", commandState)
		if info, statErr := os.Stat(rolePath); statErr == nil {
			if info.Mode().IsRegular() && info.Mode()&0o111 != 0 {
				return "exists", rolePath
			}
			return "missing", rolePath
		}
	}
	commandPath = filepath.Join(repoRoot, ".dwp", "command", commandState)
	if info, statErr := os.Stat(commandPath); statErr == nil {
		if info.Mode().IsRegular() && info.Mode()&0o111 != 0 {
			commandStatus = "exists"
		}
	}
	return commandStatus, commandPath
}
