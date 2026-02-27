package os

import (
	"os"
	"slices"

	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/strings"
)

var (
	// Args is the process argument vector.
	//
	// It is an alias of os.Args and is provided so go-service code can depend on
	// go-service packages consistently while still using the underlying OS
	// implementation.
	Args = os.Args

	// Stdout is the standard output file descriptor.
	//
	// It is an alias of os.Stdout and is typically used for writing user-facing
	// output from CLIs.
	Stdout = os.Stdout
)

// Getenv returns the value of the environment variable named by key.
func Getenv(key string) string {
	return os.Getenv(key)
}

// Setenv sets the value of the environment variable named by key to value.
func Setenv(key, value string) error {
	return os.Setenv(key, value)
}

// Unsetenv unsets the environment variable named by key.
func Unsetenv(key string) error {
	return os.Unsetenv(key)
}

// Exit causes the current program to exit with the given status code.
//
// This forwards to os.Exit. Note that deferred functions are not run.
func Exit(code int) {
	os.Exit(code)
}

// Executable returns the absolute path of the running program's executable.
//
// It forwards to os.Executable and panics if the executable path cannot be
// determined (via runtime.Must).
//
// Use this helper when inability to determine the executable is considered
// unrecoverable for your service. If you need to handle the error, call
// os.Executable directly.
func Executable() string {
	path, err := os.Executable()
	runtime.Must(err)

	return path
}

// UserHomeDir returns the current user's home directory.
//
// It forwards to os.UserHomeDir and panics if the home directory cannot be
// determined (via runtime.Must).
//
// Use this helper when inability to determine the home directory is considered
// unrecoverable for your service. If you need to handle the error, call
// os.UserHomeDir directly.
func UserHomeDir() string {
	dir, err := os.UserHomeDir()
	runtime.Must(err)

	return dir
}

// UserConfigDir returns the default root directory to use for user-specific
// configuration data.
//
// It forwards to os.UserConfigDir and panics if the configuration directory
// cannot be determined (via runtime.Must).
//
// Use this helper when inability to determine the config directory is considered
// unrecoverable for your service. If you need to handle the error, call
// os.UserConfigDir directly.
func UserConfigDir() string {
	dir, err := os.UserConfigDir()
	runtime.Must(err)

	return dir
}

// SanitizeArgs removes Go test runner flags from args.
//
// When tests execute a package as a binary, the go test harness injects flags of
// the form "-test.*" (for example "-test.v"). If you reuse an argv slice (for
// example passing it through to another command parser), those flags can be
// surprising and may break flag parsing.
//
// SanitizeArgs returns a copy of args with any element that has the "-test"
// prefix removed. Relative ordering of the remaining arguments is preserved.
func SanitizeArgs(args []string) []string {
	return slices.DeleteFunc(args, func(s string) bool { return strings.HasPrefix(s, "-test") })
}
