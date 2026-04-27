package body

import (
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
)

// NewHandler constructs request body size limiting middleware.
func NewHandler(limit int64) *Handler {
	return &Handler{limit: limit}
}

// Handler limits request bodies before downstream handlers decode them.
type Handler struct {
	limit int64
}

// ServeHTTP enforces the configured request body limit.
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if req.ContentLength > h.limit {
		_ = status.WriteError(res, &http.MaxBytesError{Limit: h.limit})
		return
	}

	data, body, err := io.ReadAll(io.LimitReader(req.Body, h.limit+1))
	if err != nil {
		_ = status.WriteError(res, status.BadRequestError(err))
		return
	}
	defer req.Body.Close()

	if int64(len(data)) > h.limit {
		_ = status.WriteError(res, &http.MaxBytesError{Limit: h.limit})
		return
	}

	req.Body = body
	next(res, req)
}
