package rpc

import (
	"github.com/alexfalkowski/go-service/net/http/content"
)

// Route for rpc.
func Route[Req any, Res any](path string, handler content.Handler[Req, Res]) {
	mux.HandleFunc("POST "+path, content.NewBodyHandler(cont, "rpc", handler))
}
