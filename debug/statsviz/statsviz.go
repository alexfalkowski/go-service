package statsviz

import (
	"github.com/alexfalkowski/go-service/v2/debug/http"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/arl/statsviz"
)

// Register for debug.
func Register(name env.Name, mux *http.ServeMux) error {
	return statsviz.Register(mux.ServeMux, statsviz.Root(http.Pattern("/debug/statsviz", name)))
}
