package rpc

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// Route registers an RPC-style HTTP POST handler under pattern.
//
// The effective route pattern passed to the router is method-qualified and has the form:
//
//	"<METHOD> <pattern>"
//
// For example:
//
//	Route("/greet.v1.Greeter/SayHello", handler) // registers "POST /greet.v1.Greeter/SayHello"
//
// The handler is constructed using [github.com/alexfalkowski/go-service/v2/net/http/content.NewRequestHandler], which:
//   - decodes the request body using Content-Type, falling back to JSON when Content-Type is absent or unknown,
//   - decodes the request body into a newly allocated request model, and
//   - encodes the returned response model using Accept, falling back to Content-Type when Accept is absent.
//
// Registration:
// The resulting handler is registered on the package-level router configured via [Register].
// [Register] must be called before Route; otherwise router/cont will be nil and this function will panic.
func Route[Req any, Res any](pattern string, handler content.RequestHandler[Req, Res]) {
	router.Handle(
		strings.Join(strings.Space, http.MethodPost, pattern),
		content.NewRequestHandler(cont, handler),
	)
}
