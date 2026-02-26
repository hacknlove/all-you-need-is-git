package statex

import (
	"strings"
	"testing"
)

func TestBuildCommitMessage(t *testing.T) {
	message := BuildCommitMessage("chore: working", "keep lease alive", []Trailer{
		{Key: "aynig-state", Value: "working"},
		{Key: "aynig-run-id", Value: "abc"},
	})

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
