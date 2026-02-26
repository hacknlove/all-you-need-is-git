package orchestrator

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"all-you-need-is-git/go/internal/config"
	"all-you-need-is-git/go/internal/gitx"
	"all-you-need-is-git/go/internal/logx"
)

type CommandParams struct {
	Config          config.Config
	BranchName      string
	IsCurrentBranch bool
	Trailers        map[string][]string
	Body            string
	CommitDate      string
	LogLevel        string
}

type Command struct {
	config          config.Config
	branchName      string
	isCurrentBranch bool
	command         string
	invalidState    string
	trailers        map[string][]string
	body            string
	commitDate      string
	logger          logx.Logger
	logLevel        string
}

func NewCommand(params CommandParams) *Command {
	state, invalidState := resolveStateTrailer(params.Trailers)
	level := params.LogLevel
	if level == "" {
		level = params.Config.LogLevel
	}
	return &Command{
		config:          params.Config,
		branchName:      params.BranchName,
		isCurrentBranch: params.IsCurrentBranch,
		command:         state,
		invalidState:    invalidState,
		trailers:        params.Trailers,
		body:            params.Body,
		commitDate:      params.CommitDate,
		logger:          logx.New(level),
		logLevel:        level,
	}
}

func (c *Command) Run() error {
	if c.invalidState != "" {
		c.logger.Warnf("Skipping branch %s (%s)", c.branchName, c.invalidState)
		return nil
	}
	if c.command == "" {
		c.logger.Debugf("Skipping branch %s (no aynig-state)", c.branchName)
		return nil
	}
	if c.command == "working" {
		c.logger.Debugf("Checking lease on branch %s", c.branchName)
		return c.checkWorking()
	}
	c.logger.Infof("Running command %s on branch %s", c.command, c.branchName)

	worktreePath, err := c.getWorkspace()
	if err != nil || worktreePath == "" {
		return err
	}
	commandPath, err := c.getCommandPath(worktreePath)
	if err != nil || commandPath == "" {
		c.logger.Debugf("Command path not found for %s", c.command)
		return err
	}
	c.logger.Debugf("Command path: %s", commandPath)

	leaseSeconds := c.config.LeaseSeconds
	if leaseSeconds <= 0 {
		leaseSeconds = 300
	}
	runID, err := uuidV4()
	if err != nil {
		return err
	}
	runnerID, err := os.Hostname()
	if err != nil {
		return err
	}

	currentCommitHash, err := gitx.RevParse(worktreePath, "HEAD")
	if err != nil {
		return err
	}

	commandLog, logPath, err := prepareCommandLogFile(worktreePath, currentCommitHash)
	if err != nil {
		return err
	}

	workingTrailers := []stateTrailer{
		{Key: "aynig-state", Value: "working"},
		{Key: "aynig-origin-state", Value: c.command},
		{Key: "aynig-run-id", Value: runID},
		{Key: "aynig-runner-id", Value: runnerID},
		{Key: "aynig-lease-seconds", Value: strconv.Itoa(leaseSeconds)},
	}
	if c.config.UseRemote != "" {
		workingTrailers = append(workingTrailers, stateTrailer{Key: "aynig-remote", Value: c.config.UseRemote})
	}
	if err := commitState(worktreePath, "chore: working", fmt.Sprintf("command %s takes control of the branch", c.command), workingTrailers); err != nil {
		return err
	}
	c.logger.Debugf("Created working commit for %s", c.branchName)

	if c.config.UseRemote != "" {
		c.logger.Infof("Pushing branch %s to %s", c.branchName, c.config.UseRemote)
		if err := gitx.Push(worktreePath, c.config.UseRemote, c.branchName); err != nil {
			return nil
		}
	}

	env := append([]string{}, os.Environ()...)
	env = append(env, "AYNIG_BODY="+c.body)
	env = append(env, "AYNIG_COMMIT_HASH="+currentCommitHash)
	if c.logLevel != "" {
		env = append(env, "AYNIG_LOG_LEVEL="+c.logLevel)
	}
	for key, values := range c.trailers {
		upperKey := strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
		envValue := strings.Join(values, ",")
		env = append(env, "AYNIG_TRAILER_"+upperKey+"="+envValue)
	}

	cmd := exec.Command(commandPath)
	cmd.Dir = worktreePath
	cmd.Env = env
	cmd.Stdout = commandLog
	cmd.Stderr = commandLog
	cmd.Stdin = nil
	setDetached(cmd)
	if err := cmd.Start(); err != nil {
		_ = commandLog.Close()
		return err
	}
	_ = commandLog.Close()
	c.logger.Infof("Launched %s in %s", c.command, worktreePath)
	c.logger.Debugf("Command log: %s", logPath)
	_ = cmd.Process.Release()
	return nil
}

func (c *Command) checkWorking() error {
	leaseSeconds := parseIntTrailer(c.trailers["aynig-lease-seconds"])
	if leaseSeconds <= 0 {
		return nil
	}
	committedAt, err := time.Parse(time.RFC3339, c.commitDate)
	if err != nil {
		return nil
	}
	if time.Now().Before(committedAt.Add(time.Duration(leaseSeconds) * time.Second)) {
		return nil
	}
	c.logger.Infof("Lease expired for branch %s", c.branchName)
	stalledRun := firstTrailer(c.trailers["aynig-run-id"], "unknown")
	worktreePath, err := c.getWorkspace()
	if err != nil || worktreePath == "" {
		return err
	}

	stalledTrailers := []stateTrailer{
		{Key: "aynig-state", Value: "stalled"},
		{Key: "aynig-stalled-run", Value: stalledRun},
	}
	if c.config.UseRemote != "" {
		stalledTrailers = append(stalledTrailers, stateTrailer{Key: "aynig-remote", Value: c.config.UseRemote})
	}
	if err := commitState(worktreePath, "chore: stalled", "Lease expired", stalledTrailers); err != nil {
		return err
	}
	if c.config.UseRemote != "" {
		c.logger.Infof("Pushing stalled state for %s to %s", c.branchName, c.config.UseRemote)
		if err := gitx.Push(worktreePath, c.config.UseRemote, c.branchName); err != nil {
			return nil
		}
	}
	return nil
}

func (c *Command) getWorkspace() (string, error) {
	if c.isCurrentBranch {
		c.logger.Debugf("Using current working directory for %s", c.branchName)
		return os.Getwd()
	}

	worktrees, err := gitx.WorktreeList(c.config.RepoRoot)
	if err != nil {
		return "", err
	}
	ref := "refs/heads/" + c.branchName
	for _, wt := range worktrees {
		if wt.Branch == ref {
			c.logger.Debugf("Using existing worktree for %s", c.branchName)
			return wt.Path, nil
		}
	}

	safeName := strings.ReplaceAll(c.branchName, "/", "_")
	hash := branchHash(c.branchName)
	baseDir := c.config.WorkTree
	if baseDir == "" {
		baseDir = "."
	}
	worktreePath := filepath.Join(c.config.RepoRoot, baseDir, "worktree-"+safeName+"-"+hash)
	if c.config.UseRemote != "" {
		if err := gitx.WorktreeAdd(c.config.RepoRoot, "-b", c.branchName, worktreePath, c.config.UseRemote+"/"+c.branchName); err != nil {
			c.logger.Warnf("Failed to create worktree for branch %s", c.branchName)
			return "", nil
		}
	} else {
		if err := gitx.WorktreeAdd(c.config.RepoRoot, worktreePath, c.branchName); err != nil {
			c.logger.Warnf("Failed to create worktree for branch %s", c.branchName)
			return "", nil
		}
	}
	c.logger.Debugf("Created worktree for %s at %s", c.branchName, worktreePath)

	return worktreePath, nil
}

func (c *Command) getCommandPath(worktreePath string) (string, error) {
	if c.command == "" {
		return "", nil
	}
	baseDir := filepath.Join(worktreePath, ".aynig", "command")
	baseDirAbs, err := filepath.Abs(baseDir)
	if err != nil {
		return "", err
	}
	commandPath := filepath.Join(baseDirAbs, filepath.FromSlash(c.command))
	if !strings.HasPrefix(commandPath, baseDirAbs+string(os.PathSeparator)) {
		return "", nil
	}
	info, err := os.Stat(commandPath)
	if err != nil {
		return "", nil
	}
	if info.Mode()&0o111 == 0 {
		return "", nil
	}
	return commandPath, nil
}

func branchHash(name string) string {
	sum := sha256.Sum256([]byte(name))
	return hex.EncodeToString(sum[:])[:8]
}

func uuidV4() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

func parseIntTrailer(values []string) int {
	if len(values) == 0 {
		return 0
	}
	value := strings.TrimSpace(values[0])
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return parsed
}

func firstTrailer(values []string, fallback string) string {
	if len(values) == 0 {
		return fallback
	}
	value := strings.TrimSpace(values[0])
	if value == "" {
		return fallback
	}
	return value
}

func resolveStateTrailer(trailers map[string][]string) (string, string) {
	values, ok := trailers["aynig-state"]
	if !ok || len(values) == 0 {
		return "", ""
	}
	if len(values) > 1 {
		return "", "multiple aynig-state trailers"
	}
	state := strings.ToLower(strings.TrimSpace(values[0]))
	if state == "" {
		return "", "empty aynig-state trailer"
	}
	return state, ""
}

func prepareCommandLogFile(worktreePath string, commitHash string) (*os.File, string, error) {
	logsDir := filepath.Join(worktreePath, ".aynig", "logs")
	if err := os.MkdirAll(logsDir, 0o755); err != nil {
		return nil, "", err
	}
	logPath := filepath.Join(logsDir, commitHash+".log")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, "", err
	}
	return file, logPath, nil
}
