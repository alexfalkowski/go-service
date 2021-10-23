package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/pkg/transport/http/breaker"
	pkgZap "github.com/alexfalkowski/go-service/pkg/transport/http/logger/zap"
	"github.com/alexfalkowski/go-service/pkg/transport/http/meta"
	"github.com/alexfalkowski/go-service/pkg/transport/http/retry"
	"github.com/alexfalkowski/go-service/pkg/transport/http/trace/opentracing"
	"go.uber.org/zap"
)

// Transport for http.Client.
var Transport = http.DefaultTransport

// NewClient for HTTP.
func NewClient(logger *zap.Logger) *http.Client {
	return NewClientWithRoundTripper(logger, Transport)
}

// NewClientWithRoundTripper for HTTP.
func NewClientWithRoundTripper(logger *zap.Logger, hrt http.RoundTripper) *http.Client {
	return &http.Client{Transport: newRoundTripper(logger, hrt)}
}

func newRoundTripper(logger *zap.Logger, hrt http.RoundTripper) http.RoundTripper {
	hrt = pkgZap.NewRoundTripper(logger, hrt)
	hrt = opentracing.NewRoundTripper(hrt)
	hrt = retry.NewRoundTripper(hrt)
	hrt = breaker.NewRoundTripper(hrt)
	hrt = meta.NewRoundTripper(hrt)

	return hrt
}
