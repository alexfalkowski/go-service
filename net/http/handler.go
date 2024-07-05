package http

// Handler for HTTP.
type Handler[Req any, Res any] interface {
	// Handle the request/response.
	Handle(ctx Context, req *Req) (*Res, error)
}
