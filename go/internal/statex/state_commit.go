package statex

import (
	"errors"
	"fmt"
	"strings"

	"all-you-need-is-git/go/internal/gitx"
)

type Trailer struct {
	Key   string
	Value string
}

func BuildCommitMessage(subject string, prompt string, trailers []Trailer) string {
	cleanSubject := strings.TrimRight(subject, "\n")
	cleanPrompt := strings.TrimRight(prompt, "\n")

	parts := []string{cleanSubject, "", cleanPrompt, ""}
	for _, trailer := range trailers {
		parts = append(parts, fmt.Sprintf("%s: %s", trailer.Key, trailer.Value))
	}
	return strings.Join(parts, "\n") + "\n"
}

func CommitState(dir string, subject string, prompt string, trailers []Trailer) error {
	message := BuildCommitMessage(subject, prompt, trailers)
	if err := gitx.Commit(dir, message, true); err != nil {
		return err
	}
	return VerifyHeadStateTrailer(dir)
}

func VerifyHeadStateTrailer(dir string) error {
	body, err := gitx.Run(dir, "show", "-s", "--format=%B", "HEAD")
	if err != nil {
		return err
	}

	trailersRaw, err := gitx.RunWithInput(dir, body, "interpret-trailers", "--parse", "--only-trailers")
	if err != nil {
		return err
	}
	parsed, err := gitx.ParseTrailersStrict(trailersRaw, nil)
	if err != nil {
		return err
	}

	stateValues := valuesForKey(parsed, "aynig-state")
	if len(stateValues) == 0 {
		return errors.New("Invalid trailer block: trailers must be contiguous at end of message")
	}
	if len(stateValues) > 1 {
		return errors.New("Invalid trailer block: multiple aynig-state trailers")
	}
	if strings.TrimSpace(stateValues[0]) == "" {
		return errors.New("Invalid trailer block: empty aynig-state trailer")
	}

	return nil
}

func valuesForKey(trailers map[string][]string, key string) []string {
	want := strings.ToLower(strings.TrimSpace(key))
	for k, values := range trailers {
		if strings.ToLower(strings.TrimSpace(k)) == want {
			return values
		}
	}
	return nil
}
