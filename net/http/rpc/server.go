package rpc

import (
	"github.com/alexfalkowski/go-service/net/http/content"
)

// Route for rpc.
func Route[Req any, Res any](path string, handler content.RequestHandler[Req, Res]) {
	mux.HandleFunc("POST "+path, content.NewRequestHandler(cont, "rpc", handler))
}
