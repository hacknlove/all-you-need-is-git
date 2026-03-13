package orchestrator

import "testing"

func TestResolveStateTrailer(t *testing.T) {
	tests := []struct {
		name       string
		trailers   map[string][]string
		wantState  string
		wantReason string
	}{
		{
			name:       "missing state",
			trailers:   map[string][]string{},
			wantState:  "",
			wantReason: "",
		},
		{
			name:       "single state",
			trailers:   map[string][]string{"dwp-state": []string{" Build "}},
			wantState:  "build",
			wantReason: "",
		},
		{
			name:       "multiple state trailers (last wins)",
			trailers:   map[string][]string{"dwp-state": []string{"build", "review"}},
			wantState:  "review",
			wantReason: "",
		},
		{
			name:       "empty state trailer",
			trailers:   map[string][]string{"dwp-state": []string{"  "}},
			wantState:  "",
			wantReason: "empty dwp-state trailer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotState, gotReason := resolveStateTrailer(tt.trailers)
			if gotState != tt.wantState {
				t.Fatalf("state mismatch: got %q want %q", gotState, tt.wantState)
			}
			if gotReason != tt.wantReason {
				t.Fatalf("reason mismatch: got %q want %q", gotReason, tt.wantReason)
			}
		})
	}
}
