package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"all-you-need-is-git/go/internal/gitx"
)

func Init() error {
	if _, err := gitx.Run("", "status"); err != nil {
		return fmt.Errorf("Error: Not a Git repository. Please run `git init` first.")
	}

	aynigDir := ".aynig"
	aynigCreated := false
	if err := os.Mkdir(aynigDir, 0o755); err != nil {
		if !os.IsExist(err) {
			return err
		}
		fmt.Printf("⊘ %s/ already exists, skipping\n", aynigDir)
	} else {
		fmt.Printf("✓ Created %s/\n", aynigDir)
		aynigCreated = true
	}

	commandDir := filepath.Join(aynigDir, "command")
	if err := os.Mkdir(commandDir, 0o755); err != nil {
		if !os.IsExist(err) {
			return err
		}
		fmt.Printf("⊘ %s/ already exists, skipping\n", commandDir)
	} else {
		fmt.Printf("✓ Created %s/\n", commandDir)
	}

	if aynigCreated {
		if err := os.WriteFile(filepath.Join(aynigDir, "COMMANDS.md"), []byte(commandsMd), 0o644); err != nil {
			return err
		}
		fmt.Println("✓ Created COMMANDS.md")
	}

	cleanPath := filepath.Join(commandDir, "clean")
	if _, err := os.Stat(cleanPath); err == nil {
		fmt.Println("⊘ clean command already exists, skipping")
	} else if !os.IsNotExist(err) {
		return err
	} else {
		if err := os.WriteFile(cleanPath, []byte(cleanScript), 0o755); err != nil {
			return err
		}
		fmt.Println("✓ Created clean command")
	}

	worktreesDir := ".worktrees"
	if err := os.Mkdir(worktreesDir, 0o755); err != nil {
		if !os.IsExist(err) {
			return err
		}
		fmt.Printf("⊘ %s/ already exists, skipping\n", worktreesDir)
	} else {
		fmt.Printf("✓ Created %s/\n", worktreesDir)
	}

	gitignorePath := ".gitignore"
	worktreesEntry := ".worktrees/"
	gitignoreContent, err := os.ReadFile(gitignorePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Error updating .gitignore: %w", err)
	}
	if stringsContainsLine(string(gitignoreContent), worktreesEntry) {
		fmt.Printf("⊘ %s already in .gitignore, skipping\n", worktreesEntry)
	} else {
		newContent := string(gitignoreContent)
		if newContent != "" && newContent[len(newContent)-1] != '\n' {
			newContent += "\n"
		}
		newContent += worktreesEntry + "\n"
		if err := os.WriteFile(gitignorePath, []byte(newContent), 0o644); err != nil {
			return fmt.Errorf("Error updating .gitignore: %w", err)
		}
		fmt.Printf("✓ Added %s to .gitignore\n", worktreesEntry)
	}

	fmt.Println("\nAYNIG initialized successfully!")
	return nil
}

func stringsContainsLine(content string, target string) bool {
	for _, line := range strings.Split(content, "\n") {
		if strings.TrimSpace(line) == target {
			return true
		}
	}
	return false
}
