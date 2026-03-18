package commands

import (
	"fmt"
	"os"
	"path/filepath"
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
	Role string
}

func Status(options StatusOptions) error {
	repoRoot, err := gitx.Run("", "rev-parse", "--show-toplevel")
	if err != nil {
		return fmt.Errorf("Error: Not a Git repository. Please run `git init` first.")
	}
	repoRoot = strings.TrimSpace(repoRoot)

	branch, err := gitx.Run("", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return err
	}
	branch = strings.TrimSpace(branch)

	headCommit, err := gitx.Run("", "rev-parse", "HEAD")
	if err != nil {
		return err
	}
	headCommit = strings.TrimSpace(headCommit)

	committerDate, err := gitx.Run("", "log", "-1", "--format=%cI")
	if err != nil {
		return err
	}
	committerDate = strings.TrimSpace(committerDate)

	fullMessage, err := gitx.Run("", "log", "-1", "--format=%B")
	if err != nil {
		return err
	}
	firstLine, body := splitCommitMessage(fullMessage)
	_ = firstLine
	trailers, err := parseTrailersFromBody(body)
	if err != nil {
		return err
	}

	state := trailerValue(trailers, "dwp-state")
	runID := trailerValue(trailers, "dwp-run-id")
	leaseSecondsRaw := trailerValue(trailers, "dwp-lease-seconds")
	originState := trailerValue(trailers, "dwp-origin-state")

	leaseStatus := leaseStatusForState(state, leaseSecondsRaw, committerDate)

	commandStatus := "missing"
	commandPath := ""
	commandState := state
	shouldResolveCommand := true
	if state == "working" && originState != "" {
		commandState = originState
	} else if state == "working" {
		shouldResolveCommand = false
	}

	roleName := strings.TrimSpace(options.Role)
	if roleName == "" {
		roleName = strings.TrimSpace(os.Getenv("AYNIG_ROLE"))
	}

	if shouldResolveCommand && commandState != "" && commandState != "working" {
		commandStatus, commandPath = resolveCommandPath(repoRoot, roleName, commandState)
	} else if !shouldResolveCommand {
		commandStatus = "lease"
	}

	fmt.Printf("branch: %s\n", branch)
	fmt.Printf("head: %s\n", headCommit)
	if state != "" {
		fmt.Printf("dwp-state: %s\n", state)
	} else {
		fmt.Printf("dwp-state: n/a\n")
	}
	if state == "working" && originState != "" {
		fmt.Printf("dwp-origin-state: %s\n", originState)
	}
	if runID != "" {
		fmt.Printf("dwp-run-id: %s\n", runID)
	} else {
		fmt.Printf("dwp-run-id: n/a\n")
	}
	fmt.Printf("lease: %s\n", leaseStatus)
	fmt.Printf("command: %s\n", commandStatus)
	if commandPath != "" {
		fmt.Printf("command-path: %s\n", commandPath)
	}
	return nil
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
