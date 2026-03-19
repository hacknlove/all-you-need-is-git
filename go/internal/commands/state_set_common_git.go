package commands

import (
	"strings"

	"all-you-need-is-git/go/internal/gitx"
)

func gitShowHeadBody() (string, error) {
	output, err := gitx.Run("", "show", "-s", "--format=%B", "HEAD")
	if err != nil && strings.Contains(err.Error(), "unknown revision") {
		return "", nil
	}
	return output, err
}
