package body

import (
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
)

// NewHandler wraps handler with request body size enforcement.
//
// It rejects requests whose body exceeds limit before calling handler. For
// accepted requests, it buffers and restores req.Body so downstream handlers can
// read it normally.
func NewHandler(handler http.Handler, limit int64) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.ContentLength > limit {
			_ = status.WriteError(res, &http.MaxBytesError{Limit: limit})
			return
		}

		if req.Body == nil || req.Body == http.NoBody {
			handler.ServeHTTP(res, req)
			return
		}

		data, body, err := io.ReadAll(io.LimitReader(req.Body, limit+1))
		if err != nil {
			_ = status.WriteError(res, status.BadRequestError(err))
			return
		}
		defer req.Body.Close()

		if int64(len(data)) > limit {
			_ = status.WriteError(res, &http.MaxBytesError{Limit: limit})
			return
		}

		req.Body = body
		handler.ServeHTTP(res, req)
	})
}
