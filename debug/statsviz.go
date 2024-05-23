package debug

import (
	"github.com/arl/statsviz"
)

// RegisterStatsviz for debug.
func RegisterStatsviz(srv *Server) {
	mux := srv.ServeMux()

	statsviz.Register(mux)
}
