package orchestrator

import "testing"

func TestFirstLowerTrailerValue(t *testing.T) {
	tests := []struct {
		name     string
		trailers map[string][]string
		want     string
	}{
		{name: "missing", trailers: map[string][]string{}, want: ""},
		{name: "case insensitive", trailers: map[string][]string{"DWP-SOURCE": []string{"git:origin"}}, want: "git:origin"},
		{name: "trims", trailers: map[string][]string{"dwp-source": []string{"  git:upstream  "}}, want: "git:upstream"},
		{name: "empty values", trailers: map[string][]string{"dwp-source": []string{}}, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := firstLowerTrailerValue(tt.trailers, "dwp-source")
			if got != tt.want {
				t.Fatalf("unexpected value: got %q want %q", got, tt.want)
			}
		})
	}
}
