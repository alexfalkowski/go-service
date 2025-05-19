package os

import "os"

// ExitFunc defines a way to exit.
type ExitFunc = func(code int)

// NewExitFunc for os.
func NewExitFunc() ExitFunc {
	return os.Exit
}
