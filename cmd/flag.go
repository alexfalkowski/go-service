package cmd

import (
	"strings"
)

// SplitFlag for cmd.
func SplitFlag(flag string) (string, string) {
	kind, name, ok := strings.Cut(flag, ":")
	if !ok {
		return "none", "not_used"
	}

	return kind, name
}
