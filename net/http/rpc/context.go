package rpc

import (
	"context"
	"net/http"
)

type (
	// Context for HTTP.
	Context interface {
		// Request of the context.
		Request() *http.Request

		// Response of the context.
		Response() http.ResponseWriter

		context.Context
	}

	//nolint:containedctx
	handlerContext struct {
		req *http.Request
		res http.ResponseWriter

		context.Context
	}
)

// NewContext for HTTP.
func NewContext(ctx context.Context, req *http.Request, res http.ResponseWriter) Context {
	return &handlerContext{req: req, res: res, Context: ctx}
}

func (c *handlerContext) Request() *http.Request {
	return c.req
}

func (c *handlerContext) Response() http.ResponseWriter {
	return c.res
}
