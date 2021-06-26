package meta

import (
	"net/http"

	"github.com/alexfalkowski/go-service/pkg/transport/meta"
	"github.com/google/uuid"
)

// NewRoundTripper for meta.
func NewRoundTripper(hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{RoundTripper: hrt}
}

// RoundTripper for meta.
type RoundTripper struct {
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

	return r.RoundTripper.RoundTrip(req.WithContext(ctx))
}
