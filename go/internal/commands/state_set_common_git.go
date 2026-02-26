package commands

import "all-you-need-is-git/go/internal/gitx"

func gitShowHeadBody() (string, error) {
	return gitx.Run("", "show", "-s", "--format=%B", "HEAD")
}
