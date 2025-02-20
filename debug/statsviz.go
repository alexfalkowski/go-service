package debug

import "github.com/arl/statsviz"

// RegisterStatsviz for debug.
func RegisterStatsviz(srv *Server) error {
	mux := srv.ServeMux()

	return statsviz.Register(mux)
}
