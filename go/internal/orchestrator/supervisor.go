package orchestrator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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
	State    string          `json:"state"`
	Subject  string          `json:"subject"`
	Body     string          `json:"body"`
	Trailers []resultTrailer `json:"trailers"`
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
	}
	if !captured.result.valid {
		supervisorLogf(stderrLog, "no valid SET_STATE line found; leaving branch in working")
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
	headCommit, err := gitx.ReadCommitInDir(opts.WorktreePath, "HEAD")
	if err != nil {
		return err
	}
	headState, invalidState := resolveStateTrailer(headCommit.Trailers)
	if invalidState != "" {
		return fmt.Errorf("invalid HEAD state trailer: %s", invalidState)
	}
	if headState != "working" {
		return nil
	}
	if strings.TrimSpace(firstTrailer(headCommit.Trailers["dwp-run-id"], "")) != strings.TrimSpace(opts.RunID) {
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
	for _, trailer := range result.Trailers {
		key := strings.TrimSpace(trailer.Key)
		if key == "" {
			return fmt.Errorf("invalid trailer in SET_STATE payload: empty key")
		}
		if strings.EqualFold(key, "dwp-state") {
			return fmt.Errorf("invalid trailer in SET_STATE payload: dwp-state is managed by state")
		}
		trailers = append(trailers, statex.Trailer{Key: key, Value: strings.TrimSpace(trailer.Value)})
	}

	if err := statex.CommitState(opts.WorktreePath, subject, result.Body, trailers); err != nil {
		return err
	}
	return pushCurrentBranchInDir(opts.WorktreePath, opts.Remote)
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
