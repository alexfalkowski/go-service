package os

import (
	"os"
	"slices"

	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/strings"
)

var (
	// Args is an alias for os.Args.
	Args = os.Args

	// Stdout is an alias for os.Stdout.
	Stdout = os.Stdout
)

// Getenv is an alias of os.Getenv.
func Getenv(key string) string {
	return os.Getenv(key)
}

// Setenv is an alias of os.Setenv.
func Setenv(key, value string) error {
	return os.Setenv(key, value)
}

// Unsetenv is an alias of os.Unsetenv.
func Unsetenv(key string) error {
	return os.Unsetenv(key)
}

// Exit is an alias for os.Exit.
func Exit(code int) {
	os.Exit(code)
}

// Executable returns the absolute path of the running application executable.
//
// It panics if the executable path cannot be determined.
func Executable() string {
	path, err := os.Executable()
	runtime.Must(err)

	return path
}

// UserHomeDir returns the home directory of the current user.
//
// It panics if the home directory cannot be determined.
func UserHomeDir() string {
	dir, err := os.UserHomeDir()
	runtime.Must(err)

	return dir
}

// UserConfigDir returns the config directory of the current user.
//
// It panics if the config directory cannot be determined.
func UserConfigDir() string {
	dir, err := os.UserConfigDir()
	runtime.Must(err)

	return dir
}

// SanitizeArgs removes testing flags (those starting with "-test") from args.
//
// This is useful when reusing a binary's arguments while avoiding go test runner flags.
func SanitizeArgs(args []string) []string {
	return slices.DeleteFunc(args, func(s string) bool { return strings.HasPrefix(s, "-test") })
}
