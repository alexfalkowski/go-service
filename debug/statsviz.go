package debug

import "github.com/arl/statsviz"

// RegisterStatsviz for debug.
func RegisterStatsviz(srv *Server) error {
	if srv == nil {
		return nil
	}

	mux := srv.ServeMux()

	return statsviz.Register(mux)
}
