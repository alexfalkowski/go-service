package statsviz

import (
	"github.com/alexfalkowski/go-service/v2/debug/http"
	"github.com/arl/statsviz"
)

// Register for debug.
func Register(mux *http.ServeMux) error {
	return statsviz.Register(mux.ServeMux)
}
