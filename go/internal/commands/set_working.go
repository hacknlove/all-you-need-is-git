package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"all-you-need-is-git/go/internal/statex"
)

type SetWorkingOptions struct {
	Subject      string
	Prompt       string
	PromptFile   string
	PromptStdin  bool
	DwpRemote    string
	LeaseSeconds int
	Trailers     []string
}

func SetWorking(opts SetWorkingOptions) error {
	fullMessage, err := gitShowHeadBody()
	if err != nil {
		return err
	}
	_, headBody := splitCommitMessage(fullMessage)
	headTrailers, err := parseTrailersFromBody(headBody)
	if err != nil {
		return err
	}

	remote := resolveDwpRemote(opts.DwpRemote, headTrailers)
	if err := validateRemoteExists(remote); err != nil {
		return err
	}

	prompt, err := readPrompt(promptOptions{Prompt: opts.Prompt, PromptFile: opts.PromptFile, PromptStdin: opts.PromptStdin}, "Lease heartbeat")
	if err != nil {
		return err
	}

	subject := strings.TrimSpace(opts.Subject)
	if subject == "" {
		subject = "chore: working"
	}

	originState := strings.TrimSpace(trailerValue(headTrailers, "dwp-origin-state"))
	if originState == "" {
		originState = strings.TrimSpace(trailerValue(headTrailers, "dwp-state"))
	}
	if originState == "" {
		return fmt.Errorf("Missing required trailer: dwp-state")
	}

	runID := strings.TrimSpace(trailerValue(headTrailers, "dwp-run-id"))
	if runID == "" {
		runID, err = randomRunID()
		if err != nil {
			return err
		}
	}

	runnerID, err := os.Hostname()
	if err != nil {
		return err
	}

	leaseSeconds := parseLeaseSeconds(headTrailers)
	if opts.LeaseSeconds > 0 {
		leaseSeconds = opts.LeaseSeconds
	}

	trailers := []statex.Trailer{
		{Key: "dwp-state", Value: "working"},
		{Key: "dwp-origin-state", Value: originState},
		{Key: "dwp-run-id", Value: runID},
		{Key: "dwp-runner-id", Value: runnerID},
		{Key: "dwp-lease-seconds", Value: strconv.Itoa(leaseSeconds)},
	}
	if remote != "" {
		trailers = append(trailers, statex.Trailer{Key: "dwp-source", Value: "git:" + remote})
	}

	reserved := map[string]struct{}{
		"dwp-state":         {},
		"dwp-origin-state":  {},
		"dwp-run-id":        {},
		"dwp-runner-id":     {},
		"dwp-lease-seconds": {},
		"dwp-source":        {},
	}
	trailers = appendDwpCopiedTrailers(trailers, headTrailers, reserved)

	for _, raw := range opts.Trailers {
		parsed, parseErr := parseTrailerArg(raw)
		if parseErr != nil {
			return parseErr
		}
		trailers = append(trailers, parsed)
	}

	if err := statex.CommitState("", subject, prompt, trailers); err != nil {
		return err
	}

	return pushCurrentBranch(remote)
}
