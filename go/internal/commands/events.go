package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"all-you-need-is-git/go/internal/gitx"
)

type EventsOptions struct {
	History bool
	Limit   int
	JSON    bool
}

type Event struct {
	Commit   string              `json:"commit"`
	Date     string              `json:"date"`
	Message  string              `json:"message"`
	State    *string             `json:"state"`
	RunID    *string             `json:"runId"`
	Origin   *string             `json:"originState"`
	Trailers map[string][]string `json:"trailers"`
}

func Events(opts EventsOptions) error {
	if _, err := gitx.Run("", "rev-parse", "--show-toplevel"); err != nil {
		return fmt.Errorf("Error: Not a Git repository. Please run `git init` first.")
	}

	limit := 1
	if opts.History {
		limit = opts.Limit
		if limit < 1 {
			limit = 1
		}
	}

	raw, err := gitx.Run("", "log", "-n", strconv.Itoa(limit), "--format=%H%x1f%cI%x1f%B%x1e")
	if err != nil {
		return err
	}

	records := strings.Split(raw, "\x1e")
	events := make([]Event, 0)
	for _, record := range records {
		if strings.TrimSpace(record) == "" {
			continue
		}
		parts := strings.Split(record, "\x1f")
		commit := ""
		date := ""
		message := ""
		if len(parts) > 0 {
			commit = parts[0]
		}
		if len(parts) > 1 {
			date = parts[1]
		}
		if len(parts) > 2 {
			message = strings.Join(parts[2:], "\x1f")
		}

		firstLine, body := splitCommitMessage(message)
		trailers, parseErr := parseTrailersFromBody(body)
		if parseErr != nil {
			return parseErr
		}
		stateValue := trailerValue(trailers, "aynig-state")
		runIDValue := trailerValue(trailers, "aynig-run-id")
		originValue := trailerValue(trailers, "aynig-origin-state")

		var statePtr *string
		if stateValue != "" {
			statePtr = &stateValue
		}
		var runIDPtr *string
		if runIDValue != "" {
			runIDPtr = &runIDValue
		}
		var originPtr *string
		if originValue != "" {
			originPtr = &originValue
		}

		events = append(events, Event{
			Commit:   commit,
			Date:     date,
			Message:  firstLine,
			State:    statePtr,
			RunID:    runIDPtr,
			Origin:   originPtr,
			Trailers: trailers,
		})
	}

	if opts.JSON {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(events)
	}

	if len(events) == 0 {
		fmt.Println("No commits found.")
		return nil
	}

	for _, event := range events {
		shortCommit := "unknown"
		if event.Commit != "" {
			if len(event.Commit) > 7 {
				shortCommit = event.Commit[:7]
			} else {
				shortCommit = event.Commit
			}
		}
		state := "n/a"
		if event.State != nil {
			state = *event.State
		}
		runID := "n/a"
		if event.RunID != nil {
			runID = *event.RunID
		}
		origin := ""
		if event.Origin != nil {
			origin = " origin=" + *event.Origin
		}
		fmt.Printf("%s %s state=%s run=%s%s %s\n", shortCommit, event.Date, state, runID, origin, event.Message)
	}

	return nil
}
