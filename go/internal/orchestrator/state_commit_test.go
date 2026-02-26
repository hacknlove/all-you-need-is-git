package orchestrator

import (
	"strings"
	"testing"
)

func TestBuildStateCommitMessage(t *testing.T) {
	message := buildStateCommitMessage(
		"chore: working",
		"keep lease alive",
		[]stateTrailer{
			{Key: "aynig-state", Value: "working"},
			{Key: "aynig-run-id", Value: "abc"},
		},
	)

	expected := []string{
		"chore: working",
		"",
		"keep lease alive",
		"",
		"aynig-state: working",
		"aynig-run-id: abc",
		"",
	}
	if message != strings.Join(expected, "\n") {
		t.Fatalf("unexpected message:\n%s", message)
	}
}
