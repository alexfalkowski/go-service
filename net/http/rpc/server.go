package rpc

import (
	"net/http"

	"github.com/alexfalkowski/go-service/net/http/content"
)

// Route for rpc.
func Route[Req any, Res any](path string, handler content.RequestHandler[Req, Res]) {
	pattern := http.MethodPost + " " + path

	mux.HandleFunc(pattern, content.NewRequestHandler(cont, "rpc", handler))
}
