package runtime

import "go.uber.org/automaxprocs/maxprocs"

// RegisterMaxProcs for runtime.
func RegisterMaxProcs() error {
	_, err := maxprocs.Set()

	return err
}
