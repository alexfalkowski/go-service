package os

import (
	"os"
	"strings"
)

// GetFromEnv parses the value in the format of env:VARIABLE and retrieves it, or just returns the value.
func GetFromEnv(value string) string {
	s := strings.Split(value, ":")

	if len(s) != 2 || s[0] != "env" {
		return value
	}

	return os.Getenv(s[1])
}
