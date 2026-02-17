package statsviz

import (
	"github.com/alexfalkowski/go-service/v2/debug/http"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/arl/statsviz"
)

// Register installs statsviz handlers under /debug/statsviz.
func Register(name env.Name, mux *http.ServeMux) error {
	return statsviz.Register(mux.ServeMux, statsviz.Root(http.Pattern(name, "/debug/statsviz")))
}
