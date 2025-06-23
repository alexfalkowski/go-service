package os

import (
	"os"
	"slices"

	"github.com/alexfalkowski/go-service/v2/strings"
)

var (
	// Args is an alias for os.Args.
	Args = os.Args

	// Exit is an alias for os.Exit.
	Exit = os.Exit

	// Stdout is an alias for os.Stdout.
	Stdout = os.Stdout

	// Getenv is an alias of os.Getenv.
	Getenv = os.Getenv

	// Setenv is an alias of os.Setenv.
	Setenv = os.Setenv

	// Unsetenv is an alias of os.Unsetenv.
	Unsetenv = os.Unsetenv
)

// Executable of the running application.
func Executable() string {
	path, _ := os.Executable()

	return path
}

// UserHomeDir of the current user.
func UserHomeDir() string {
	home, _ := os.UserHomeDir()

	return home
}

// SanitizeArgs removes all flags that start with -test.
func SanitizeArgs(args []string) []string {
	return slices.DeleteFunc(args, func(s string) bool {
		return strings.HasPrefix(s, "-test")
	})
}
