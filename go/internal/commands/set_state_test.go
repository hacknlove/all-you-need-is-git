package commands

import "testing"

func TestSetStateRequiresState(t *testing.T) {
	err := SetState(SetStateOptions{})
	if err == nil || err.Error() != "Missing required flag: --aynig-state" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSetStateRejectsWorking(t *testing.T) {
	err := SetState(SetStateOptions{State: "working"})
	if err == nil || err.Error() != "Invalid aynig-state: working (use aynig set-working)" {
		t.Fatalf("unexpected error: %v", err)
	}
}
