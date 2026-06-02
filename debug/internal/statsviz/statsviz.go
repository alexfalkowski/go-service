package statsviz

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/debug/http"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/env"
	statsviz "github.com/arl/statsviz"
)

// Register installs statsviz handlers on mux.
//
// The handler tree is mounted at "/debug/statsviz" (namespaced by service name via [http.Pattern]), which
// serves statsviz's runtime visualization UI and endpoints.
//
// This registration is intended to be composed into the go-service debug server wiring. The underlying
// statsviz server is closed when the application lifecycle stops.
func Register(lc di.Lifecycle, name env.Name, mux *http.ServeMux) error {
	server, err := statsviz.NewServer(statsviz.Root(http.Pattern(name, "/debug/statsviz")))
	if err != nil {
		return err
	}

	server.Register(mux.ServeMux)
	lc.Append(di.Hook{
		OnStop: func(context.Context) error {
			return server.Close()
		},
	})

	return nil
}
