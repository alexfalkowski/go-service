package debug

import "github.com/arl/statsviz"

// RegisterStatsviz for debug.
func RegisterStatsviz(mux *ServeMux) error {
	return statsviz.Register(mux.ServeMux)
}
