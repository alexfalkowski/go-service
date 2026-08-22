package http

import (
	"net/http"
	"net/http/httptrace"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/http/telemetry"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/time"
)

// ClosingRoundTripper adapts a function to RoundTripper while making request-body ownership explicit.
//
// The function's third return value reports whether ClosingRoundTripper should close req.Body before returning.
// Return true when the function rejects the request locally without delegating to another RoundTripper.
// Return false after delegating because the delegated RoundTripper owns the request body.
type ClosingRoundTripper func(req *Request) (*Response, error, bool)

// RoundTrip calls s and closes req.Body when s asks it to.
func (s ClosingRoundTripper) RoundTrip(req *Request) (*Response, error) {
	res, err, closeBody := s(req)
	if closeBody && req != nil && req.Body != nil && req.Body != NoBody {
		_ = req.Body.Close()
	}

	return res, err
}

// NewClient constructs an HTTP client with a request timeout.
//
// When tracing or metrics are enabled, the returned client wraps the provided RoundTripper with a telemetry
// transport. When tracing is enabled, it also installs an httptrace-based client trace derived from the
// request context.
//
// The provided timeout is assigned to [http.Client.Timeout] (total time limit for requests, including
// connection time, redirects, and reading the response body).
func NewClient(rt http.RoundTripper, timeout time.Duration) *http.Client {
	var transport http.RoundTripper

	if metrics.IsEnabled() || tracer.IsEnabled() {
		options := []telemetry.Option{}
		if tracer.IsEnabled() {
			options = append(options, telemetry.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
				return telemetry.NewClientTrace(ctx)
			}))
		}

		transport = telemetry.NewTransport(rt, options...)
	} else {
		transport = rt
	}

	return &http.Client{
		Transport: transport,
		Timeout:   timeout.Duration(),
	}
}
