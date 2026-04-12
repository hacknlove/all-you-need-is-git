package orchestrator

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"all-you-need-is-git/go/internal/gitx"
	"all-you-need-is-git/go/internal/statex"
)

const setStatePrefix = "SET_STATE "

type SuperviseOptions struct {
	WorktreePath  string
	BranchName    string
	CommandPath   string
	OriginState   string
	RunID         string
	Remote        string
	StdoutLogPath string
	StderrLogPath string
}

type commandResult struct {
	State        string          `json:"state"`
	Subject      string          `json:"subject"`
	Body         string          `json:"body"`
	KeepTrailers bool            `json:"keep_trailers"`
	Trailers     []resultTrailer `json:"trailers"`
}

type resultTrailer struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func Supervise(opts SuperviseOptions) error {
	if strings.TrimSpace(opts.WorktreePath) == "" {
		return fmt.Errorf("missing required flag: --worktree-path")
	}
	if strings.TrimSpace(opts.CommandPath) == "" {
		return fmt.Errorf("missing required flag: --command-path")
	}
	if strings.TrimSpace(opts.RunID) == "" {
		return fmt.Errorf("missing required flag: --run-id")
	}
	if strings.TrimSpace(opts.StdoutLogPath) == "" {
		return fmt.Errorf("missing required flag: --stdout-log-path")
	}
	if strings.TrimSpace(opts.StderrLogPath) == "" {
		return fmt.Errorf("missing required flag: --stderr-log-path")
	}
	if err := os.MkdirAll(filepath.Dir(opts.StdoutLogPath), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(opts.StderrLogPath), 0o755); err != nil {
		return err
	}

	stdoutLog, err := os.OpenFile(opts.StdoutLogPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer stdoutLog.Close()

	stderrLog, err := os.OpenFile(opts.StderrLogPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer stderrLog.Close()

	stdoutReader, stdoutWriter := io.Pipe()
	resultCh := make(chan captureResult, 1)
	go func() {
		result, captureErr := captureSetState(stdoutReader, stdoutLog)
		resultCh <- captureResult{result: result, err: captureErr}
	}()

	cmd := exec.Command(opts.CommandPath)
	cmd.Dir = opts.WorktreePath
	cmd.Env = os.Environ()
	cmd.Stdout = stdoutWriter
	cmd.Stderr = stderrLog
	cmd.Stdin = nil

	if err := cmd.Start(); err != nil {
		_ = stdoutWriter.Close()
		_ = stdoutReader.Close()
		supervisorLogf(stderrLog, "failed to start %s: %v", opts.CommandPath, err)
		<-resultCh
		return err
	}

	waitErr := cmd.Wait()
	_ = stdoutWriter.Close()

	captured := <-resultCh
	if captured.err != nil {
		supervisorLogf(stderrLog, "failed to capture stdout: %v", captured.err)
		return captured.err
	}
	if waitErr != nil {
		supervisorLogf(stderrLog, "command exited with error: %v", waitErr)
		if err := applyFailedCommandResult(opts, waitErr, opts.StdoutLogPath, opts.StderrLogPath); err != nil {
			supervisorLogf(stderrLog, "failed to mark branch as stalled: %v", err)
			return err
		}
		return nil
	}
	if !captured.result.valid {
		supervisorLogf(stderrLog, "no valid SET_STATE line found; refreshing working state")
		if err := refreshWorkingState(opts); err != nil {
			supervisorLogf(stderrLog, "failed to refresh working state: %v", err)
			return err
		}
		return nil
	}

	if err := applyCommandResult(opts, captured.result.value); err != nil {
		supervisorLogf(stderrLog, "failed to apply SET_STATE result: %v", err)
		return err
	}
	return nil
}

type captureResult struct {
	result parsedResult
	err    error
}

type parsedResult struct {
	valid bool
	value commandResult
}

func captureSetState(src io.Reader, log io.Writer) (parsedResult, error) {
	reader := bufio.NewReader(src)
	var latest parsedResult
	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			if _, writeErr := io.WriteString(log, line); writeErr != nil {
				return parsedResult{}, writeErr
			}
			if parsed, ok, parseErr := parseSetStateLine(line); parseErr != nil {
				_, _ = fmt.Fprintf(log, "aynig __supervise: invalid SET_STATE line: %v\n", parseErr)
			} else if ok {
				latest = parsedResult{valid: true, value: parsed}
			}
		}
		if err == io.EOF {
			return latest, nil
		}
		if err != nil {
			return parsedResult{}, err
		}
	}
}

func parseSetStateLine(line string) (commandResult, bool, error) {
	trimmed := strings.TrimRight(line, "\r\n")
	if !strings.HasPrefix(trimmed, setStatePrefix) {
		return commandResult{}, false, nil
	}
	payload := strings.TrimSpace(strings.TrimPrefix(trimmed, setStatePrefix))
	if payload == "" {
		return commandResult{}, true, fmt.Errorf("empty payload")
	}

	var result commandResult
	if err := json.Unmarshal([]byte(payload), &result); err != nil {
		return commandResult{}, true, err
	}
	result.State = strings.ToLower(strings.TrimSpace(result.State))
	result.Subject = strings.TrimSpace(result.Subject)
	return result, true, nil
}

func applyCommandResult(opts SuperviseOptions, result commandResult) error {
	headCommit, ok, err := currentWorkingHead(opts)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	state := strings.ToLower(strings.TrimSpace(result.State))
	if state == "" {
		return fmt.Errorf("missing state in SET_STATE payload")
	}
	if state == "working" {
		return fmt.Errorf("invalid state in SET_STATE payload: working")
	}

	subject := result.Subject
	if subject == "" {
		subject = fmt.Sprintf("chore: set %s", state)
	}

	trailers := []statex.Trailer{{Key: "dwp-state", Value: state}}
	if strings.TrimSpace(opts.Remote) != "" {
		trailers = append(trailers, statex.Trailer{Key: "dwp-source", Value: "git:" + strings.TrimSpace(opts.Remote)})
	}
	if result.KeepTrailers {
		reserved := map[string]struct{}{
			"dwp-state":         {},
			"dwp-source":        {},
			"dwp-origin-state":  {},
			"dwp-run-id":        {},
			"dwp-runner-id":     {},
			"dwp-lease-seconds": {},
		}
		trailers = appendCopiedDwpTrailers(trailers, headCommit.Trailers, reserved)
	}
	for _, trailer := range result.Trailers {
		key := strings.TrimSpace(trailer.Key)
		if key == "" {
			return fmt.Errorf("invalid trailer in SET_STATE payload: empty key")
		}
		if _, blocked := reservedPayloadTrailerKeys()[strings.ToLower(key)]; blocked {
			return fmt.Errorf("invalid trailer in SET_STATE payload: %s is managed by aynig", key)
		}
		parsed := statex.Trailer{Key: key, Value: strings.TrimSpace(trailer.Value)}
		if err := statex.ValidateTrailer(parsed); err != nil {
			return err
		}
		trailers = append(trailers, parsed)
	}

	if err := statex.CommitState(opts.WorktreePath, subject, result.Body, trailers); err != nil {
		return err
	}
	return pushCurrentBranchInDir(opts.WorktreePath, opts.Remote)
}

func reservedPayloadTrailerKeys() map[string]struct{} {
	return map[string]struct{}{
		"dwp-state":         {},
		"dwp-source":        {},
		"dwp-origin-state":  {},
		"dwp-run-id":        {},
		"dwp-runner-id":     {},
		"dwp-lease-seconds": {},
		"dwp-stalled-run":   {},
	}
}

func applyFailedCommandResult(opts SuperviseOptions, waitErr error, stdoutLogPath string, stderrLogPath string) error {
	headCommit, ok, err := currentWorkingHead(opts)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	exitCode := exitCodeFromError(waitErr)
	body := buildFailureBody(exitCode, stdoutLogPath, stderrLogPath)
	originState := strings.TrimSpace(firstTrailer(headCommit.Trailers["dwp-origin-state"], opts.OriginState))
	trailers := []statex.Trailer{
		{Key: "dwp-state", Value: "stalled"},
		{Key: "dwp-stalled-run", Value: strings.TrimSpace(opts.RunID)},
	}
	if originState != "" {
		trailers = append(trailers, statex.Trailer{Key: "dwp-origin-state", Value: originState})
	}
	if strings.TrimSpace(opts.Remote) != "" {
		trailers = append(trailers, statex.Trailer{Key: "dwp-source", Value: "git:" + strings.TrimSpace(opts.Remote)})
	}

	if err := statex.CommitState(opts.WorktreePath, "chore: stalled", body, trailers); err != nil {
		return err
	}
	return pushCurrentBranchInDir(opts.WorktreePath, opts.Remote)
}

func refreshWorkingState(opts SuperviseOptions) error {
	headCommit, ok, err := currentWorkingHead(opts)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	trailers := copyAllTrailers(headCommit.Trailers)
	body := "command ended without changing the state, waiting in case it spawns any other process that will eventually change the state"
	if err := statex.CommitState(opts.WorktreePath, "chore: working", body, trailers); err != nil {
		return err
	}
	return pushCurrentBranchInDir(opts.WorktreePath, opts.Remote)
}

func currentWorkingHead(opts SuperviseOptions) (gitx.CommitMessage, bool, error) {
	headCommit, err := gitx.ReadCommitInDir(opts.WorktreePath, "HEAD")
	if err != nil {
		return gitx.CommitMessage{}, false, err
	}
	headState, invalidState := resolveStateTrailer(headCommit.Trailers)
	if invalidState != "" {
		return gitx.CommitMessage{}, false, fmt.Errorf("invalid HEAD state trailer: %s", invalidState)
	}
	if headState != "working" {
		return gitx.CommitMessage{}, false, nil
	}
	if strings.TrimSpace(firstTrailer(headCommit.Trailers["dwp-run-id"], "")) != strings.TrimSpace(opts.RunID) {
		return gitx.CommitMessage{}, false, nil
	}
	return headCommit, true, nil
}

func exitCodeFromError(err error) int {
	if err == nil {
		return 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}
	return -1
}

func buildFailureBody(exitCode int, stdoutLogPath string, stderrLogPath string) string {
	lines := []string{
		fmt.Sprintf("command exited with code %d", exitCode),
		"",
		"last 5 stdout lines:",
		tailLinesOrPlaceholder(stdoutLogPath, 5),
		"",
		"last 5 stderr lines:",
		tailLinesOrPlaceholder(stderrLogPath, 5),
	}
	return strings.Join(lines, "\n")
}

func tailLinesOrPlaceholder(path string, n int) string {
	lines, err := tailFileLines(path, n)
	if err != nil {
		return fmt.Sprintf("(unable to read log: %v)", err)
	}
	if len(lines) == 0 {
		return "(empty)"
	}
	return strings.Join(lines, "\n")
}

func tailFileLines(path string, n int) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	normalized := strings.ReplaceAll(string(content), "\r\n", "\n")
	normalized = strings.TrimRight(normalized, "\n")
	if normalized == "" {
		return nil, nil
	}
	lines := strings.Split(normalized, "\n")
	if len(lines) <= n {
		return lines, nil
	}
	return lines[len(lines)-n:], nil
}

func copyAllTrailers(trailers map[string][]string) []statex.Trailer {
	keys := make([]string, 0, len(trailers))
	for key := range trailers {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return strings.ToLower(strings.TrimSpace(keys[i])) < strings.ToLower(strings.TrimSpace(keys[j]))
	})

	out := make([]statex.Trailer, 0)
	for _, key := range keys {
		trimmedKey := strings.TrimSpace(key)
		if trimmedKey == "" {
			continue
		}
		for _, value := range trailers[key] {
			out = append(out, statex.Trailer{Key: trimmedKey, Value: strings.TrimSpace(value)})
		}
	}
	return out
}

func appendCopiedDwpTrailers(out []statex.Trailer, headTrailers map[string][]string, reserved map[string]struct{}) []statex.Trailer {
	keys := make([]string, 0, len(headTrailers))
	for key := range headTrailers {
		lower := strings.ToLower(strings.TrimSpace(key))
		if !strings.HasPrefix(lower, "dwp-") {
			continue
		}
		if _, blocked := reserved[lower]; blocked {
			continue
		}
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return strings.ToLower(strings.TrimSpace(keys[i])) < strings.ToLower(strings.TrimSpace(keys[j]))
	})
	for _, key := range keys {
		normalizedKey := strings.ToLower(strings.TrimSpace(key))
		for _, value := range headTrailers[key] {
			if strings.TrimSpace(value) == "" {
				continue
			}
			out = append(out, statex.Trailer{Key: normalizedKey, Value: strings.TrimSpace(value)})
		}
	}
	return out
}

func pushCurrentBranchInDir(dir string, remote string) error {
	if strings.TrimSpace(remote) == "" {
		return nil
	}
	branch, err := gitx.BranchCurrent(dir)
	if err != nil {
		return err
	}
	return gitx.Push(dir, remote, branch)
}

func supervisorLogf(log io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(log, "aynig __supervise: "+format+"\n", args...)
}
