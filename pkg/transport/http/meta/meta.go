package meta

import (
	"net/http"

	"github.com/alexfalkowski/go-service/pkg/transport/meta"
	"github.com/google/uuid"
)

// NewRoundTripper for meta.
func NewRoundTripper(userAgent string, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{userAgent: userAgent, RoundTripper: hrt}
}

// RoundTripper for meta.
type RoundTripper struct {
	userAgent string
	http.RoundTripper
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	requestID := meta.RequestID(ctx)
	if requestID == "" {
		requestID = uuid.New().String()
	}

	req.Header.Set("Request-ID", requestID)
	ctx = meta.WithRequestID(ctx, requestID)

	req.Header.Set("User-Agent", r.userAgent)
	ctx = meta.WithUserAgent(ctx, r.userAgent)

	return r.RoundTripper.RoundTrip(req.WithContext(ctx))
}
