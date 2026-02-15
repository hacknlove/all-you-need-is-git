package gitx

import "strings"

type WorktreeEntry struct {
	Path   string
	Branch string
}

func parseWorktreeList(output string) []WorktreeEntry {
	entries := []WorktreeEntry{}
	current := WorktreeEntry{}
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			if current.Path != "" {
				entries = append(entries, current)
			}
			current = WorktreeEntry{}
			continue
		}
		if strings.HasPrefix(line, "worktree ") {
			current.Path = strings.TrimPrefix(line, "worktree ")
		} else if strings.HasPrefix(line, "branch ") {
			current.Branch = strings.TrimPrefix(line, "branch ")
		}
	}
	if current.Path != "" {
		entries = append(entries, current)
	}
	return entries
}
