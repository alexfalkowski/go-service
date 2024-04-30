package cmd

import (
	"strings"
)

// SplitFlag for cmd.
func SplitFlag(flag string) (string, string) {
	c := strings.Split(flag, ":")

	if len(c) != 2 {
		return "env", "CONFIG_FILE"
	}

	return c[0], c[1]
}

// NewReadWriter for cmd.
func NewReadWriter(kind, location string) ReaderWriter {
	if kind == "file" {
		return NewFile(location)
	}

	return NewENV(location)
}
