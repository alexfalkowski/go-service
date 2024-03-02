package runtime

import (
	"go.uber.org/automaxprocs/maxprocs"
)

// RegisterMaxProcs for runtime.
func RegisterMaxProcs() {
	maxprocs.Set()
}
