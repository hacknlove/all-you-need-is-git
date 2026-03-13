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
	DwpRemote   string
	Trailers    []string
}

func SetState(opts SetStateOptions) error {
	state := strings.ToLower(strings.TrimSpace(opts.State))
	if state == "" {
		return fmt.Errorf("Missing required flag: --dwp-state")
	}
	if state == "working" {
		return fmt.Errorf("Invalid dwp-state: working (use aynig set-working)")
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

	remote := resolveDwpRemote(opts.DwpRemote, headTrailers)
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

	trailers := []statex.Trailer{{Key: "dwp-state", Value: state}}
	if remote != "" {
		trailers = append(trailers, statex.Trailer{Key: "dwp-source", Value: "git:" + remote})
	}

	for _, raw := range opts.Trailers {
		parsed, parseErr := parseTrailerArg(raw)
		if parseErr != nil {
			return parseErr
		}
		if strings.EqualFold(strings.TrimSpace(parsed.Key), "dwp-state") {
			return fmt.Errorf("Invalid trailer: dwp-state is managed by --dwp-state")
		}
		trailers = append(trailers, parsed)
	}

	if err := statex.CommitState("", subject, prompt, trailers); err != nil {
		return err
	}

	return pushCurrentBranch(remote)
}
