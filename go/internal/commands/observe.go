package commands

import (
	"strings"

	"all-you-need-is-git/go/internal/gitx"
)

func splitCommitMessage(full string) (string, string) {
	normalized := strings.ReplaceAll(full, "\r\n", "\n")
	lines := strings.Split(normalized, "\n")
	firstLine := ""
	if len(lines) > 0 {
		firstLine = lines[0]
		lines = lines[1:]
	}
	body := strings.Join(lines, "\n")
	body = strings.TrimLeft(body, "\n")
	return firstLine, body
}

func parseTrailersFromMessage(fullMessage string) (map[string][]string, error) {
	if strings.TrimSpace(fullMessage) == "" {
		return map[string][]string{}, nil
	}
	return gitx.ParseCommitTrailers("", fullMessage)
}

func trailerValue(trailers map[string][]string, key string) string {
	if trailers == nil {
		return ""
	}
	want := strings.ToLower(strings.TrimSpace(key))
	for k, values := range trailers {
		if strings.ToLower(strings.TrimSpace(k)) != want {
			continue
		}
		if len(values) == 0 {
			return ""
		}
		return strings.TrimSpace(values[len(values)-1])
	}
	return ""
}
