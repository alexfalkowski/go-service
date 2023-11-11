package debug

import (
	"net/http"

	"github.com/alexfalkowski/go-service/env"
	sh "github.com/alexfalkowski/go-service/transport/http"
	"github.com/arl/statsviz"
)

// This is not exposed by statsviz.
const path = "/debug/statsviz"

// Register debug.
func Register(env env.Environment, server *sh.Server) error {
	if !env.IsDevelopment() {
		return nil
	}

	srv, err := statsviz.NewServer(statsviz.Root(path))
	if err != nil {
		return err
	}

	idx := srv.Index()
	index := func(w http.ResponseWriter, r *http.Request, p map[string]string) {
		idx.ServeHTTP(w, r)
	}

	if err := server.Mux.HandlePath("GET", path, index); err != nil {
		return err
	}

	ws := srv.Ws()
	websocket := func(w http.ResponseWriter, r *http.Request, p map[string]string) {
		ws.ServeHTTP(w, r)
	}

	if err := server.Mux.HandlePath("GET", path+"/ws", websocket); err != nil {
		return err
	}

	return nil
}
