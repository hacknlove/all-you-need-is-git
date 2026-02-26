package commands

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"all-you-need-is-git/go/internal/config"
	"all-you-need-is-git/go/internal/gitx"
	"all-you-need-is-git/go/internal/statex"
)

type promptOptions struct {
	Prompt      string
	PromptFile  string
	PromptStdin bool
}

func readPrompt(opts promptOptions, fallback string) (string, error) {
	count := 0
	if strings.TrimSpace(opts.Prompt) != "" {
		count++
	}
	if strings.TrimSpace(opts.PromptFile) != "" {
		count++
	}
	if opts.PromptStdin {
		count++
	}
	if count > 1 {
		return "", fmt.Errorf("Only one prompt source is allowed: --prompt | --prompt-file | --prompt-stdin")
	}

	if strings.TrimSpace(opts.Prompt) != "" {
		return opts.Prompt, nil
	}
	if strings.TrimSpace(opts.PromptFile) != "" {
		content, err := os.ReadFile(opts.PromptFile)
		if err != nil {
			return "", err
		}
		return string(content), nil
	}
	if opts.PromptStdin {
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		return string(content), nil
	}
	return fallback, nil
}

func parseTrailerArg(raw string) (statex.Trailer, error) {
	idx := strings.Index(raw, ":")
	if idx <= 0 {
		return statex.Trailer{}, fmt.Errorf("Invalid trailer format: %q (expected key:value)", raw)
	}
	key := strings.TrimSpace(raw[:idx])
	value := strings.TrimSpace(raw[idx+1:])
	if key == "" {
		return statex.Trailer{}, fmt.Errorf("Invalid trailer format: %q (empty key)", raw)
	}
	return statex.Trailer{Key: key, Value: value}, nil
}

func resolveAynigRemote(cliRemote string, headTrailers map[string][]string) string {
	if strings.TrimSpace(cliRemote) != "" {
		return strings.TrimSpace(cliRemote)
	}
	return strings.TrimSpace(trailerValue(headTrailers, "aynig-remote"))
}

func validateRemoteExists(remote string) error {
	if strings.TrimSpace(remote) == "" {
		return nil
	}
	if _, err := gitx.Run("", "remote", "get-url", remote); err != nil {
		return fmt.Errorf("Unknown remote %q", remote)
	}
	return nil
}

func pushCurrentBranch(remote string) error {
	if strings.TrimSpace(remote) == "" {
		return nil
	}
	branch, err := gitx.BranchCurrent("")
	if err != nil {
		return err
	}
	return gitx.Push("", remote, branch)
}

func randomRunID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return "run-" + hex.EncodeToString(buf), nil
}

func parseLeaseSeconds(headTrailers map[string][]string) int {
	raw := strings.TrimSpace(trailerValue(headTrailers, "aynig-lease-seconds"))
	if raw == "" {
		return config.Default().LeaseSeconds
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return config.Default().LeaseSeconds
	}
	return value
}

func appendAynigCopiedTrailers(out []statex.Trailer, headTrailers map[string][]string, reserved map[string]struct{}) []statex.Trailer {
	keys := make([]string, 0)
	for key := range headTrailers {
		lower := strings.ToLower(strings.TrimSpace(key))
		if !strings.HasPrefix(lower, "aynig-") {
			continue
		}
		if _, blocked := reserved[lower]; blocked {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value := trailerValue(headTrailers, key)
		if strings.TrimSpace(value) == "" {
			continue
		}
		out = append(out, statex.Trailer{Key: strings.ToLower(strings.TrimSpace(key)), Value: value})
	}
	return out
}
