package debug

import (
	"net/http"

	"github.com/arl/statsviz"
)

// RegisterStatsviz for debug.
func RegisterStatsviz(mux *http.ServeMux) {
	statsviz.Register(mux)
}
