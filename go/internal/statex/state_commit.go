package statex

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"all-you-need-is-git/go/internal/gitx"
)

type Trailer struct {
	Key   string
	Value string
}

var trailerKeyPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`)

func ValidateTrailer(trailer Trailer) error {
	key := strings.TrimSpace(trailer.Key)
	if key == "" {
		return errors.New("Invalid trailer: empty key")
	}
	if !trailerKeyPattern.MatchString(key) {
		return fmt.Errorf("Invalid trailer key: %q", trailer.Key)
	}
	if strings.Contains(trailer.Value, "\n") || strings.Contains(trailer.Value, "\r") {
		return fmt.Errorf("Invalid trailer value for %q: newlines are not allowed", key)
	}
	return nil
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
	fullMessage, err := gitx.Run(dir, "show", "-s", "--format=%B", "HEAD")
	if err != nil {
		return err
	}

	parsed, err := gitx.ParseCommitTrailers(dir, fullMessage)
	if err != nil {
		return err
	}

	stateValues := valuesForKey(parsed, "dwp-state")
	if len(stateValues) == 0 {
		return errors.New("Invalid trailer block: missing dwp-state trailer")
	}
	last := stateValues[len(stateValues)-1]
	if strings.TrimSpace(last) == "" {
		return errors.New("Invalid trailer block: empty dwp-state trailer")
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
