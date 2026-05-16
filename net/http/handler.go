package http

import (
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-sync"
	snoop "github.com/felixge/httpsnoop"
)

// NotFoundHandler handles an unmatched mux request.
//
// It returns true when it has written a response. Returning false lets the next handler decide.
type NotFoundHandler func(res ResponseWriter, req *Request) bool

// NewNotFoundHandler wraps mux and lets handlers replace generated 404 responses.
//
// Non-404 mux responses are flushed unchanged. This preserves method handling such as 405 Method Not
// Allowed while still allowing MVC, REST, or RPC packages to provide typed not-found responses.
func NewNotFoundHandler(mux *ServeMux, pool *sync.BufferPool, handlers ...NotFoundHandler) Handler {
	if len(handlers) == 0 {
		return mux
	}

	return &notFoundHandler{mux: mux, pool: pool, handlers: handlers}
}

type notFoundHandler struct {
	mux      *ServeMux
	pool     *sync.BufferPool
	handlers []NotFoundHandler
}

func (h *notFoundHandler) ServeHTTP(res ResponseWriter, req *Request) {
	handler, pattern := h.mux.Handler(req)
	if !strings.IsEmpty(pattern) {
		h.mux.ServeHTTP(res, req)
		return
	}

	response := &bufferedWriter{
		header:   Header{},
		response: res,
		buffer:   h.pool.Get(),
		pool:     h.pool,
	}
	metrics := snoop.CaptureMetricsFn(response, func(res ResponseWriter) {
		handler.ServeHTTP(res, req)
	})
	defer response.Close()

	if metrics.Code != StatusNotFound {
		response.Flush()
		return
	}

	for _, handler := range h.handlers {
		if handler(res, req) {
			return
		}
	}

	response.Flush()
}
