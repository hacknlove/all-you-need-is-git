package commands

import (
	"fmt"
	"os"
	"os/exec"
)

const installCommand = "curl -fsSL https://aynig.org/install.sh | bash"

func Update() error {
	cmd := exec.Command("sh", "-c", installCommand)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	return nil
}
