package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/pkg/meta"
	"github.com/google/uuid"
)

type metaRoundTripper struct {
	http.RoundTripper
}

func (r *metaRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	requestID := meta.Attribute(ctx, meta.RequestID)
	if requestID == "" {
		requestID = uuid.New().String()
	}

	req.Header.Set("Request-ID", requestID)

	ctx = meta.WithAttribute(ctx, meta.RequestID, requestID)

	return r.RoundTripper.RoundTrip(req.WithContext(ctx))
}
