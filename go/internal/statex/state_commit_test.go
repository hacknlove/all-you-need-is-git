package statex

import (
	"strings"
	"testing"
)

func TestBuildCommitMessage(t *testing.T) {
	message := BuildCommitMessage("chore: working", "keep lease alive", []Trailer{
		{Key: "dwp-state", Value: "working"},
		{Key: "dwp-run-id", Value: "abc"},
	})

	expected := []string{
		"chore: working",
		"",
		"keep lease alive",
		"",
		"dwp-state: working",
		"dwp-run-id: abc",
		"",
	}
	if message != strings.Join(expected, "\n") {
		t.Fatalf("unexpected message:\n%s", message)
	}
}
