package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"all-you-need-is-git/go/internal/statex"
)

type SetWorkingOptions struct {
	Subject     string
	Prompt      string
	PromptFile  string
	PromptStdin bool
	AynigRemote string
	Trailers    []string
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

	remote := resolveAynigRemote(opts.AynigRemote, headTrailers)
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

	originState := strings.TrimSpace(trailerValue(headTrailers, "aynig-origin-state"))
	if originState == "" {
		originState = strings.TrimSpace(trailerValue(headTrailers, "aynig-state"))
	}
	if originState == "" {
		return fmt.Errorf("Missing required trailer: aynig-state")
	}

	runID := strings.TrimSpace(trailerValue(headTrailers, "aynig-run-id"))
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

	trailers := []statex.Trailer{
		{Key: "aynig-state", Value: "working"},
		{Key: "aynig-origin-state", Value: originState},
		{Key: "aynig-run-id", Value: runID},
		{Key: "aynig-runner-id", Value: runnerID},
		{Key: "aynig-lease-seconds", Value: strconv.Itoa(leaseSeconds)},
	}
	if remote != "" {
		trailers = append(trailers, statex.Trailer{Key: "aynig-remote", Value: remote})
	}

	reserved := map[string]struct{}{
		"aynig-state":         {},
		"aynig-origin-state":  {},
		"aynig-run-id":        {},
		"aynig-runner-id":     {},
		"aynig-lease-seconds": {},
		"aynig-remote":        {},
	}
	trailers = appendAynigCopiedTrailers(trailers, headTrailers, reserved)

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
