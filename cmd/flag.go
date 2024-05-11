package cmd

import (
	"strings"
)

// SplitFlag for cmd.
func SplitFlag(flag string) (string, string) {
	c := strings.Split(flag, ":")

	if len(c) != 2 {
		return "none", "not_used"
	}

	return c[0], c[1]
}
