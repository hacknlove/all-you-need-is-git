package gitx

import "testing"

func TestParseTrailersStrictBasic(t *testing.T) {
	parsed, err := ParseTrailersStrict("Token: value", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if parsed["Token"][0] != "value" {
		t.Fatalf("expected value, got %q", parsed["Token"][0])
	}
}

func TestParseTrailersStrictRepeated(t *testing.T) {
	parsed, err := ParseTrailersStrict("Reviewed-by: A\nReviewed-by: B", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(parsed["Reviewed-by"]) != 2 {
		t.Fatalf("expected 2 values, got %d", len(parsed["Reviewed-by"]))
	}
}

func TestParseTrailersStrictContinuation(t *testing.T) {
	parsed, err := ParseTrailersStrict("Notes: first line\n  second line", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if parsed["Notes"][0] != "first line\nsecond line" {
		t.Fatalf("unexpected continuation: %q", parsed["Notes"][0])
	}
}

func TestParseTrailersStrictBlankLine(t *testing.T) {
	if _, err := ParseTrailersStrict("A: 1\n\nB: 2", nil); err == nil {
		t.Fatalf("expected error on blank line")
	}
}

func TestParseTrailersStrictCustomSeparators(t *testing.T) {
	parsed, err := ParseTrailersStrict("Key=left:right", &TrailerOptions{Separators: []string{":", "="}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if parsed["Key"][0] != "left:right" {
		t.Fatalf("expected left:right, got %q", parsed["Key"][0])
	}
}
