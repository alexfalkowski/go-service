package statsviz

import (
	"github.com/alexfalkowski/go-service/v2/debug/http"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/arl/statsviz"
)

// Register installs statsviz handlers on mux.
//
// The handler tree is mounted at "/debug/statsviz" (namespaced by service name via http.Pattern), which
// serves statsviz's runtime visualization UI and endpoints.
//
// This registration is intended to be composed into the go-service debug server wiring. Any error
// returned by the underlying statsviz.Register call is returned to the caller.
func Register(name env.Name, mux *http.ServeMux) error {
	return statsviz.Register(mux.ServeMux, statsviz.Root(http.Pattern(name, "/debug/statsviz")))
}
