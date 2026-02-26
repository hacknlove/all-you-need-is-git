package orchestrator

import (
	"fmt"
	"strings"

	"all-you-need-is-git/go/internal/gitx"
)

type stateTrailer struct {
	Key   string
	Value string
}

func buildStateCommitMessage(subject string, prompt string, trailers []stateTrailer) string {
	cleanSubject := strings.TrimRight(subject, "\n")
	cleanPrompt := strings.TrimRight(prompt, "\n")

	parts := []string{cleanSubject, "", cleanPrompt, ""}
	for _, trailer := range trailers {
		parts = append(parts, fmt.Sprintf("%s: %s", trailer.Key, trailer.Value))
	}
	return strings.Join(parts, "\n") + "\n"
}

func commitState(worktreePath string, subject string, prompt string, trailers []stateTrailer) error {
	message := buildStateCommitMessage(subject, prompt, trailers)
	return gitx.Commit(worktreePath, message, true)
}
