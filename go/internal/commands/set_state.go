package commands

import (
	"fmt"
	"strings"

	"all-you-need-is-git/go/internal/statex"
)

type SetStateOptions struct {
	State       string
	Subject     string
	Prompt      string
	PromptFile  string
	PromptStdin bool
	AynigRemote string
	Trailers    []string
}

func SetState(opts SetStateOptions) error {
	state := strings.ToLower(strings.TrimSpace(opts.State))
	if state == "" {
		return fmt.Errorf("Missing required flag: --aynig-state")
	}
	if state == "working" {
		return fmt.Errorf("Invalid aynig-state: working (use aynig set-working)")
	}

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

	prompt, err := readPrompt(promptOptions{Prompt: opts.Prompt, PromptFile: opts.PromptFile, PromptStdin: opts.PromptStdin}, "")
	if err != nil {
		return err
	}

	subject := strings.TrimSpace(opts.Subject)
	if subject == "" {
		subject = fmt.Sprintf("chore: set %s", state)
	}

	trailers := []statex.Trailer{{Key: "aynig-state", Value: state}}
	if remote != "" {
		trailers = append(trailers, statex.Trailer{Key: "aynig-remote", Value: remote})
	}

	for _, raw := range opts.Trailers {
		parsed, parseErr := parseTrailerArg(raw)
		if parseErr != nil {
			return parseErr
		}
		if strings.EqualFold(strings.TrimSpace(parsed.Key), "aynig-state") {
			return fmt.Errorf("Invalid trailer: aynig-state is managed by --aynig-state")
		}
		trailers = append(trailers, parsed)
	}

	if err := statex.CommitState("", subject, prompt, trailers); err != nil {
		return err
	}

	return pushCurrentBranch(remote)
}
