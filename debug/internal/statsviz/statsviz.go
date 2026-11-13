package statsviz

import (
	"github.com/alexfalkowski/go-service/v3/context"
	"github.com/alexfalkowski/go-service/v3/debug/http"
	"github.com/alexfalkowski/go-service/v3/env"
	"github.com/alexfalkowski/go-signal/v2"
	statsviz "github.com/arl/statsviz"
)

// Register installs statsviz handlers on mux.
//
// The handler tree is mounted at "/debug/statsviz" (namespaced by service name via [http.Pattern]), which
// serves statsviz's runtime visualization UI and endpoints.
//
// This registration is intended to be composed into the go-service debug server wiring. The underlying
// statsviz server is closed when the application lifecycle stops.
func Register(lc signal.Lifecycle, name env.Name, mux *http.ServeMux) error {
	server, err := statsviz.NewServer(statsviz.Root(http.Pattern(name, "/debug/statsviz")))
	if err != nil {
		return err
	}

	server.Register(mux.ServeMux)
	lc.Append(signal.Hook{
		OnStop: func(context.Context) error {
			return server.Close()
		},
	})

	return nil
}
