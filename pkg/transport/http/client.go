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

// ClientParams for HTTP.
type ClientParams struct {
	Logger       *zap.Logger
	RoundTripper http.RoundTripper
}

// Transport for http.Client.
var Transport = http.DefaultTransport

// NewClient for HTTP.
func NewClient(params *ClientParams) *http.Client {
	hrt := params.RoundTripper
	if hrt == nil {
		hrt = http.DefaultTransport
	}

	return &http.Client{Transport: newRoundTripper(params.Logger, hrt)}
}

func newRoundTripper(logger *zap.Logger, hrt http.RoundTripper) http.RoundTripper {
	hrt = pkgZap.NewRoundTripper(logger, hrt)
	hrt = opentracing.NewRoundTripper(hrt)
	hrt = retry.NewRoundTripper(hrt)
	hrt = breaker.NewRoundTripper(hrt)
	hrt = meta.NewRoundTripper(hrt)

	return hrt
}
