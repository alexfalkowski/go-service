package runtime

import (
	"go.uber.org/automaxprocs/maxprocs"
)

// Register runtime.
func Register() {
	maxprocs.Set()
}
