package debug

import (
	"github.com/arl/statsviz"
)

// RegisterStatsviz for debug.
func RegisterStatsviz(server *Server) {
	statsviz.Register(server.ServeMux())
}
