package http

import (
	"context"
)

// Client for HTTP.
type Client[Req any, Res any] interface {
	// Call with request.
	Call(ctx context.Context, req *Req) (*Res, error)
}
