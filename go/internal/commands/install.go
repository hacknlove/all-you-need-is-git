package commands

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"all-you-need-is-git/go/internal/gitx"
)

func Install(repo string, ref string, subfolder string) error {
	repo = normalizeRepo(repo)
	tmpDir, err := os.MkdirTemp("", "aynig-install-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	cloneArgs := []string{"clone", "--depth", "1"}
	if ref != "" {
		cloneArgs = append(cloneArgs, "--branch", ref)
		fmt.Printf("Cloning %s at %s...\n", repo, ref)
	} else {
		fmt.Printf("Cloning %s (default branch)...\n", repo)
	}
	cloneArgs = append(cloneArgs, repo, tmpDir)
	if _, err := gitx.Run("", cloneArgs...); err != nil {
		return err
	}

	sourcePath := filepath.Join(tmpDir, ".aynig")
	if subfolder != "" {
		sourcePath = filepath.Join(tmpDir, subfolder)
	}
	if _, err := os.Stat(sourcePath); err != nil {
		refInfo := ""
		if ref != "" {
			refInfo = " at " + ref
		}
		subfolderInfo := ""
		if subfolder != "" {
			subfolderInfo = " in " + subfolder
		}
		return fmt.Errorf("Error: No .aynig directory found in %s%s%s", repo, refInfo, subfolderInfo)
	}

	destPath := ".aynig"
	statusBefore, err := gitx.StatusPorcelain("", destPath)
	if err != nil {
		return err
	}
	if strings.TrimSpace(statusBefore) != "" {
		return fmt.Errorf("\nError: .aynig has uncommitted changes. Please commit or stash them first.")
	}

	fmt.Println("Copying workflows...")
	if err := copyDir(sourcePath, destPath); err != nil {
		return err
	}

	statusAfter, err := gitx.StatusPorcelain("", destPath)
	if err != nil {
		return err
	}
	modified, untracked := parseStatus(statusAfter)

	overwritten := filterFiles(modified, ".aynig/COMMANDS.md")
	newFiles := filterFiles(untracked, ".aynig/COMMANDS.md")

	if len(newFiles) > 0 {
		fmt.Printf("✓ Installed new: %s\n", strings.Join(newFiles, ", "))
	}
	if len(overwritten) > 0 {
		fmt.Printf("⚠  Overwritten: %s\n", strings.Join(overwritten, ", "))
	}

	commandsMdModified := contains(modified, ".aynig/COMMANDS.md")
	if commandsMdModified && len(overwritten) > 0 {
		diff, err := gitx.DiffFile("", ".aynig/COMMANDS.md")
		if err != nil {
			return err
		}
		hasOpencode := commandExists("opencode")
		hasClaude := commandExists("claude")

		if hasOpencode || hasClaude {
			tool := "claude"
			if hasOpencode {
				tool = "opencode"
			}
			fmt.Printf("\n🤖 Calling %s to merge COMMANDS.md...\n", tool)
			prompt := "Some commands were just installed/overwritten. Please merge this COMMANDS.md file intelligently by looking at the git diff below. Keep the most accurate and valid description for each command, removing any duplicates.\n\nGit diff:\n" + diff
			destCommandsPath := filepath.Join(destPath, "COMMANDS.md")
			var cmd *exec.Cmd
			if tool == "opencode" {
				cmd = exec.Command(tool, destCommandsPath, "-p", prompt)
			} else {
				cmd = exec.Command(tool, "-p", prompt, destCommandsPath)
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Printf("⚠  Error calling %s: %v\n", tool, err)
				fmt.Println("⚠  Please manually review and merge COMMANDS.md")
			} else {
				fmt.Printf("✓ COMMANDS.md merged by %s\n", tool)
			}
		} else {
			fmt.Println("\n⚠  Warning: COMMANDS.md was modified. Please manually review the changes.")
			preview := diffPreview(diff, 20)
			fmt.Printf("\nGit diff preview:\n%s\n", preview)
			if len(strings.Split(diff, "\n")) > 20 {
				fmt.Println("... (run `git diff .aynig/COMMANDS.md` to see full diff)")
			}
		}
	} else if commandsMdModified {
		fmt.Println("✓ COMMANDS.md updated")
	}

	fmt.Println("\n✓ Workflows installed successfully!")
	return nil
}

func normalizeRepo(repo string) string {
	if isFullRepoURL(repo) {
		return repo
	}
	if isShorthandRepo(repo) {
		return "https://github.com/" + repo + ".git"
	}
	return repo
}

func isFullRepoURL(repo string) bool {
	return strings.HasPrefix(repo, "http://") ||
		strings.HasPrefix(repo, "https://") ||
		strings.HasPrefix(repo, "git@") ||
		strings.HasPrefix(repo, "ssh://")
}

func isShorthandRepo(repo string) bool {
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return false
	}
	return isRepoSlug(parts[0]) && isRepoSlug(parts[1])
}

func isRepoSlug(part string) bool {
	if part == "" {
		return false
	}
	for _, r := range part {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == '.' {
			continue
		}
		return false
	}
	return true
}

func copyDir(src string, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dest, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		return copyFileMode(path, target, info.Mode())
	})
}

func copyFileMode(src string, dest string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func parseStatus(status string) ([]string, []string) {
	modified := []string{}
	untracked := []string{}
	for _, line := range strings.Split(status, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.HasPrefix(line, "?? ") {
			untracked = append(untracked, strings.TrimSpace(line[3:]))
			continue
		}
		if len(line) >= 4 {
			path := strings.TrimSpace(line[3:])
			modified = append(modified, path)
		}
	}
	return modified, untracked
}

func filterFiles(files []string, excluded string) []string {
	out := []string{}
	for _, file := range files {
		if file != excluded {
			out = append(out, filepath.Base(file))
		}
	}
	return out
}

func contains(files []string, target string) bool {
	for _, file := range files {
		if file == target {
			return true
		}
	}
	return false
}

func commandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func diffPreview(diff string, maxLines int) string {
	lines := strings.Split(diff, "\n")
	if len(lines) <= maxLines {
		return diff
	}
	return strings.Join(lines[:maxLines], "\n")
}
