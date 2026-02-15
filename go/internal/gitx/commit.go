package gitx

import (
	"strings"
)

type CommitMessage struct {
	Message  string
	Body     string
	Date     string
	Trailers map[string][]string
}

func ReadCommit(branch string) (CommitMessage, error) {
	out, err := Run("", "log", branch, "-1", "--pretty=format:%s%x1f%b%x1f%cI")
	if err != nil {
		return CommitMessage{}, err
	}
	parts := strings.Split(out, "\x1f")
	if len(parts) < 3 {
		return CommitMessage{}, nil
	}

	message := strings.TrimRight(parts[0], "\n")
	body := strings.TrimRight(parts[1], "\n")
	date := strings.TrimSpace(parts[2])

	trailers := map[string][]string{}
	if body != "" {
		trailersRaw, err := RunWithInput("", body, "interpret-trailers", "--parse", "--only-trailers")
		if err != nil {
			return CommitMessage{}, err
		}
		parsed, err := ParseTrailersStrict(trailersRaw, nil)
		if err != nil {
			return CommitMessage{}, err
		}
		trailers = parsed
	}

	return CommitMessage{
		Message:  message,
		Body:     body,
		Date:     date,
		Trailers: trailers,
	}, nil
}
