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
	"all-you-need-is-git/go/internal/statex"
)

var currentExecutablePath = os.Executable

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
		c.logger.Debugf("Skipping branch %s (no dwp-state)", c.branchName)
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
		c.logger.Warnf("Command path not found for %s on branch %s", c.command, c.branchName)
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

	stdoutLogPath, stderrLogPath, err := prepareCommandLogPaths(worktreePath, currentCommitHash)
	if err != nil {
		return err
	}

	workingTrailers := []statex.Trailer{
		{Key: "dwp-state", Value: "working"},
		{Key: "dwp-origin-state", Value: c.command},
		{Key: "dwp-run-id", Value: runID},
		{Key: "dwp-runner-id", Value: runnerID},
		{Key: "dwp-lease-seconds", Value: strconv.Itoa(leaseSeconds)},
	}
	if c.config.UseRemote != "" {
		workingTrailers = append(workingTrailers, statex.Trailer{Key: "dwp-source", Value: "git:" + c.config.UseRemote})
	}
	if err := statex.CommitState(worktreePath, "chore: working", fmt.Sprintf("command %s takes control of the branch", c.command), workingTrailers); err != nil {
		return err
	}
	c.logger.Debugf("Created working commit for %s", c.branchName)

	if c.config.UseRemote != "" {
		c.logger.Infof("Pushing branch %s to %s", c.branchName, c.config.UseRemote)
		if err := gitx.Push(worktreePath, c.config.UseRemote, c.branchName); err != nil {
			return nil
		}
	}

	env := c.commandEnv(currentCommitHash, stdoutLogPath, stderrLogPath, worktreePath)
	executablePath, err := currentExecutablePath()
	if err != nil {
		return err
	}

	cmdArgs := []string{
		"__supervise",
		"--worktree-path", worktreePath,
		"--branch", c.branchName,
		"--command-path", commandPath,
		"--origin-state", c.command,
		"--run-id", runID,
		"--stdout-log-path", stdoutLogPath,
		"--stderr-log-path", stderrLogPath,
	}
	if c.config.UseRemote != "" {
		cmdArgs = append(cmdArgs, "--remote", c.config.UseRemote)
	}

	cmd := exec.Command(executablePath, cmdArgs...)
	cmd.Dir = worktreePath
	cmd.Env = env
	cmd.Stdin = nil
	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	cmd.Stdout = devNull
	cmd.Stderr = devNull
	setDetached(cmd)
	if err := cmd.Start(); err != nil {
		_ = devNull.Close()
		return err
	}
	_ = devNull.Close()
	c.logger.Infof("Launched %s in %s", c.command, worktreePath)
	c.logger.Debugf("Command stdout log: %s", stdoutLogPath)
	c.logger.Debugf("Command stderr log: %s", stderrLogPath)
	_ = cmd.Process.Release()
	return nil
}

func (c *Command) checkWorking() error {
	leaseSeconds := parseIntTrailer(c.trailers["dwp-lease-seconds"])
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
	stalledRun := firstTrailer(c.trailers["dwp-run-id"], "unknown")
	originState := firstTrailer(c.trailers["dwp-origin-state"], "")
	worktreePath, err := c.getWorkspace()
	if err != nil || worktreePath == "" {
		return err
	}

	stalledTrailers := []statex.Trailer{
		{Key: "dwp-state", Value: "stalled"},
		{Key: "dwp-stalled-run", Value: stalledRun},
	}
	if originState != "" {
		stalledTrailers = append(stalledTrailers, statex.Trailer{Key: "dwp-origin-state", Value: originState})
	}
	if c.config.UseRemote != "" {
		stalledTrailers = append(stalledTrailers, statex.Trailer{Key: "dwp-source", Value: "git:" + c.config.UseRemote})
	}
	if err := statex.CommitState(worktreePath, "chore: stalled", "Lease expired", stalledTrailers); err != nil {
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
		c.logger.Debugf("Using repository root for current branch %s", c.branchName)
		return c.config.RepoRoot, nil
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
	commandPath, err := c.findCommandPath(worktreePath, c.command)
	if err != nil || commandPath == "" {
		return "", err
	}
	return commandPath, nil
}

// commandsBase returns the checkout that .aynig commands are resolved from:
// the pinned commands checkout when configured, or the event branch worktree.
func (c *Command) commandsBase(worktreePath string) string {
	if strings.TrimSpace(c.config.CommandsRoot) != "" {
		return c.config.CommandsRoot
	}
	return worktreePath
}

func (c *Command) findCommandPath(worktreePath string, commandName string) (string, error) {
	commandsBase := c.commandsBase(worktreePath)
	roleName := strings.TrimSpace(c.config.Role)
	roleEnv := strings.TrimSpace(os.Getenv("ROLE"))
	if roleName == "" {
		roleName = roleEnv
	}
	if roleName != "" {
		rolesRoot := filepath.Join(commandsBase, ".aynig", "roles")
		roleDir, ok, err := containedPath(rolesRoot, roleName+"/command")
		if err != nil {
			return "", err
		}
		if !ok {
			c.logger.Warnf("Ignoring role %q (resolves outside %s)", roleName, rolesRoot)
		} else {
			rolePath, err := c.resolveCommandPath(roleDir, commandName)
			if err != nil {
				return "", err
			}
			if rolePath != "" {
				return rolePath, nil
			}
		}
	}
	baseDir := filepath.Join(commandsBase, ".aynig", "command")
	return c.resolveCommandPath(baseDir, commandName)
}

// containedPath joins rel (slash-separated) under baseDir and reports whether
// the result stays inside baseDir after path normalization.
func containedPath(baseDir string, rel string) (string, bool, error) {
	baseDirAbs, err := filepath.Abs(baseDir)
	if err != nil {
		return "", false, err
	}
	joined := filepath.Join(baseDirAbs, filepath.FromSlash(rel))
	if !strings.HasPrefix(joined, baseDirAbs+string(os.PathSeparator)) {
		return joined, false, nil
	}
	return joined, true, nil
}

func (c *Command) resolveCommandPath(baseDir string, commandName string) (string, error) {
	commandPath, ok, err := containedPath(baseDir, commandName)
	if err != nil {
		return "", err
	}
	c.logger.Debugf("Trying command path: %s", commandPath)
	if !ok {
		c.logger.Infof("Command path not found at %s (outside base directory %s)", commandPath, baseDir)
		return "", nil
	}
	info, err := os.Stat(commandPath)
	if err != nil {
		if os.IsNotExist(err) {
			c.logger.Infof("Command path not found at %s", commandPath)
			return "", nil
		}
		return "", nil
	}
	if info.Mode()&0o111 == 0 {
		c.logger.Infof("Command path not found at %s (not executable)", commandPath)
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
	values, ok := trailers["dwp-state"]
	if !ok || len(values) == 0 {
		return "", ""
	}
	state := strings.ToLower(strings.TrimSpace(values[len(values)-1]))
	if state == "" {
		return "", "empty dwp-state trailer"
	}
	return state, ""
}

func prepareCommandLogPaths(worktreePath string, commitHash string) (string, string, error) {
	logsDir := filepath.Join(worktreePath, ".aynig", "logs")
	if err := os.MkdirAll(logsDir, 0o755); err != nil {
		return "", "", err
	}
	stdoutLogPath := filepath.Join(logsDir, commitHash+".stdout.log")
	stderrLogPath := filepath.Join(logsDir, commitHash+".stderr.log")
	return stdoutLogPath, stderrLogPath, nil
}

func (c *Command) commandEnv(commitHash, stdoutLogPath, stderrLogPath, worktreePath string) []string {
	env := append([]string{}, os.Environ()...)
	envNames := envNameSet(env)
	env = append(env, "BODY="+c.body)
	envNames["BODY"] = struct{}{}
	env = append(env, "COMMIT_HASH="+commitHash)
	envNames["COMMIT_HASH"] = struct{}{}
	env = append(env, "STDOUT_LOG_PATH="+stdoutLogPath)
	envNames["STDOUT_LOG_PATH"] = struct{}{}
	env = append(env, "STDERR_LOG_PATH="+stderrLogPath)
	envNames["STDERR_LOG_PATH"] = struct{}{}
	env = append(env, "WORKTREE_PATH="+worktreePath)
	envNames["WORKTREE_PATH"] = struct{}{}
	env = append(env, "COMMANDS_PATH="+filepath.Join(c.commandsBase(worktreePath), ".aynig"))
	envNames["COMMANDS_PATH"] = struct{}{}
	envNames["ROLE"] = struct{}{}
	if c.logLevel != "" {
		env = append(env, "LOG_LEVEL="+c.logLevel)
		envNames["LOG_LEVEL"] = struct{}{}
	}
	if role := strings.TrimSpace(c.config.Role); role != "" {
		env = append(env, "ROLE="+role)
	}
	for key, values := range c.trailers {
		upperKey := strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
		if _, exists := envNames[upperKey]; exists {
			continue
		}
		envValue := strings.Join(values, ",")
		env = append(env, upperKey+"="+envValue)
		envNames[upperKey] = struct{}{}
	}
	return env
}

func envNameSet(env []string) map[string]struct{} {
	names := make(map[string]struct{}, len(env))
	for _, entry := range env {
		name, _, found := strings.Cut(entry, "=")
		if !found {
			continue
		}
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		names[name] = struct{}{}
	}
	return names
}
