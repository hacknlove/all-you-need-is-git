package gitx

import (
	"strings"
)

func RepoRoot() (string, error) {
	out, err := Run("", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func Fetch(dir string) error {
	_, err := Run(dir, "fetch")
	return err
}

func BranchCurrent(dir string) (string, error) {
	out, err := Run(dir, "branch", "--show-current")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func BranchListLocal(dir string) ([]string, error) {
	out, err := Run(dir, "branch", "--format=%(refname:short)")
	if err != nil {
		return nil, err
	}
	return splitLines(out), nil
}

func BranchListRemote(dir string) ([]string, error) {
	out, err := Run(dir, "branch", "-r", "--format=%(refname:short)")
	if err != nil {
		return nil, err
	}
	return splitLines(out), nil
}

func BranchUpstream(dir string, branch string) (string, error) {
	if strings.TrimSpace(branch) == "" {
		return "", nil
	}
	out, err := Run(dir, "rev-parse", "--abbrev-ref", branch+"@{upstream}")
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(out), nil
}

func WorktreeList(dir string) ([]WorktreeEntry, error) {
	out, err := Run(dir, "worktree", "list", "--porcelain")
	if err != nil {
		return nil, err
	}
	return parseWorktreeList(out), nil
}

func WorktreeAdd(dir string, args ...string) error {
	_, err := Run(dir, append([]string{"worktree", "add"}, args...)...)
	return err
}

func RevParse(dir string, arg string) (string, error) {
	out, err := Run(dir, "rev-parse", arg)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func Commit(dir string, message string, allowEmpty bool) error {
	args := []string{"commit"}
	if allowEmpty {
		args = append(args, "--allow-empty")
	}
	args = append(args, "-F", "-")
	_, err := RunWithInput(dir, message, args...)
	return err
}

func Push(dir string, remote string, branch string) error {
	_, err := Run(dir, "push", remote, branch)
	return err
}

func StatusPorcelain(dir string, pathspec string) (string, error) {
	args := []string{"status", "--porcelain"}
	if pathspec != "" {
		args = append(args, "--", pathspec)
	}
	return Run(dir, args...)
}

func DiffFile(dir string, path string) (string, error) {
	return Run(dir, "diff", "--", path)
}

func splitLines(out string) []string {
	lines := []string{}
	for _, line := range strings.Split(out, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			lines = append(lines, trimmed)
		}
	}
	return lines
}
