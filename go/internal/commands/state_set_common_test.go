package commands

import "testing"

func TestResolveDwpRemote(t *testing.T) {
	trailers := map[string][]string{
		"dwp-source": {" git:origin "},
	}
	got := resolveDwpRemote("", trailers)
	if got != "origin" {
		t.Fatalf("unexpected remote: %q", got)
	}

	got = resolveDwpRemote("upstream", trailers)
	if got != "upstream" {
		t.Fatalf("cli remote should win, got %q", got)
	}

	got = resolveDwpRemote("", map[string][]string{"dwp-source": {"http:https://example.com"}})
	if got != "" {
		t.Fatalf("expected empty for non-git locator, got %q", got)
	}
}

func TestParseLeaseSeconds(t *testing.T) {
	trailers := map[string][]string{"dwp-lease-seconds": {"120"}}
	if got := parseLeaseSeconds(trailers); got != 120 {
		t.Fatalf("unexpected lease seconds: %d", got)
	}
}

func TestAppendDwpCopiedTrailers(t *testing.T) {
	reserved := map[string]struct{}{"dwp-state": {}}
	trailers := map[string][]string{
		"dwp-note":  {"hello"},
		"dwp-state": {"build"},
		"custom":    {"skip"},
	}
	got := appendDwpCopiedTrailers(nil, trailers, reserved)
	if len(got) != 1 || got[0].Key != "dwp-note" || got[0].Value != "hello" {
		t.Fatalf("unexpected copied trailers: %+v", got)
	}
}
