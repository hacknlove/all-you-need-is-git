package gitx

import (
	"errors"
	"strconv"
	"strings"
)

type TrailerOptions struct {
	Separators   []string
	LowerCaseKey bool
}

func ParseTrailersStrict(raw string, opts *TrailerOptions) (map[string][]string, error) {
	if opts == nil {
		opts = &TrailerOptions{}
	}
	separators := opts.Separators
	if len(separators) == 0 {
		separators = []string{":", "="}
	}

	if raw == "" {
		return map[string][]string{}, nil
	}

	normalized := strings.ReplaceAll(raw, "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")
	lines := strings.Split(normalized, "\n")
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) == 0 {
		return map[string][]string{}, nil
	}

	out := map[string][]string{}
	currentKey := ""

	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			return nil, errors.New("Invalid trailers: blank line at " + strconv.Itoa(i+1))
		}

		if len(line) > 0 && (line[0] == ' ' || line[0] == '\t') {
			if currentKey == "" {
				return nil, errors.New("Invalid trailers: continuation without a preceding trailer at " + strconv.Itoa(i+1))
			}
			cont := strings.TrimLeft(line, " \t")
			values := out[currentKey]
			if len(values) == 0 {
				return nil, errors.New("Invalid trailers: continuation without a preceding trailer at " + strconv.Itoa(i+1))
			}
			values[len(values)-1] = values[len(values)-1] + "\n" + cont
			out[currentKey] = values
			continue
		}

		sepIdx := -1
		sep := ""
		for _, s := range separators {
			idx := strings.Index(line, s)
			if idx != -1 && (sepIdx == -1 || idx < sepIdx) {
				sepIdx = idx
				sep = s
			}
		}
		if sepIdx <= 0 {
			return nil, errors.New("Invalid trailers: expected \"token" + strings.Join(separators, "|") + "value\" at " + strconv.Itoa(i+1) + ": \"" + line + "\"")
		}

		key := strings.TrimSpace(line[:sepIdx])
		if key == "" {
			return nil, errors.New("Invalid trailers: empty token at " + strconv.Itoa(i+1))
		}
		if opts.LowerCaseKey {
			key = strings.ToLower(key)
		}
		value := line[sepIdx+len(sep):]
		value = strings.TrimLeft(value, " \t")

		currentKey = key
		out[key] = append(out[key], value)
	}

	return out, nil
}
